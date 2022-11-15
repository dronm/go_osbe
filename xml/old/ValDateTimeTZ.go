package xml

import (
	"encoding/xml"
	
	"osbe/fields"
)
type ValDateTimeTZ struct {
	fields.ValDateTimeTZ
}

func (v *ValDateTimeTZ) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	TextFieldToXML(v, e, start)
	return nil
}


