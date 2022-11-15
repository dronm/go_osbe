package model

import (
	"osbe/fields"
)

//Complete model
type Complete_Model struct {	
	//Pattern fields.ValText `json:"pattern" length:500`
	Count fields.ValInt `json:"count" default:10`	
	Ic fields.ValInt `json:"ic" default:1 minValue:0 maxValue:1`
	Mid fields.ValInt `json:"mid" default:1 minValue:0 maxValue:1`
	Ord_directs fields.ValText `json:"ord_directs" length:500`
	Field_sep fields.ValText `json:"field_sep" length:2`
}
