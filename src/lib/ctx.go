package trident

import (
	pf "trident.li/pitchfork/lib"
)

type TriCtx struct {
	pfctx       pf.PfCtx
	sel_vouchee TriUser /* Selected User target of a vouch*/
}

func NewTriCtx() pf.PfCtx {
	pctx := pf.NewPfCtx(NewTriUser, NewTriGroup, TriMenuOverride, nil, nil)
	tctx := &TriCtx{pfctx: pctx}
	pctx.SetAppData(tctx)
	return pctx
}

func TriGetCtx(ctx pf.PfCtx) (tctx *TriCtx) {
	tctx = ctx.GetAppData().(*TriCtx)
	tctx.pfctx = ctx
	return
}

func (tctx *TriCtx) TriTheUser() (user TriUser) {
	return tctx.pfctx.TheUser().(TriUser)
}

func (tctx *TriCtx) TriSelectedUser() (user TriUser) {
	return tctx.pfctx.SelectedUser().(TriUser)
}

func (tctx *TriCtx) TriSelectedGroup() (grp TriGroup) {
	return tctx.pfctx.SelectedGroup().(TriGroup)
}

func (tctx *TriCtx) SelectedVouchee() (user TriUser) {
	/* Return a copy, not a reference */
	/* XXX: verify */
	return tctx.sel_vouchee
}

func (tctx *TriCtx) HasSelectedVouchee() bool {
	return tctx.sel_vouchee != nil
}

func (tctx *TriCtx) SelectVouchee(username string, perms pf.Perm) (err error) {
	tctx.pfctx.PDbgf("SelectVouchee", perms, "%q", username)

	/* Nothing to select, always works */
	if username == "" {
		tctx.sel_vouchee = nil
		return nil
	}

	tctx.sel_vouchee = tctx.pfctx.NewUser().(*TriUserS)
	err = tctx.sel_vouchee.Select(tctx.pfctx, username, perms)
	if err != nil {
		tctx.sel_vouchee = nil
	}

	return
}

func TriMenuOverride(ctx pf.PfCtx, menu *pf.PfMenu) {
	loc := ctx.GetLoc()

	switch loc {
	case "group":
		group_menu(ctx, menu)
		break

	case "user":
		user_menu(ctx, menu)
		break

	case "user password":
		user_pw_menu(ctx, menu)
		break
	}

	return
}
