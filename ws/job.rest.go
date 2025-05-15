package ws

type RestJob struct {
	Api     string            `json:"api,omitempty"`
	Schema  string            `json:"schema,omitempty"`
	Method  string            `json:"method,omitempty"`
	Host    string            `json:"host,omitempty"`
	Path    string            `json:"path,omitempty"`
	Query   map[string]string `json:"query,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Cookies map[string]string `json:"cookies,omitempty"`
}

var RestFunc func(rj *RestJob, respPtr any) error

func (rj *RestJob) Fetch(respPtr any) (err error) {
	if RestFunc == nil {
		return err
	}
	err = RestFunc(rj, respPtr)
	return
}
