package TriUI

import (
	"errors"
	"html/template"
	pf "trident.li/pitchfork/lib"
	pu "trident.li/pitchfork/ui"
)

func h_recover(cui pu.PfUI) {
	var msg string
	var err error
	var usr string
	var usp string
	var nmp string

	if cui.IsPOST() {
		var err2 error
		var err3 error

		cmd := "user password recover"
		arg := []string{"", "", ""}

		usr, err = cui.FormValue("username")
		usp, err2 = cui.FormValue("user")
		nmp, err3 = cui.FormValue("nominator")
		pw1, err4 := cui.FormValue("password")
		pw2, err5 := cui.FormValue("passwordr")
		arg[1] = usp + nmp

		if err != nil && err2 != nil && err3 != nil || err4 != nil || err5 != nil {
			if pw1 != "" {
				if pw1 == pw2 {
					msg, err = cui.HandleCmd(cmd, arg)
					if err == nil {
						/* Reset form values */
						usp = ""
						nmp = ""
					}
				} else {
					err = errors.New("Passwords did not match")
				}
			}
		}
	}

	var errmsg = ""

	if err != nil {
		/* Failed */
		errmsg = err.Error()
	} else {
		/* Success */
	}

	/* Output the page */
	type rec struct {
		Username  string `label:"Username" hint:"Your username" pfmin:"3" pfreq:"yes"`
		User      string `label:"Token: User Portion" hint:"The user portion of the recovery token" pfmin:"16" pfmax:"16" pfreq:"yes"`
		Nominator string `label:"Token: Nominator Portion" hint:"The nominator portion of the recovery token" pfmin:"8" pfmax:"16" pfreq:"yes"`
		Password  string `label:"New Password" hint:"The new password to set" pfmin:"6" pftype:"password" pfreq:"yes"`
		PasswordR string `label:"Repeat Password" hint:"Repeat the password so that one knows it is the same as the other" pfmin:"6" pftype:"password" pfreq:"yes"`
		Button    string `label:"Recover password" pftype:"submit"`
	}

	type Page struct {
		*pu.PfPage
		Intro   template.HTML
		Recover rec
		Message string
		Error   string
	}

	intro := pf.HEB("<p>\n" +
		"This form can be used after receiving both the user and nominator portions of the password recovery procedure.\n" +
		"</p>\n")

	p := Page{cui.Page_def(), intro, rec{usr, usp, nmp, "", "", ""}, msg, errmsg}
	cui.Page_show("misc/recover.tmpl", p)
}
