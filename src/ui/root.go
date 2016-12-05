package TriUI

import (
	pu "trident.li/pitchfork/ui"
)

func h_root(cui pu.PfUI, menu *pu.PfUIMenu) {
	menu.Add(pu.PfUIMentry{"recover", "Password Recover", pu.PERM_NONE | pu.PERM_HIDDEN, h_recover, nil})
}
