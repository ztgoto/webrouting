package utils

import "regexp"

var (
	// SpaceRegexp 匹配空白字符正则
	SpaceRegexp = regexp.MustCompile(`\s`)
)
