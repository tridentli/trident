package TriUI

import (
	pf "trident.li/pitchfork/lib"
	pu "trident.li/pitchfork/ui"
	"trident.li/trident/src/lib"
)

func h_user_vouches(cui pu.PfUI) {
	tcui := TriGetUI(cui)

	var vouch trident.TriVouch
	var isedit bool
	var msg string
	var err error
	var canvouch bool
	var errmsg string

	theuser := tcui.TriTheUser()
	user := tcui.TriSelectedUser()
	grp := tcui.TriSelectedGroup()

	/* SysAdmin and User-Self can edit */
	isedit = cui.IsSysAdmin() || cui.SelectedSelf()

	if isedit && cui.IsPOST() {
		msg, err = pu.H_user_profile_post(cui)
	}

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
		User         trident.TriUser
		TrustGroup   trident.TriGroup
		VouchIn      []trident.TriVouch
		VouchOut     []trident.TriVouch
		Details      []pf.PfUserDetail
		Languages    []pf.PfUserLanguage
		Attestations []trident.TriGroupAttestation
		IsEdit       bool
		CanVouch     bool
		MyUserName   string
		IsAdmin      bool
		GroupMember  trident.TriGroupMember
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

	/* Set the last link nicer */
	cui.AddCrumb("", "Profile", user.GetFullName()+" ("+user.GetUserName()+")'s Profile")

	p := Page{cui.Page_def(), msg, errmsg, user, grp, vi, vo, nil, nil, nil, isedit,
		canvouch, theuser.GetUserName(), isadmin, nil}

	/*
	 * We search for user's username thus will only get one result back
	 * Hence why grpusers[0] is the user we are looking for.
	 */
	grpusers, err := grp.GetMembers(user.GetUserName(), user.GetUserName(), 0, 1, true, true)
	if err == nil {
		p.GroupMember = grpusers[0].(trident.TriGroupMember)

		p.Details, err = user.GetDetails()
	}

	if err == nil {
		p.Languages, err = user.GetLanguages()
	}

	if err == nil {
		p.Attestations, err = grp.GetAttestations()
	}

	if err != nil {
		errmsg = err.Error()
	}

	cui.Page_show("user/profile_with_group.tmpl", p)
}

func h_user(cui pu.PfUI, menu *pu.PfUIMenu) {
	menu.Add(pu.PfUIMentry{"vouches", "Vouches", pf.PERM_GROUP_MEMBER, h_user_vouches, nil})
}
