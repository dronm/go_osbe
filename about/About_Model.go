package about

/**
 * Andrey Mikhalevich 15/12/21
 * This file is part of the OSBE framework
 */

import (
	"osbe/fields"
)

//
type About struct {
	Author fields.ValText `json:"author"`
	Tech_mail fields.ValText `json:"tech_mail"`
	App_name fields.ValText `json:"app_name"`
	Fw_version fields.ValText `json:"fw_version"`
	App_version fields.ValText `json:"app_version"`
	Db_name fields.ValText `json:"db_name"`
}

