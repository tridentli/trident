package TriUI

import (
	"strconv"
	"trident.li/keyval"
	pf "trident.li/pitchfork/lib"
	pu "trident.li/pitchfork/ui"
	tr "trident.li/trident/src/lib"
)

func h_group_members(cui pu.PfUI) {
	path := cui.GetPath()

	if len(path) != 0 && path[0] != "" {
		pu.H_group_member_profile(cui)
		return
	}

	var err error

	tctx := tr.TriGetCtx(cui)
	total := 0
	offset := 0

	offset_v, err := cui.FormValue("offset")
	if err == nil && offset_v != "" {
		offset, _ = strconv.Atoi(offset_v)
	}

	search, err := cui.FormValue("search")
	if err != nil {
		search = ""
	}

	grp := tctx.TriSelectedGroup()

	total, err = grp.GetMembersTot(search)
	if err != nil {
		cui.Err("error: " + err.Error())
		return
	}

	members, err := grp.GetMembers(search, cui.TheUser().GetUserName(), offset, 10, false, cui.IAmGroupAdmin(), false)
	if err != nil {
		cui.Err(err.Error())
		return
	}

	/* Output the page */
	type Page struct {
		*pu.PfPage
		Group        pf.PfGroup
		GroupMembers []pf.PfGroupMember
		PagerOffset  int
		PagerTotal   int
		Search       string
		IsAdmin      bool
	}
	isadmin := cui.IAmGroupAdmin()

	p := Page{cui.Page_def(), grp, members, offset, total, search, isadmin}
	cui.Page_show("group/members.tmpl", p)
}

type NominateAdd struct {
	group        tr.TriGroup
	Action       string          `label:"Action" pftype:"hidden"`
	Vouchee      string          `label:"Username" pfset:"nobody" pfget:"none"`
	Comment      string          `label:"Vouch comment" pftype:"text" hint:"Vouch description for this user" pfreq:"yes"`
	Attestations map[string]bool `label:"Attestations (all required)" hint:"Attestations for this user" options:"GetAttestationOpts" pfcheckboxmode:"yes"`
	Button       string          `label:"Nominate" pftype:"submit"`
	Error        string          /* Used by pfform() */
}

func (na *NominateAdd) GetAttestationOpts(obj interface{}) (kvs keyval.KeyVals, err error) {
	return na.group.GetAttestationsKVS()
}

func h_group_nominate_existing(cui pu.PfUI) {
	errmsg := ""
	tctx := tr.TriGetCtx(cui)
	grp := tctx.TriSelectedGroup()

	vouchee_name, err := cui.FormValue("vouchee")
	if err != nil {
		pu.H_errtxt(cui, "No valid vouchee")
		return
	}

	err = tctx.SelectVouchee(vouchee_name, pu.PERM_USER_NOMINATE)
	if err != nil {
		pu.H_errtxt(cui, "Vouchee unselectable")
		return
	}

	if cui.IsPOST() {
		action, err := cui.FormValue("action")
		if err == nil && action == "nominate" {
			err = vouch_nominate(cui)
			if err != nil {
				errmsg = err.Error()
			}
		}
	}

	vouchee := tctx.SelectedVouchee()

	type Page struct {
		*pu.PfPage
		Vouchee     string
		GroupName   string
		NominateAdd *NominateAdd
	}

	na := &NominateAdd{group: grp, Vouchee: vouchee.GetUserName(), Action: "nominate", Error: errmsg}

	p := Page{cui.Page_def(), vouchee.GetUserName(), grp.GetGroupName(), na}
	cui.Page_show("group/nominate_existing.tmpl", p)
}

type NominateNew struct {
	group        tr.TriGroup
	Action       string          `label:"Action" pftype:"hidden"`
	Email        string          `label:"Email address of nominee" pfset:"none"`
	FullName     string          `label:"Full Name" hint:"Full Name of this user" pfreq:"yes"`
	Affiliation  string          `label:"Affiliation" hint:"Who the user is affiliated to" pfreq:"yes"`
	Biography    string          `label:"Biography" pftype:"text" hint:"Biography for this user" pfreq:"yes"`
	Vouch        string          `label:"Vouch" pftype:"text" hint:"Vouch for this user" pfreq:"yes"`
	Attestations map[string]bool `label:"Attestations (all required)" hint:"Attestations for this user" options:"GetAttestationOpts" pfcheckboxmode:"yes"`
	Button       string          `label:"Nominate" pftype:"submit"`
}

func (na *NominateNew) GetAttestationOpts(obj interface{}) (kvs keyval.KeyVals, err error) {
	return na.group.GetAttestationsKVS()
}

func h_group_nominate(cui pu.PfUI) {
	tctx := tr.TriGetCtx(cui)

	var list []pf.PfUser
	notfound := false

	user := tctx.TriSelectedUser()
	grp := tctx.TriSelectedGroup()

	if cui.IsPOST() {
		action, err := cui.FormValue("action")
		if err == nil && action == "nominate" {
			vouch_nominate_new(cui)
			return
		}
	}

	search, err := cui.FormValue("search")
	if err != nil {
		search = ""
	}

	if search != "" {
		list, err = user.GetList(cui, search, 0, 0)
		if err != nil {
			return
		}
		if len(list) == 0 {
			notfound = true
		}
	}

	type Page struct {
		*pu.PfPage
		Search    string
		GroupName string
		Message   string
		Error     string
		Users     []pf.PfUser
		NotFound  bool
		NewForm   *NominateNew
	}

	newform := &NominateNew{group: grp, Action: "nominate", Email: search}

	p := Page{cui.Page_def(), search, grp.GetGroupName(), "", "", list, notfound, newform}
	cui.Page_show("group/nominate.tmpl", p)
}

func h_group(cui pu.PfUI, menu *pu.PfUIMenu) {
	menu.Replace("members", h_group_members)

	m := []pu.PfUIMentry{
		{"nominate", "Nominate", pf.PERM_GROUP_MEMBER, h_group_nominate, nil},
		{"nominate_existing", "Nominate existing user", pf.PERM_GROUP_MEMBER | pf.PERM_HIDDEN, h_group_nominate_existing, nil},
	}

	menu.Add(m...)
}
