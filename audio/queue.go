package audio

import (
	// "Waveform/audio/tools"
	"Waveform/audio/tools"
	"cmp"
	"errors"
	"io"
	"log"
	"math/rand"
	"slices"

	// "sort"
	"time"
)

type SongTimer struct {
	elapsed         int
	pause, shutdown chan bool
	output          chan int
}

type Queue struct {
	ranomized bool
	current   *SongData
	position  int
	timer     *SongTimer
	inner     []SongData
}

func (self *Queue) PlayNext() {
	self.position = (self.position + 1) % (len(self.inner))
	if self.current != nil {

		self.reset_current()
		self.current = &self.inner[self.position]
		_, err := self.current.Player.Seek(0, io.SeekStart)
		if err != nil {
			log.Println("Could Not seek start")
		}
		self.current.Player.Play()

	} else {
		log.Println("Initialize the queue first!")
	}
	log.Printf("Currently Playing: %s", self.Current().Title)

}
func (self *Queue) GetTrackList() []string {
	return tools.Map(self.inner, func(a SongData) string { return a.Title })
}

func (self *Queue) reset_current() {
	self.current.Player.Pause()
	self.current.Player.Seek(0, io.SeekStart)
}
func (self *Queue) PlayPrevious() {
	if self.position == 0 {
		_, err := self.current.Player.Seek(0, io.SeekStart)
		if err != nil {
			log.Println("Could Not seek start")
		}

	} else {
		self.reset_current()
		self.position -= 1
		self.current = &self.inner[self.position]
		_, err := self.current.Player.Seek(0, io.SeekStart)
		if err != nil {
			log.Println("Could Not seek start")
		}
	}
	self.current.Player.Play()
	log.Printf("Currently Playing: %s", self.Current().Title)
	log.Println(self.position)

}

func (self *Queue) Init() {

	self.current = &self.inner[0]
	self.current.Player.SetVolume(0.15)
	self.timer = new_timer()
}

func (self *Queue) Sort(field func(a *SongData) string) {
	slices.SortStableFunc(self.inner, func(a, b SongData) int {

		return cmp.Compare(field(&a), (field(&b)))
	})
}
func (self *Queue) Length() int {
	return len(self.inner)
}
func (self *Queue) Current() *SongData {
	return self.current
}

// This function shuffles the queue randomly and puts the current song on top. This is how all players do that.
func (self *Queue) Randomize() {
	self.inner[0], self.inner[self.position] = self.inner[self.position], self.inner[0]
	tmp := self.inner[1:]

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(tmp), func(i, j int) {
		tmp[i], tmp[j] = tmp[j], tmp[i]
	})
	self.inner = slices.Insert(tmp, 0, self.inner[0])
	self.current = &self.inner[0]
	self.ranomized = true

}

func (self *Queue) Pause() {
	self.current.Player.Pause()
}

func (self *Queue) is_playing() bool {
	return self.current.Player.IsPlaying()
}

func (self *Queue) Resume() {
	self.current.Player.Play()
}

func (self *Queue) GetCurrentCover() (*[]byte, error) {

	if self.current.Cover == nil {
		return nil, errors.New("Unexpected error while loading cover.")
	}
	return self.current.Cover, nil
}
func (self *Queue) PlayPause() {

	if !self.is_playing() {

		self.Current().Player.Play()
	} else {
		self.Current().Player.Pause()

	}
}

func (self *SongTimer) ResetTimer() {
	self.shutdown <- true
	self.elapsed = 0
}
func (self *SongTimer) SetTimer(until int) {
	go func() {
		defer func() {
			self.output <- -1
		}()
		elapsed := 0
		for elapsed <= until {
			select {
			case self.output <- elapsed:
				elapsed++
				time.Sleep(time.Second)
			case <-self.pause:
				<-self.pause
			case <-self.shutdown:
				return
			default:
				log.Fatalln("This code should be unreachable! Looks like the timer channel has been overflown.")

			}

		}

	}()

}
func (self *Queue) DecreaseVolume() {
	if vol := self.Current().Player.Volume(); vol >= 0.05 {

		self.Current().Player.SetVolume(vol - 0.05)
	} else {
		self.Current().Player.SetVolume(0)
	}
}
func (self *Queue) IncreaseVolume() {
	if vol := self.Current().Player.Volume(); vol <= 0.95 {

		self.Current().Player.SetVolume(vol + 0.05)
	} else {

		self.Current().Player.SetVolume(1)
	}
}
func (self *Queue) Volume() float64 {
	return self.current.Player.Volume()
}
func (self *Queue) IsRandomized() bool {
	return self.ranomized
}
func new_timer() *SongTimer {
	return &SongTimer{

		elapsed: 0,
		output:  make(chan int, 1),
		pause:   make(chan bool),
	}
}
