package contact

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

type Contact_upsert struct {
	Name fields.ValText `json:"name" required:"true"`
	Tel fields.ValText `json:"tel" required:"true"`
	Email fields.ValText `json:"email"`
	Tel_ext fields.ValText `json:"tel_ext"`
}

type Contact_upsert_argv struct {
	Argv *Contact_upsert `json:"argv"`	
}

