package etc

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/jeremywohl/flatten"
	"github.com/spf13/viper"
	"regexp"
	"strings"
)

type Loader interface {
	Load(target interface{}) error
	SetEnvPrefix(name string) Loader
	SetFileSearcher(configName string, searchPath ...string) Loader
	SetEnvSearcher() Loader
}

type loader struct {
	v               *viper.Viper
	name            string
	envSearchEnable bool
}

func New() (o Loader) {
	config := viper.New()
	t := &loader{
		v: config,
	}
	o = t
	return
}

func (t *loader) SetEnvPrefix(name string) Loader {
	if len(name) <= 0 {
		return t
	}
	matched, _ := regexp.MatchString(`[a-zA-Z]+`, name)
	if !matched {
		return t
	}
	t.name = name
	return t
}

func (t *loader) SetFileSearcher(configName string, searchPath ...string) Loader {

	t.v.SetConfigName(configName)
	if len(searchPath) > 0 {
		for _, p := range searchPath {
			t.v.AddConfigPath(p)
		}
	}
	e := t.v.ReadInConfig()
	if e != nil {
		fmt.Println(e)
	}
	return t
}

func (t *loader) SetEnvSearcher() Loader {
	t.envSearchEnable = true
	return t
}

func (t *loader) prepareEnv(target interface{}) Loader {
	t.v.AutomaticEnv()
	if len(t.name) > 0 {
		t.v.SetEnvPrefix(strings.ToUpper(t.name))
	} else {
		t.v.SetEnvPrefix("")
	}
	t.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	confMap := structs.Map(target)
	flat, err := flatten.Flatten(confMap, "", flatten.DotStyle)
	if err != nil {
		return t
	}

	for key := range flat {
		err = t.v.BindEnv(key)
		if err != nil {
			return t
		}
	}
	return t
}

func (t *loader) Load(target interface{}) (err error) {

	if len(t.name) > 0 {
		_ = t.v.ReadInConfig()
	}

	if t.envSearchEnable {
		t.prepareEnv(target)
	}

	err = t.v.Unmarshal(target)
	return
}
