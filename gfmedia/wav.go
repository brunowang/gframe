package gfmedia

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	maxSubChunkCount = 20
	maxSubChunkSize  = 1 << 20
)

var (
	ErrNonstandardWavFile = fmt.Errorf("nonstandard wav file")
)

type WavAudio struct {
	raw  []byte
	head *WavHead
}

func NewWavAudio(raw []byte) *WavAudio {
	return &WavAudio{raw: raw}
}

func (a *WavAudio) GetHead() (*WavHead, error) {
	if a.head != nil {
		return a.head, nil
	}
	head := WavHead{}
	if err := head.Deserialize(a.raw); err != nil {
		return nil, err
	}
	a.head = &head
	return a.head, nil
}

func (a WavAudio) ToPcm() (*PcmAudio, error) {
	head, err := a.GetHead()
	if err != nil {
		return nil, err
	}
	pcm := make([]byte, len(a.raw)-head.Len())
	copy(pcm, a.raw[head.Len():])

	return &PcmAudio{raw: pcm}, nil
}

type WavHead struct {
	Riff
	Fmt
	Extras
	Data
}

func NewWavHead(numChannels uint16, sampleRate uint32, bitsPerSample uint16,
	bodyLen uint32, extras ...Extra) *WavHead {
	return &WavHead{
		Riff:   newRiff(36 + bodyLen + uint32(Extras(extras).Len())),
		Fmt:    newFmt(numChannels, sampleRate, bitsPerSample),
		Extras: extras,
		Data:   newData(bodyLen),
	}
}

func (h *WavHead) ExtendFmtChunk(ext FmtExtra) *WavHead {
	h.UnknownFmtBytes = append(h.UnknownFmtBytes, ext.UnknownFmtBytes...)
	h.FmtChunkSize += uint32(len(ext.UnknownFmtBytes))
	h.ChunkSize += uint32(len(ext.UnknownFmtBytes))
	return h
}

func (h *WavHead) Serialize() ([]byte, error) {
	ext, err := h.Extras.Marshal()
	if err != nil {
		return nil, err
	}
	buf, byteOrder := &bytes.Buffer{}, binary.LittleEndian

	if err := binary.Write(buf, byteOrder, h.Riff); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, byteOrder, h.FmtStandard); err != nil {
		return nil, err
	}
	if _, err := buf.Write(h.FmtExtra.UnknownFmtBytes); err != nil {
		return nil, err
	}
	if _, err := buf.Write(ext); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, byteOrder, h.Data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (h *WavHead) Deserialize(raw []byte) error {
	wav, head := bytes.NewReader(raw), WavHead{}
	if err := readRiff(wav, &head); err != nil {
		return err
	}

	for i := 0; i < maxSubChunkCount; i++ {
		chunk, err := readChunk(wav)
		if err != nil {
			return err
		}

		subChunkID := string(chunk.SubChunkID[:])
		if subChunkID == "data" {
			head.DataChunkID = chunk.SubChunkID
			head.DataChunkSize = chunk.SubChunkSize
			break
		}

		if subChunkID == "fmt " {
			head.FmtChunkID = chunk.SubChunkID
			head.FmtChunkSize = chunk.SubChunkSize
			if head.FmtChunkSize < 16 {
				return ErrChunkSizeInvalid
			}
			subChunk, byteOrder := chunk.SubChunkData, binary.LittleEndian

			head.AudioFormat = byteOrder.Uint16(subChunk[:2])
			head.NumChannels = byteOrder.Uint16(subChunk[2:4])
			head.SampleRate = byteOrder.Uint32(subChunk[4:8])
			head.ByteRate = byteOrder.Uint32(subChunk[8:12])
			head.BlockAlign = byteOrder.Uint16(subChunk[12:14])
			head.BitsPerSample = byteOrder.Uint16(subChunk[14:16])
			head.UnknownFmtBytes = subChunk[16:]
			continue
		}

		ext := Extra{
			SubChunk: SubChunk{
				SubChunkID:   chunk.SubChunkID,
				SubChunkSize: chunk.SubChunkSize,
				SubChunkData: chunk.SubChunkData,
			},
		}
		head.Extras = append(head.Extras, ext)
	}
	*h = head
	return nil
}

func (h *WavHead) Len() int {
	// ChunkSize不包含ChunkID和ChunkSize这8个字节，需要加回来
	return int(h.ChunkSize + 8 - h.DataChunkSize)
}

func read(r io.Reader, s []byte) error {
	n, e := r.Read(s)
	if e != nil {
		return e
	} else if n < len(s) {
		return io.EOF
	}
	return nil
}

func readRiff(r io.Reader, head *WavHead) error {
	chunkSize := [4]byte{}
	if err := read(r, head.ChunkID[:]); err != nil {
		return err
	}
	if err := read(r, chunkSize[:]); err != nil {
		return err
	}
	head.ChunkSize = binary.LittleEndian.Uint32(chunkSize[:])

	if err := read(r, head.Format[:]); err != nil {
		return err
	}
	if string(head.Format[:]) != "WAVE" {
		return ErrNonstandardWavFile
	}
	return nil
}

func readChunk(r io.Reader) (*SubChunk, error) {
	chunk := &SubChunk{}
	if err := read(r, chunk.SubChunkID[:]); err != nil {
		return nil, err
	}
	subChunkSize := [4]byte{}
	if err := read(r, subChunkSize[:]); err != nil {
		return nil, err
	}
	chunk.SubChunkSize = binary.LittleEndian.Uint32(subChunkSize[:])

	if string(chunk.SubChunkID[:]) == "data" {
		return chunk, nil
	}
	if chunk.SubChunkSize > maxSubChunkSize {
		return nil, ErrChunkSizeTooLarge
	}

	chunk.SubChunkData = make([]byte, chunk.SubChunkSize)
	if err := read(r, chunk.SubChunkData[:]); err != nil {
		return nil, err
	}
	return chunk, nil
}
