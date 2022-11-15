package xml

import (
	"encoding/xml"
	
	"osbe/fields"
)

type ValText struct {
	fields.ValText
}

func (v *ValText) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	TextFieldToXML(v, e, start)
	return nil
}


