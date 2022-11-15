package xml

import (
	"encoding/xml"
	
	"osbe/fields"
)

type ValJSON struct {
	fields.ValJSON
}

func (v *ValJSON) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	TextFieldToXML(v, e, start)
	return nil
}


