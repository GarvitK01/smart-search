package parser

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fumiama/go-docx"
	xtract "github.com/sassoftware/pdf-xtract"
)

func extractJSON(filePath string) (string, error) {

	// Replace all newlines with spaces
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(content), "\n", " "), nil
}

func extractPDF(filePath string) (string, error) {

	cfg := xtract.NewDefaultConfig()
	proc := xtract.NewProcessor(cfg)
	text, _, err := proc.Extract(context.Background(), filePath)
	if err != nil {
		return "", err
	}
	return text, nil
}

func extractDOCX(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open docx: %w", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return "", fmt.Errorf("stat docx: %w", err)
	}

	doc, err := docx.Parse(f, info.Size())
	if err != nil {
		return "", fmt.Errorf("parse docx: %w", err)
	}

	var paragraphs []string
	for _, item := range doc.Document.Body.Items {
		if p, ok := item.(*docx.Paragraph); ok {
			paragraphs = append(paragraphs, p.String())
		}
	}
	return strings.Join(paragraphs, "\n"), nil
}

func extractTXT(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func ExtractText(filePath string) (string, error) {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".txt":
		return extractTXT(filePath)
	case ".json":
		return extractJSON(filePath)
	case ".pdf":
		return extractPDF(filePath)
	case ".docx":
		return extractDOCX(filePath)
	default:
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
}
