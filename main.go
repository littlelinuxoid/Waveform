package main

import (
	audio "Waveform/audio"
	"Waveform/audio/tools"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/ebitengine/oto/v3"

	// "github.com/ebitengine/oto/v3"

	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func WF_play_button(player *oto.Player) *widget.Button {

	bt := widget.NewButton("  ", nil)
	bt.OnTapped = func() {
		if !player.IsPlaying() {
			player.Play()
			bt.Text = "  "
			bt.Refresh()
		} else {
			player.Pause()
			bt.Text = "  "
			bt.Refresh()
		}
	}
	return bt
}
func WF_next_button(q *audio.Queue) *widget.Button {
	bt := widget.NewButton(" 󰒭 ", nil)
	bt.OnTapped = func() {

		q.PlayNext()
	}
	return bt
}
func WF_prev_button(q *audio.Queue) *widget.Button {

	return widget.NewButton(" 󰒮 ", func() {
		q.PlayPrevious()

	})
}
func main() {

	app := app.New()
	window := app.NewWindow("Waveform")
	app.Settings().SetTheme(theme.DefaultTheme())
	window.Resize(fyne.NewSize(1920/2, 1000))
	ctx := tools.NewContext()
	q := audio.FromDirectory(ctx, "./resources")
	newsong := q.Current()

	log.Println(newsong)
	timer := widget.NewSlider(0, 100)
	// image := canvas.NewImageFromResource(fyne.NewStaticResource("huh", *newsong.Cover))
	play_bt := WF_play_button(newsong.Player)
	prev_bt := WF_prev_button(q)
	next_bt := WF_next_button(q)

	inner := container.NewGridWithRows(3,
		// image,
		timer,
		container.New(layout.NewGridLayoutWithColumns(3),
			prev_bt,
			play_bt,
			next_bt,
		),
	)
	// image.FillMode = canvas.ImageFillContain
	// image.Resize(fyne.NewSize(100, 500))
	// image.Refresh()

	mainscreen := container.New(layout.NewCustomPaddedLayout(100, 100, 50, 50),
		inner,
	)

	window.SetContent(mainscreen)
	window.ShowAndRun()

}
