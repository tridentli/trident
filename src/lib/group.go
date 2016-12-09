package trident

import (
	"errors"
	"time"
	"trident.li/keyval"
	pf "trident.li/pitchfork/lib"
)

type TriGroup interface {
	pf.PfGroup
	Add_default_attestations(ctx pf.PfCtx) (err error)
	GetVouch_adminonly() bool
	GetAttestations() (output []TriGroupAttestation, err error)
	GetAttestationsKVS() (kvs keyval.KeyVals, err error)
}

type TriGroupS struct {
	pf.PfGroup
	Please_vouch    bool   `label:"Please Vouch" pfset:"group_admin" hint:"Members must vouch before becoming active"`
	Vouch_adminonly bool   `label:"Vouch group admins only" pfset:"group_admin" hint:"Only adminstators may Vouvh"`
	Min_invouch     int    `label:"Minimum Inbound Vouches" pfset:"group_admin" hint:"Number of incoming vouches required to vett."`
	Min_outvouch    int    `label:"Minimum Outbound Vouches" pfset:"group_admin" hint:"Number of outgoing vouches required"`
	Target_invouch  int    `label:"Target Invouches" pfset:"group_admin"`
	Max_inactivity  string `label:"Maximum Inactivity" pfset:"group_admin" coalesce:"30 days"`
	Can_time_out    bool   `label:"Can Time Out" pfset:"group_admin"`
	Max_vouchdays   int    `label:"Maximum Vouch Days" pfset:"group_admin"`
	Idle_guard      string `label:"Idle Guard" pfset:"group_admin" coalesce:"30 days"`
	Nom_enabled     bool   `label:"Nominations Enabled" pfset:"group_admin"`
}

/* Don't call directly, use ctx.NewGroup() */
func NewTriGroup() pf.PfGroup {
	return &TriGroupS{PfGroup: pf.NewPfGroup()}
}

func (grp *TriGroupS) fetch(group_name string, nook bool) (err error) {
	/* Make sure the name is mostly sane */
	group_name, err = pf.Chk_ident("Group Name", group_name)
	if err != nil {
		return
	}

	p := []string{"ident"}
	v := []string{group_name}
	err = pf.StructFetchA(grp, "trustgroup", "", p, v, "", true)
	if nook && err == pf.ErrNoRows {
		/* No rows is sometimes okay */
	} else if err != nil {
		pf.Log("Group:fetch() " + err.Error() + " '" + group_name + "'")
	}

	return
}

func (grp *TriGroupS) Select(ctx pf.PfCtx, group_name string, perms pf.Perm) (err error) {
	err = grp.fetch(group_name, false)
	if err != nil {
		/* Failed to fetch */
		return
	}

	return grp.PfGroup.Select(ctx, group_name, perms)
}

func (grp *TriGroupS) Exists(group_name string) (exists bool) {
	err := grp.fetch(group_name, true)
	if err == pf.ErrNoRows {
		return false
	}

	return true
}

func (grp *TriGroupS) Refresh() (err error) {
	return grp.fetch(grp.GetGroupName(), false)
}

func (grp *TriGroupS) GetVouch_adminonly() bool {
	return grp.Vouch_adminonly
}

