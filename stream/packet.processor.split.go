package stream

// SplitPacketProcessor 分隔符分包器, 每个包用Split标记结尾
type SplitPacketProcessor struct {
	split    byte
	callback func([]byte)
}

func NewSplitPP(split byte) (pp PacketProcessor) {
	pp = &SplitPacketProcessor{split: split}
	return
}

func (spp *SplitPacketProcessor) Data(data []byte) (n int) {
	// find split
	var to = -1
	var from = 0
	for i, d := range data {
		if d == spp.split {
			to = i + 1
			size := spp.handlePacket(data, from, to)
			n = n + size
			from = i + 1
		}
	}

	return
}

func (spp *SplitPacketProcessor) handlePacket(data []byte, from, to int) (size int) {

	packet := data[from:to]
	size = len(packet)
	if spp.callback != nil {
		spp.callback(packet)
	}
	return
}

func (spp *SplitPacketProcessor) Listen(onPacket func(p []byte)) {
	spp.callback = onPacket
	return
}
