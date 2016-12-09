package TriUI

import (
	"strconv"
	"strings"
	"trident.li/keyval"
	pf "trident.li/pitchfork/lib"
	pu "trident.li/pitchfork/ui"
	tr "trident.li/trident/src/lib"
)

func h_group_member(cui pu.PfUI) {
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

	total, err = grp.ListGroupMembersTot(search)
	if err != nil {
		cui.Err("error: " + err.Error())
		return
	}

	members, err := grp.ListGroupMembers(search, cui.TheUser().GetUserName(), offset, 10, false, cui.IAmGroupAdmin(), false)
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
	Message      string          /* Used by pfform() */
	Error        string          /* Used by pfform() */
}

func (na *NominateAdd) GetAttestationOpts(obj interface{}) (kvs keyval.KeyVals, err error) {
	return na.group.GetAttestationsKVS()
}

func h_group_nominate_existing(cui pu.PfUI) {
	msg := ""
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
			msg, err = vouch_nominate(cui)
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

	na := &NominateAdd{group: grp, Vouchee: vouchee.GetUserName(), Action: "nominate", Message: msg, Error: errmsg}

	p := Page{cui.Page_def(), vouchee.GetUserName(), grp.GetGroupName(), na}
	cui.Page_show("group/nominate_existing.tmpl", p)
}

type NominateNew struct {
	group        tr.TriGroup
	Action       string          `label:"Action" pftype:"hidden"`
	Search       string          `label:"Search" pftype:"hidden"`
	Email        string          `label:"Email address of nominee" pfset:"none"`
	FullName     string          `label:"Full Name" hint:"Full Name of this user" pfreq:"yes"`
	Affiliation  string          `label:"Affiliation" hint:"Who the user is affiliated to" pfreq:"yes"`
	Biography    string          `label:"Biography" pftype:"text" hint:"Biography for this user" pfreq:"yes"`
	Comment      string          `label:"Vouch Comment" pftype:"text" hint:"Vouch for this user" pfreq:"yes"`
	Attestations map[string]bool `label:"Attestations (all required)" hint:"Attestations for this user" options:"GetAttestationOpts" pfcheckboxmode:"yes"`
	Button       string          `label:"Nominate" pftype:"submit"`
	Message      string          /* Used by pfform() */
	Error        string          /* Used by pfform() */
}

func (na *NominateNew) GetAttestationOpts(obj interface{}) (kvs keyval.KeyVals, err error) {
	return na.group.GetAttestationsKVS()
}

func h_group_nominate(cui pu.PfUI) {
	var msg string
	var err error
	var errmsg string
	var list []pf.PfUser
	var search string

	tctx := tr.TriGetCtx(cui)
	user := tctx.TriSelectedUser()
	grp := tctx.TriSelectedGroup()

	/* Something posted? */
	if cui.IsPOST() {
		/* An action to perform? */
		action, err := cui.FormValue("action")
		if err == nil && action == "nominate" {
			msg, err = vouch_nominate_new(cui)
			if err != nil {
				errmsg += err.Error()
			}
		}

		/* Search field? */
		search, err = cui.FormValue("search")
		if err != nil {
			search = ""
		}

		/* Simple 'is there an @ sign, it must be an email address' check */
		if strings.Index(search, "@") == -1 {
			/* Not an email, do not allow searching */
			search = ""
		}
	}

	/* Need to search the list? */
	notfound := true
	if search != "" {
		/* Get list of users matching the given search query */
		list, err = user.GetList(cui, search, 0, 0, true)
		if err != nil {
			cui.Errf("Listing users failed: %s", err.Error())
			pu.H_error(cui, pu.StatusBadRequest)
			return
		}

		if len(list) != 0 {
			notfound = false
		}
	}

	type Page struct {
		*pu.PfPage
		Search    string
		GroupName string
		Users     []pf.PfUser
		NotFound  bool
		NewForm   *NominateNew
	}

	/* Re-fill in the form (for people who do not enable the attestations) */
	descr, _ := cui.FormValue("fullname")
	affil, _ := cui.FormValue("affiliation")
	bio, _ := cui.FormValue("biography")
	comment, _ := cui.FormValue("comment")

	newform := &NominateNew{group: grp, Action: "nominate", Email: search, Message: msg, Error: errmsg, Search: search, FullName: descr, Affiliation: affil, Biography: bio, Comment: comment}

	p := Page{cui.Page_def(), search, grp.GetGroupName(), list, notfound, newform}
	cui.Page_show("group/nominate.tmpl", p)
}

func h_vouches_csv(cui pu.PfUI) {
	grp := cui.SelectedGroup()

	vouches, err := tr.Vouches_Get(cui, grp.GetGroupName())
	if err != nil {
		pu.H_errmsg(cui, err)
		return
	}

	csv := ""

	for _, v := range vouches {
		csv += v.Vouchor + "," + v.Vouchee + "," + v.Entered.Format(pf.Config.DateFormat) + "\n"
	}

	fname := grp.GetGroupName() + ".csv"

	cui.SetContentType("text/vcard")
	cui.SetFileName(fname)
	cui.SetExpires(60)
	cui.SetRaw([]byte(csv))
	return
}

func h_vouches(cui pu.PfUI) {
	fmt := cui.GetArg("format")

	if fmt == "csv" {
		h_vouches_csv(cui)
		return
	}

	grp := cui.SelectedGroup()
	vouches, err := tr.Vouches_Get(cui, grp.GetGroupName())
	if err != nil {
		pu.H_errmsg(cui, err)
		return
	}

	/* Output the page */
	type Page struct {
		*pu.PfPage
		Vouches []tr.Vouch
	}

	p := Page{cui.Page_def(), vouches}
	cui.Page_show("group/vouches.tmpl", p)
}

func h_group(cui pu.PfUI, menu *pu.PfUIMenu) {
	menu.Replace("member", h_group_member)

	m := []pu.PfUIMentry{
		{"nominate", "Nominate", pf.PERM_GROUP_MEMBER, h_group_nominate, nil},
		{"nominate_existing", "Nominate existing user", pf.PERM_GROUP_MEMBER | pf.PERM_HIDDEN, h_group_nominate_existing, nil},
		{"vouches", "Vouches", pf.PERM_GROUP_MEMBER, h_vouches, nil},
	}

	menu.Add(m...)
}
