package xml

import (
	"encoding/xml"
	
	"osbe/fields"
)
type ValArray struct {
	fields.ValArray
}

func (v *ValArray) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	TextFieldToXML(v, e, start)
	return nil
}


