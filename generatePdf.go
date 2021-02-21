package pdfGenerator

//
//import (
//	"bytes"
//	"fmt"
//	"github.com/jung-kurt/gofpdf"
//	"go/build"
//	"os"
//	"strings"
//)
//
//const (
//	COLOR_BLACK  = "black"
//	COLOR_DARK_BLUE  = "dark_blue"
//	COLOR_LIGHT_BLUE = "light_blue"
//	COLOR_LIGHT_BLUE_TEXT = "light_blue_text"
//	COLOR_WHITE      = "white"
//	BORDER_WIDTH     = 0.1
//	START_X = 4
//	LINE_HEIGHT = 3.8
//)
//
//func generatePdf() ([]byte, string, error) {
//
//	var fileName string
//
//	// находим путь для папки с шрифтами
//	gopath := os.Getenv("GOPATH")
//	fmt.Printf("gopath %s\n", gopath)
//	if gopath == "" {
//		gopath = build.Default.GOPATH
//	}
//	fontPath := fmt.Sprintf("%s/src/github.com/pepelazz/pdfGenerator/font/", gopath)
//	fmt.Printf("fontPath %s\n", fontPath)
//
//	//pwd, err := os.Getwd()
//	//fmt.Printf("pwd %s\n", pwd)
//	//if err != nil {
//	//	return nil, fileName, err
//	//}
//	//pdf := gofpdf.New("P", "mm", "A4", pwd+"/forPrintPdf/font/")
//	pdf := gofpdf.New("P", "mm", "A4", fontPath)
//	//fmt.Printf("w: %v h: %v\n", pW, pH)
//	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
//	pdf.SetFont("Helvetica", "", 8)
//	tr := pdf.UnicodeTranslatorFromDescriptor("cp1251")
//
//	// печать страниц
//	pdf.SetFooterFunc(func() {
//		pdf.SetY(-10)
//		pdf.SetFont("Arial", "", 8)
//		pdf.SetTextColor(0, 0, 0)
//		pdf.CellFormat(0, 10, fmt.Sprintf("%v", pdf.PageNo()),
//			"", 0, "R", false, 0, "")
//	})
//
//	pdf.SetLineWidth(BORDER_WIDTH)
//
//	// первая страница
//	pW, _ := pdf.GetPageSize()
//	pdf.AddPage()
//
//
//	// печать текста
//	printText := func(text []string, w float64, alignStr string) {
//		startX := pdf.GetX()
//		for _, line := range text {
//			pdf.SetXY(startX, pdf.GetY())
//			pdf.CellFormat(w, LINE_HEIGHT, tr(line), "", 1, alignStr, false, 0, "")
//		}
//	}
//
//	pdf.ImageOptions("./forPrintPdf/footage/global_survey_logo.jpg", pW*3/11, pdf.GetY(), pW*5/11, 0, false, gofpdf.ImageOptions{ImageType: "jpg"}, 0, "")
//
//	setXStart(pdf)
//	addY(pdf, 30)
//
//	printText([]string{
//		"ООО «ГЛОБАЛ СЮРВЕЙ» ИНН 5047216117 тел.: 8 (495) 005-42-25 Россия, Московская область, г. Химки, Ленинградская",
//		"улица, д. 29, офис 911/2 www.globalsurvey.ru e-mail: info@globalsurvey.ruGLOBAL SURVEY LLC Tel.: 8 (495) 005-42-25 Russia,",
//		"Moscow region, Khimki, Leningradskaya street, 29, office 911/2 www.globalsurvey.ru e-mail: info@globalsurvey.ru.",
//	}, pW, "CM")
//
//	setXStart(pdf)
//	addY(pdf, 10)
//
//	pdf.SetFontSize(13)
//	printText([]string{"СЮРВЕЙЕРСКИЙ ОТЧЕТ / SURVEY REPORT"}, pW, "CM")
//
//
//	// добавление перевернутой страницы
//	//pdf.AddPageFormat("L", gofpdf.SizeType{pW, pH})
//
//	buf := new(bytes.Buffer)
//	err := pdf.Output(buf)
//	if err != nil {
//		return nil, fileName, err
//	}
//
//	fileName = "report.pdf"
//
//	return buf.Bytes(), fileName, nil
//}
//
//func setXStart(pdf *gofpdf.Fpdf)  {
//	pdf.SetXY(0, pdf.GetY())
//}
//
//func addY(pdf *gofpdf.Fpdf, v float64) {
//	pdf.SetXY(pdf.GetX(), pdf.GetY()+v)
//}
//
//func getTextLines(vRu, vEn string) (linesRu, linesEn []string, h int) {
//	linesRu = strings.Split(vRu, "\n")
//	linesEn = strings.Split(vEn, "\n")
//	h = len(linesRu)
//	if len(linesEn) > len(linesRu) {
//		h = len(linesEn)
//	}
//	return
//}
//
//func getTextWithHyper(pdf *gofpdf.Fpdf, str string, wd float64) (height float64, lines [][]byte) {
//	lines = pdf.SplitLines([]byte(str), wd)
//	_, lineHt := pdf.GetFontSize()
//	lineHt = lineHt + 1
//	ht := float64(len(lines)) * lineHt
//	if ht == 0 {
//		ht = lineHt
//	}
//	//y := (297.0 - ht) / 2.0
//	return ht, lines
//}
