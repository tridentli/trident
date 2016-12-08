package trident

import (
	pf "trident.li/pitchfork/lib"
)

type TriGroupMember interface {
	pf.PfGroupMember
	GetVouchesFor() int
	GetVouchesBy() int
	GetVouchesForMe() int
	GetVouchesByMe() int
}

type TriGroupMemberS struct {
	pf.PfGroupMember
	VouchesFor   int
	VouchesBy    int
	VouchesForMe int
	VouchesByMe  int
}

func (o *TriGroupMemberS) GetVouchesFor() int {
	return o.VouchesFor
}

func (o *TriGroupMemberS) GetVouchesBy() int {
	return o.VouchesBy
}

func (o *TriGroupMemberS) GetVouchesForMe() int {
	return o.VouchesForMe
}

func (o *TriGroupMemberS) GetVouchesByMe() int {
	return o.VouchesByMe
}

func NewTriGroupMember() TriGroupMember {
	pfg := pf.NewPfGroupMember()
	return &TriGroupMemberS{PfGroupMember: pfg}
}
