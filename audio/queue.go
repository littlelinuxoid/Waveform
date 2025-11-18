package audio

import (
	// "Waveform/audio/tools"
	"cmp"
	"errors"
	"io"
	"log"
	"math/rand"
	"slices"
	// "sort"
	"time"
)

type Queue struct {
	inner    []SongData
	current  *SongData
	position int
	len      int
}

func (self *Queue) PlayNext() {
	self.position = (self.position + 1) % (self.len)
	if self.current != nil {

		self.ResetCurrent()
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
func (self *Queue) ResetCurrent() {
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
		self.ResetCurrent()
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
}

func (self *Queue) Sort(field func(a *SongData) string) {
	slices.SortStableFunc(self.inner, func(a, b SongData) int {

		return cmp.Compare(field(&a), (field(&b)))
	})
}
func (self *Queue) Length() int {
	return self.len
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

}

func (self *Queue) Pause() {
	self.current.Player.Pause()
}

func (self *Queue) IsPlaying() bool {
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
