package pdf

import (
	"bytes"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/skylibdrvlz/20.11.2025-links-checker/models"
)

func GeneratePDF(linkSets []*models.LinkSet) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 10)
	for _, linkSet := range linkSets {
		if pdf.GetY() > 250 {
			pdf.AddPage()
		}

		pdf.Cell(0, 8, fmt.Sprintf("Check #%d:", linkSet.ID))
		pdf.Ln(8)

		for link, status := range linkSet.Links {
			pdf.MultiCell(0, 6, fmt.Sprintf("%s - %s", link, status), "", "", false)
		}
		pdf.Ln(8)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("PDF output failed: %w", err)
	}
	return buf.Bytes(), nil
}

func GenerateErrorPDF(errorMsg string) []byte {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, errorMsg) // просто текст ошибки
	var buf bytes.Buffer
	pdf.Output(&buf)
	return buf.Bytes()
}
