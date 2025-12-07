package document

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
)

// ExtractPDFText extracts text content from a PDF file
func ExtractPDFText(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	totalPages := r.NumPage()
	var textBuilder strings.Builder

	// Limit to first 50 pages to avoid excessive processing
	maxPages := totalPages
	if maxPages > 50 {
		maxPages = 50
	}

	for pageNum := 1; pageNum <= maxPages; pageNum++ {
		p := r.Page(pageNum)
		if p.V.IsNull() {
			continue
		}

		text, err := p.GetPlainText(nil)
		if err != nil {
			continue // Skip pages with errors
		}

		textBuilder.WriteString(fmt.Sprintf("\n--- Page %d ---\n", pageNum))
		textBuilder.WriteString(text)
	}

	result := textBuilder.String()
	
	// Limit total text length
	if len(result) > 10000 {
		result = result[:10000] + "\n\n[內容過長，已截斷]"
	}

	return result, nil
}

// ReadTextFile reads a text file
func ReadTextFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	text := string(data)
	
	// Limit text length
	if len(text) > 10000 {
		text = text[:10000] + "\n\n[內容過長，已截斷]"
	}

	return text, nil
}
