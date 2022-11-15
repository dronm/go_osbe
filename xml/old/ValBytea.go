package xml

import (
	"encoding/xml"
	
	"osbe/fields"
)

type ValBytea struct {
	fields.ValBytea
}

func (v *ValBytea) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	TextFieldToXML(v, e, start)
	return nil
}


