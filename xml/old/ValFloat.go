package xml

import (
	"encoding/xml"
	
	"osbe/fields"
)
type ValFloat struct {
	fields.ValFloat
}

func (v *ValFloat) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	TextFieldToXML(v, e, start)
	return nil
}


