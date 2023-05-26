package office

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func buildFileHeader(base string, path string, info fs.FileInfo) (header *zip.FileHeader, err error) {
	p, err := filepath.Rel(base, path)
	if err != nil {
		return nil, err
	}
	header, err = zip.FileInfoHeader(info)
	if err != nil {
		return nil, err
	}
	header.Name = p
	header.Method = zip.Deflate
	if info.IsDir() {
		header.Name += string(os.PathSeparator)
		//log.Info("add folder:%s", header.Name)
	}
	header.Modified = time.Unix(info.ModTime().Unix(), 0)
	return header, nil
}

func cp(w io.Writer, src string) (err error) {
	if err != nil {
		return
	}

	rc, err := os.Open(src)
	if err != nil {
		return
	}
	defer func() {
		_ = rc.Close()
	}()

	_, err = io.Copy(w, rc)
	if err != nil {
		return
	}

	return
}
