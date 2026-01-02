package main

import (
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	mainWindow  fyne.Window
	mainMenu    *fyne.Container
	mainButtons *fyne.Container
	new         int
	active      Deal
)

type Deal struct {
	Number   int
	Pair     string
	StopLoss string
	Amount   string
}

type SeaTheme struct{}

func (s SeaTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	colors := map[fyne.ThemeColorName]color.Color{
		theme.ColorNamePrimary:    color.RGBA{64, 224, 208, 255}, 
		theme.ColorNameBackground: color.RGBA{10, 40, 70, 255},    
		theme.ColorNameForeground: color.RGBA{230, 245, 255, 255}, 
		theme.ColorNameButton:     color.RGBA{0, 105, 148, 255},   
	}
	if c, ok := colors[n]; ok {
		return c
	}
	return theme.DarkTheme().Color(n, v)
}
func (s SeaTheme) Icon(n fyne.ThemeIconName) fyne.Resource { return theme.DarkTheme().Icon(n) }
func (s SeaTheme) Font(sf fyne.TextStyle) fyne.Resource    { return theme.DarkTheme().Font(sf) }
func (s SeaTheme) Size(n fyne.ThemeSizeName) float32       { return theme.DarkTheme().Size(n) }

func main() {
	app := app.New()
	mainWindow = app.NewWindow("helper v0")
	app.Settings().SetTheme(&SeaTheme{})
	mainWindow.Resize(fyne.NewSize(900, 600))
	UI()
	new = 0
	mainWindow.ShowAndRun()
}

func UI() {
	// 1. Все лейблы
	MainLabel := canvas.NewText("Хелпер v0", color.White)
	MainLabel.TextSize = 30
	MainLabel.TextStyle = fyne.TextStyle{Bold: true}

	DealLabel := canvas.NewText("Открываем сделку", color.White)
	DealLabel.TextSize = 30
	DealLabel.TextStyle = fyne.TextStyle{Bold: true}

	ActiveLabel := canvas.NewText("Активные сделки", color.White)
	ActiveLabel.TextSize = 30
	ActiveLabel.TextStyle = fyne.TextStyle{Bold: true}

	CheckLabel := canvas.NewText("Проверьте данные сделки!", color.White)
	CheckLabel.TextSize = 30
	CheckLabel.TextStyle = fyne.TextStyle{Bold: true}

	SuccessLabel := canvas.NewText("Сделка успешно открыта!", color.White)
	SuccessLabel.TextSize = 30
	SuccessLabel.TextStyle = fyne.TextStyle{Bold: true}

	Menu1 := canvas.NewText("Открыть сделку", color.White)
	Menu1.TextSize = 25
	Menu1.TextStyle = fyne.TextStyle{Bold: true}

	Menu2 := canvas.NewText("Активные сделки", color.White)
	Menu2.TextSize = 25
	Menu2.TextStyle = fyne.TextStyle{Bold: true}

	Menu3 := canvas.NewText("Выход", color.White)
	Menu3.TextSize = 25
	Menu3.TextStyle = fyne.TextStyle{Bold: true}

	// 2. Все вкладки и кнопки к ним

	pair := widget.NewEntry()
	stop := widget.NewEntry()
	amount := widget.NewEntry()

	successDeal := container.NewVBox(
		layout.NewSpacer(),
		container.NewCenter(SuccessLabel),
		widget.NewButton("Вернуться в главное меню", func() {
			mainWindow.SetContent(container.NewBorder(
				mainMenu,
				nil,
				nil,
				nil,
				mainButtons,
			))
		}),
		layout.NewSpacer(),
	)

	CheckButton := widget.NewButton("Подтверждаю, что всё верно", func() {
		active = Deal{100, pair.Text, stop.Text, amount.Text}
		new++
		mainWindow.SetContent(successDeal)
	})

	checkDealMenu := container.NewVBox(
		container.NewCenter(CheckLabel),
		container.NewGridWithColumns(1,
			widget.NewLabel("Торговая пара:"),
			pair,
			widget.NewLabel("Стоп-лосс:"),
			stop,
			widget.NewLabel("Сумма сделки:"),
			amount,
		),
		container.NewStack(CheckButton),
	)

	dealForm := container.NewVBox(
		container.NewGridWithColumns(1,
			widget.NewLabel("Торговая пара:"),
			pair,
			widget.NewLabel("Стоп-лосс:"),
			stop,
			widget.NewLabel("Сумма:"),
			amount),
		widget.NewButton("Продолжить", func() {
			mainWindow.SetContent(checkDealMenu)
		},
		),
	)

	activeDealsMenu := container.NewVBox(
		container.NewCenter(ActiveLabel),
		container.NewHBox(
			widget.NewButton("Назад в меню", func() {
				mainWindow.SetContent(container.NewBorder(
					mainMenu,
					nil,
					nil,
					nil,
					mainButtons,
				))
			}),
		),
		container.NewGridWithColumns(3,
			widget.NewLabel("Тут будут активные сделки"),
			widget.NewLabel("Но уже после"),
			widget.NewLabel("Нового года"),
		),
	)

	dealMenu := container.NewVBox(
		container.NewCenter(DealLabel),
		layout.NewSpacer(),
	)

	OpenDealButton := widget.NewButton("", func() {
		mainWindow.SetContent(container.NewBorder(
			dealMenu,
			nil,
			nil,
			nil,
			dealForm,
		))
	})
	ActiveDealsButton := widget.NewButton("", func() {
		mainWindow.SetContent(activeDealsMenu)
	})
	ExitButton := widget.NewButton("", func() {
		os.Exit(0)
	})

	mainMenu = container.NewVBox(
		container.NewCenter(MainLabel),
	)
	mainButtons = container.NewGridWithColumns(3,
		container.NewStack(OpenDealButton,
			container.NewCenter(Menu1)),
		container.NewStack(ActiveDealsButton,
			container.NewCenter(Menu2)),
		container.NewStack(ExitButton,
			container.NewCenter(Menu3)),
	)
	mainWindow.SetContent(container.NewBorder(
		mainMenu,
		nil,
		nil,
		nil,
		mainButtons,
	))
}
