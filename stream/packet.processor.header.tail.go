package stream

// HeaderTailPacketProcessor 固定包头包尾分包器
type HeaderTailPacketProcessor struct {
	head     []byte
	tail     []byte
	callback func([]byte)
}

func NewHeadTailPP(head []byte, tail []byte) (pp PacketProcessor) {
	pp = &HeaderTailPacketProcessor{head: head, tail: tail}
	return
}

func (spp *HeaderTailPacketProcessor) Data(data []byte) (n int) {
	// find split
	var to = -1
	var from = -1
	for i := range data {
		if spp.findSub(i, spp.head, data) {
			from = i
			continue
		}
		if spp.findSub(i, spp.tail, data) {
			to = i + len(spp.tail)
			size := spp.handlePacket(data, from, to)
			n = n + size
			continue
		}
	}

	return
}
func (spp *HeaderTailPacketProcessor) findSub(offset int, compare []byte, data []byte) (find bool) {
	size := len(data)
	if offset+len(compare) > size {
		find = false
		return
	}
	if string(compare) != string(data[offset:offset+len(compare)]) {
		find = false
		return
	}

	find = true
	return
}

func (spp *HeaderTailPacketProcessor) handlePacket(data []byte, from, to int) (size int) {

	packet := data[from:to]
	size = len(packet)
	if spp.callback != nil {
		spp.callback(packet)
	}
	return
}

func (spp *HeaderTailPacketProcessor) Listen(onPacket func(p []byte)) {
	spp.callback = onPacket
	return
}
