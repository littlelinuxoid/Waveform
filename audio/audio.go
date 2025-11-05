package audio

import (
	"Waveform/audio/tools"
	"log"
	"os"
	"path"
	"strings"

	"github.com/bogem/id3v2/v2"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

type SongData struct {
	Artist     string
	Album      string
	Length_sec int
	Cover      *[]byte
	Player     *oto.Player
	Genre      string
	Title      string
}

const (
	// This assumes 320kbps bitrate as it is the most common nowadays and it is not encoded into the file data
	// I am aware of TLEN tag in the ID3 standart, but it is literally NOWHERE to be found in about 700 mp3 files i tested.
	BITRATE      int64 = 320
	LENGTH_COEFF int64 = (BITRATE * 1000)
)

func FromDirectory(ctx *oto.Context, filepath string) *Queue {
	// this should get SongData objects for all valid mp3 files in a given directory and smash them into a Queue object
	// log.Printf("FromDirectory is called on %s\n", filepath)

	contents, err := os.ReadDir(filepath)
	if err != nil {
		log.Println(err.Error())
	}

	contents = tools.Filter(contents, func(de os.DirEntry) bool {
		return strings.ToLower(path.Ext(de.Name())) == ".mp3"
	})

	log.Println(contents)
	answer := make([]SongData, len(contents))

	for i, entry := range contents {
		f, err := os.Open(filepath + "/" + entry.Name())
		if err != nil {
			log.Println(err.Error())
		}
		sdata, err := extract_data(ctx, f)
		if err != nil {
			log.Println(err.Error())
		}
		answer[i] = sdata
		log.Println(answer)

	}

	q := &Queue{
		inner:    answer,
		len:      len(answer),
		position: 0,
		current:  nil,
	}
	return q
}

func extract_data(ctx *oto.Context, file *os.File) (SongData, error) {
	// avoid calling id3v2.Open since it opens the file itself, do not want to open the same file twice.
	log.Printf("Extracting data from %s", file.Name())
	data, err := id3v2.Open(file.Name(), id3v2.Options{Parse: true})
	if err != nil {
		log.Printf("Could not get information about file %s: "+err.Error(), file.Name())
		return SongData{}, err
	}

	info, err := file.Stat()
	if err != nil {
		log.Printf("Could not get information about file %s: %s", file.Name(), err.Error())
		return SongData{}, err
	}
	cover, err := tools.Extract_cover(data)
	if err != nil {
		log.Printf("File %s has no picture attached to it, using default", file.Name())
		temp, _ := os.Open("./resources/default/music_icon.png")
		if _, err := temp.Read(cover); err != nil {
			log.Println("The default icon is not present. Using no icon")
		}

	}
	// subtract size of a picture to get approximate the size of audio.
	// File also might contain other fields like lyrics, year and others, but i will worry later.
	secs := int((info.Size() - int64(len(cover))) * 8 / (LENGTH_COEFF))
	log.Println("Success")
	return SongData{
		Artist:     data.Artist(),
		Album:      data.Album(),
		Length_sec: secs,
		Player:     OneSong(ctx, file.Name()),
		Cover:      &cover,
		Genre:      data.Genre(),
		Title:      data.Title(),
	}, nil

}

func OneSong(ctx *oto.Context, filepath string) *oto.Player {
	// log.Printf("OneSong is called on %s\n", filepath)
	file, err := os.Open(filepath)
	if err != nil {
		panic("I just shit myself" + err.Error())
	}

	decoded, err := mp3.NewDecoder(file)
	if err != nil {
		panic("I Just Decoded Myself " + err.Error())
	}

	player := ctx.NewPlayer(decoded)

	return player
}
