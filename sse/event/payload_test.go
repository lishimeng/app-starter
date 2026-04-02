package event

import (
	"encoding/json"
	"testing"
)

func TestPayloadMarshall(t *testing.T) {

	var data = make(map[string]interface{})
	data["a"] = 13
	data["b"] = "sample"
	bs, _ := json.Marshal(data)
	var p = Payload{
		Data:  bs,
		Event: "",
		ID:    "24352343424",
		Time:  0,
	}
	s := p.Marshall()
	t.Log(s)
}
