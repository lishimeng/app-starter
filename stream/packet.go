package stream

type PacketProcessor interface {
	Listen(onPacket func(p []byte))
	Data(data []byte) (n int)
}
