package github

import "encoding/json"

func UnmarshalFlogoActivity(data []byte) (FlogoActivity, error) {
	var r FlogoActivity
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *FlogoActivity) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type FlogoActivity struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Ref         string `json:"ref"`
	Version     string `json:"version"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}
