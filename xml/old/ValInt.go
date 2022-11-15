package xml

import (
	"encoding/xml"
	
	"osbe/fields"
)

type ValInt struct {
	fields.ValInt
}

func (v *ValInt) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	TextFieldToXML(v, e, start)
	return nil
}


