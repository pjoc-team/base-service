package url

import "strings"

func CompactUrl(domain string, path string, parameters string) string {
	builder := strings.Builder{}
	if domain != "" {
		builder.WriteString(strings.TrimSuffix(domain, "/"))
		builder.WriteString("/")
	}
	if path != "" {
		path = strings.TrimPrefix(path, "/")
		path = strings.TrimSuffix(path, "?")
		builder.WriteString(path)
	}
	if parameters != "" {
		builder.WriteString("?")
		builder.WriteString(strings.TrimPrefix(parameters, "?"))
	}
	return builder.String()
}
