package msocks

//implements FrameReceiver interface
func (c *Conn) ReceiveFrame(f Frame) (err error) {
	switch ft := f.(type) {
	default:
		err = ErrUnexpectedPkg
		log.Errorf("%s", err)
		c.Close()
		return
	case *FrameSynResult:
		return c.inSynResult(ft.Errno)
	case *FrameData:
		return c.inData(ft)
	case *FrameFin:
		return c.inFin(ft)
	case *FrameRst:
		log.Debugf("reset %s.", c.String())
		c.final()
	}
	return
}

func (c *Conn) inSynResult(errno uint32) (err error) {
	c.statusLock.Lock()
	defer c.statusLock.Unlock()

	if c.status != ST_SYN_SENT {
		return ErrNotSyn
	}

	if errno == ERR_NONE {
		c.status = ST_EST
	} else {
		c.final()
	}

	select {
	case c.chSynResult <- errno:
	default:
	}
	return
}

func (c *Conn) inData(ft *FrameData) (err error) {
	log.Infof("%s recved %d bytes.", c.String(), len(ft.Data))
	err = c.rqueue.Push(ft.Data)
	if err != nil {
		return
	}
	return
}

func (c *Conn) inFin(ft *FrameFin) (err error) {
	// always need to close read pipe
	// coz fin means remote will never send data anymore
	c.rqueue.Close()

	c.statusLock.Lock()
	defer c.statusLock.Unlock()

	c.final()
	return
}


func (c *Conn) CloseFrame() error {
	// maybe conn closed
	c.final()
	return nil
}