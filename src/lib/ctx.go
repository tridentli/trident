package trident

import (
	pf "trident.li/pitchfork/lib"
)

type TriCtx interface {
	pf.PfCtx /* Pitchfork Context */

	TriTheUser() (user TriUser)
	TriSelectedUser() (user TriUser)
	TriSelectedGroup() (grp TriGroup)
	SelectedVouchee() (user TriUser)
	HasSelectedVouchee() bool
	SelectVouchee(username string, perms pf.Perm) (err error)
}

type TriCtxS struct {
	pf.PfCtx
	sel_vouchee TriUser /* Selected User target of a vouch*/
}

func NewTriCtx() pf.PfCtx {
	pctx := pf.NewPfCtx(NewTriUser, NewTriGroup, TriMenuOverride, nil, nil)
	tctx := &TriCtxS{PfCtx: pctx}
	return tctx
}

func TriGetCtx(ctx pf.PfCtx) (tctx *TriCtxS) {
	tctxp, ok := ctx.(*TriCtxS)
	if !ok {
		panic("Not a TriCtx")
	}
	return tctxp
}

func (ctx *TriCtxS) TriTheUser() (user TriUser) {
	return ctx.PfCtx.TheUser().(TriUser)
}

func (ctx *TriCtxS) TriSelectedUser() (user TriUser) {
	return ctx.PfCtx.SelectedUser().(TriUser)
}

func (ctx *TriCtxS) TriSelectedGroup() (grp TriGroup) {
	return ctx.PfCtx.SelectedGroup().(TriGroup)
}

func (ctx *TriCtxS) SelectedVouchee() (user TriUser) {
	/* Return a copy, not a reference */
	/* XXX: verify */
	return ctx.sel_vouchee
}

func (ctx *TriCtxS) HasSelectedVouchee() bool {
	return ctx.sel_vouchee != nil
}

func (ctx *TriCtxS) SelectVouchee(username string, perms pf.Perm) (err error) {
	ctx.PDbgf("SelectVouchee", perms, "%q", username)

	/* Nothing to select, always works */
	if username == "" {
		ctx.sel_vouchee = nil
		return nil
	}

	ctx.sel_vouchee = ctx.NewUser().(*TriUserS)
	err = ctx.sel_vouchee.Select(ctx, username, perms)
	if err != nil {
		ctx.sel_vouchee = nil
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
