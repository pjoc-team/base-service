package sign

import "github.com/pjoc-team/base-service/pkg/util"

func GenerateMd5Key(length int) string {
	return util.RandString(length)
}

func GenerateMd5KeyWith32Word() string {
	return GenerateMd5Key(32)
}
