package audio

import (
	"Waveform/audio/tools"
	"cmp"
	"errors"
	"io"
	"log"
	"math/rand"
	"slices"
	"sync"
	"time"
)

const DEBUG = true

type Queue struct {
	ranomized      bool
	current        *SongData
	inner          []SongData
	Pause, Changed *sync.Cond
	position       int
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
	if DEBUG {
		log.Println(self.position)
	}
}

func (self *Queue) Init() {
	go func() {
		self.current = &self.inner[0]
		self.current.Player.SetVolume(0.15)
		self.TrackTime()
	}()
}

func (self *Queue) Sort(field func(a *SongData) string) {
	slices.SortStableFunc(self.inner, func(a, b SongData) int {
		return cmp.Compare(field(&a), (field(&b)))
	})
	self.ranomized = false
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

func (self *Queue) is_playing() bool {
	return self.current.Player.IsPlaying()
}

func (self *Queue) GetCurrentCover() (*[]byte, error) {
	if self.current.Cover == nil {
		return nil, errors.New("Unexpected error while loading cover.")
	}
	return self.current.Cover, nil
}

func (self *Queue) PlayPause() {
	go func() {

		if !self.is_playing() {
			log.Println("Resume")
			self.current.Player.Play()
		} else {
			log.Println("Pause")
			self.current.Player.Pause()
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

func (self *Queue) TrackTime() {
	go func() {

		self.Pause.L.Lock()
		a := time.NewTicker(time.Second)
		defer func() {
			self.Changed.L.Unlock()
			a.Stop()
			self.current.Elapsed = 0
		}()

		for range a.C {
			if self.current.Elapsed <= self.current.Length_sec {
				break
			}
			
			
		}
	}()
}
