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
	LimitCount int //max count that can be served at one go, see sql_limit.go for details
	LimitConstant string //deprecated, not used
	DocPerPageCount int //default document per page
	mx sync.RWMutex
	FieldList string
	CopyFieldList string
	FieldDefOrder *string
}
func (m *ModelMD) GetFields() fields.FieldCollection {
	return m.Fields
}

//does both: makes field list for select (comma separated list m.FieldList) and makes m.FieldDefOrder (comma separated list for ORDER BY)
func (m *ModelMD) initFieldOrder(encryptkey string, copyMode bool) {
	if (m.FieldList == "" && !copyMode) || (m.CopyFieldList == "" && copyMode) || m.FieldDefOrder == nil {
		var l_sel []string
		var l_ord []string
		m.mx.Lock()
		if (m.FieldList == "" && !copyMode) || (m.CopyFieldList == "" && copyMode) {
			l_sel = make([]string, len(m.Fields))
		}
		if m.FieldDefOrder == nil {
			l_ord = make([]string, len(m.Fields))
		}
		for _, fld := range m.Fields {
			fld_id := fld.GetId()
			if l_sel != nil {
				if copyMode && fld.GetNoValueOnCopy() {
					l_sel[fld.GetOrderInList()] = "NULL AS " + fld_id //do not copy value!
				
				}else if encryptkey != "" && fld.GetEncrypted() {				
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
		if l_sel != nil && !copyMode {
			m.FieldList = strings.Join(l_sel, ",")
			
		}else if l_sel != nil && copyMode {
			m.CopyFieldList = strings.Join(l_sel, ",")
		}
		
		if l_ord != nil {
			_s := ""
			m.FieldDefOrder = &_s //initialize
			//var sb strings.Builder
			//sb_first := true
			for _, o := range l_ord {
				if o != "" {
					if _s != "" {
						_s+= ","
					}
					/*if !sb_first {
						sb.WriteString(",")
					}else{
						sb_first = true
					}*/
					_s+= o
					//sb.WriteString(o)
				}
			}
			//_s = sb.String()
		}		
		m.mx.Unlock()		
	}
}

//fields as comma separated list for sql used in select query
func (m *ModelMD) GetFieldList(encryptkey string) string {
	m.initFieldOrder(encryptkey, false)
	return m.FieldList
}

func (m *ModelMD) GetCopyFieldList(encryptkey string) string {
	m.initFieldOrder(encryptkey, true)
	return m.CopyFieldList
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
	m.initFieldOrder(encryptkey, false)
	return *m.FieldDefOrder
}

