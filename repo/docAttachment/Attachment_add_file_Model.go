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

type Attachment_add_file struct {
	Ref Ref_Type `json:"ref" required:"true"`
	Content_data fields.ValBytea `json:"content_data"`
	Content_info Content_info_Type `json:"content_info" required:"true"`
}
type Attachment_add_file_argv struct {
	Argv *Attachment_add_file `json:"argv"`	
}

