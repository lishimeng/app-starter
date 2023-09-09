package cms

import "testing"

func TestInit(t *testing.T) {
	ws := WebSiteInfo{
		Name:      "test dummy",
		BaseUrl:   "http://localhost:80",
		Copyright: "copyright @ me",
		Icp:       "icp format",
		Favicon:   "icon url",
		Logo:      "logo url",
	}
	Init(WithName("dummy"), WithRedis(), WithDatabase(), WithConfigFile(ws))
	if c.Logo != ws.Logo {
		t.Fatal("logo not match")
	}
}
