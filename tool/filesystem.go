package tool

import (
	"io/fs"
	"os"
	"path/filepath"
)

func Walk(root string, depth int, fun func(path string, info os.FileInfo)) (err error) {
	var folders []string
	err = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if root == path {
			return nil
		}
		if info.IsDir() {
			folders = append(folders, path)
		}
		fun(path, info)
		if depth > 0 && len(folders) > 0 {
			nextDepth := depth - 1
			for _, p := range folders {
				e := Walk(p, nextDepth, fun)
				if e != nil {
					return e
				}
			}
		}
		if info.IsDir() {
			return filepath.SkipDir
		} else {
			return nil
		}
	})
	return
}
