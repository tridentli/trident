package TriUI

import (
	"errors"
	"net/http"
	"strings"
	pf "trident.li/pitchfork/lib"
	pu "trident.li/pitchfork/ui"
	tr "trident.li/trident/src/lib"
)

func attestations_get(cui pu.PfUI, grp tr.TriGroup) (attestations string, err error) {
	required_attestations, err := grp.GetAttestations()
	if err != nil {
		return
	}

	/* Walk required_attestations verify all are present */
	var atts []string
	for _, att := range required_attestations {
		value, e := cui.FormValue("attestation-" + att.Ident)
		if e != nil || value != "on" {
			err = errors.New("Incomplete attestations: " + att.Ident)
			return
		}

		atts = append(atts, att.Ident)
	}

	attestations = strings.Join(atts, ",")
	return
}

func h_vouch_edit(cui pu.PfUI) {
	tcui := TriGetUI(cui)
	user := tcui.TriSelectedUser()
	grp := tcui.TriSelectedGroup()
	vouchee := tcui.SelectedVouchee()
	comment, err := tcui.FormValue("comment")

	if err != nil {
		cmd := "group vouch update"
		arg := []string{grp.GetGroupName(), user.GetUserName(), vouchee.GetUserName(), comment}
		_, err = cui.HandleCmd(cmd, arg)

		if err != nil {
			pu.H_errmsg(cui, err)
			return
		}
	}

	cui.SetRedirect("/user/"+user.GetUserName()+"/group/"+grp.GetGroupName()+
		"/member/"+vouchee.GetUserName(), http.StatusSeeOther)
	return
}

func h_vouch_add(cui pu.PfUI) {
	tcui := TriGetUI(cui)

	user := tcui.TriSelectedUser()
	grp := tcui.TriSelectedGroup()
	vouchee := tcui.SelectedVouchee()
	comment, err := tcui.FormValue("comment")

	var attestations string

	if err != nil {
		attestations, err = attestations_get(cui, grp)
	}

	if err != nil {
		cmd := "group vouch add"
		arg := []string{grp.GetGroupName(), user.GetUserName(), vouchee.GetUserName(), comment, attestations}

		_, err = cui.HandleCmd(cmd, arg)

		if err != nil {
			pu.H_errmsg(cui, err)
			return
		}
	}

	cui.SetRedirect("/group/"+grp.GetGroupName()+"/member/"+user.GetUserName(), http.StatusSeeOther)
	return
}

func h_vouch_remove(cui pu.PfUI) {
	tcui := TriGetUI(cui)
	user := tcui.TriSelectedUser()
	grp := tcui.TriSelectedGroup()
	vouchee := tcui.SelectedVouchee()

	cmd := "group vouch remove"
	arg := []string{grp.GetGroupName(), user.GetUserName(), vouchee.GetUserName()}

	_, err := cui.HandleCmd(cmd, arg)

	if err != nil {
		pu.H_errmsg(cui, err)
		return
	}

	cui.SetRedirect("/group/"+grp.GetGroupName()+"/member/"+user.GetUserName(), http.StatusSeeOther)
	return
}

func h_vouch_nominate_form(cui pu.PfUI) {
	tcui := TriGetUI(cui)
	user := tcui.TriSelectedUser()
	grp := tcui.TriSelectedGroup()
	vouchee := tcui.SelectedVouchee()

	type Page struct {
		*pu.PfPage
		GroupName    string
		Vouchor      string
		Vouchee      string
		Attestations []tr.TriGroupAttestation
	}

	attestations, err := grp.GetAttestations()
	if err != nil {
		pu.H_errtxt(cui, "Incomplete attestations")
		return
	}

	p := Page{cui.Page_def(), grp.GetGroupName(), user.GetUserName(), vouchee.GetUserName(), attestations}
	cui.Page_show("vouch/nominate_form.tmpl", p)
}

