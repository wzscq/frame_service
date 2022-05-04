package definition

import (
	"strings"
)

func HasRight(roles *interface{},userRoles string)(bool){
	if roles == nil {
		return false
	}

	userRoles=","+userRoles+","
	rolesStr,ok:=(*roles).(string)
	if ok {
		if rolesStr == "*" {
			return true
		}

		if strings.Contains(userRoles,","+rolesStr+",") {
			return true
		}
		
		return false
	}

	rolesArr,ok:=(*roles).([]interface{})
	if ok {
		for idx:=range rolesArr {
			rolesStr,ok:=(rolesArr[idx]).(string)
			if ok {
				if strings.Contains(userRoles,","+rolesStr+",") {
					return true
				}
			}
		}
	}

	return false
}