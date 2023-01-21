package gfmedia

import (
	"encoding/binary"
	"fmt"
)

var (
	ErrChunkSizeTooLarge = fmt.Errorf("riff chunk size too large")
	ErrChunkSizeInvalid  = fmt.Errorf("riff chunk size invalid")
	ErrExtraChunkInvalid = fmt.Errorf("extra riff chunk invalid")
)

type Riff struct {
	ChunkID   [4]byte
	ChunkSize uint32
	Format    [4]byte
}

func newRiff(chunkSize uint32) Riff {
	ret := Riff{
		ChunkSize: chunkSize,
	}
	copy(ret.ChunkID[:], "RIFF")
	copy(ret.Format[:], "WAVE")
	return ret
}

func (r Riff) Len() int {
	return 16
}

type FmtStandard struct {
	FmtChunkID    [4]byte
	FmtChunkSize  uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}

type FmtExtra struct {
	UnknownFmtBytes []byte
}

type Fmt struct {
	FmtStandard
	FmtExtra
}

func newFmt(numChannels uint16, sampleRate uint32, bitsPerSample uint16) Fmt {
	ret := Fmt{
		FmtStandard: FmtStandard{
			FmtChunkSize:  16,
			AudioFormat:   1, // 1为PCM编码格式
			NumChannels:   numChannels,
			SampleRate:    sampleRate,
			ByteRate:      sampleRate * uint32(numChannels) * uint32(bitsPerSample) / 8,
			BlockAlign:    numChannels * bitsPerSample / 8,
			BitsPerSample: bitsPerSample,
		},
	}
	copy(ret.FmtChunkID[:], "fmt ")
	return ret
}

func (f Fmt) Len() int {
	return int(8 + f.FmtChunkSize)
}

type Data struct {
	DataChunkID   [4]byte
	DataChunkSize uint32
}

func newData(chunkSize uint32) Data {
	ret := Data{
		DataChunkSize: chunkSize,
	}
	copy(ret.DataChunkID[:], "data")
	return ret
}

type Extra struct {
	SubChunk
}

func (e Extra) Len() int {
	return e.SubChunk.Len()
}

func (e Extra) IsValid() bool {
	return e.SubChunkSize == uint32(len(e.SubChunkData))
}

type Extras []Extra

func (e Extras) Len() int {
	sum := 0
	for _, v := range e {
		sum += v.Len()
	}
	return sum
}

func (e Extras) Marshal() ([]byte, error) {
	bs := make([]byte, 0, e.Len())
	for _, v := range e {
		if !v.IsValid() {
			return nil, ErrExtraChunkInvalid
		}
		bs = append(bs, v.SubChunkID[:]...)
		subChunkSize := [4]byte{}
		binary.LittleEndian.PutUint32(subChunkSize[:], v.SubChunkSize)
		bs = append(bs, subChunkSize[:]...)
		bs = append(bs, v.SubChunkData...)
	}
	return bs, nil
}

type SubChunk struct {
	SubChunkID   [4]byte
	SubChunkSize uint32
	SubChunkData []byte
}

func (s SubChunk) Len() int {
	return int(8 + s.SubChunkSize)
}
