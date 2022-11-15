package fields

/*import (
	"encoding/json"
)*/

type Ref struct {
	Keys map[string]interface{} `json:"keys"`
	Descr string `json:"descr"`
	DataType string `json:"dataType"`
}
/*
func (r *Ref) UnmarshalJSON(d []byte) error {
	return json.Unmarshal(d, r)
}

func (r *Ref) MarshalJSON() ([]byte, error) {
	return json.Marshal(r)
}*/
