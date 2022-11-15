package stat

/**
 * Andrey Mikhalevich 15/12/21
 * This file is part of the OSBE framework
 *
 */

import (
	"reflect"
	
		
	"osbe/fields"
	"osbe/model"
)

//
type SrvStatistics struct {
	Name fields.ValText `json:"name"`
	Max_client_count fields.ValInt `json:"max_client_count"`
	Client_count fields.ValInt `json:"client_count"`
	Downloaded_bytes fields.ValUint `json:"downloaded_bytes"`
	Uploaded_bytes fields.ValUint `json:"uploaded_bytes"`
	Handshakes fields.ValUint `json:"handshakes"`
	Run_seconds fields.ValUint `json:"run_seconds"`
}

func NewModelMD_SrvStatistics() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(SrvStatistics{})),
		ID: "SrvStatistics_Model",
	}
}

