package TriUI

import (
	pu "trident.li/pitchfork/ui"
	tf "trident.li/trident/src/lib"
)

func NewTriUI() pu.PfUI {
	tctx := tf.NewTriCtx()
	pfui := pu.NewPfUI(tctx, nil, TriUIMenuOverride, nil)
	return pfui
}

func TriUIMenuOverride(cui pu.PfUI, menu *pu.PfUIMenu) {
	path := cui.GetCrumbParts()

	lp := len(path)

	if lp == 0 {
		h_root(cui, menu)
		return
	}

	if lp >= 1 {
		switch path[0] {
		case "group":
			if lp == 2 {
				/* group/<groupname>/ */
				h_group(cui, menu)

			} else if lp >= 4 && path[2] == "member" {
				/* /group/<groupname>/member/<username> */
				menu.Replace("", h_user_vouches)
				menu.Replace("profile", h_user_vouches)
				menu.Filter([]string{"", "profile", "pwreset", "log", "pgp_keys", "email"})
			}
			break
		}
	}

	return
}
