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

type Attachment_delete_file struct {
	Ref Ref_Type `json:"ref" required:"true"`
	Content_id fields.ValText `json:"content_id" required:"true"`
}
type Attachment_delete_file_argv struct {
	Argv *Attachment_delete_file `json:"argv"`	
}

