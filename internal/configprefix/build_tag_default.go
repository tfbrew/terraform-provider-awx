//go:build !repoAWX && !repoAAP
// +build !repoAWX,!repoAAP

package configprefix

import "strings"

const Prefix = "aap"

func ReplaceText(input string) string {
	return strings.ReplaceAll(input, "awx_", "aap_")
}
