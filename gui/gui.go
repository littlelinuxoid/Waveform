package gui

import (
	audio "Waveform/audio"
	"Waveform/audio/tools"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
)

func get_mappings(q *audio.Queue) func(r rune) {

	return func(r rune) {

		switch r {
		case ' ':
			q.PlayPause()
		case '+':
			q.IncreaseVolume()
			log.Printf("Volume Increased, new volume: %.1f", q.Volume()*100)
		case '-':
			q.DecreaseVolume()
			log.Printf("Volume Decreased, new volume: %.1f", q.Volume()*100)
		case 'n':
			q.PlayNext()
			log.Println("Next Song")
		case 'p':
			q.PlayPrevious()
			log.Println("Previous Song")

		case 'r':
			if !q.IsRandomized() {

				q.Randomize()
				log.Println("Randomized!")
			} else {
				log.Println("Sorted Back!")

				q.Sort(func(a *audio.SongData) string { return a.Title })

			}
			log.Println("New Order:", q.GetTrackList())

		}

	}

}
func Run() {

	app := app.New()

	window := app.NewWindow("Waveform")
	app.Settings().SetTheme(theme.DefaultTheme())
	window.Resize(fyne.NewSize(1920/2, 1000))
	q := audio.FromDirectory(tools.NewContext(), "./resources")
	q.Init()

	window.Canvas().SetOnTypedRune(get_mappings(q))

	btplay := widget.NewButton("  ", nil)
	btplay.OnTapped = func() {
	}
	btnext := widget.NewButton(" 󰒭 ", nil)
	btnext.OnTapped = func() {
	}
	btprev := widget.NewButton(" 󰒮 ", func() {

	})

	inner := container.NewGridWithRows(3,

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
