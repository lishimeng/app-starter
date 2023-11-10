package stream

func (s *SessionCtx) Write(p []byte) (size int, err error) {
	size, err = s.proxy.Write(p)
	return
}
