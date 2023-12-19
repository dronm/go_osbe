package userOperation

import (
	"reflect"	
		
	"osbe/fields"
	"osbe/model"
)

type UserOperationDialog struct {
	User_id fields.ValInt `json:"user_id" required:"true" primaryKey:"true"`
	Operation_id fields.ValText `json:"operation_id" required:"true" primaryKey:"true"`
	Operation fields.ValText `json:"operation"`
	Status fields.ValText `json:"status"`
	Date_time fields.ValDateTimeTZ `json:"date_time"`
	Error_text fields.ValText `json:"error_text"`
	Comment_text fields.ValText `json:"comment_text"`
	Date_time_end fields.ValDateTimeTZ `json:"date_time_end"`
	End_wal_lsn fields.ValText `json:"end_wal_lsn"`
}

func (o *UserOperationDialog) SetNull() {
	o.User_id.SetNull()
	o.Operation_id.SetNull()
	o.Operation.SetNull()
	o.Status.SetNull()
	o.Date_time.SetNull()
	o.Error_text.SetNull()
	o.Comment_text.SetNull()
	o.Date_time_end.SetNull()
	o.End_wal_lsn.SetNull()
}

func NewModelMD_UserOperationDialog() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(UserOperationDialog{})),
		ID: "UserOperationDialog_Model",
		Relation: "user_operations_dialog",
		AggFunctions: []*model.AggFunction{
			&model.AggFunction{Alias: "totalCount", Expr: "count(*)"},
		},
	}
}
