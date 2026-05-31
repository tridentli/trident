package trident

import (
	"trident.li/keyval"
	pf "trident.li/pitchfork/lib"
)

type TriGroupAttestation struct {
	Ident      string
	Descr      string
	TrustGroup string
}

func (grpa *TriGroupAttestation) toString() (out string) {
	out = "{" + grpa.Ident + "} " + grpa.Descr
	return
}

func (grp *TriGroupS) GetAttestations() (output []TriGroupAttestation, err error) {
	q := "SELECT at.ident, " +
		"at.descr, " +
		"at.trustgroup " +
		"FROM attestations at " +
		"WHERE at.trustgroup = $1"
	rows, err := pf.DB.Query(q, grp.GetGroupName())
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var at TriGroupAttestation
		err = rows.Scan(&at.Ident, &at.Descr, &at.TrustGroup)
		if err != nil {
			return
		}

		output = append(output, at)
	}

	return
}

func (grp *TriGroupS) GetAttestationsKVS() (kvs keyval.KeyVals, err error) {
	attestations, err := grp.GetAttestations()
	if err != nil {
		return
	}

	for _, att := range attestations {
		kvs.Add(att.Ident, att.Descr)
	}

	return
}

func grp_attestation_list(ctx pf.PfCtx, args []string) (err error) {
	groupname := args[0]

	tctx := TriGetCtx(ctx)

	err = ctx.SelectGroup(groupname, pf.PERM_GROUP_MEMBER)
	if err != nil {
		return
	}

	grp := tctx.TriSelectedGroup()

	ats, err := grp.GetAttestations()
	var at TriGroupAttestation
	for _, at = range ats {
		ctx.OutLn(at.toString())
	}

	return
}

func grp_attestation_add(ctx pf.PfCtx, args []string) (err error) {
	groupname := args[0]
	ident := args[1]
	descr := args[2]

	err = ctx.SelectGroup(groupname, pf.PERM_GROUP_ADMIN)
	if err != nil {
		return
	}

	tctx := TriGetCtx(ctx)
	grp := tctx.TriSelectedGroup()

	q := "INSERT INTO attestations " +
		"(ident, descr, trustgroup) " +
		"VALUES($1, $2, $3)"
	err = pf.DB.Exec(ctx,
		"Added new attestation ($1,$2,$3)",
		1, q,
		ident, descr, grp.GetGroupName())
	if err != nil {
		return
	}

	return
}

func grp_attestation_delete(ctx pf.PfCtx, args []string) (err error) {
	groupname := args[0]
	ident := args[1]
	descr := args[2]

	err = ctx.SelectGroup(groupname, pf.PERM_GROUP_ADMIN)
	if err != nil {
		return
	}

	tctx := TriGetCtx(ctx)
	grp := tctx.TriSelectedGroup()

	q := "DELETE FROM attestations " +
		"WHERE ident = $1 " +
		"AND descr = $2 " +
		"AND trustgroup = $3"
	err = pf.DB.Exec(ctx,
		"Removed attestation ($1,$2,$3)",
		1, q,
		ident, descr, grp.GetGroupName())
	if err != nil {
		return
	}

	return
}

func attestation_menu(ctx pf.PfCtx, args []string) (err error) {
	menu := pf.NewPfMenu([]pf.PfMEntry{
		{Cmd: "list", Fun: grp_attestation_list, Args_min: 1, Args_max: 1, Args: []string{"groupname"}, Perms: pf.PERM_USER, Desc: "List attestations"},
		{Cmd: "add", Fun: grp_attestation_add, Args_min: 3, Args_max: 3, Args: []string{"groupname", "ident", "Description"}, Perms: pf.PERM_GROUP_ADMIN, Desc: "Add an attestation"},
		{Cmd: "delete", Fun: grp_attestation_delete, Args_min: 2, Args_max: 2, Args: []string{"groupname", "ident"}, Perms: pf.PERM_GROUP_ADMIN, Desc: "Delete an attestation"},
	})

	err = ctx.Menu(args, menu)
	return
}
