package etc

import "github.com/lishimeng/go-libs/etc"

func Load(config interface{}, name string, path ...string) (err error) {
	_, err = etc.LoadEnvs(name, path, config)

	return
}
