package about

/**
 * Andrey Mikhalevich 15/12/21
 * This file is part of the OSBE framework
 */

import (
	"reflect"
	
		
	"osbe/fields"
	"osbe/model"
)

//const About_DATA_TABLE = ""

//Exported model metadata
var (
	About_md fields.FieldCollection	
)

func About_Model_init() {	
	About_md = fields.GenModelMD(reflect.ValueOf(About{}))
}

//
type About struct {
	Author fields.ValText `json:"author"`
	Tech_mail fields.ValText `json:"tech_mail"`
	App_name fields.ValText `json:"app_name"`
	Fw_version fields.ValText `json:"fw_version"`
	App_version fields.ValText `json:"app_version"`
	Db_name fields.ValText `json:"db_name"`
}
/*func (m *About) GetDataTable() string{
	return About_DATA_TABLE
}*/
func (m *About) GetID() model.ModelID{
	return model.ModelID("About_Model")
}