func h_vouch_nominate_new(cui pu.PfUI) {
	tcui := TriGetUI(cui)

	var cmd string
	var args []string

	user := tcui.TriSelectedUser()
	grp := tcui.TriSelectedGroup()

	address, err := cui.FormValue("address")
	descr, err2 := cui.FormValue("descr")
	comment, err3 := cui.FormValue("comment")
	bio_info, err4 := cui.FormValue("bio_info")
	affil, err5 := cui.FormValue("affiliation")
	attestations, err6 := attestations_get(cui, grp)

	if err != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
		pu.H_errtxt(cui, "Invalid parameters")
		return
	}

	/* Generate a username */
	vouchee_ident, err := pf.Fullname_to_ident(descr)

	cmd = "user nominate"
	args = []string{vouchee_ident, address, bio_info, affil, descr}

	_, err = cui.HandleCmd(cmd, args)

	if err != nil {
		pu.H_errmsg(cui, err)
		return
	}

	cmd = "user email confirm_force"
	args = []string{vouchee_ident, address}

	_, err = cui.HandleCmd(cmd, args)

	if err != nil {
		pu.H_errmsg(cui, err)
		return
	}

	cmd = "group member nominate"
	args = []string{grp.GetGroupName(), vouchee_ident}

	_, err = cui.HandleCmd(cmd, args)

	if err != nil {
		pu.H_errmsg(cui, err)
		return
	}

	cmd = "group vouch add"
	args = []string{grp.GetGroupName(), user.GetUserName(), vouchee_ident, comment, attestations}

	_, err = cui.HandleCmd(cmd, args)

	if err != nil {
		pu.H_errmsg(cui, err)
		return
	}

	cui.SetRedirect("/user/"+user.GetUserName()+"/group/"+grp.GetGroupName()+"/members", http.StatusSeeOther)
	return
}

func h_vouch_nominate(cui pu.PfUI) {
	tcui := TriGetUI(cui)

	var cmd string
	var arg []string

	user := tcui.TriSelectedUser()
	grp := tcui.TriSelectedGroup()
	vouchee := tcui.SelectedVouchee()
	comment, err := tcui.FormValue("comment")

	if err != nil {
		pu.H_errtxt(cui, "Invalid parameters")
		return
	}

	attestations, err := attestations_get(cui, grp)

	cmd = "group member add"
	arg = []string{grp.GetGroupName(), vouchee.GetUserName()}

	_, err = cui.HandleCmd(cmd, arg)

	if err != nil {
		pu.H_errmsg(cui, err)
		return
	}

	cmd = "group vouch add"
	arg = []string{grp.GetGroupName(), user.GetUserName(), vouchee.GetUserName(), comment, attestations}

	_, err = cui.HandleCmd(cmd, arg)

	if err != nil {
		pu.H_errmsg(cui, err)
		return
	}

	cui.SetRedirect("/user/"+user.GetUserName()+"/group/"+grp.GetGroupName()+"/members", http.StatusSeeOther)
	return
}

func h_vouch(cui pu.PfUI) {
	tcui := TriGetUI(cui)

	path := tcui.GetPath()

	if len(path) == 0 || path[0] == "" {
		pu.H_group_member_profile(cui)
		return
	}

	if path[0] == "nominate_new" {
		h_vouch_nominate_new(cui)
		return
	}

	if !cui.HasSelectedGroup() {
		pu.H_errtxt(cui, "Vouch: No group selected")
		return
	}
	if !cui.HasSelectedUser() {
		pu.H_errtxt(cui, "Vouch: No User Selected")
		return
	}

	/* Check member access to group */
	err := tcui.SelectVouchee(path[0], pf.PERM_GROUP_MEMBER|pf.PERM_USER_VIEW)
	if err != nil {
		tcui.Errf("Selecting Vouchee: %s", err.Error())
		pu.H_error(tcui, http.StatusNotFound)
		return
	}

	vouchee := tcui.SelectedVouchee()

	tcui.AddCrumb(path[0], vouchee.GetUserName(), vouchee.GetFullName())

	tcui.SetPath(path[1:])

	p := pf.PERM_GROUP_MEMBER | pf.PERM_HIDDEN

	menu := pu.NewPfUIMenu([]pu.PfUIMentry{
		{"edit", "Edit Vouch", p, h_vouch_edit, nil},
		{"add", "Create Vouch", p, h_vouch_add, nil},
		{"remove", "Remove Vouch", p, h_vouch_remove, nil},
		{"nominate", "Nominate and Vouch", p, h_vouch_nominate, nil},
		{"nominate_form", "Create form to nominate existing user", p, h_vouch_nominate_form, nil},
	})

	cui.UIMenu(menu)
}
