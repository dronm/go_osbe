package xml

import (
	"encoding/xml"
	
	"osbe/fields"
)

type ValAssocArray struct {
	fields.ValAssocArray
}

func (v *ValAssocArray) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if !v.GetIsNull() {
		//e.EncodeElement(v.String(), start)
		/*for el_k, el_v := range v.TypedValue {
			tag := xml.StartElement{Name: xml.Name{"", el_k}, Attr: nil}
			xml.EscapeText(, []byte(el_v))
			e.EncodeElement(el_v, tag)
		}*/
		
	}else{
		EncodeEmptyElement(e, start)
	}

	return nil
}


