package otelawslambda

func contentLength(maybeBase64content string, isBase64 bool) int {
	if isBase64 {
		padding := 0
		for i := len(maybeBase64content) - 1; i >= 0 && i > len(maybeBase64content)-3; i-- {
			if maybeBase64content[i] == '=' {
				padding++
			} else {
				break
			}
		}
		return (len(maybeBase64content) - padding) * 3 / 4
	}
	return len(maybeBase64content)
}
