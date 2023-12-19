package userOperation

import (
	"reflect"	
		
	"osbe/fields"
	"osbe/model"
)

type UserOperation struct {
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

func (o *UserOperation) SetNull() {
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

func NewModelMD_UserOperation() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(UserOperation{})),
		ID: "UserOperation_Model",
		Relation: "user_operations",
		AggFunctions: []*model.AggFunction{
			&model.AggFunction{Alias: "totalCount", Expr: "count(*)"},
		},
	}
}
//for insert
type UserOperation_argv struct {
	Argv *UserOperation `json:"argv"`	
}

//Keys for delete/get object
type UserOperation_keys struct {
	User_id fields.ValInt `json:"user_id"`
	Operation_id fields.ValText `json:"operation_id"`
	Mode string `json:"mode" openMode:"true"` //open mode insert|copy|edit
}
type UserOperation_keys_argv struct {
	Argv *UserOperation_keys `json:"argv"`	
}

//old keys for update
type UserOperation_old_keys struct {
	Old_user_id fields.ValInt `json:"old_user_id" required:"true"`
	User_id fields.ValInt `json:"user_id"`
	Old_operation_id fields.ValText `json:"old_operation_id" required:"true"`
	Operation_id fields.ValText `json:"operation_id"`
	Operation fields.ValText `json:"operation"`
	Status fields.ValText `json:"status"`
	Date_time fields.ValDateTimeTZ `json:"date_time"`
	Error_text fields.ValText `json:"error_text"`
	Comment_text fields.ValText `json:"comment_text"`
	Date_time_end fields.ValDateTimeTZ `json:"date_time_end"`
	End_wal_lsn fields.ValText `json:"end_wal_lsn"`
}

type UserOperation_old_keys_argv struct {
	Argv *UserOperation_old_keys `json:"argv"`	
}

