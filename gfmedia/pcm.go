package gfmedia

type PcmAudio struct {
	raw []byte
}

func NewPcmAudio(raw []byte) *PcmAudio {
	return &PcmAudio{raw: raw}
}

func (a PcmAudio) ToWav(h *WavHead) (*WavAudio, error) {
	rawLen := len(a.raw)

	head, err := h.Serialize()
	if err != nil {
		return nil, err
	}
	wav := make([]byte, len(head)+rawLen)
	copy(wav, head)
	copy(wav[len(head):], a.raw)

	return &WavAudio{raw: wav}, nil
}
