package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("工具库")
	w.SetContent(container.NewVBox(
		widget.NewButtonWithIcon("文件分割", theme.SettingsIcon(), func() {
			createFileSpliterWindow(a)
		}),
		widget.NewButtonWithIcon("文件Hash", theme.InfoIcon(), func() {
			createFileHashWindow(a)
		}),
	))
	w.Resize(fyne.NewSize(400, 300))
	w.ShowAndRun()

}
