module github.com/lishimeng/app-starter

go 1.14

require (
	github.com/fatih/structs v1.1.0
	github.com/go-redis/cache/v8 v8.4.0
	github.com/go-redis/redis/v8 v8.11.4
	github.com/jeremywohl/flatten v1.0.1
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/kataras/iris/v12 v12.2.0-alpha9
	github.com/klauspost/compress v1.15.0 // indirect
	github.com/lishimeng/go-app-shutdown v1.0.1
	github.com/lishimeng/go-log v1.0.0
	github.com/lishimeng/go-orm v1.1.1
	github.com/spf13/viper v1.10.1
)

replace (
	github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/go-yaml/yaml/v2 => github.com/go-yaml/yaml/v2 v2.2.8
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.2.8
)