func (grp *TriGroupS) ListGroupMembers(search string, username string, offset int, max int, nominated bool, inclhidden bool, exact bool) (members []pf.PfGroupMember, err error) {
	var rows *pf.Rows

	grpname := grp.GetGroupName()

	members = nil

	ord := "ORDER BY m.descr"

	m := pf.NewPfGroupMember()

	q := m.SQL_Selects() + ", " +
		"COALESCE(for_vouches.num, 0) AS vouches_for, " +
		"COALESCE(for_me_vouches.num, 0) AS vouches_for_me, " +
		"COALESCE(by_vouches.num, 0) AS vouches_by, " +
		"COALESCE(by_me_vouches.num, 0) AS vouches_by_me " +
		m.SQL_Froms() + " " +
		"LEFT OUTER JOIN ( " +
		"  SELECT 'for' AS dir, mv.vouchee AS member, COUNT(*) AS num " +
		"  FROM member_vouch mv " +
		"  WHERE mv.trustgroup = $1 " +
		"  AND mv.positive " +
		"  GROUP BY mv.vouchee " +
		") as for_vouches on (m.ident = for_vouches.member) " +
		"LEFT OUTER JOIN ( " +
		"  SELECT 'by' AS dir, mv.vouchor AS member, COUNT(*) AS num " +
		"  FROM member_vouch mv " +
		"  WHERE mv.trustgroup = $1 " +
		"  AND mv.positive " +
		"  GROUP BY mv.vouchor " +
		") as by_vouches on (m.ident = by_vouches.member) " +
		"LEFT OUTER JOIN ( " +
		"  SELECT 'for_me' AS dir, mv.vouchor AS member, COUNT(*) AS num " +
		"  FROM member_vouch mv " +
		"  WHERE ROW(mv.trustgroup, mv.vouchee) = ROW($1, $2) " +
		"  AND mv.positive " +
		"  GROUP BY mv.vouchor " +
		") as for_me_vouches on (m.ident = for_me_vouches.member) " +
		"LEFT OUTER JOIN ( " +
		"  SELECT 'by_me' AS dir, mv.vouchee AS member, COUNT(*) AS num " +
		"  FROM member_vouch mv " +
		"  WHERE ROW(mv.trustgroup, mv.vouchor) = ROW($1, $2) " +
		"  AND mv.positive " +
		"  GROUP BY mv.vouchee " +
		") as by_me_vouches on (m.ident = by_me_vouches.member) " +
		"WHERE grp.ident = $1 " +
		"AND me.email = mt.email "

	if inclhidden {
		if nominated {
			q += "AND ms.ident = 'nominated' "
		}
	} else {
		if nominated {
			q += "AND (NOT ms.hidden OR ms.ident = 'nominated') "
		} else {
			q += "AND NOT ms.hidden "
		}
	}

	if search == "" {
		if max != 0 {
			q += ord + " LIMIT $4 OFFSET $3"
			rows, err = pf.DB.Query(q, grpname, username, offset, max)
		} else {
			q += ord
			rows, err = pf.DB.Query(q, grpname, username)
		}
	} else {
		if exact {
			q += "AND (m.ident ~* $3) " +
				ord

		} else {
			q += "AND (m.ident ~* $3 " +
				"OR m.descr ~* $3 " +
				"OR m.affiliation ~* $3) " +
				ord
		}

		if max != 0 {
			q += " LIMIT $5 OFFSET $4"
			rows, err = pf.DB.Query(q, grpname, username, search, offset, max)
		} else {
			rows, err = pf.DB.Query(q, grpname, username, search)
		}
	}

	if err != nil {
		pf.Log("Query failed: " + err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		var fullname string
		var username string
		var affiliation string
		var groupname string
		var groupdesc string
		var groupadmin bool
		var groupstate string
		var groupcansee bool
		var email string
		var pgpkey_id string
		var entered string
		var activity string
		var tel string
		var sms string
		var airport string

		member := NewTriGroupMember().(*TriGroupMemberS)

		err = rows.Scan(
			&username,
			&fullname,
			&affiliation,
			&groupname,
			&groupdesc,
			&groupadmin,
			&groupstate,
			&groupcansee,
			&email,
			&pgpkey_id,
			&entered,
			&activity,
			&tel,
			&sms,
			&airport,
			&member.VouchesFor,
			&member.VouchesForMe,
			&member.VouchesBy,
			&member.VouchesByMe)
		if err != nil {
			pf.Log("Error listing members: " + err.Error())
			return nil, err
		}

		member.Set(groupname, groupdesc, username, fullname, affiliation, groupadmin, groupstate, groupcansee, email, pgpkey_id, entered, activity, sms, tel, airport)
		members = append(members, member)
	}

	return members, nil
}

func (grp *TriGroupS) Add_default_attestations(ctx pf.PfCtx) (err error) {
	att := make(map[string]string)
	att["met"] = "I have met them in person more than once"
	att["trust"] = "I trust them to take action"
	att["fate"] = "I will share membership fate with them"

	for a, descr := range att {
		q := "INSERT INTO attestations " +
			"(ident, descr, trustgroup) " +
			"VALUES($1, $2, $3)"
		err = pf.DB.Exec(ctx,
			"Added default attestation $1 to group $3",
			1, q,
			a, descr, grp.GetGroupName())

		if err != nil {
			return
		}
	}

	return
}

func (grp *TriGroupS) Add_default_mailinglists(ctx pf.PfCtx) (err error) {
	err = grp.PfGroup.Add_default_mailinglists(ctx)
	if err != nil {
		return
	}

	mls := make(map[string]string)
	mls["vetting"] = "Vetting and Vouching"

	for lhs, descr := range mls {
		err = pf.Ml_addv(ctx, grp.PfGroup, lhs, descr, true, true, true)
		if err != nil {
			return
		}
	}

	return
}

func group_add(ctx pf.PfCtx, args []string) (err error) {
	var group_name string

	/* Make sure the name is mostly sane */
	group_name, err = pf.Chk_ident("Group Name", args[0])
	if err != nil {
		return
	}

	d_maxin := 180 * 24 * time.Hour
	i_maxin := d_maxin.Seconds()

	d_guard := 7 * 24 * time.Hour
	i_guard := d_guard.Seconds()

	grp := ctx.NewGroup().(TriGroup)
	exists := grp.Exists(group_name)
	if exists {
		err = errors.New("Group already exists")
		return
	}

	q := "INSERT INTO trustgroup " +
		"(ident, descr, shortname, min_invouch, pgp_required, " +
		" please_vouch, vouch_adminonly, min_outvouch, max_inactivity, can_time_out, " +
		" max_vouchdays, idle_guard, nom_enabled, target_invouch, has_wiki) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) "
	err = pf.DB.Exec(ctx,
		"Created group $1",
		1, q,
		group_name, group_name, group_name, 0, false,
		true, false, 0, i_maxin, false,
		0, i_guard, true, 0, false)

	if err != nil {
		err = errors.New("Group creation failed")
		return
	}

	err = ctx.SelectGroup(group_name, pf.PERM_SYS_ADMIN)
	if err != nil {
		err = errors.New("Group creation failed")
		return
	}

	/* Fetch our newly created group */
	tctx := TriGetCtx(ctx)
	grp = tctx.TriSelectedGroup()

	/* Select yourself */
	ctx.SelectMe()
	if err != nil {
		return
	}

	err = grp.Add_default_attestations(ctx)
	if err != nil {
		return
	}

	err = grp.Add_default_mailinglists(ctx)
	if err != nil {
		return
	}

	grp.Member_add(ctx)
	grp.Member_set_state(ctx, pf.GROUP_STATE_APPROVED)
	grp.Member_set_admin(ctx, true)

	/* All worked */
	ctx.OutLn("Creation of group %s complete", group_name)
	return
}

func group_member_nominate(ctx pf.PfCtx, args []string) (err error) {
	grp := ctx.SelectedGroup()
	user := args[1]
	err = ctx.SelectUser(user, pf.PERM_USER_NOMINATE)
	if err != nil {
		err = errors.New("User selection failed")
		return
	}
	return grp.Member_add(ctx)
}

func group_menu(ctx pf.PfCtx, menu *pf.PfMenu) {
	menu.Replace("add", group_add)

	m := []pf.PfMEntry{
		{"vouch", vouch_menu, 0, -1, nil, pf.PERM_USER, "Vouch Commands"},
		{"nominate", group_member_nominate, 2, 2, []string{"group", "username"}, pf.PERM_GROUP_MEMBER, "Nominate a member for a group"},
	}

	menu.Add(m...)
}
