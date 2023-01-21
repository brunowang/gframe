package gfmedia

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestWavToPcm(t *testing.T) {
	assertPcm, err := os.ReadFile("input.pcm")
	if err != nil {
		panic(err)
	}

	wav, err := os.ReadFile("input.wav")
	if err != nil {
		panic(err)
	}
	pcm, err := NewWavAudio(wav).ToPcm()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, len(pcm.raw), len(assertPcm))

	pcmf, err := os.Create("output.pcm")
	if err != nil {
		panic(err)
	}
	defer pcmf.Close()
	if _, err := pcmf.Write(pcm.raw); err != nil {
		panic(err)
	}
}

func TestPcmToWav(t *testing.T) {
	assertWav, err := os.ReadFile("input.wav")
	if err != nil {
		panic(err)
	}
	head, err := NewWavAudio(assertWav).GetHead()
	if err != nil {
		panic(err)
	}

	pcm, err := os.ReadFile("input.pcm")
	if err != nil {
		panic(err)
	}
	wav, err := NewPcmAudio(pcm).ToWav(head)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, len(wav.raw), len(assertWav))

	wavf, err := os.Create("output.wav")
	if err != nil {
		panic(err)
	}
	defer wavf.Close()
	if _, err := wavf.Write(wav.raw); err != nil {
		panic(err)
	}
}
