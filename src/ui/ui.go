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
				h_group(cui, menu)
			} else if lp == 3 {
				menu.Add(pu.PfUIMentry{"vouches", "Vouches", pu.PERM_GROUP_MEMBER, h_user_vouches, nil})
			}
			break
		}
	}

	return
}
