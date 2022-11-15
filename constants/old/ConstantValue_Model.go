package constants

/**
 * Andrey Mikhalevich 16/12/21
 * This file is part of the OSBE framework
 */

import (
	"osbe/fields"
	"osbe/model"
)

//
type ConstantValue struct {
	Id fields.ValText `json:"id" primaryKey:"true"`
	Val fields.ValText `json:"val"`
	Val_type fields.ValText `json:"val_type"`
}
func (m *ConstantValue) GetID() model.ModelID{
	return model.ModelID("ConstantValueList_Model")
}

