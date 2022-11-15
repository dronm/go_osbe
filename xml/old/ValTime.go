package xml

import (
	"encoding/xml"
	
	"osbe/fields"
)

type ValTime struct {
	fields.ValTime
}

func (v *ValTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	TextFieldToXML(v, e, start)
	return nil
}


