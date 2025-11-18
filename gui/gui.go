package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	audio "Waveform/audio"
	// "Waveform/audio/tools"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
)

func getimage(q *audio.Queue) *canvas.Image {

	cover, err := q.GetCurrentCover()
	if err != nil {
		log.Println(err.Error())
	}
	image := canvas.NewImageFromResource(fyne.NewStaticResource("huh", *cover))
	return image
}
func Run() {

	app := app.New()
	window := app.NewWindow("Waveform")
	app.Settings().SetTheme(theme.DefaultTheme())
	window.Resize(fyne.NewSize(1920/2, 1000))

	timer := widget.NewSlider(0, 100)

	// image := getimage(q)
	btplay := widget.NewButton("  ", nil)
	btplay.OnTapped = func() {
	}
	btnext := widget.NewButton(" 󰒭 ", nil)
	btnext.OnTapped = func() {

	}
	btprev := widget.NewButton(" 󰒮 ", func() {

	})

	inner := container.NewGridWithRows(3,

		timer,
		container.New(layout.NewGridLayoutWithColumns(3),
			btprev,
			btplay,
			btnext,
		),
	)

	mainscreen := container.New(layout.NewCustomPaddedLayout(100, 100, 50, 50),
		inner,
	)

	window.SetContent(mainscreen)
	window.ShowAndRun()

}
