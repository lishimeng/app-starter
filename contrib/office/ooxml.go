package office

import (
	"github.com/lishimeng/go-log"
	"os"
)

func ApplyWorkspace() (ws string, err error) {
	ws, err = os.MkdirTemp(os.TempDir(), "ooxml")
	return
}

func ClearWorkspace(ws string) {
	err := os.RemoveAll(ws)
	if err != nil {
		log.Info(err)
	}
}

// Marshall 编译成open xml文件
func Marshall(sourceFolder string, destFile string) (err error) {
	d := OpenXmlMarshaller{}
	err = d.Marshall(sourceFolder, destFile)
	return
}

// UnMarshall open xml文件反编译成文件明细
func UnMarshall(sourceFile string, destFolder string) (err error) {
	d := OpenXmlMarshaller{}
	err = d.UnMarshall(sourceFile, destFolder)
	return
}
