package osbe

// This file describes Enum object.
// Enums are lists of fixed predefined values.
// Enum must implement Enumer interface.
// Enum may have descriptions in many languages.

type EnumDescrCollection map[string]string 	//holds language specific description of enum values.
						//key is enum ID

type Enum map[string]EnumDescrCollection

// CheckValue checks if a given value is in enum value list of values.
func (e *Enum) CheckValue(v string) bool {
	_, ok := (*e)[v]
	return ok
}
// GetDescription retrieves enum description for a specific language.
func (e *Enum) GetDescription(v string, localeID string) string{
	if descr, ok := (*e)[v]; ok {
		if descr_v, descr_v_ok := descr[localeID]; descr_v_ok {
			return descr_v
		}		
	}
	return ""
}

// Enumer 
type Enumer interface {
	CheckValue(string) bool
	GetDescription(string, string) string
}

// EnumCollection is an application enum list.
type EnumCollection map[string] Enumer
