package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"image/color"
	"math"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

//go:embed *.json
var jsonFiles embed.FS

var ( // Global variables
	mainWindow      fyne.Window
	UnPassLabel     *canvas.Text
	MainLabel       *canvas.Text
	PriceLabel      *canvas.Text
	ListLabel       *canvas.Text
	DealLabel       *canvas.Text
	Red             color.Color
	Green           color.Color
	CandleColor     color.Color
	PrevCandleColor color.Color
	scroll          string
	pass            string
	PairText        string
	priceText       string
	pairs           []Pair
	candles         []Candle
	CandleY         float32
	WickY           float32
	WickSizeUp      float32
	WickSizeDown    float32
	CandleSize      float32
	CandleScale     float32
	pricePos        float32
)

type Pair struct { // Struct for JSON
	Pair string `json:"pair"`
}

type Candle struct { // Struct for candles
	Id    int     `json:"candle_id"`
	Open  float32 `json:"open"`
	Close float32 `json:"close"`
	High  float32 `json:"high"`
	Low   float32 `json:"low"`
}

func Init() { // Variables initialization
	Red = color.NRGBA{255, 0, 0, 255}
	Green = color.NRGBA{0, 255, 0, 255}
	scroll = ""

	MainLabel = canvas.NewText("Trade & Stage", color.White)
	MainLabel.TextSize = 25
	MainLabel.TextStyle = fyne.TextStyle{Bold: true}

	UnPassLabel = canvas.NewText("401 =(", color.White)
	UnPassLabel.TextSize = 25
	UnPassLabel.TextStyle = fyne.TextStyle{Bold: true}

	ListLabel = canvas.NewText("Список", color.White)
	ListLabel.TextSize = 25
	ListLabel.TextStyle = fyne.TextStyle{Bold: true}

	DealLabel = canvas.NewText(PairText, color.White)
	DealLabel.TextSize = 25
	DealLabel.TextStyle = fyne.TextStyle{Bold: true}

	for range 84 {
		scroll += "A"
	}
}

func Start() {
	Init()
	Load()
	Authorize()
}

func Authorize() { // Window for sign in
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
		container.NewCenter(MainLabel),
	)
	mainWindow.SetContent(codeBox)
}

func PairsList() { // Window for decide a pair
	buttonContainer := container.NewVBox()

	for idx := range pairs {
		// Важно: создаем локальную копию переменной
		pairIdx := idx
		button := widget.NewButton(pairs[pairIdx].Pair, func() {
			PairText = pairs[pairIdx].Pair
			Deal()
		})
		buttonContainer.Add(button)
	}

	listBox := container.NewVBox(
		container.NewCenter(ListLabel),
		buttonContainer,
	)

	mainWindow.SetContent(container.NewBorder(
		listBox,
		nil,
		nil,
		nil,
	))
}

func Pass(append string) { // Password handle
	if len(pass) <= 3 {
		pass += append
		if pass == "pas" && len(pass) == 3 {
			PairsList()
		} else if len(pass) == 3 {
			mainWindow.SetContent(container.NewBorder(
				container.NewCenter(UnPassLabel),
				nil,
				nil,
				nil,
				widget.NewButton("Try again...", Start),
			))
			pass = ""
		}
	}
}

func Deal() {
	DealLabel.Text = PairText
	RectContainer := container.NewWithoutLayout()

	if len(candles) == 0 {
		fmt.Println("Нет данных свечей!")
		return
	}

	CandleY = candles[0].Open
	WickY = candles[0].Low
	CandleScale = 10

	for i := 0; i < len(candles) && i < 30; i++ {
		WickSizeUp = candles[i].High - candles[i].Close
		WickSizeDown = candles[i].Open - candles[i].Low
		CandleSize = float32(math.Abs(float64(candles[i].Close - candles[i].Open)))

		if candles[i].Close >= candles[i].Open {
			CandleColor = Green
		} else {
			CandleColor = Red
		}

		if PrevCandleColor == CandleColor {
			if CandleColor == Green {
				CandleY -= (CandleSize + CandleScale)
			} else {
				CandleY += (CandleSize + CandleScale)
			}
		}

		rect := canvas.NewRectangle(CandleColor)
		rect.Resize(fyne.NewSize(11, CandleSize+CandleScale))
		rect.Move(fyne.NewPos(float32(i*15), CandleY))

		wickUp := canvas.NewRectangle(CandleColor)
		wickUp.Resize(fyne.NewSize(1, WickSizeUp+40))
		wickUp.Move(fyne.NewPos(float32(i*15)+5, CandleY-WickSizeUp-CandleScale))

		wickDown := canvas.NewRectangle(CandleColor)
		wickDown.Resize(fyne.NewSize(1, WickSizeDown+10))
		wickDown.Move(fyne.NewPos(float32(i*15)+5, CandleY+CandleScale+CandleSize))

		price := canvas.NewRectangle(CandleColor)
		price.Resize(fyne.NewSize(52, 4))
		priceText = strconv.FormatFloat(float64(candles[i].Close), 'f', 6, 32)
		priceLabel := canvas.NewText(strconv.FormatFloat(float64(candles[i].Close), 'f', 3, 32), CandleColor)

		if CandleColor == Green {
			pricePos = CandleY - WickSizeUp - CandleScale
			priceLabel.Move(fyne.NewPos(float32(i*15)+5, pricePos-20))
		} else {
			pricePos = CandleY + WickSizeDown + CandleScale + CandleSize + 10
			priceLabel.Move(fyne.NewPos(float32(i*15)+5, pricePos+20))
		}
		price.Move(fyne.NewPos(float32(i*15)+5, pricePos))

		RectContainer.Add(wickUp)
		RectContainer.Add(wickDown)
		RectContainer.Add(rect)

		if candles[i].Id == 1 || candles[i].Id == 2 {
			RectContainer.Add(price)
			RectContainer.Add(priceLabel)
		}

		PrevCandleColor = CandleColor
	}

	RectContainer.Add(container.NewCenter(canvas.NewText(scroll, color.NRGBA{0, 0, 0, 0})))
	RectContainer.Resize(fyne.NewSize(400, 350))
	RectContainer.Move(fyne.NewPos(0, 100))

	back := widget.NewButton("Назад", func() { PairsList() })
	exit := widget.NewButton("Выход", func() { os.Exit(0) })

	mainWindow.SetContent(container.NewBorder(
		container.NewGridWithColumns(3,
			back,
			container.NewCenter(DealLabel),
			exit,
		),
		nil,
		nil,
		nil,
		container.NewScroll(RectContainer),
	))
}

func Load() {
	data, err := jsonFiles.ReadFile("pairs.json")
	if err != nil {
		fmt.Println("Ошибка чтения pairs.json:", err)
		panic(err)
	}

	if err := json.Unmarshal(data, &pairs); err != nil {
		panic(err)
	}

	data2, err := jsonFiles.ReadFile("1H.json")
	if err != nil {
		fmt.Println("Ошибка чтения 1H.json:", err)
		panic(err)
	}

	if err := json.Unmarshal(data2, &candles); err != nil {
		panic(err)
	}
}

func main() {
	pass = ""
	app := app.NewWithID("com.yourcompany.yourapp")
	mainWindow = app.NewWindow("T&S")
	mainWindow.Resize(fyne.NewSize(400, 500))
	Start()
	mainWindow.ShowAndRun()
}
