package pdfGenerator

import (
	"bytes"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"os"
	"strings"
)

type (
	PdfDoc struct {
		Pdf *gofpdf.Fpdf
		Tr  func(string) string
	}

	Image struct {
		Src  string
		X    float64
		Y    float64
		W    float64
		H    float64
		Type string
	}

	Text struct {
		Data       string
		Align      string
		Width      float64
		LineHeight float64
		FontSize   float64
		AddY       float64
		X          float64
		IsBold     bool
		Height     float64 // высота текста
		IsBorder   bool
		PaddingTop float64
	}
)

func Init() (*PdfDoc, error) {
	pd := PdfDoc{}
	// находим путь для папки с шрифтами
	//gopath := os.Getenv("GOPATH")
	//if gopath == "" {
	//	gopath = build.Default.GOPATH
	//}
	//fontPath := fmt.Sprintf("%s/src/github.com/pepelazz/pdfGenerator/font/", gopath)
	//
	//pd.Pdf = gofpdf.New("P", "mm", "A4", fontPath)
	////fmt.Printf("w: %v h: %v\n", pW, pH)
	//pd.Pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	//pd.Pdf.AddFont("Roboto-Regular", "", "Roboto-Regular.json")
	//pd.Pdf.AddFont("Roboto-Regular", "B", "Roboto-Bold.json")
	//pd.Pdf.AddFont("Roboto-Regular", "I", "Roboto-Italic.json")
	//pd.Pdf.AddFont("Roboto-Regular", "BI", "Roboto-BoldItalic.json")
	//pd.Pdf.SetFont("Roboto-Regular", "", 8)
	//pd.Tr = pd.Pdf.UnicodeTranslatorFromDescriptor("cp1251")
	return &pd, nil
}

func (pd *PdfDoc) Print(fileName string) error {

	buf := new(bytes.Buffer)
	err := pd.Pdf.Output(buf)
	if err != nil {
		return err
	}

	f, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	defer f.Close()

	_, err = f.Write(buf.Bytes())
	return err
}

func (pd *PdfDoc) PrintToByte() ([]byte, error) {

	buf := new(bytes.Buffer)
	err := pd.Pdf.Output(buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

func (pd *PdfDoc) AddPage() {
	pd.Pdf.AddPage()
}

func (pd *PdfDoc) GetY() float64 {
	return pd.Pdf.GetY()
}

func (pd *PdfDoc) AddImage(img Image) {
	if len(img.Src) == 0 {
		fmt.Printf("AddImage missed Src")
		return
	}
	if len(img.Type) == 0 {
		img.Type = "jpg"
	}
	pd.Pdf.ImageOptions(img.Src, img.X, pd.GetY()+img.Y, img.W, img.H, false, gofpdf.ImageOptions{ImageType: img.Type}, 0, "")

}

// печать текста
func (pd *PdfDoc) PrintText(t Text) {
	if t.Width == 0 {
		t.Width, _ = pd.Pdf.GetPageSize()
	}
	if t.LineHeight == 0 {
		t.LineHeight = LINE_HEIGHT
	}
	if len(t.Align) == 0 {
		t.Align = "CM"
	}
	if t.FontSize > 0 {
		pd.Pdf.SetFontSize(t.FontSize)
	}
	if t.IsBold {
		pd.Pdf.SetFontStyle("B")
	}
	pd.AddY(t.AddY)
	// рисуем контур
	if t.IsBorder {
		pd.Pdf.Rect(t.X, pd.Pdf.GetY(), t.Width-BORDER_WIDTH, t.Height, "FD")
	}
	pd.AddY(t.PaddingTop)
	for _, line := range strings.Split(t.Data, "\n") {
		pd.Pdf.SetXY(t.X, pd.GetY())
		pd.Pdf.CellFormat(t.Width, t.LineHeight, pd.Tr(line), "", 1, t.Align, false, 0, "")
	}
	// в конце сбрасываем стиль
	pd.Pdf.SetFontStyle("")
}

func (pd *PdfDoc) AddY(v float64) {
	pd.Pdf.SetXY(pd.Pdf.GetX(), pd.Pdf.GetY()+v)
}

func (pd *PdfDoc) SetY(v float64) {
	pd.Pdf.SetXY(pd.Pdf.GetX(), v)
}
func (pd *PdfDoc) SetX(v float64) {
	pd.Pdf.SetXY(v, pd.Pdf.GetY())
}

func (pd *PdfDoc) GetFullWidth() float64 {
	pW, _ := pd.Pdf.GetPageSize()
	leftMargin, _, _, _ := pd.Pdf.GetMargins()
	return pW - leftMargin*2
}

func (t *Text) CalcHeight(pd *PdfDoc) {
	fontSize, _ := pd.Pdf.GetFontSize()
	if t.FontSize > 0 {
		fontSize = t.FontSize
	}
	h := len(strings.Split(t.Data, "\n"))
	t.Height = fontSize * float64(h) * 0.5

}
