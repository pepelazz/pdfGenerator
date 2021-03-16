package pdfGenerator

import (
	"fmt"
	"strings"
)

const (
	COLOR_BLACK           = "black"
	COLOR_DARK_BLUE       = "dark_blue"
	COLOR_LIGHT_BLUE      = "light_blue"
	COLOR_LIGHT_BLUE_TEXT = "light_blue_text"
	COLOR_WHITE           = "white"
	BORDER_WIDTH          = 0.1
	START_X               = 4
	LINE_HEIGHT           = 3.8
)

type (
	Table struct {
		Rows          []TableRow
		PaddingTop    float64
		PaddingBottom float64
		RowWidth      []float64 // в процентах
		IsBold        bool
	}

	TableRow struct {
		Cells      []TableCell
		PaddingTop float64
	}

	TableCell struct {
		Data       Text
		PaddingTop float64
		Colspan int
	}
)

func (tbl *Table) AddRow(r TableRow) {
	tbl.Rows = append(tbl.Rows, r)
}

// helper для создания простой строки из массива текстовых блоков
func (tbl *Table) AddRowSimple(pd *PdfDoc, txtList []string) {
	r := TableRow{Cells: []TableCell{}}
	for i, txt := range txtList {
		colWidth := tbl.RowWidth[i]
		// если текст превышает ширину, то разделяем его
			if pd.Pdf.GetStringWidth(txt) >= colWidth {
				lines := pd.Pdf.SplitLines([]byte(txt), colWidth)
				txtArr := make([]string, len(lines))
				for j, ln := range lines {
					txtArr[j] = fmt.Sprintf("%s", ln)
				}
				txt = strings.Join(txtArr, "\n")
		}
		r.Cells = append(r.Cells, TableCell{Data: Text{Data: txt}})
	}
	tbl.Rows = append(tbl.Rows, r)
}

// helper для создания простой строки из массива текстовых блоков без автоматических переносов
func (tbl *Table) AddRowSimpleWithoutHypernation(pd *PdfDoc, txtList []string) {
	r := TableRow{Cells: []TableCell{}}
	for _, txt := range txtList {
		r.Cells = append(r.Cells, TableCell{Data: Text{Data: txt}})
	}
	tbl.Rows = append(tbl.Rows, r)
}

func (r *TableRow) AddCell(c TableCell) {
	r.Cells = append(r.Cells, c)
}

func (tbl *Table) GetHeight(pd *PdfDoc) float64 {
	h := 0.0
	for _, r := range tbl.Rows {
		h = h + r.getHeight(pd, *tbl)
	}
	return h
}

// вычисление - выходит за границу или нет
func (tbl *Table) IsTableOutOfPage(pd *PdfDoc, headerHeight float64) bool {
	_, pH := pd.Pdf.GetPageSize()
	_, _, _, bMargin := pd.Pdf.GetMargins()
	return pd.Pdf.GetY()+tbl.GetHeight(pd)+headerHeight > (pH - bMargin)
}

func (pd *PdfDoc) DrawTable(tbl Table) error {
	pd.Pdf.SetLineWidth(BORDER_WIDTH)
	startX := pd.Pdf.GetX()
	for _, r := range tbl.Rows {
		drawTableRow(pd, tbl, r)
		pd.Pdf.SetXY(startX, pd.Pdf.GetY()+r.getHeight(pd, tbl))
	}
	return nil
}

func drawTableRow(pd *PdfDoc, tbl Table, row TableRow) {
	h := row.getHeight(pd, tbl)

	// если выходит за границы страницы, то добавляем страницу
	_, pH := pd.Pdf.GetPageSize()
	_, _, _, bMargin := pd.Pdf.GetMargins()
	if pd.Pdf.GetY()+h > (pH - bMargin) {
		pd.AddPage()
	}

	borderWidth := pd.Pdf.GetLineWidth()
	yStart := pd.Pdf.GetY()
	xStart := pd.Pdf.GetX()
	pd.Pdf.SetDrawColor(31, 46, 60)
	for i, c := range row.Cells {
		// проставляем ширину колонки
		if len(tbl.RowWidth) > i {
			c.Data.Width = tbl.RowWidth[i]
		}
		if c.Colspan > 1 {
			j := 1
			for  {
				if len(tbl.RowWidth) == i+j {
					break
				}
				c.Data.Width = c.Data.Width + tbl.RowWidth[i+j] - borderWidth
				j++
				if j == c.Colspan {
					break
				}
			}
		}
		drawTableCell(pd, tbl, row, c, h)
		xStart = xStart + c.Data.Width - borderWidth
		pd.Pdf.SetXY(xStart, yStart)
	}
}

