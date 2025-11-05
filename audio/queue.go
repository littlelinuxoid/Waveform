package audio

import (
	"io"
	"log"
)

type Queue struct {
	inner    []SongData
	current  *SongData
	position int
	len      int
}

func (self *Queue) PlayNext() {
	if self.current != nil {
		self.current.Player.Reset()
		if a := self.position + 1; a != self.len {

			self.current = &self.inner[a]
			_, err := self.current.Player.Seek(0, io.SeekStart)
			if err != nil {
				log.Println("Could Not seek start")
			}
			self.current.Player.Play()
		} else {
			log.Println("The queue is Done!")
		}
	} else {
		// Start playing if not playing yet.
		log.Println("I AM HERE")
		self.current = &self.inner[0]
		self.current.Player.Play()
	}
	self.position += 1

}

func (self *Queue) Length() int {
	return self.len
}
func (self *Queue) Current() *SongData {
	return self.current
}
func (self *Queue) Pause() {
	self.current.Player.Pause()
}
func (self *Queue) PlayPrevious() {
	if self.position == 0 {
		_, err := self.current.Player.Seek(0, io.SeekStart)
		if err != nil {
			log.Println("Could Not seek start")
		}
	} else {
		self.position -= 1
		self.current = &self.inner[self.position]
		_, err := self.current.Player.Seek(0, io.SeekStart)
		if err != nil {
			log.Println("Could Not seek start")
		}
		self.current.Player.Play()
	}

}
