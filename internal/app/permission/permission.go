package permission

import "strings"

func Has(permissions []string, target string) bool {
	target = strings.TrimSpace(target)
	if target == "" {
		return true
	}

	for _, permission := range permissions {
		if strings.TrimSpace(permission) == target {
			return true
		}
	}

	return false
}
