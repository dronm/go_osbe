package xml

import (
	"encoding/xml"
	
	"osbe/fields"
)
type ValDate struct {
	fields.ValDate
}

func (v *ValDate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	TextFieldToXML(v, e, start)
	return nil
}


