package files

import "strings"

func EllipsisFront(path string, maxLength int) string {
	if len(path) <= maxLength {
		return path
	}
	parts := strings.Split(path, "/")
	
	currentParts := []string{}
	currentPartsLength := 0
	for i := len(parts) - 1; i >= 0; i-- {
		newPart := parts[i]
		if currentPartsLength+len(newPart)+1 > maxLength {
			break
		}
		currentParts = append(currentParts, newPart)
		currentPartsLength += len(newPart) + 1
	}

	for i, j := 0, len(currentParts)-1; i < j; i, j = i+1, j-1 {
		currentParts[i], currentParts[j] = currentParts[j], currentParts[i]
	}

	joinedPath := strings.Join(currentParts, "/")

	ellipsis := ".../"
	if len(joinedPath) > maxLength {
		joinedPath = joinedPath[len(joinedPath)-maxLength:]
	}
	return ellipsis + joinedPath
}
