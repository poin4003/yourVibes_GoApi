package truncate

import "strings"

func TruncateContent(content string, wordLimit int) string {
	words := strings.Fields(content)
	if len(words) <= wordLimit {
		return content
	}

	return strings.Join(words[:wordLimit], " ") + "..."
}
