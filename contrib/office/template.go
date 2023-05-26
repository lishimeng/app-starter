package office

import (
	"bytes"
	"github.com/lishimeng/go-log"
	"io"
	"os"
	"text/template"
)

func RenderTemplate(file string, data any) (err error) {
	tpl, err := readTpl(file)
	err = render(tpl, data, file)

	if err != nil {
		log.Info(err)
		return
	}
	return
}

func render(tpl *template.Template, data any, dest string) (err error) {

	log.Info("render template:%s", dest)
	wc, err := os.Create(dest)
	if err != nil {
		log.Info(err)
		return
	}
	defer func() {
		_ = wc.Close()
	}()

	err = tpl.Execute(wc, data)
	return
}

func readTpl(tplPath string) (tpl *template.Template, err error) {
	buf := new(bytes.Buffer)
	rc, err := os.Open(tplPath)
	if err != nil {
		log.Info(err)
		return
	}
	defer func() {
		_ = rc.Close()
	}()
	_, err = io.Copy(buf, rc)
	if err != nil {
		log.Info(err)
		return
	}
	tpl, err = template.New("bill").Parse(buf.String())
	if err != nil {
		log.Info(err)
		return
	}
	return
}
