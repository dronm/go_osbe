package osbe

type EnumDescrCollection map[string]string
type Enum map[string]EnumDescrCollection

func (e *Enum) CheckValue(v string) bool{
	_, ok := (*e)[v]
	return ok
}
func (e *Enum) GetDescription(v string, localeID string) string{
	if descr, ok := (*e)[v]; ok {
		if descr_v, descr_v_ok := descr[localeID]; descr_v_ok {
			return descr_v
		}		
	}
	return ""
}

type Enumer interface {
	CheckValue(string) bool
	GetDescription(string, string) string
}

type EnumCollection map[string] Enumer
