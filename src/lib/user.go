package trident

import (
	"errors"
	pf "trident.li/pitchfork/lib"
)

type TriUser interface {
	pf.PfUser
	IsNominator(ctx pf.PfCtx, nom_name string) (ok bool)
	BestNominator(ctx pf.PfCtx) (nom_name string, err error)
}

type TriUserS struct {
	pf.PfUser `pfset:"self" pfget:"self"`
}

func NewTriUser() pf.PfUser {
	return &TriUserS{PfUser: pf.NewPfUser(nil, nil)}
}

func (user *TriUserS) IsNominator(ctx pf.PfCtx, nom_name string) (ok bool) {
	cnt := 0

	q := "SELECT COUNT(*) " +
		"FROM member_vouch mv " +
		"JOIN member_email me ON mv.vouchor = me.member " +
		"WHERE vouchee = $1 " +
		"AND mv.positive " +
		"AND me.pgpkey_id IS NOT NULL " +
		"WHERE vouchor = $2"
	_ = pf.DB.QueryRow(q, user.GetUserName(), nom_name).Scan(&cnt)

	return cnt == 0
}

func (user *TriUserS) BestNominator(ctx pf.PfCtx) (nom_name string, err error) {
	q := "SELECT vouchor " +
		"FROM member_vouch mv " +
		"JOIN member_email me ON mv.vouchor = me.member " +
		"WHERE vouchee = $1 " +
		"AND mv.positive " +
		"AND me.pgpkey_id IS NOT NULL " +
		"ORDER BY mv.entered " +
		"LIMIT 1"
	err = pf.DB.QueryRow(q, user.GetUserName()).Scan(&nom_name)
	return
}

/* TODO: Verify: Only show member of groups my user is associated with and are non-anonymous */
func (user *TriUserS) GetList(ctx pf.PfCtx, search string, offset int, max int) (users []pf.PfUser, err error) {
	users = nil

	/* The fields we match on */
	matches := []string{"ident", "me.email", "descr", "affiliation", "bio_info", "comment", "d.value"}

	var p []string
	var t []pf.DB_Op
	var v []interface{}

	for _, m := range matches {
		p = append(p, m)
		t = append(t, pf.DB_OP_ILIKE)
		v = append(v, "%"+search+"%")
	}

	j := "INNER JOIN member_email me ON member.ident = me.member " +
		"LEFT OUTER JOIN member_details d ON d.member = member.ident " +
		"LEFT OUTER JOIN member_vouch v ON v.vouchee = member.ident"

	o := "GROUP BY member.ident " +
		"ORDER BY member.ident"

	objs, err := pf.StructFetchMulti(ctx.NewUserI, "member", j, pf.DB_OP_OR, p, t, v, o, offset, max)
	if err != nil {
		return
	}

	/* Get the groups these folks are in */
	for _, o := range objs {
		u := o.(pf.PfUser)
		users = append(users, u)
	}

	return users, err
}

func user_pw_send(ctx pf.PfCtx, is_reset bool, nom_name string) (err error) {
	var user_email pf.PfUserEmail
	var nom_email pf.PfUserEmail
	var pw pf.PfPass
	var user_portion string
	var nom_portion string

	theuser := ctx.SelectedUser()

	username := theuser.GetUserName()

	/* Make sure the name is mostly sane */
	nom_name, err = pf.Chk_ident("UserName", nom_name)
	if err != nil {
		return
	}

	if nom_name == username {
		err = errors.New("Nominator cannot be the same as the user")
		return
	}

	err = ctx.SelectUser(nom_name, pf.PERM_USER_NOMINATE)
	if err != nil {
		return
	}

	nom_user := ctx.SelectedUser()
	nom_email, err = nom_user.GetPriEmail(ctx, false)
	if err != nil {
		return
	}

	/* Reselect the user, this is the one it is all about */
	err = ctx.SelectUser(username, pf.PERM_USER)
	if err != nil {
		return
	}

	user_email, err = theuser.GetPriEmail(ctx, true)
	if err != nil {
		return
	}

	user_portion, err = pw.GenPass(16)
	if err != nil {
		return
	}

	err = Mail_PassResetUser(ctx, user_email, is_reset, nom_email, user_portion)
	if err != nil {
		return
	}

	nom_portion, err = pw.GenPass(16)
	if err != nil {
		return
	}

	err = Mail_PassResetNominator(ctx, nom_email, is_reset, user_email, nom_portion)
	if err != nil {
		return
	}

	err = theuser.SetRecoverToken(ctx, user_portion+nom_portion)

	return
}

func user_pw_reset(ctx pf.PfCtx, args []string) (err error) {

	username := args[0]
	nom_name := ""

	err = ctx.SelectUser(username, pf.PERM_USER_SELF)
	if err != nil {
		return
	}

	tctx := TriGetCtx(ctx)
	user := tctx.TriSelectedUser()

	/*
	 * Note that when the user does not have a valid nominator
	 * the password can't be reset either
	 */
	if len(args) >= 2 {
		nom_name = args[1]
		if !user.IsNominator(ctx, nom_name) {
			err = errors.New(nom_name + " is not a nominator for this user")
			return
		}
	} else {
		nom_name, err = user.BestNominator(ctx)
		if err != nil {
			err = errors.New("No nominator with valid PGP key")
			return
		}
	}

	/* Send out the new password */
	err = user_pw_send(ctx, true, nom_name)
	if err == nil {
		ctx.OutLn("Recovery Passwords sent to user and " + nom_name)
	}
	return
}

func user_nominate(ctx pf.PfCtx, args []string) (err error) {
	username := args[0]
	email := args[1]
	bio_info := args[2]
	affiliation := args[3]
	descr := args[4]

	return pf.User_new(ctx, username, email, bio_info, affiliation, descr)
}

func user_merge(ctx pf.PfCtx, args []string) (err error) {
	u_new := args[0]
	u_old := args[1]

	err = pf.DB.TxBegin(ctx)
	if err != nil {
		return err
	}

	/* No error yet */
	err = nil

	q := ""

	if err == nil {
		q = "UPDATE member_vouch " +
			"SET vouchor = $1 " +
			"WHERE vouchor = $2"
		err = pf.DB.Exec(ctx,
			"Update Vouches $2 to $1",
			-1, q,
			u_new, u_old)
	}

	if err == nil {
		q = "UPDATE member_vouch " +
			"SET vouchee = $1" +
			"WHERE vouchee = $2"
		err = pf.DB.Exec(ctx,
			"Update Vouchee $2 to $1",
			-1, q,
			u_new, u_old)
	}

	return pf.User_merge(ctx, u_new, u_old, err)
}

func user_pw_menu(ctx pf.PfCtx, menu *pf.PfMenu) {
	m := []pf.PfMEntry{
		{"reset", user_pw_reset, 1, 2, []string{"username", "nominator"}, pf.PERM_USER, "Send a recovery password split between the user and a nominator"},
	}

	menu.Add(m...)
}

func user_menu(ctx pf.PfCtx, menu *pf.PfMenu) {
	m := []pf.PfMEntry{
		{"nominate", user_nominate, 5, 5, []string{"username", "email", "bio_info", "affiliation", "descr"}, pf.PERM_USER, "Nominate New User"},
	}

	menu.Add(m...)
}
