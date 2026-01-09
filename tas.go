package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	mainWindow  fyne.Window
	UnPassLabel *canvas.Text
	pass        string
)

func main() {
	pass = ""
	app := app.New()
	mainWindow = app.NewWindow("T&S")
	mainWindow.Resize(fyne.NewSize(500, 400))
	UI()
	mainWindow.ShowAndRun()
}

func UI() {
	MainLabel := canvas.NewText("Welcome to Trade & Stage", color.White)
	MainLabel.TextSize = 25
	MainLabel.TextStyle = fyne.TextStyle{Bold: true}
	UnPassLabel = canvas.NewText("OOps =(", color.White)
	UnPassLabel.TextSize = 25
	UnPassLabel.TextStyle = fyne.TextStyle{Bold: true}
	Authorize(*MainLabel)
}

func Authorize(lab canvas.Text) {
	code1st := container.NewGridWithRows(3,
		widget.NewButton("?", func() { Pass("p") }),
		widget.NewButton("?", func() { Pass("0") }),
		widget.NewButton("?", func() { Pass("0") }),
	)
	code2nd := container.NewGridWithRows(3,
		widget.NewButton("?", func() { Pass("0") }),
		widget.NewButton("?", func() { Pass("a") }),
		widget.NewButton("?", func() { Pass("0") }),
	)
	code3rd := container.NewGridWithRows(3,
		widget.NewButton("?", func() { Pass("0") }),
		widget.NewButton("?", func() { Pass("0") }),
		widget.NewButton("?", func() { Pass("s") }),
	)
	codeBox := container.NewGridWithRows(3,
		layout.NewSpacer(),
		container.NewGridWithColumns(5,
			layout.NewSpacer(),
			code1st,
			code2nd,
			code3rd,
			layout.NewSpacer(),
		),
		container.NewCenter(&lab),
	)

	mainWindow.SetContent(codeBox)
}

func Menu() {
	list := container.NewGridWithRows(10,
		widget.NewButton("IMOEXF", func() {}),
		widget.NewButton("GLDRUBF", func() {}),
		widget.NewButton("BR-1.26", func() {}),
		widget.NewButton("NASD-3.26", func() {}),
		widget.NewButton("SBRF-3.26", func() {}),
		widget.NewButton("GAZP-3.26", func() {}),
		widget.NewButton("VTBR-3.26", func() {}),
		widget.NewButton("YDEX-3.26", func() {}),
		widget.NewButton("BTC", func() {}),
		widget.NewButton("ETH", func() {}),
	)
	mainWindow.SetContent(container.NewBorder(
		container.NewCenter(canvas.NewText("List...", color.White)),
		nil,
		nil,
		nil,
		list,
	))
}

func Pass(append string) {
	if len(pass) <= 3 {
		pass += append
		if pass == "pas" && len(pass) == 3 {
			Menu()
		} else if len(pass) == 3 {
			mainWindow.SetContent(container.NewBorder(
				container.NewCenter(UnPassLabel),
				nil,
				nil,
				nil,
				widget.NewButton("Try again...", UI),
			))
			pass = ""
		}
	}
}