func drawTableCell(pd *PdfDoc, tbl Table, row TableRow, c TableCell, height float64) {
	pd.Pdf.SetFillColor(255, 255, 255)
	borderWidth := pd.Pdf.GetLineWidth()
	x := pd.Pdf.GetX() + BORDER_WIDTH
	pd.Pdf.Rect(x, pd.Pdf.GetY(), c.Data.Width-borderWidth, height, "FD")
	c.Data.X = x
	// определяем padding
	paddingTop := 1.0
	if tbl.PaddingTop != 0 {
		paddingTop = tbl.PaddingTop
	}
	if row.PaddingTop != 0 {
		paddingTop = row.PaddingTop
	}
	if c.PaddingTop != 0 {
		paddingTop = c.PaddingTop
	}
	pd.AddY(paddingTop)
	if tbl.IsBold {
		c.Data.IsBold = tbl.IsBold
	}
	pd.PrintText(c.Data)
}

//func draw2cells(pdf *gofpdf.Fpdf, Tr func(string) string, nameRu, nameEn string, bgColor, borderColor, textColor string) {
//	pdf.SetX(START_X)
//	pdf.SetTextColor(255, 255, 255)
//	switch bgColor {
//	case COLOR_WHITE:
//		pdf.SetFillColor(255, 255, 255)
//		pdf.SetTextColor(0, 0, 0)
//	case COLOR_DARK_BLUE:
//		pdf.SetFillColor(31, 46, 60)
//	case COLOR_LIGHT_BLUE:
//		pdf.SetFillColor(55, 82, 105)
//	}
//	switch borderColor {
//	case COLOR_WHITE:
//		pdf.SetDrawColor(255, 255, 255)
//	case COLOR_DARK_BLUE:
//		pdf.SetDrawColor(31, 46, 60)
//	case COLOR_LIGHT_BLUE:
//		pdf.SetDrawColor(55, 82, 105)
//	}
//	pdf.SetLineWidth(BORDER_WIDTH)
//	//linesRu, linesEn, maxStr := getTextLines(nameRu, nameEn)
//
//	// в случае если текст без переносов и одной строкой
//	if maxStr == 1 {
//		pdf.CellFormat(70, 6, Tr(nameRu), "1", 0, "CM", true, 0, "")
//		if textColor == COLOR_LIGHT_BLUE_TEXT {
//			pdf.SetTextColor(118, 182, 224)
//			pdf.CellFormat(70, 6, Tr(nameEn), "1", 1, "CM", true, 0, "")
//			pdf.SetTextColor(255, 255, 255)
//		} else {
//			pdf.CellFormat(70, 6, Tr(nameEn), "1", 1, "CM", true, 0, "")
//		}
//	} else {
//		// если текст из нескольких строк
//		fSize, _ := pdf.GetFontSize()
//		htMax := fSize * float64(maxStr) * 0.5
//		yStart := pdf.GetY()
//		pdf.Rect(pdf.GetX()+BORDER_WIDTH, pdf.GetY(), 70-BORDER_WIDTH, htMax, "FD")
//		_, fontSize := pdf.GetFontSize()
//		//addY(pdf, 1)
//		pdf.SetXY(START_X, yStart+1)
//		for _, line := range linesRu {
//			pdf.CellFormat(70.0, fontSize, Tr(line), "", 1, "CM", false, 0, "")
//		}
//
//		pdf.SetX(70 + START_X)
//		pdf.Rect(pdf.GetX(), yStart, 70-BORDER_WIDTH, htMax, "FD")
//		pdf.SetXY(70 + START_X, yStart+1)
//		// печатаем строки с текстом
//		if textColor == COLOR_LIGHT_BLUE_TEXT {
//			pdf.SetTextColor(118, 182, 224)
//		}
//		for _, line := range linesEn {
//			pdf.SetX(70 + START_X)
//			pdf.CellFormat(70.0, fontSize, string(line), "", 1, "CM", false, 0, "")
//		}
//		pdf.SetTextColor(255, 255, 255)
//		pdf.SetXY(START_X, yStart+htMax)
//	}
//}

func (r *TableRow) getHeight(pd *PdfDoc, tbl Table) float64 {
	maxHeight := 0.0
	for _, c := range r.Cells {
		if c.Data.Height == 0 {
			c.Data.CalcHeight(pd)
		}
		h := c.Data.Height
		// добавляем padding
		pt := tbl.PaddingTop
		if r.PaddingTop > 0 {
			pt = r.PaddingTop
		}
		if c.PaddingTop > 0 {
			pt = c.PaddingTop
		}
		h = h + pt


		if h > maxHeight {
			maxHeight = h
		}
	}
	maxHeight = maxHeight + tbl.PaddingBottom
	return maxHeight
}
