/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package db

func getWhereQuery(condition string, args ...interface{}) string {
	withoutDeletedAtIsNull := false

	for _, arg := range args {
		switch v := arg.(type) {
		case FetchOpt:
			if v.WithoutDeletedAtIsNull {
				withoutDeletedAtIsNull = true
				break
			}

		case *FetchOpt:
			if v != nil {
				if v.WithoutDeletedAtIsNull {
					withoutDeletedAtIsNull = true
					break
				}
			}
		}
	}

	whereQuery := ""
	if !withoutDeletedAtIsNull {
		whereQuery = "WHERE deleted_at IS NULL"
		if condition != "" {
			whereQuery += " AND " + condition
		}
	} else if condition != "" {
		whereQuery += "WHERE " + condition
	}

	return whereQuery
}

func getEndQuery(args ...interface{}) string {
	endQuery := ""

	for _, arg := range args {
		switch v := arg.(type) {
		case FetchOpt:
			endQuery = appendEndQuery(endQuery, v.EndQuery)
		case *FetchOpt:
			if v != nil {
				endQuery = appendEndQuery(endQuery, v.EndQuery)
			}
		}
	}

	return endQuery
}

func appendEndQuery(base, add string) string {
	if add == "" {
		return base
	}

	if base != "" {
		base += " "
	}
	return base + add
}

func getPars(args ...interface{}) []interface{} {
	pars := make([]interface{}, 0)

	for _, arg := range args {
		switch v := arg.(type) {
		case FetchOpt:
			if v.ForceRW {
				pars = append(pars, srFetchOpt{
					rwForce: true,
				})
			}

		case *FetchOpt:
			if v.ForceRW {
				pars = append(pars, srFetchOpt{
					rwForce: true,
				})
			}

		default:
			pars = append(pars, arg)
		}
	}

	return pars
}

func parsAndOthers(args ...interface{}) (bool, []interface{}) {
	var (
		rwForce = false
		pars    = make([]interface{}, 0)
	)

	for _, arg := range args {
		switch v := arg.(type) {
		case FetchOpt:
			if v.ForceRW {
				rwForce = v.ForceRW
			}

		case *FetchOpt:
			if v.ForceRW {
				rwForce = v.ForceRW
			}

		case srFetchOpt:
			if v.rwForce {
				rwForce = v.rwForce
			}

		case *srFetchOpt:
			if v.rwForce {
				rwForce = v.rwForce
			}

		default:
			pars = append(pars, arg)
		}
	}

	return rwForce, pars
}
