package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("File Spliter")

	splitValueEntry := widget.NewEntry()
	splitValueEntry.SetPlaceHolder("整数")

	unitSelect := widget.NewSelect([]string{UNIT_KB, UNIT_MB}, func(value string) {})

	splitType := widget.NewSelect([]string{SPLIT_BY_LINE, SPLIT_BY_SIZE}, func(value string) {
		if value == SPLIT_BY_LINE {
			splitValueEntry.SetText(SPLIT_LINE_DEFAULT)
			unitSelect.Hide()
		} else {
			splitValueEntry.SetText(SPLIT_SIZE_DEFAULT)
			unitSelect.SetSelectedIndex(1)
			unitSelect.Show()
		}
	})
	splitType.SetSelectedIndex(0)

	firstLineCheck := widget.NewCheck("首行列名", func(checked bool) {})

	fileLabel := widget.NewLabel("")
	fileInfoLabel := widget.NewLabel("")
	var filePath string
	openButton := widget.NewButton("打开", func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader != nil {
				filePath = reader.URI().Path()
				fileLabel.SetText(filePath)

				fileLines, fileSize, _ := getFileInfo(filePath)
				fileInfoLabel.SetText(fmt.Sprintf("行数：%d,  大小: %s", fileLines, fileSize))
			}
		}, w)
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".csv", ".txt"}))
		fileDialog.Show()
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "文件", Widget: openButton},
			{Text: "", Widget: fileLabel},
			{Text: "文件信息", Widget: fileInfoLabel},
			{Text: "分割方式", Widget: splitType},
			{Text: "行数/大小", Widget: splitValueEntry},
			{Text: "单位", Widget: unitSelect},
			{Text: "", Widget: firstLineCheck},
		},
		SubmitText: "执行",
		OnSubmit: func() {
			if splitValueEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("请输入行数或文件大小"), w)
				return
			}
			num, err := strconv.Atoi(splitValueEntry.Text)
			if err != nil {
				dialog.ShowError(fmt.Errorf("不是整数: %s", splitValueEntry.Text), w)
				return
			}
			if filePath == "" {
				dialog.ShowError(fmt.Errorf("没有选择文件"), w)
				return
			}
			splitFile(filePath, splitType.Selected, num, unitSelect.Selected, firstLineCheck.Checked, w)
			dialog.ShowInformation("INFO", "分割完成", w)
		},
	}

	w.SetContent(container.NewVBox(form))
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}

// getFileInfo 获取文件信息
func getFileInfo(filename string) (fileLine int, fileSize string, err error) {
	// 获取文件信息
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return
	}
	fileSize = bytesToString(fileInfo.Size())

	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 4096) // 设置缓冲区大小
	for {
		_, err := reader.ReadBytes('\n') // 读取到下一个换行符
		if err != nil {
			if err == io.EOF {
				break // 文件结束
			}
			return fileLine, fileSize, err
		}
		fileLine++
	}
	return
}
