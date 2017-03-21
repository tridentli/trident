package trident

import (
	"errors"
	"strings"
	"time"
	pf "trident.li/pitchfork/lib"
)

type TriVouch struct {
	Vouchor      string    `label:"Voucher" pfset:"self" pfcol:"voucher"`
	Vouchee      string    `label:"Vouchee" pfset:"self" pfcol:"vouchee"`
	GroupName    string    `label:"TrustGroup" pfset:"self" pfcol:"trustgroup"`
	Comment      string    `label:"Comment" pfset:"self" pfcol:"comment"`
	Entered      time.Time `label:"Entered" pfset:"nobody"`
	Positive     bool      `label:"Positive" pfset:"self" pfcol:"positive" hint:"Is the vouch of a positive nature"`
	State        string    `label:"Member state in trustgroup" pfset:"nobody"`
	Affiliation  string    `label:"Vouchee Afiliation" pfset:"nobody"`
	MyVouch      bool
	Attestations []TriGroupAttestation
}

func (vouch *TriVouch) String() (out string) {
	var msg string

	out = "Group: " + vouch.GroupName + "\n"
	out += "   " + vouch.Vouchor + " -> " + vouch.Vouchee + "\n"
	out += "   Entered: " + vouch.Entered.String() + "\n"
	msgs := strings.Split(vouch.Comment, "\n")

	for _, msg = range msgs {
		out += "   " + msg + "\n"
	}
	out += "\n"

	return
}

/*
 * Generates a set of vouches that were issued
 * by the provided user in a given group.
 * If username is populated any vouches issued
 * by that user will have MyVouch set true
 */

