package office

import (
	"archive/zip"
	"bytes"
	"github.com/lishimeng/go-log"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type OpenXmlMarshaller struct {
	buf *bytes.Buffer
}

func (d *OpenXmlMarshaller) Marshall(inputFolder string, outFile string) (err error) {
	// 将一个文件夹压缩
	d.buf = new(bytes.Buffer)

	err = d.compact(inputFolder)
	if err != nil {
		log.Info(err)
		return
	}
	d.release(outFile)
	return
}

func (d *OpenXmlMarshaller) compact(inputFolder string) (err error) {
	w := zip.NewWriter(d.buf)

	defer func() {
		_ = w.Close()
	}()

	err = filepath.Walk(inputFolder, func(path string, info fs.FileInfo, e error) error {
		if inputFolder == path {
			return nil
		}
		header, e := buildFileHeader(inputFolder, path, info)
		if e != nil {
			log.Info(e)
			return e
		}
		if info.IsDir() {
			return nil
		}
		writer, e := w.CreateHeader(header)
		if e != nil {
			log.Info(e)
			return e
		}
		//log.Info("add file:%s", header.Name)
		e = cp(writer, path)
		if e != nil {
			log.Info(e)
			return e
		}
		return nil
	})
	if err != nil {
		log.Info(err)
		return
	}
	return
}

func (d *OpenXmlMarshaller) release(out string) {
	f, err := os.Create(out)
	if err != nil {
		log.Info(err)
		return
	}
	_, err = io.Copy(f, d.buf)
	//_, err = d.buf.WriteTo(f)
	if err != nil {
		log.Info(err)
		return
	}
}

func (d *OpenXmlMarshaller) UnMarshall(inputFile string, outFolder string) (err error) {

	d.buf = new(bytes.Buffer)
	return d.unCompact(inputFile, outFolder)
}

func (d *OpenXmlMarshaller) unCompact(inputFile string, outFolder string) (err error) {
	r, err := zip.OpenReader(inputFile)
	if err != nil {
		log.Info(err)
		return
	}
	abs, err := filepath.Abs(outFolder)
	if err != nil {
		log.Info(err)
		return
	}

	for _, item := range r.File {
		//log.Info("extract: %s", item.Name)
		p := filepath.Join(abs, item.Name)
		err = extractFile(item, p)
		if err != nil {
			log.Info(err)
			return
		}
	}
	return
}

func extractFile(item *zip.File, dest string) (err error) {
	if item.FileInfo().IsDir() {
		err = os.MkdirAll(dest, 755)
		if err != nil {
			log.Info(err)
			return
		}
	}
	rc, err := item.Open()
	if err != nil {
		log.Info(err)
		return
	}
	defer func() { _ = rc.Close() }()

	dir := filepath.Dir(dest)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		log.Info(err)
		return
	}
	wc, err := os.Create(dest)
	if err != nil {
		log.Info(err)
		return
	}
	_, err = io.Copy(wc, rc)
	return
}
