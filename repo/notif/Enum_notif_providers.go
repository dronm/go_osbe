package notif

import (
	"osbe/fields"
)

type ValEnum_notif_providers struct {
	fields.ValText
}

func (e *ValEnum_notif_providers) GetValues() []string {
	return []string{ "email", "sms", "wa", "tm", "vb" }
}

//func (e *ValEnum_notif_providers) GetDescriptions() map[string]map[string]string {
//	return make(map[string]{ "email", "sms", "wa", "tm", "vb" }
//}

