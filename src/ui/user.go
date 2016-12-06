package TriUI

import (
	pu "trident.li/pitchfork/ui"
	tr "trident.li/trident/src/lib"
)

func h_user_vouches(cui pu.PfUI) {
	tctx := tr.TriGetCtx(cui)

	var vouch tr.TriVouch
	var isedit bool
	var msg string
	var err error
	var canvouch bool
	var errmsg string

	theuser := tctx.TriTheUser()
	user := tctx.TriSelectedUser()
	grp := tctx.TriSelectedGroup()

	/* SysAdmin and User-Self can edit */
	isedit = cui.IsSysAdmin() || cui.SelectedSelf()

	if err != nil {
		/* Failed */
		errmsg = err.Error()
	} else {
		/* Success */
	}

	/* Refresh updated version */
	err = user.Refresh(cui)
	if err != nil {
		errmsg += err.Error()
	}

	/* Output the page */
	type Page struct {
		*pu.PfPage
		Message      string
		Error        string
		User         tr.TriUser
		Group        tr.TriGroup
		VouchIn      []tr.TriVouch
		VouchOut     []tr.TriVouch
		Attestations []tr.TriGroupAttestation
		IsEdit       bool
		CanVouch     bool
		MyUserName   string
		IsAdmin      bool
		GroupMember  tr.TriGroupMember
	}

	vi, err := vouch.ListFor(user, grp, theuser.GetUserName())
	if err != nil {
		cui.Err(err.Error())
		return
	}

	vo, err := vouch.ListBy(user, grp, theuser.GetUserName())
	if err != nil {
		cui.Err(err.Error())
		return
	}

	canvouch = false
	if user.GetUserName() != theuser.GetUserName() {
		if grp.GetVouch_adminonly() == false {
			/* Non-admin can vouch */
			canvouch = true

			/* Unless the user has already vouched */
			for _, v := range vi {
				if v.Vouchor == theuser.GetUserName() {
					canvouch = false
				}
			}
		}
	}

	isadmin := cui.IAmGroupAdmin()

	p := Page{cui.Page_def(), msg, errmsg, user, grp, vi, vo, nil, isedit,
		canvouch, theuser.GetUserName(), isadmin, nil}

	/*
	 * We search for user's username thus will only get one result back
	 * Hence why grpusers[0] is the user we are looking for.
	 */
	grpusers, err := grp.GetMembers(user.GetUserName(), user.GetUserName(), 0, 1, true, true)
	if err == nil {
		p.GroupMember = grpusers[0].(tr.TriGroupMember)
	}

	if err == nil {
		p.Attestations, err = grp.GetAttestations()
	}

	if err != nil {
		errmsg = err.Error()
	}

	cui.Page_show("user/profile_vouches.tmpl", p)
}
