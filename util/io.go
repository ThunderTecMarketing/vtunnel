package util

import "io"

func Pipe(dst, src io.ReadWriteCloser) {
	go func() {
		defer src.Close()
		io.Copy(src, dst)
	}()
	defer dst.Close()
	io.Copy(dst, src)
}


