package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash/crc32"
	"io"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// 文件Hash工具
func createFileHashWindow(a fyne.App) {
	w := a.NewWindow("文件Hash工具")

	fileLabel := widget.NewLabel("")
	fileInfoLabel := widget.NewLabel("")
	fileMd5Text := widget.NewEntry()
	fileSha1Text := widget.NewEntry()
	fileSha256Text := widget.NewEntry()
	fileSha512Text := widget.NewEntry()
	fileCrcText := widget.NewEntry()
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

				fileHash, err := computeFileHash(filePath)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				fileMd5Text.SetText(fileHash.MD5)
				fileSha1Text.SetText(fileHash.SHA1)
				fileSha256Text.SetText(fileHash.SHA256)
				fileSha512Text.SetText(fileHash.SHA512)
				fileCrcText.SetText(fmt.Sprintf("%d", fileHash.CRC32))
			}
		}, w)

		fileDialog.Show()
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "文件", Widget: openButton},
			{Text: "", Widget: fileLabel},
			{Text: "文件信息", Widget: fileInfoLabel},
			{Text: "MD5", Widget: fileMd5Text},
			{Text: "SHA1", Widget: fileSha1Text},
			{Text: "SHA256", Widget: fileSha256Text},
			{Text: "SHA512", Widget: fileSha512Text},
			{Text: "CRC32", Widget: fileCrcText},
		},
	}

	w.SetContent(container.NewVBox(form))
	w.Resize(fyne.NewSize(600, 400))
	w.Show()
}

// FileHash 结构体存储所有哈希值
type FileHash struct {
	MD5    string
	SHA1   string
	SHA256 string
	SHA512 string
	CRC32  uint32
}

// computeFileHash 计算文件哈希值
func computeFileHash(filename string) (*FileHash, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 小文件使用常规方法
	return conventionalCompute(file)
}

func conventionalCompute(file *os.File) (*FileHash, error) {
	md5Hash := md5.New()
	sha1Hash := sha1.New()
	sha256Hash := sha256.New()
	sha512Hash := sha512.New()
	crc32Hash := crc32.NewIEEE()

	multiWriter := io.MultiWriter(md5Hash, sha1Hash, sha256Hash, sha512Hash, crc32Hash)

	if _, err := io.Copy(multiWriter, file); err != nil {
		return nil, err
	}

	return &FileHash{
		MD5:    fmt.Sprintf("%x", md5Hash.Sum(nil)),
		SHA1:   fmt.Sprintf("%x", sha1Hash.Sum(nil)),
		SHA256: fmt.Sprintf("%x", sha256Hash.Sum(nil)),
		SHA512: fmt.Sprintf("%x", sha512Hash.Sum(nil)),
		CRC32:  crc32Hash.Sum32(),
	}, nil
}
