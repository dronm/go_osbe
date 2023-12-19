package bank

import (
	"reflect"	
		
	"osbe/fields"
	"osbe/model"
)

type BankList struct {
	Bik fields.ValText `json:"bik" primaryKey:"true" alias:"БИК" defOrder:"ASC"`
	Codegr fields.ValText `json:"codegr" defOrder:"ASC"`
	Gr_descr fields.ValText `json:"gr_descr" alias:"Регион"`
	Name fields.ValText `json:"name" alias:"Наименование"`
	Korshet fields.ValText `json:"korshet" alias:"Кoр.счет"`
	Adres fields.ValText `json:"adres" alias:"Адрес"`
	Gor fields.ValText `json:"gor" alias:"Город"`
	Tgroup fields.ValBool `json:"tgroup"`
}

func (o *BankList) SetNull() {
	o.Bik.SetNull()
	o.Codegr.SetNull()
	o.Gr_descr.SetNull()
	o.Name.SetNull()
	o.Korshet.SetNull()
	o.Adres.SetNull()
	o.Gor.SetNull()
	o.Tgroup.SetNull()
}

func NewModelMD_BankList() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(BankList{})),
		ID: "BankList_Model",
		Relation: "banks.banks_list",
		AggFunctions: []*model.AggFunction{
			&model.AggFunction{Alias: "totalCount", Expr: "count(*)"},
		},
	}
}
