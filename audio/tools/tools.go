package tools

import (
	"errors"
	"strconv"

	"github.com/bogem/id3v2/v2"
	"github.com/ebitengine/oto/v3"
)

func NewContext() *oto.Context {
	op := &oto.NewContextOptions{}
	op.SampleRate = 44100
	op.ChannelCount = 2
	op.Format = oto.FormatSignedInt16LE
	ctx, channel, err := oto.NewContext(op)
	if err != nil {
		panic("Could Not create a context, Please report this issue!")
	}
	<-channel
	return ctx

}
func Extract_cover(Data *id3v2.Tag) ([]byte, error) {
	cover := Data.GetLastFrame("APIC")
	if cover == nil {
		return nil, errors.New("No picture attached")

	}
	return cover.(id3v2.PictureFrame).Picture, nil

}

func Filter[T any](ts []T, op func(T) bool) []T {
	result := make([]T, len(ts))
	counter := 0
	for _, val := range ts {
		if op(val) {
			result[counter] = val
			counter++
		}
	}
	return result[:counter]
}
func Format_lenth(input int) string {
	var answer string
	answer += strconv.Itoa(input / 60) //minutes
	answer += ":"
	answer += strconv.Itoa(input % 60)

	return answer
}
func Map[T, U any](ts []T, op func(T) U) []U {
	result := make([]U, len(ts))
	for i, t := range ts {
		result[i] = op(t)
	}
	return result

}
