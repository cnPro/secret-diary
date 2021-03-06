package main

import (
	"bufio"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func (s *myWindow) setHeader(level int) {
	var cfmt = gui.NewQTextCharFormat()
	switch level {
	case 1:
		cfmt.SetFontPointSize(18)
	case 2:
		cfmt.SetFontPointSize(16)
	case 3:
		cfmt.SetFontPointSize(14)
	default:
		cfmt.SetFontPointSize(12)
	}

	cfmt.SetForeground(gui.NewQBrush3(gui.NewQColor2(core.Qt__blue), core.Qt__SolidPattern))
	s.mergeFormatOnLineOrSelection(cfmt)
}

func (s *myWindow) setStandard() {
	var cfmt = gui.NewQTextCharFormat()
	cfmt.SetFontPointSize(14)
	cfmt.SetForeground(gui.NewQBrush3(gui.NewQColor2(core.Qt__black), core.Qt__SolidPattern))

	s.mergeFormatOnLineOrSelection(cfmt)
}

func (s *myWindow) textStyle(styleIndex int) {
	var cursor = s.editor.TextCursor()

	if styleIndex != 0 {

		var style = gui.QTextListFormat__ListDisc

		switch styleIndex {
		case 1:
			{
				style = gui.QTextListFormat__ListDisc
			}

		case 2:
			{
				style = gui.QTextListFormat__ListCircle
			}

		case 3:
			{
				style = gui.QTextListFormat__ListSquare
			}

		case 4:
			{
				style = gui.QTextListFormat__ListDecimal
			}

		case 5:
			{
				style = gui.QTextListFormat__ListLowerAlpha
			}

		case 6:
			{
				style = gui.QTextListFormat__ListUpperAlpha
			}

		case 7:
			{
				style = gui.QTextListFormat__ListLowerRoman
			}

		case 8:
			{
				style = gui.QTextListFormat__ListUpperRoman
			}
		}

		cursor.BeginEditBlock()

		var (
			blockFmt = cursor.BlockFormat()
			listFmt  = gui.NewQTextListFormat()
		)

		if cursor.CurrentList().Pointer() != nil {
			listFmt = gui.NewQTextListFormatFromPointer(cursor.CurrentList().Format().Pointer())
		} else {
			listFmt.SetIndent(blockFmt.Indent() + 1)
			blockFmt.SetIndent(0)
			cursor.SetBlockFormat(blockFmt)
		}

		listFmt.SetStyle(style)
		cursor.CreateList(listFmt)

		cursor.EndEditBlock()

	} else {
		var bfmt = gui.NewQTextBlockFormat()
		bfmt.SetObjectIndex(-1)
		cursor.MergeBlockFormat(bfmt)
	}
}

func (s *myWindow) textColor() {
	var col = widgets.QColorDialog_GetColor(s.editor.TextColor(), s.editor, "", 0)
	if !col.IsValid() {
		return
	}
	var cfmt = gui.NewQTextCharFormat()
	cfmt.SetForeground(gui.NewQBrush3(col, core.Qt__SolidPattern))
	s.mergeFormatOnLineOrSelection(cfmt)
}

func (s *myWindow) textBgColor() {
	var col = widgets.QColorDialog_GetColor(s.editor.TextColor(), s.editor, "", 0)
	if !col.IsValid() {
		return
	}
	var cfmt = gui.NewQTextCharFormat()
	cfmt.SetBackground(gui.NewQBrush3(col, core.Qt__SolidPattern))
	s.mergeFormatOnLineOrSelection(cfmt)
}

func (s *myWindow) textBold() {
	var afmt = gui.NewQTextCharFormat()
	var fw = gui.QFont__Normal
	if s.actionTextBold.IsChecked() {
		fw = gui.QFont__Bold
	}
	afmt.SetFontWeight(int(fw))
	s.mergeFormatOnLineOrSelection(afmt)
}

func (s *myWindow) textUnderline() {
	var afmt = gui.NewQTextCharFormat()
	afmt.SetFontUnderline(s.actionTextUnderline.IsChecked())
	s.mergeFormatOnLineOrSelection(afmt)
}

func (s *myWindow) textStrikeOut() {
	var afmt = gui.NewQTextCharFormat()
	afmt.SetFontStrikeOut(s.actionStrikeOut.IsChecked())
	s.mergeFormatOnLineOrSelection(afmt)
}

func (s *myWindow) textItalic() {
	var afmt = gui.NewQTextCharFormat()
	afmt.SetFontItalic(s.actionTextItalic.IsChecked())
	s.mergeFormatOnLineOrSelection(afmt)
}

func (s *myWindow) insertImage() {
	filename := widgets.QFileDialog_GetOpenFileName(s.window, "select a file", ".", "Image (*.png *.jpg)", "Image (*.png *.jpg)", widgets.QFileDialog__ReadOnly)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	img := gui.NewQImage()
	ok := img.LoadFromData(data, len(data), "")
	if !ok {
		return
	}

	uri := core.NewQUrl3("rc://"+filename, core.QUrl__TolerantMode)

	img = s.scaleImage(img)

	s.editor.Document().AddResource(int(gui.QTextDocument__ImageResource), uri, img.ToVariant())
	url := uri.Url(core.QUrl__None)
	cursor := s.editor.TextCursor()
	cursor.InsertImage4(img, url)

	ba := core.NewQByteArray()

	iod := core.NewQBuffer2(ba, nil)

	iod.Open(core.QIODevice__WriteOnly)

	ok = img.Save2(iod, filepath.Ext(filename)[1:], -1)
	//fmt.Println(filepath.Ext(filename))
	if ok {
		s.document.Images[url] = []byte(ba.Data())
	}

	//fmt.Println("save image:", ok)
}

func (s *myWindow) scaleImage(src *gui.QImage) (res *gui.QImage) {
	dlg := widgets.NewQDialog(s.window, core.Qt__Dialog)
	dlg.SetWindowTitle(T("Scale Image Size"))

	grid := widgets.NewQGridLayout(dlg)

	width := widgets.NewQLabel2(strconv.Itoa(src.Width())+" =>", dlg, core.Qt__Widget)
	grid.AddWidget(width, 0, 0, 0)
	scaledW := src.Width()
	scaledH := src.Height()
	delta := 30
	if scaledW > s.editor.Width()-delta {
		scaledW = s.editor.Geometry().Width() - delta
		scaledH = int(float64(src.Height()) * float64(scaledW) / float64(src.Width()))
	}
	wValidor := gui.NewQIntValidator(dlg)
	wValidor.SetRange(10, scaledW)
	hValidor := gui.NewQIntValidator(dlg)
	hValidor.SetRange(10, scaledH)

	widthInput := widgets.NewQLineEdit(dlg)
	widthInput.SetText(strconv.Itoa(scaledW))
	widthInput.SetValidator(wValidor)
	grid.AddWidget(widthInput, 0, 1, 0)

	height := widgets.NewQLabel2(strconv.Itoa(src.Height())+" =>", dlg, core.Qt__Widget)

	grid.AddWidget(height, 1, 0, 0)

	heightInput := widgets.NewQLineEdit(dlg)
	heightInput.SetText(strconv.Itoa(scaledH))
	heightInput.SetValidator(hValidor)
	grid.AddWidget(heightInput, 1, 1, 0)

	btb := widgets.NewQGridLayout(nil)

	okBtn := widgets.NewQPushButton2(T("OK"), dlg)
	btb.AddWidget(okBtn, 0, 0, 0)

	cancelBtn := widgets.NewQPushButton2(T("Cancel"), dlg)
	btb.AddWidget(cancelBtn, 0, 1, 0)

	grid.AddLayout2(btb, 2, 0, 1, 2, 0)

	dlg.SetLayout(grid)

	widthInput.ConnectKeyReleaseEvent(func(e *gui.QKeyEvent) {
		w, err := strconv.Atoi(widthInput.Text())
		if err != nil {
			return
		}
		w0 := float64(src.Width())
		h0 := float64(src.Height())
		h := float64(w) * h0 / w0
		heightInput.SetText(strconv.Itoa(int(h)))
	})
	heightInput.ConnectKeyReleaseEvent(func(e *gui.QKeyEvent) {
		h, err := strconv.Atoi(heightInput.Text())
		if err != nil {
			return
		}
		w0 := float64(src.Width())
		h0 := float64(src.Height())
		w := float64(h) * w0 / h0
		widthInput.SetText(strconv.Itoa(int(w)))
	})

	okBtn.ConnectClicked(func(b bool) {
		w, err := strconv.Atoi(widthInput.Text())
		if err != nil {
			res = src
		}
		h, err := strconv.Atoi(heightInput.Text())
		if err != nil {
			res = src
		}
		res = src.Scaled2(w, h, core.Qt__KeepAspectRatioByExpanding, core.Qt__SmoothTransformation)
		dlg.Hide()
		dlg.Destroy(true, true)
	})

	cancelBtn.ConnectClicked(func(b bool) {
		res = src
		dlg.Hide()
		dlg.Destroy(true, true)
	})

	dlg.Exec()
	return
}

func (s *myWindow) getImageList(html string) []string {
	r := strings.NewReader(html)
	bufr := bufio.NewReader(r)
	reg1, err := regexp.Compile(`<img[^><]+/>`)
	if err != nil {
		//fmt.Println(err)
		return nil
	}
	reg2, err := regexp.Compile(`src="([^"]+)"`)
	if err != nil {
		//fmt.Println(err)
		return nil
	}
	imgs := []string{}
	for line, _, err := bufr.ReadLine(); err == nil; line, _, err = bufr.ReadLine() {
		line1 := string(line)
		res1 := reg1.FindAllString(line1, -1)
		imgs = append(imgs, res1...)
	}

	res := []string{}
	for _, img := range imgs {
		res2 := reg2.FindStringSubmatch(img)
		res = append(res, res2[1])
	}
	//fmt.Println(res)
	return res
}

func (s *myWindow) insertTable() {
	dlg := widgets.NewQDialog(s.window, core.Qt__Dialog)
	dlg.SetWindowTitle(T("Table Rows and Columns"))

	grid := widgets.NewQGridLayout(dlg)

	row := widgets.NewQLabel2(T("Rows:"), dlg, core.Qt__Widget)
	grid.AddWidget(row, 0, 0, 0)

	rowInput := widgets.NewQLineEdit(dlg)
	rowInput.SetText("3")
	rowInput.SetValidator(gui.NewQIntValidator(dlg))
	grid.AddWidget(rowInput, 0, 1, 0)

	col := widgets.NewQLabel2(T("Columns:"), dlg, core.Qt__Widget)

	grid.AddWidget(col, 1, 0, 0)

	colInput := widgets.NewQLineEdit(dlg)
	colInput.SetText("3")
	colInput.SetValidator(gui.NewQIntValidator(dlg))
	grid.AddWidget(colInput, 1, 1, 0)

	btb := widgets.NewQGridLayout(nil)

	okBtn := widgets.NewQPushButton2(T("OK"), dlg)
	btb.AddWidget(okBtn, 0, 0, 0)

	cancelBtn := widgets.NewQPushButton2(T("Cancel"), dlg)
	btb.AddWidget(cancelBtn, 0, 1, 0)

	grid.AddLayout2(btb, 2, 0, 1, 2, 0)

	dlg.SetLayout(grid)

	okBtn.ConnectClicked(func(b bool) {
		cursor := s.editor.TextCursor()
		r, err := strconv.Atoi(rowInput.Text())
		if err != nil {
			return
		}
		c, err := strconv.Atoi(colInput.Text())
		if err != nil {
			return
		}
		tbl := cursor.InsertTable2(r, c)
		tbl.Format().SetBorderBrush(gui.NewQBrush2(core.Qt__SolidPattern))
		dlg.Hide()
		dlg.Destroy(true, true)
	})

	cancelBtn.ConnectClicked(func(b bool) {
		dlg.Hide()
		dlg.Destroy(true, true)
	})

	dlg.SetModal(true)
	dlg.Show()
}

func (s *myWindow) addJustifyActions(tb *widgets.QToolBar) {
	rsrcPath := ":/qml/icons"
	var leftIcon = gui.QIcon_FromTheme2("format-justify-left", gui.NewQIcon5(rsrcPath+"/textleft.png"))
	actionAlignLeft := tb.AddAction2(leftIcon, "&Left")
	actionAlignLeft.SetPriority(widgets.QAction__LowPriority)
	actionAlignLeft.ConnectTriggered(func(b bool) {
		s.textAlign(1)
	})

	var centerIcon = gui.QIcon_FromTheme2("format-justify-center", gui.NewQIcon5(rsrcPath+"/textcenter.png"))
	actionAlignCenter := tb.AddAction2(centerIcon, "C&enter")
	actionAlignCenter.SetPriority(widgets.QAction__LowPriority)
	actionAlignCenter.ConnectTriggered(func(b bool) {
		s.textAlign(2)
	})

	var rightIcon = gui.QIcon_FromTheme2("format-justify-right", gui.NewQIcon5(rsrcPath+"/textright.png"))
	actionAlignRight := tb.AddAction2(rightIcon, "&Right")
	actionAlignRight.SetPriority(widgets.QAction__LowPriority)
	actionAlignRight.ConnectTriggered(func(b bool) {
		s.textAlign(3)
	})

	var fillIcon = gui.QIcon_FromTheme2("format-justify-fill", gui.NewQIcon5(rsrcPath+"/textjustify.png"))
	actionAlignJustify := tb.AddAction2(fillIcon, "&Justify")
	actionAlignJustify.SetPriority(widgets.QAction__LowPriority)
	actionAlignJustify.ConnectTriggered(func(b bool) {
		s.textAlign(4)
	})
}

func (s *myWindow) textAlign(n int) {
	switch n {
	case 1:

		s.editor.SetAlignment(core.Qt__AlignLeft | core.Qt__AlignAbsolute)

	case 2:

		s.editor.SetAlignment(core.Qt__AlignHCenter)

	case 3:

		s.editor.SetAlignment(core.Qt__AlignRight | core.Qt__AlignAbsolute)

	case 4:

		s.editor.SetAlignment(core.Qt__AlignJustify)

	}
}

func (s *myWindow) getTable() (t *gui.QTextTable, cell *gui.QTextTableCell) {
	cursor := s.editor.TextCursor()
	blk := s.editor.Document().FindBlock(cursor.Position())
	for _, frame := range blk.Document().RootFrame().ChildFrames() {
		//fmt.Println("table cell:", frame.FrameFormat().IsTableCellFormat(), "table:", frame.FrameFormat().IsTableFormat())
		if frame.FrameFormat().IsTableFormat() {
			table := gui.NewQTextTableFromPointer(frame.Pointer())
			cell := table.CellAt2(cursor.Position())
			//fmt.Println(cell.Row(), cell.Column())
			return table, cell
		}
	}

	return nil, nil
}
