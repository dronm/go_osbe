package bank

import (
	"reflect"	
		
	"osbe/fields"
	"osbe/model"
)

type Bank struct {
	Bik fields.ValText `json:"bik" primaryKey:"true" defOrder:"ASC"`
	Codegr fields.ValText `json:"codegr" defOrder:"ASC"`
	Name fields.ValText `json:"name"`
	Korshet fields.ValText `json:"korshet"`
	Adres fields.ValText `json:"adres"`
	Gor fields.ValText `json:"gor"`
	Tgroup fields.ValBool `json:"tgroup"`
}

func (o *Bank) SetNull() {
	o.Bik.SetNull()
	o.Codegr.SetNull()
	o.Name.SetNull()
	o.Korshet.SetNull()
	o.Adres.SetNull()
	o.Gor.SetNull()
	o.Tgroup.SetNull()
}

func NewModelMD_Bank() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(Bank{})),
		ID: "Bank_Model",
		Relation: "banks.banks",
		AggFunctions: []*model.AggFunction{
			&model.AggFunction{Alias: "totalCount", Expr: "count(*)"},
		},
	}
}
//for insert
type Bank_argv struct {
	Argv *Bank `json:"argv"`	
}

//Keys for delete/get object
type Bank_keys struct {
	Bik fields.ValText `json:"bik"`
	Mode string `json:"mode" openMode:"true"` //open mode insert|copy|edit
}
type Bank_keys_argv struct {
	Argv *Bank_keys `json:"argv"`	
}

//old keys for update
type Bank_old_keys struct {
	Old_bik fields.ValText `json:"old_bik"`
	Bik fields.ValText `json:"bik"`
	Codegr fields.ValText `json:"codegr"`
	Name fields.ValText `json:"name"`
	Korshet fields.ValText `json:"korshet"`
	Adres fields.ValText `json:"adres"`
	Gor fields.ValText `json:"gor"`
	Tgroup fields.ValBool `json:"tgroup"`
}

type Bank_old_keys_argv struct {
	Argv *Bank_old_keys `json:"argv"`	
}

