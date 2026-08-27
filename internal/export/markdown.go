package export

import (
	"fmt"
	"os"
	"strings"

	"github.com/mooreceipts/stendoclip/internal/store"
)

func ToMarkdown(stack *store.ClippingStack, outputPath string) error {
	sb := strings.Builder{}
	sb.WriteString("# Stendoclip Clipboard History\n\n")
	sb.WriteString("| Pinned | Text | Timestamp |\n")
	sb.WriteString("|--------|------|----------|\n")

	for i := 0; i < stack.Len(); i++ {
		entry := stack.Get(i)
		pinned := "✓"
		if !entry.Pinned {
			pinned = ""
		}

		// Preserve full text while keeping a valid Markdown table.
		text := strings.ReplaceAll(entry.Text, "|", "\\|")
		text = strings.ReplaceAll(text, "\r\n", "<br>")
		text = strings.ReplaceAll(text, "\n", "<br>")

		ts := entry.Timestamp.Format("2006-01-02 15:04:05")
		sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", pinned, text, ts))
	}

	return os.WriteFile(outputPath, []byte(sb.String()), 0o644)
}
