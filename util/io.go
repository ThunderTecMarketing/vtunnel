package util

import (
	"io"
	"os/exec"
	"os"
	"fmt"
	"path/filepath"
)

func Pipe(dst, src io.ReadWriteCloser) {
	go func() {
		defer src.Close()
		io.Copy(src, dst)
	}()
	defer dst.Close()
	io.Copy(dst, src)
}



func GetCurrentExecDir() (dir string, err error) {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		fmt.Printf("exec.LookPath(%s), err: %s\n", os.Args[0], err)
		return "", err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("filepath.Abs(%s), err: %s\n", path, err)
		return "", err
	}

	dir = filepath.Dir(absPath)

	return dir, nil
}