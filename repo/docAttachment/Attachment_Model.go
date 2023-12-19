package docAttachment

/**
 * Andrey Mikhalevich 15/12/21
 * This file is part of the OSBE framework
 *
 * THIS FILE IS GENERATED FROM TEMPLATE build/templates/models/Model.go.tmpl
 * ALL DIRECT MODIFICATIONS WILL BE LOST WITH THE NEXT BUILD PROCESS!!!
 */

import (
	"reflect"	
		
	"osbe/fields"
	"osbe/model"
)

type Attachment struct {
	Id fields.ValInt `json:"id" primaryKey:"true" autoInc:"true"`
	Date_time fields.ValDateTimeTZ `json:"date_time"`
	Ref fields.ValJSON `json:"ref"`
	Content_info fields.ValJSON `json:"content_info"`
	Content_data fields.ValBytea `json:"content_data"`
	Content_preview fields.ValBytea `json:"content_preview"`
}

func (o *Attachment) SetNull() {
	o.Id.SetNull()
	o.Date_time.SetNull()
	o.Ref.SetNull()
	o.Content_info.SetNull()
	o.Content_data.SetNull()
	o.Content_preview.SetNull()
}

func NewModelMD_Attachment() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(Attachment{})),
		ID: "Attachment_Model",
		Relation: "attachments",
	}
}
//for insert
type Attachment_argv struct {
	Argv *Attachment `json:"argv"`	
}

//Keys for delete/get object
type Attachment_keys struct {
	Id fields.ValInt `json:"id"`
	Mode string `json:"mode" openMode:"true"` //open mode insert|copy|edit
}
type Attachment_keys_argv struct {
	Argv *Attachment_keys `json:"argv"`	
}

//old keys for update
type Attachment_old_keys struct {
	Old_id fields.ValInt `json:"old_id"`
	Id fields.ValInt `json:"id"`
	Date_time fields.ValDateTimeTZ `json:"date_time"`
	Ref fields.ValJSON `json:"ref"`
	Content_info fields.ValJSON `json:"content_info"`
	Content_data fields.ValBytea `json:"content_data"`
	Content_preview fields.ValBytea `json:"content_preview"`
}

type Attachment_old_keys_argv struct {
	Argv *Attachment_old_keys `json:"argv"`	
}

