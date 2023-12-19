package docAttachment

import(
	"errors"
	
	"osbe/fields"
)

type HttpInt int
func (v *HttpInt) UnmarshalJSON(data []byte) error {
	
	if fields.ExtValIsNull(data){
		*v = 0
		return nil
	}
	
	v_str := fields.ExtRemoveQuotes(data)
	temp, err := fields.StrToInt(v_str)
	if err != nil {
		return errors.New(fields.ER_UNMARSHAL_INT + err.Error())
	}
	*v = HttpInt(temp)
	
	return nil	
}

type AttachmentKey struct {
	Id HttpInt `json:"id"`
}

//Attachment reference
type Ref_Type struct {
	Keys AttachmentKey `json:"keys"`
	DataType string `json:"dataType"`
}
func NewRef_Type(id int64, dataType string) *Ref_Type {
	return &Ref_Type{Keys: AttachmentKey{Id: HttpInt(id)}, DataType: dataType}
}

//Attachment info
type Content_info_Type struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Size int64 `json:"size"`
}
func NewContent_info_Type(id, name string) *Content_info_Type {
	return &Content_info_Type{Id: id, Name: name}
}

