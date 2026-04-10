package reader

import (
	"fmt"
	"strings"
)

func BuildPrompt(files []File, context string) string {
	var sb strings.Builder

	if context != "" {
		fmt.Fprintf(&sb, "## Product Context\n%s\n\n", context)
	}

	fmt.Fprintf(&sb, "## Project Files\n\n")

	for _, f := range files {
		fmt.Fprintf(&sb, "### %s\n```\n%s\n```\n\n", f.Path, f.Content)
	}

	fmt.Fprintf(&sb, `## Task
Based on the project files above, suggest:
1. New features that would add value to this product
2. Technical improvements (performance, security, refactoring)
3. UX improvements

For each suggestion, explain why it makes sense based on what you saw in the code.`)

	return sb.String()
}
