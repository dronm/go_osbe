package notif

import (
	"reflect"	
	"osbe/fields"
	"osbe/model"
)

type NotifTemplate struct {
	Id fields.ValInt `json:"id" primaryKey:"true" autoInc:"true"`
	Notif_provider ValEnum_notif_providers `json:"notif_provider" required:"true" alias:"Способ доставки" defOrder:"ASC"`
	Notif_type fields.ValText `json:"notif_type" required:"true" alias:"Тип сообщения" defOrder:"ASC"`
	Template fields.ValText `json:"template" required:"true" alias:"Шаблон"`
	Comment_text fields.ValText `json:"comment_text" required:"true" alias:"Комментарий"`
	Fields fields.ValJSON `json:"fields" required:"true" alias:"Поля"`
	Provider_values fields.ValJSON `json:"provider_values" alias:"Специфичные для провайдера поля, например subject для писем, фичи оформления для месенджеров"`
}

func (o *NotifTemplate) SetNull() {
	o.Id.SetNull()
	o.Notif_provider.SetNull()
	o.Notif_type.SetNull()
	o.Template.SetNull()
	o.Comment_text.SetNull()
	o.Fields.SetNull()
	o.Provider_values.SetNull()
}

func NewModelMD_NotifTemplate() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(NotifTemplate{})),
		ID: "NotifTemplate_Model",
		Relation: "notif_templates",
		AggFunctions: []*model.AggFunction{
			&model.AggFunction{Alias: "totalCount", Expr: "count(*)"},
		},
	}
}
//for insert
type NotifTemplate_argv struct {
	Argv *NotifTemplate `json:"argv"`	
}

//Keys for delete/get object
type NotifTemplate_keys struct {
	Id fields.ValInt `json:"id"`
	Mode string `json:"mode" openMode:"true"` //open mode insert|copy|edit
}
type NotifTemplate_keys_argv struct {
	Argv *NotifTemplate_keys `json:"argv"`	
}

//old keys for update
type NotifTemplate_old_keys struct {
	Old_id fields.ValInt `json:"old_id"`
	Id fields.ValInt `json:"id"`
	Notif_provider ValEnum_notif_providers `json:"notif_provider" alias:"Способ доставки"`
	Notif_type fields.ValText `json:"notif_type" alias:"Тип сообщения"`
	Template fields.ValText `json:"template" alias:"Шаблон"`
	Comment_text fields.ValText `json:"comment_text" alias:"Комментарий"`
	Fields fields.ValJSON `json:"fields" alias:"Поля"`
	Provider_values fields.ValJSON `json:"provider_values" alias:"Специфичные для провайдера поля, например subject для писем, фичи оформления для месенджеров"`
}

type NotifTemplate_old_keys_argv struct {
	Argv *NotifTemplate_old_keys `json:"argv"`	
}

