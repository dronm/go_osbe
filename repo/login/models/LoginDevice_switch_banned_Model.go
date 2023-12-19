package models

/**
 * Andrey Mikhalevich 16/12/21
 * This file is part of the OSBE framework
 *
 * THIS FILE IS GENERATED FROM TEMPLATE build/templates/models/Model.go.tmpl
 * ALL DIRECT MODIFICATIONS WILL BE LOST WITH THE NEXT BUILD PROCESS!!!
 */

//Controller method model
import (
		
	"osbe/fields"
)

type LoginDevice_switch_banned struct {
	Banned fields.ValBool `json:"banned" required:"true"`
	Hash fields.ValText `json:"hash" required:"true"`
	User_id fields.ValInt `json:"user_id" required:"true"`
}
type LoginDevice_switch_banned_argv struct {
	Argv *LoginDevice_switch_banned `json:"argv"`	
}