func (vouch *TriVouch) ListBy(user TriUser, tg TriGroup, username string) (vouches []TriVouch, err error) {
	q := "SELECT mv.vouchor, mv.comment, mt.trustgroup, " +
		"DATE_TRUNC('seconds', mv.entered) as entered, " +
		"vm.affiliation, mv.vouchee, mt.state " +
		"FROM member_vouch mv " +
		"JOIN member vm ON (mv.vouchor = vm.ident) " +
		"JOIN member m ON (m.ident = mv.vouchor) " +
		"JOIN member_trustgroup mt ON " +
		" ROW(m.ident, mv.trustgroup) = ROW(mt.member, mt.trustgroup) " +
		"WHERE ROW(mv.vouchor, mv.trustgroup) = ROW($1, $2) " +
		"AND mv.positive " +
		"ORDER BY mv.entered DESC"

	rows, err := pf.DB.Query(q, user.GetUserName(), tg.GetGroupName())
	if err != nil {
		err = errors.New("Could not retrieve vouches for user " + err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		var v TriVouch

		err = rows.Scan(&v.Vouchor, &v.Comment, &v.GroupName,
			&v.Entered, &v.Affiliation, &v.Vouchee, &v.State)
		if err != nil {
			vouches = nil
			return
		}

		v.MyVouch = false
		if username == v.Vouchor {
			v.MyVouch = true
		}

		vouches = append(vouches, v)
	}

	return
}

/*
 * Generates a set of vouches that were issued
 * by other users for the provided user in a given group.
 * If username is populated any vouches issued
 * by that user will have MyVouch set true
 */

func (vouch *TriVouch) ListFor(user TriUser, tg TriGroup, username string) (vouches []TriVouch, err error) {
	q := "SELECT mv.vouchor, mv.comment, mt.trustgroup, " +
		"DATE_TRUNC('seconds', mv.entered) as entered, " +
		"vm.affiliation, mv.vouchee, mt.state " +
		"FROM member_vouch mv " +
		"JOIN member vm ON (mv.vouchor = vm.ident) " +
		"JOIN member m ON (m.ident = mv.vouchee) " +
		"JOIN member_trustgroup mt ON " +
		" ROW(m.ident, mv.trustgroup) = ROW(mt.member, mt.trustgroup) " +
		"WHERE ROW(mv.vouchee, mv.trustgroup) = ROW($1, $2) " +
		"AND mv.positive " +
		"ORDER BY mv.entered DESC"

	rows, err := pf.DB.Query(q, user.GetUserName(), tg.GetGroupName())
	if err != nil {
		err = errors.New("Could not retrieve vouches for user " + err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		var v TriVouch

		err = rows.Scan(&v.Vouchor, &v.Comment, &v.GroupName,
			&v.Entered, &v.Affiliation, &v.Vouchee, &v.State)
		if err != nil {
			vouches = nil
			return
		}

		v.MyVouch = false
		if username == v.Vouchor {
			v.MyVouch = true
		}

		vouches = append(vouches, v)
	}

	return
}

func vouches_member(ctx pf.PfCtx, args []string, is_by bool) (err error) {
	group_name := args[0]
	user_name := args[1]

	var tv TriVouch
	var vouches []TriVouch

	if !ctx.IsLoggedIn() {
		err = errors.New("Not authenticated")
		return
	}

	theuser := ctx.TheUser()

	err = ctx.SelectUser(user_name, pf.PERM_USER_SELF)
	if err != nil {
		return
	}

	err = ctx.SelectGroup(group_name, pf.PERM_GROUP_MEMBER)
	if err != nil {
		return
	}

	tctx := TriGetCtx(ctx)
	user := tctx.TriSelectedUser()
	tg := tctx.TriSelectedGroup()

	if is_by {
		vouches, err = tv.ListBy(user, tg, theuser.GetUserName())
	} else {
		vouches, err = tv.ListFor(user, tg, theuser.GetUserName())
	}

	for _, tv = range vouches {
		ctx.OutLn(tv.String())
	}

	return
}

func vouches_by_member(ctx pf.PfCtx, args []string) (err error) {
	return vouches_member(ctx, args, true)
}

func vouches_for_member(ctx pf.PfCtx, args []string) (err error) {
	return vouches_member(ctx, args, false)
}

func vouch_add(ctx pf.PfCtx, args []string) (err error) {
	group_name := args[0]
	user_name := args[1]
	vouchee_name := args[2]
	comment := args[3]
	attestations := args[4]

	err = ctx.SelectUser(user_name, pf.PERM_USER_SELF)
	if err != nil {
		return
	}

	err = ctx.SelectGroup(group_name, pf.PERM_GROUP_MEMBER)
	if err != nil {
		return
	}

	tctx := TriGetCtx(ctx)
	user := tctx.TriSelectedUser()
	tg := tctx.TriSelectedGroup()

	/* Get Attestations */
	required_attestations, err := tg.GetAttestations()
	if err != nil {
		return
	}

	/* Parse attestations into structure. (Only those marked true) */
	attestation_array := strings.Split(attestations, ",")

	/* Walk required_attestations verify all are present */
	attest_okay := true

	for _, att := range required_attestations {
		attest_inner := false

		for _, i_label := range attestation_array {
			if i_label == att.Ident {
				attest_inner = true
			}
		}

		if !attest_inner {
			attest_okay = false
		}
	}

	if !attest_okay {
		err = errors.New("Incomplete attestations")
		return
	}

	q := "INSERT INTO member_vouch " +
		"(vouchor, vouchee, trustgroup, comment, entered, positive) " +
		"VALUES($1, $2, $3, $4, now(), 't') "
	err = pf.DB.Exec(ctx,
		"Created Vouch $1 -> $2, group: $3",
		1, q,
		user.GetUserName(), vouchee_name, tg.GetGroupName(), comment)
	return
}

/*
 * We're going to limit the scope of this to only changing the comment.
 * Entered at will update automaticall
 */
func vouch_update(ctx pf.PfCtx, args []string) (err error) {
	group_name := args[0]
	user_name := args[1]
	vouchee_name := args[2]
	comment := args[3]

	err = ctx.SelectUser(user_name, pf.PERM_USER_SELF)
	if err != nil {
		return
	}

	err = ctx.SelectGroup(group_name, pf.PERM_GROUP_MEMBER)
	if err != nil {
		return
	}

	tctx := TriGetCtx(ctx)
	user := tctx.TriSelectedUser()
	tg := tctx.TriSelectedGroup()

	q := "UPDATE member_vouch " +
		"SET comment = $1, " +
		"entered = now() " +
		"WHERE vouchor = $2 " +
		"AND vouchee = $3 " +
		"AND trustgroup = $4"
	err = pf.DB.Exec(ctx,
		"Updated Vouch $1 -> $2, group: $3",
		-1, q,
		comment, user.GetUserName(), vouchee_name, tg.GetGroupName())
	return
}

func vouch_remove(ctx pf.PfCtx, args []string) (err error) {
	group_name := args[0]
	user_name := args[1]
	vouchee_name := args[2]

	err = ctx.SelectUser(user_name, pf.PERM_USER_SELF)
	if err != nil {
		return
	}

	err = ctx.SelectGroup(group_name, pf.PERM_GROUP_MEMBER)
	if err != nil {
		return
	}

	tctx := TriGetCtx(ctx)
	user := tctx.TriSelectedUser()
	tg := tctx.TriSelectedGroup()

	q := "DELETE FROM member_vouch " +
		"WHERE vouchor = $1 " +
		"AND vouchee = $2 " +
		"AND trustgroup = $3 " +
		"AND positive = true"
	err = pf.DB.Exec(ctx,
		"Deleted Vouch $1 -> $2, group: $3",
		-1, q,
		user.GetUserName(), vouchee_name, tg.GetGroupName())
	return
}

type Vouch struct {
	Vouchor string
	Vouchee string
	Entered time.Time
}

func Vouches_Get(ctx pf.PfCtx, group_name string) (v []Vouch, err error) {
	err = ctx.SelectGroup(group_name, pf.PERM_GROUP_MEMBER)
	if err != nil {
		return
	}

	q := "SELECT " +
		"mv.vouchor, mv.vouchee, mv.entered " +
		"FROM member_vouch mv " +
		"JOIN member m1 ON (mv.vouchor = m1.ident) " +
		"JOIN member m2 ON (mv.vouchee = m2.ident) " +
		"JOIN member_trustgroup mt1 ON " +
		"ROW(mv.vouchor, mv.trustgroup) = " +
		"ROW(mt1.member, mt1.trustgroup) " +
		"JOIN member_trustgroup mt2 ON " +
		"ROW(mv.vouchee, mv.trustgroup) = " +
		"ROW(mt2.member, mt2.trustgroup) " +
		"JOIN member_state ms1 ON (mt1.state = ms1.ident) " +
		"JOIN member_state ms2 ON (mt2.state = ms2.ident) " +
		"WHERE mv.trustgroup = $1 " +
		"AND ms1.can_login " +
		"AND ms2.can_login " +
		"AND mv.positive"

	rows, err := pf.DB.Query(q, group_name)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var vouch Vouch

		err = rows.Scan(&vouch.Vouchor, &vouch.Vouchee, &vouch.Entered)
		if err != nil {
			return
		}

		v = append(v, vouch)
	}

	return
}

func vouches_emit(ctx pf.PfCtx, args []string) (err error) {
	group_name := args[0]

	vouches, err := Vouches_Get(ctx, group_name)
	if err != nil {
		return
	}

	for _, v := range vouches {
		ctx.Outf("%s,%s,%s\n", v.Vouchor, v.Vouchee, v.Entered)
	}

	return
}

func vouch_menu(ctx pf.PfCtx, args []string) (err error) {
	perms := pf.PERM_USER

	var menu = pf.NewPfMenu([]pf.PfMEntry{
		{"add", vouch_add, 5, 5, []string{"group", "vouchor", "vouchee", "comment", "attestations"}, pf.PERM_USER, "Add a vouch"},
		{"update", vouch_update, 4, 4, []string{"group", "vouchor", "vouchee", "comment"}, perms, "Update vouch comment"},
		{"remove", vouch_remove, 3, 3, []string{"group", "vouchor", "vouchee"}, perms, "Remove a vouch"},
		{"list_for", vouches_for_member, 2, 2, []string{"group", "username"}, perms, "List vouches for a member"},
		{"list_by", vouches_by_member, 2, 2, []string{"group", "username"}, perms, "List vouches by a member"},
		{"emit", vouches_emit, 1, 1, []string{"group"}, perms, "Emit Vouches"},
	})

	return ctx.Menu(args, menu)
}
