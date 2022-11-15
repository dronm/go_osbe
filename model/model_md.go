package model

//

import(
	"sync"	
	"strings"
	"fmt"
	
	"osbe/fields"
)
//aggregation function
type AggFunction struct {	
	Alias string
	Expr string
}

//
type ModelMD struct {	
	Fields fields.FieldCollection
	ID string
	Relation string
	AggFunctions []*AggFunction
	LimitCount int
	LimitConstant string
	mx sync.RWMutex
	FieldList string
	FieldDefOrder *string
}
func (m *ModelMD) GetFields() fields.FieldCollection {
	return m.Fields
}

//does both: makes field list for select (comma separated list m.FieldList) and makes m.FieldDefOrder (comma separated list for ORDER BY)
func (m *ModelMD) initFieldOrder(encryptkey string) {
	if m.FieldList == "" || m.FieldDefOrder == nil {
		var l_sel []string
		var l_ord []string
		m.mx.Lock()
		if m.FieldList == "" {
			l_sel = make([]string, len(m.Fields))
		}
		if m.FieldDefOrder == nil {
			l_ord = make([]string, len(m.Fields))
		}
		for _, fld := range m.Fields {
			fld_id := fld.GetId()
			if l_sel != nil {
				if encryptkey != "" && fld.GetEncrypted() {				
					l_sel[fld.GetOrderInList()] = fmt.Sprintf(`PGP_SYM_DECRYPT(%s, "%s") AS %s`, fld_id, encryptkey, fld_id)
				}else{
					l_sel[fld.GetOrderInList()] = fld_id
				}
			}
			if l_ord != nil {
				if  o := fld.GetDefOrder(); o.IsSet {
					ord := ""
					if o.Value {
						ord = "ASC"
					}else{
						ord = "DESC"
					}
					l_ord[fld.GetDefOrderIndex()] = fld_id +" " + ord
				}
			}
		}
		if l_sel != nil {
			m.FieldList = strings.Join(l_sel, ",")
		}
		if l_ord != nil {
			_s := ""
			m.FieldDefOrder = &_s //initialize
			for _, o := range l_ord {
				if o != "" {
					if _s != "" {
						_s+= ","
					}
					_s+= o
				}
			}
		}		
		m.mx.Unlock()		
	}
}

//fields as comma separated list for sql used in select query
func (m *ModelMD) GetFieldList(encryptkey string) string {
	/*
	if m.FieldList == "" {
		m.mx.Lock()
		l := make([]string, len(m.Fields))
		for _, fld := range m.Fields {
			fld_id := fld.GetId()
			if encryptkey != "" && fld.GetEncrypted() {				
				l[fld.GetOrderInList()] = fmt.Sprintf(`PGP_SYM_DECRYPT(%s, "%s") AS %s`, fld_id, encryptkey, fld_id)
			}else{
				l[fld.GetOrderInList()] = fld_id
			}
		}
		m.FieldList = strings.Join(l, ",")
		m.mx.Unlock()		
	}
	*/
	m.initFieldOrder(encryptkey)
	return m.FieldList
}

//ORDER BY FIELD1 DIR1, FIELD2 DIR2, ...
func (m *ModelMD) GetFieldDefOrder(encryptkey string) string {
	/*
	if m.FieldDefOrder == nil {
		*m.FieldDefOrder = "" //initialize
		m.mx.Lock()
		l := make([]string, len(m.Fields))
		for _, fld := range m.Fields {
			if  o := fld.GetDefOrder(); o.IsSet {
				ord := ""
				if o.GetValue() {
					ord = "ASC"
				}else{
					ord = "DESC"
				}
				l[fld.GetDefOrderIndex()] = fld.GetId() +" " + ord
			}
		}
		for i, o := range l {
			if o != "" {
				if *m.FieldDefOrder != "" {
					*m.FieldDefOrder+= ","
				}
				*m.FieldDefOrder+= o
			}
		}
		m.mx.Unlock()		
	}
	*/
	m.initFieldOrder(encryptkey)
	return *m.FieldDefOrder
}

