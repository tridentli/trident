package TriUI

import (
	"strconv"
	pf "trident.li/pitchfork/lib"
	pu "trident.li/pitchfork/ui"
)

func h_group_vcp(cui pu.PfUI) {
	criterias := []string{"Unmarked", "Dunno", "Vouched"}
	limits := []int{10, 25, 50}

	criteria, err := cui.FormValue("criteria")
	if err != nil || (criteria != "Unmarked" && criteria != "Dunno" && criteria != "Vouched") {
		criteria = "Unmarked"
	}

	limit := 25
	limit_v, err := cui.FormValue("limit")
	if err == nil {
		limit_t, err2 := strconv.Atoi(limit_v)
		if err2 == nil && limit_t < 51 {
			limit = limit_t
		}
	}

	offset := 0
	offset_v, err := cui.FormValue("offset")
	if err == nil {
		offset_t, err2 := strconv.Atoi(offset_v)
		if err2 == nil && offset_t < 50 {
			offset = offset_t
		}
	}

	if cui.IsPOST() {
		marks, err := cui.FormValueM("marked[]")
		if err == nil {
			for _, m := range marks {
				q := ""

				if criteria == "Unmarked" {
					q = "INSERT INTO member_vouch " +
						"(vouchor, vouchee, trustgroup, positive) " +
						"VALUES ($1, $2, $3, FALSE)"
				} else {
					positivity := "positive"

					if criteria == "Dunno" {
						positivity = "NOT positive"
					} else if criteria == "Vouched" {
						positivity = "positive"
					} else {
						cui.Errf("Unknown Criteria: %s", criteria)
						pu.H_errtxt(cui, "Invalid")
						return
					}

					q = "DELETE FROM member_vouch mv " +
						"WHERE mv.vouchor = $1 " +
						"AND mv.vouchee = $2 " +
						"AND mv.trustgroup = $3 " +
						"AND " + positivity
				}

				err = pf.DB.Exec(cui, "Update vouch: vouchor: $1, vouchee: $2, group: $3", 1, q, cui.TheUser().GetUserName(), m, cui.SelectedGroup().GetGroupName())
			}
		}
	}

	type VCP struct {
		UserName    string
		FullName    string
		Affiliation string
	}

	/* Criteria "Unmarked" */
	and_where := "AND mv.vouchee IS NULL"
	action := "Dunno"

	switch criteria {
	case "Dunno":
		and_where = "AND mv.vouchee IS NOT NULL AND NOT mv.positive"
		action = "Reconsider"
		break

	case "Vouched":
		and_where = "AND mv.vouchee IS NOT NULL AND mv.positive"
		action = "Unvouch"
		break
	}

	total := 0
	q := "SELECT COUNT(*) " +
		"FROM member m " +
		"JOIN member_trustgroup mt ON (mt.member = m.ident " +
		"  AND mt.trustgroup = $1 " +
		"  AND mt.member <> $2) " +
		"JOIN member_state ms ON (ms.ident = mt.state) " +
		"LEFT OUTER JOIN member_vouch mv " +
		"ON (mv.trustgroup = mt.trustgroup " +
		" AND mv.vouchee = m.ident " +
		" AND mv.vouchor = $2) " +
		"WHERE NOT ms.hidden " + and_where
	err = pf.DB.QueryRow(q, cui.SelectedGroup().GetGroupName(), cui.TheUser().GetUserName()).Scan(&total)
	if err != nil {
		pu.H_errtxt(cui, "Query broken")
		return
	}

	q = "SELECT m.ident, m.descr, m.affiliation " +
		"FROM member m " +
		"JOIN member_trustgroup mt ON (mt.member = m.ident " +
		"  AND mt.trustgroup = $1 " +
		"  AND mt.member <> $2) " +
		"JOIN member_state ms ON (ms.ident = mt.state) " +
		"LEFT OUTER JOIN member_vouch mv " +
		"ON (mv.trustgroup = mt.trustgroup " +
		" AND mv.vouchee = m.ident " +
		" AND mv.vouchor = $2) " +
		"WHERE NOT ms.hidden " + and_where + " " +
		"ORDER BY mt.entered " +
		"LIMIT $3 OFFSET $4"
	rows, err := pf.DB.Query(q, cui.SelectedGroup().GetGroupName(), cui.TheUser().GetUserName(), limit, offset)
	if err != nil {
		pu.H_errmsg(cui, err)
		return
	}

	defer rows.Close()

	var vcps []VCP
	for rows.Next() {
		var vcp VCP

		err = rows.Scan(&vcp.UserName, &vcp.FullName, &vcp.Affiliation)
		if err != nil {
			return
		}

		vcps = append(vcps, vcp)
	}

	/* Output the page */
	type Page struct {
		*pu.PfPage
		PagerOffset int
		PagerTotal  int
		GroupName   string
		Members     []VCP
		Action      string
		Criteria    string
		Criterias   []string
		Limit       int
		Limits      []int
	}

	p := Page{cui.Page_def(), offset, total, cui.SelectedGroup().GetGroupName(), vcps, action, criteria, criterias, limit, limits}
	cui.Page_show("group/vcp.tmpl", p)
}
