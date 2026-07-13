package lab

import "regexp"

var sensitiveValue = regexp.MustCompile(`(?i)(token|secret|private[_-]?key|api[_-]?key|rpc[_-]?url)=([^\s]+)`)

func SanitizeLog(value string) string {
	return sensitiveValue.ReplaceAllString(value, "$1=[REDACTED]")
}
