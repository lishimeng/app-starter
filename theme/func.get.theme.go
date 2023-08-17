package theme

import (
	"encoding/json"
	"fmt"
	"github.com/lishimeng/app-starter/factory"
	"github.com/lishimeng/app-starter/persistence"
	"github.com/lishimeng/go-log"
	"strconv"
)

//GetPageTheme 开放方法：获取页面主题配置。从缓存中获取
func GetPageTheme(page string) (configs []themeConfig) {
	if !existPageCache(page) {
		return
	}
	configs, err := getPageCache(page)
	if err == nil {
		return
	}
	log.Debug(err)
	configs = GetPageThemeSkipCache(page)
	err = setPageCache(page, configs)
	if err != nil {
		log.Debug(err)
		return
	}
	return
}

//GetPageThemeSkipCache 开放方法：获取页面主题配置。跳过缓存
func GetPageThemeSkipCache(page string) (configs []themeConfig) {
	var themeConfigs []AppThemeConfig
	_, err := persistence.New().Context.QueryTable(new(AppThemeConfig)).
		Filter("AppName", AppName).
		Filter("ConfigPage", page).
		All(&themeConfigs)
	if err != nil {
		log.Debug(err)
		return
	}
	if len(themeConfigs) == 0 {
		return
	}
	for _, item := range themeConfigs {
		var configValue interface{}
		switch item.ConfigContentType {
		case string(NumberConfigContentType):
			configValue, _ = strconv.Atoi(item.ConfigContent)
		case string(BooleanConfigContentType):
			configValue, _ = strconv.ParseBool(item.ConfigContent)
		default:
			configValue = item.ConfigContent
		}
		if configValue == nil {
			return
		}
		configs = append(configs, themeConfig{
			Id:                item.Id,
			AppName:           item.AppName,
			ConfigPage:        item.ConfigPage,
			ConfigName:        item.ConfigName,
			ConfigContent:     configValue,
			ConfigContentType: item.ConfigContentType,
		})
	}
	return
}

//FormatPageTheme 格式化->map
func FormatPageTheme(configs []themeConfig) (mapConfigs map[string]interface{}) {
	mapConfigs = make(map[string]interface{})
	for _, config := range configs {
		mapConfigs[config.ConfigName] = config.ConfigContent
	}
	return
}

func setPageCache(page string, configs []themeConfig) (err error) {
	if factory.GetCache() == nil {
		return
	}
	data, err := json.Marshal(configs)
	if err != nil {
		return
	}
	return factory.GetCache().Set(pageKey(page), data)
}

func existPageCache(page string) bool {
	if factory.GetCache() == nil {
		return false
	}
	return factory.GetCache().Exists(pageKey(page))
}

func getPageCache(page string) (configs []themeConfig, err error) {
	if factory.GetCache() == nil {
		return
	}
	data := make([]byte, 0)
	err = factory.GetCache().GetSkipLocal(pageKey(page), &data)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &configs)
	if err != nil {
		return
	}
	return
}

func pageKey(page string) string {
	return fmt.Sprintf(pageThemeCacheKeyTpl, AppName, page)
}
