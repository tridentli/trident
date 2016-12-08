package TriUI

import (
	"errors"
	"trident.li/keyval"
	pf "trident.li/pitchfork/lib"
	pu "trident.li/pitchfork/ui"
	tr "trident.li/trident/src/lib"
)

type VouchAdd struct {
	group        tr.TriGroup
	Action       string          `label:"Action" pftype:"hidden"`
	Group        string          `label:"Group" pftype:"hidden"`
	Vouchee      string          `label:"Username" pftype:"hidden"`
	Comment      string          `label:"Vouch" pftype:"text" hint:"Vouch for this user" pfreq:"yes"`
	Attestations map[string]bool `label:"Attestations (all required)" hint:"Attestations for this user" options:"GetAttestationOpts" pfcheckboxmode:"yes"`
	Button       string          `label:"Vouch" pftype:"submit"`
}

func (va *VouchAdd) GetAttestationOpts(obj interface{}) (kvs keyval.KeyVals, err error) {
	return va.group.GetAttestationsKVS()
}

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

	if cui.IsPOST() {
		action, err := cui.FormValue("action")
		if err == nil {
			switch action {
			case "vouch_add":
				err = vouch_add(cui)
				break

			case "vouch_edit":
				err = vouch_edit(cui)
				break

			case "vouch_remove":
				err = vouch_remove(cui)
				break

			default:
				cui.Errf("Unknown action %q", action)
				err = errors.New("Unknown action provided")
				break
			}
		}
	}

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

	details, err := user.GetDetails()
	if err != nil {
		cui.Errf("Failed to GetDetails(): %s", err.Error())
		pu.H_error(cui, pu.StatusBadRequest)
		return
	}

	languages, err := user.GetLanguages()
	if err != nil {
		cui.Errf("Failed to GetLanguages(): %s", err.Error())
		pu.H_error(cui, pu.StatusBadRequest)
		return
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
		IsAdmin      bool
		GroupMember  tr.TriGroupMember
		Details      []pf.PfUserDetail
		Languages    []pf.PfUserLanguage
		VouchAdd     *VouchAdd
	}

	vf := &VouchAdd{group: grp, Group: grp.GetGroupName(), Vouchee: user.GetUserName(), Action: "vouch_add"}

	p := Page{cui.Page_def(), msg, errmsg, user, grp, vi, vo, nil, isedit,
		canvouch, isadmin, nil, details, languages, vf}

	/*
	 * We search for user's username thus will only get one result back
	 * Hence why grpusers[0] is the user we are looking for.
	 *
	 * TODO: Update SQL for the case when there are no vouches for the user
	 *       we now generate an empty GroupMember object instead
	 */
	grpusers, err := grp.GetMembers("", user.GetUserName(), 0, 1, true, true, true)
	if err == nil && len(grpusers) > 0 {
		p.GroupMember = grpusers[0].(tr.TriGroupMember)
	} else {
		p.GroupMember = tr.NewTriGroupMember()
	}

	if err == nil {
		p.Attestations, err = grp.GetAttestations()
	}

	if err != nil {
		p.Error += err.Error()
	}

	cui.Page_show("user/profile_vouches.tmpl", p)
}
