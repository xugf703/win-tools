package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func createFileSpliterWindow(a fyne.App) {
	w := a.NewWindow("文件分割工具")

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
	openButton := widget.NewButtonWithIcon("点击打开文件", theme.FolderOpenIcon(), func() {
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
	w.Show()
}

// splitFile 分割文件到多个小文件
func splitFile(filePath string, splitType string, splitValue int, unit string, firstLine bool, w fyne.Window) error {
	file, err := os.Open(filePath)
	if err != nil {
		dialog.ShowInformation("ERROR", "打开文件失败", w)
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	extName := filepath.Ext(file.Name())

	// 读取首行
	var header string
	if scanner.Scan() {
		header = scanner.Text()
	}

	// 按行数分割
	if splitType == SPLIT_BY_LINE {
		return splitByLines(file, extName, splitValue, header, firstLine)
	} else {
		return splitBySize(file, extName, splitValue, unit, header, firstLine)
	}

}

// splitByLines 按行数分割文件
func splitByLines(file *os.File, extName string, splitValue int, header string, firstLine bool) error {
	scanner := bufio.NewScanner(file)
	lineCount := 0
	fileIndex := 0
	fileBaseName := strings.TrimSuffix(file.Name(), extName)

	// 初始化第一个文件
	nextFileName := fmt.Sprintf("%s_%d.%s", fileBaseName, fileIndex, extName)
	nextFile, err := os.Create(nextFileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer nextFile.Close()

	writer := bufio.NewWriter(nextFile)

	// 写入首行（如果需要）
	if firstLine {
		if _, err := writer.WriteString(header + "\n"); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
	}

	for scanner.Scan() {
		line := scanner.Text()

		// 写入行
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write line: %w", err)
		}

		lineCount++
		if lineCount == splitValue {
			// 刷新缓冲区
			if err := writer.Flush(); err != nil {
				return fmt.Errorf("failed to flush buffer: %w", err)
			}

			// 关闭当前文件
			if err := nextFile.Close(); err != nil {
				return fmt.Errorf("failed to close file: %w", err)
			}

			// 准备下一个文件
			fileIndex++
			nextFileName = fmt.Sprintf("%s_%d.%s", fileBaseName, fileIndex, extName)
			nextFile, err = os.Create(nextFileName)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			defer nextFile.Close()

			writer = bufio.NewWriter(nextFile)

			// 写入首行（如果需要）
			if firstLine {
				if _, err := writer.WriteString(header + "\n"); err != nil {
					return fmt.Errorf("failed to write header: %w", err)
				}
			}
			lineCount = 0
		}
	}

	// 处理可能的剩余行
	if lineCount > 0 {
		if err := writer.Flush(); err != nil {
			return fmt.Errorf("failed to flush buffer: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	return nil
}

// splitBySize 按文件大小分割文件
func splitBySize(file *os.File, extName string, splitValue int, unit string, header string, firstLine bool) error {
	reader := bufio.NewReader(file)
	fileIndex := 0
	currentSize := 0
	fileBaseName := strings.TrimSuffix(file.Name(), extName)

	// 初始化第一个文件
	nextFileName := fmt.Sprintf("%s_%d.%s", fileBaseName, fileIndex, extName)
	nextFile, err := os.Create(nextFileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer nextFile.Close()

	writer := bufio.NewWriter(nextFile)

	// 写入首行（如果需要）
	if firstLine {
		if _, err := writer.WriteString(header + "\n"); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading file: %w", err)
		}

		lineBytes := []byte(line)

		// 写入行
		if _, err := writer.WriteString(line); err != nil {
			return fmt.Errorf("failed to write line: %w", err)
		}

		currentSize += len(lineBytes)
		if currentSize >= convertToBytes(splitValue, unit) {
			// 刷新缓冲区
			if err := writer.Flush(); err != nil {
				return fmt.Errorf("failed to flush buffer: %w", err)
			}

			// 关闭当前文件
			if err := nextFile.Close(); err != nil {
				return fmt.Errorf("failed to close file: %w", err)
			}

			// 准备下一个文件
			fileIndex++
			nextFileName = fmt.Sprintf("%s_%d.%s", fileBaseName, fileIndex, extName)
			nextFile, err = os.Create(nextFileName)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			defer nextFile.Close()

			writer = bufio.NewWriter(nextFile)

			// 写入首行（如果需要）
			if firstLine {
				if _, err := writer.WriteString(header + "\n"); err != nil {
					return fmt.Errorf("failed to write header: %w", err)
				}
			}

			currentSize = 0
		}
	}

	// 处理可能的剩余行
	if currentSize > 0 {
		if err := writer.Flush(); err != nil {
			return fmt.Errorf("failed to flush buffer: %w", err)
		}
	}
	return nil
}
