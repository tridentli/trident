package TriUI

import (
	pu "trident.li/pitchfork/ui"
	tf "trident.li/trident/src/lib"
)

type TriUI interface {
	pu.PfUIi
	tf.TriCtx
}

type TriUIi interface {
	pu.PfUIi
}

type TriUIS struct {
	pu.PfUI
}

func NewTriUI() pu.PfUI {
	tctx := tf.NewTriCtx()
	pfui := pu.NewPfUI(tctx, nil, TriUIMenuOverride, nil)
	tcui := &TriUIS{PfUI: pfui}
	return tcui
}

func TriGetUI(cui pu.PfUI) TriUI {
	return cui.(TriUI)
}

func TriUIMenuOverride(cui pu.PfUI, menu *pu.PfUIMenu) {
	loc := cui.GetLoc()

	switch loc {
	case "":
		h_root(cui, menu)

	case "group":
		h_group(cui, menu)
		break

	case "user":
		h_user(cui, menu)
		break
	}

	return
}
