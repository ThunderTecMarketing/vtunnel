package msocks

//implements FrameSender interface
func (s *Session) SendFrame(f Frame) (err error) {
	err = s.conn.Write(f)
	if err != nil {
		return
	}

	return
}
