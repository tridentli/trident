package TriUI

import (
	"strconv"
	pf "trident.li/pitchfork/lib"
	pu "trident.li/pitchfork/ui"
	tr "trident.li/trident/src/lib"
)

func h_group_members(cui pu.PfUI) {
	var err error

	tcui := TriGetUI(cui)

	total := 0
	offset := 0

	offset_v := tcui.GetArg("offset")
	if offset_v != "" {
		offset, _ = strconv.Atoi(offset_v)
	}

	search := tcui.GetArg("search")

	grp := tcui.TriSelectedGroup()

	total, err = grp.GetMembersTot(search)
	if err != nil {
		cui.Err("error: " + err.Error())
		cui.Err("error: " + err.Error())
		return
	}

	members, err := grp.GetMembers(search, cui.TheUser().GetUserName(), offset, 10, false, false)
	if err != nil {
		cui.Err(err.Error())
		return
	}

	/* Output the page */
	type Page struct {
		*pu.PfPage
		Group       pf.PfGroup
		Users       []pf.PfGroupMember
		PagerOffset int
		PagerTotal  int
		Search      string
		IsAdmin     bool
	}
	isadmin := cui.IAmGroupAdmin()

	p := Page{cui.Page_def(), grp, members, offset, total, search, isadmin}
	cui.Page_show("group/members.tmpl", p)
}

func h_group_nominate(cui pu.PfUI) {
	tcui := TriGetUI(cui)

	var err error
	var list []pf.PfUser
	var notfound bool

	user := tcui.TriSelectedUser()
	grp := tcui.TriSelectedGroup()

	attestations, err := grp.GetAttestations()
	if err != nil {
		cui.Err(err.Error())
		return
	}

	search := cui.GetArg("search")
	notfound = false
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
		Search       string
		Ident        string
		GroupName    string
		Message      string
		Error        string
		Users        []pf.PfUser
		NotFound     bool
		Attestations []tr.TriGroupAttestation
	}

	p := Page{cui.Page_def(), search, user.GetUserName(), grp.GetGroupName(), "", "", list, notfound, attestations}
	cui.Page_show("group/nominate.tmpl", p)
}

func h_group(cui pu.PfUI, menu *pu.PfUIMenu) {
	menu.Replace("members", h_group_members)
	menu.Add(pu.PfUIMentry{"vouch", "Vouch", pf.PERM_GROUP_MEMBER | pf.PERM_HIDDEN, h_vouch, nil})
	menu.Add(pu.PfUIMentry{"nominate", "Nominate", pf.PERM_GROUP_MEMBER, h_group_nominate, nil})
}
