package socket

import (
	"encoding/gob"
			
	"osbe/sql"
)

const SESS_PRESET_FILTER = "PRESET_FILTER"

var ps_registered bool

//Individual global filter
type PresetFilter map[string]sql.FilterCondCollection

func (f *PresetFilter) Add(modelID string, conditions sql.FilterCondCollection) {
	(*f)[modelID] = conditions
}
func (f *PresetFilter) Get(modelID string) sql.FilterCondCollection {
	if v, ok := (*f)[modelID]; ok {
		return v
	}
	return nil
}

func RegisterPresetFilter(){
	if !ps_registered {
		gob.Register(PresetFilter{})	
		ps_registered = true
	}
}

func NewPresetFilter() PresetFilter{
	RegisterPresetFilter()
	return make(PresetFilter,0)
}

