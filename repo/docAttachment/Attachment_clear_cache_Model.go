package docAttachment

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

type Attachment_clear_cache struct {
	Ref Ref_Type `json:"ref" required:"true"`
	Content_id fields.ValText `json:"content_id" required:"true"`
}

type Attachment_clear_cache_argv struct {
	Argv *Attachment_clear_cache `json:"argv"`	
}

