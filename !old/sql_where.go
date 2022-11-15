package osbe

import (
	"reflect"
	"errors"
	"strings"
	"fmt"
	
	"osbe/fields"
	"osbe/model"
	"osbe/sql"
)

const (
	SGN_PAR_E = "e"			//equal
	SGN_PAR_L = "l"			//less
	SGN_PAR_G = "g"			//greater
	SGN_PAR_LE = "le"		//less and equal
	SGN_PAR_GE = "ge"		//greater and equal
	SGN_PAR_LK = "lk"		//like
	SGN_PAR_NE = "ne"		//not equal
	SGN_PAR_I = "i"			// IS
	SGN_PAR_IN = "in"		// in
	SGN_PAR_INCL = "incl"		//include
	SGN_PAR_ANY = "any"		//Any
	SGN_PAR_OVERLAP = "overlap"	//overlap
	
)

type conditionJoin int
const (
	CONDITION_JOIN_AND conditionJoin = iota
	CONDITION_JOIN_OR
)

type argConditions struct {
	Fields []string
	Signs []sql.SQLCondition
	Vals []interface{}
	InsCases []bool
	Joins []conditionJoin
}
func (c *argConditions) sql(ind int) string {
	switch c.Joins[ind] {
	case CONDITION_JOIN_AND:
		return "AND"
	case CONDITION_JOIN_OR:
		return "OR"
	}
	return "UNDEFIND_JOIN"
}

//parses reflect.Value, extracts data from cond_fields, cond_sgns, cond_ic, cond_vals
//returns
//cond_fields - slice of string
//cond_sgns - slice of sql.SQLCondition
//cond_vals - slice of interface{}
//cond_ic - slice of []bool
func parseSQLWhereFromArgs(rfltArgs reflect.Value, fieldSep string, modelMetadata fields.FieldCollection) ([]string, []sql.SQLCondition, []interface{}, []bool, error) {		
	if ids := GetTextArgValByName(rfltArgs, "Cond_fields", ""); ids != "" {		
		//fields
		fields_s := strings.Split(ids, fieldSep) //fld_t.GetValue()
		f_cnt := len(fields_s)
		if f_cnt == 0 {
			return nil, nil, nil, nil, nil
		}
		
		//signs
		var sgns_s []sql.SQLCondition
		if sgns := GetTextArgValByName(rfltArgs, "Cond_sgns", ""); sgns != "" {		
			sgns_str := strings.Split(sgns, fieldSep)
			sgns_s = make([]sql.SQLCondition, len(sgns_str))
			for ind, sgn := range sgns_str {
				switch sgn {
					case SGN_PAR_E:
						sgns_s[ind] = sql.SGN_SQL_E
					case SGN_PAR_L:
						sgns_s[ind] = sql.SGN_SQL_L
					case SGN_PAR_G:
						sgns_s[ind] = sql.SGN_SQL_G
					case SGN_PAR_LE:
						sgns_s[ind] = sql.SGN_SQL_LE
					case SGN_PAR_GE:
						sgns_s[ind] = sql.SGN_SQL_GE
					case SGN_PAR_LK:
						sgns_s[ind] = sql.SGN_SQL_LK
					case SGN_PAR_NE:
						sgns_s[ind] = sql.SGN_SQL_NE
					case SGN_PAR_I:
						sgns_s[ind] = sql.SGN_SQL_I
					case SGN_PAR_IN:
						sgns_s[ind] = sql.SGN_SQL_IN
					case SGN_PAR_INCL:
						sgns_s[ind] = sql.SGN_SQL_INCL
					case SGN_PAR_ANY:
						sgns_s[ind] = sql.SGN_SQL_ANY
					case SGN_PAR_OVERLAP:
						sgns_s[ind] = sql.SGN_SQL_OVERLAP							
					default:
						return nil, nil, nil, nil, errors.New(fmt.Sprintf(ER_SQL_WHERE_UNKNOWN_COND, sgn))
				}
			}
		}
		if f_cnt > len(sgns_s) {
			//field count mismatch
			return nil, nil, nil, nil, errors.New("1 "+ER_SQL_WHERE_FILED_CNT_MISMATCH)
		}
				
		//ics
		var ic_s []bool
		if ics := GetTextArgValByName(rfltArgs, "Cond_ic", ""); ics != "" {		
			ics_str := strings.Split(ics, fieldSep)
			ic_s = make([]bool, len(ics_str))
			for i,ic := range ics_str {
				ic_s[i],_ = fields.StrToBool(ic)					
			}		
		}
		if ic_s == nil {
			ic_s = make([]bool, f_cnt)
		}
		if f_cnt > len(ic_s) {
			//missing some ics					
			for i:= len(ic_s); i<f_cnt; i++ {
				ic_s = append(ic_s, false)
			}
		}
		//values
		var vals_s []interface{}
		if vals := GetTextArgValByName(rfltArgs, "Cond_vals", ""); vals != "" {		
			vals_str := strings.Split(vals, fieldSep)
			if f_cnt < len(vals_str) {
				//field count mismatch
				return nil, nil, nil, nil, errors.New("2 "+ER_SQL_WHERE_FILED_CNT_MISMATCH)
			}
			//cast string value to real field type value
			valid_err := ""
			var md_field_ids map[string]fields.Fielder //case insensetive field ids
			for ind, val_str := range vals_str {
				if len(val_str) == 0 {
					appendError(&valid_err, "field value not set")
					continue
				}
				/*
				if fields_s[ind] == "custom" {
					//Controller defind condition - skeep here
					fields_s = append(fields_s[:ind], fields_s[ind+1:]...)
					sgns_s = append(sgns_s[:ind], sgns_s[ind+1:]...)
					ic_s = append(ic_s[:ind], ic_s[ind+1:]...)
					f_cnt--
					continue
				}*/
				
				//in most cases first letter is capitalized
				id := strings.ToUpper(string(fields_s[ind][:1])) + string(fields_s[ind][1:])										
				model_f, ok := modelMetadata[id]
				if !ok {
					//case insensetive check!!!	
					if md_field_ids == nil {
						md_field_ids = make(map[string]fields.Fielder, len(modelMetadata))
						for _, m_f := range modelMetadata {
							m_f_id := m_f.GetId()
							if !ok && m_f_id == fields_s[ind] &&
							len(fields_s) == 1 {
								model_f = m_f
								ok = true
								break
								
							}else if !ok && m_f_id == fields_s[ind] {
								model_f = m_f
								ok = true
							}
							md_field_ids[m_f_id] = m_f
						}
					}
					if !ok {
						if model_f_lc, ok_lc := md_field_ids[fields_s[ind]]; ok_lc {
							model_f = model_f_lc
							ok = true
						}
					}						
				}
				if ok {						
					var err error
					var val_i interface{}
					
					//might be wild char signs % -at the begining and at the end of the val_str!!!
					if val_str[0:1]=="%" || val_str[len(val_str)-1:] == "%" {
						//treat as string
						//@ToDo validate for injections!
						val_i = val_str						
					}else{
						switch model_f.GetDataType() {
						case fields.FIELD_TYPE_FLOAT:
							//str to float64
							var tp_v float64
							tp_v, err = fields.StrToFloat(val_str)
							if err == nil {
								err = fields.ValidateFloat(model_f.(fields.FielderFloat), tp_v)
								if err == nil {
									val_i = tp_v
								}
							}
							if ic_s[ind] {
								ic_s[ind] = false
							}
						case fields.FIELD_TYPE_INT:
							var tp_v int64
							tp_v, err = fields.StrToInt(val_str)
							if err == nil {
								err = fields.ValidateInt(model_f.(fields.FielderInt), tp_v)				
								if err == nil {
									val_i = tp_v
								}
							}
							if ic_s[ind] {
								ic_s[ind] = false
							}
							
						case fields.FIELD_TYPE_BOOL:
							tp_v,_ := fields.StrToBool(val_str)
							val_i = tp_v
							if ic_s[ind] {
								ic_s[ind] = false
							}
														
						case fields.FIELD_TYPE_TEXT:
							err = fields.ValidateText(model_f.(fields.FielderText), val_str)
							if err == nil {
								val_i = val_str
							}
						case fields.FIELD_TYPE_DATE:
							err = fields.ValidateDate(model_f.(fields.Fielder), val_str)
							if err == nil {
								val_i = val_str
							}
						case fields.FIELD_TYPE_DATETIME:
							err = fields.ValidateDateTime(model_f.(fields.Fielder), val_str)
							if err == nil {
								val_i = val_str
							}
						case fields.FIELD_TYPE_DATETIMETZ:
							err = fields.ValidateDateTimeTZ(model_f.(fields.Fielder), val_str)
							if err == nil {
								val_i = val_str
							}

							
						default:
							err = errors.New(fmt.Sprintf("'%s' unsupported condition field type",fields_s[ind])) 
						}
					}
					if err != nil {
						appendError(&valid_err, err.Error() ) 
					}else{
						vals_s = append(vals_s, val_i)
					}						
				}else{
					return nil, nil, nil, nil, errors.New(fmt.Sprintf("parseSQLWhereFromArgs(): field %s not found in model", id))
				}	
			}
			if valid_err != "" {
				return nil, nil, nil, nil, errors.New(valid_err)
			}
		}
//			}
//		}
//fmt.Println("vals_s=",vals_s, "Len=", len(vals_s))
//fmt.Println("f_cnt=",f_cnt)
//fmt.Println("fields_s=",fields_s)
//can be nil if cystom is set
		if vals_s == nil || f_cnt > len(vals_s) {
			return nil, nil, nil, nil, errors.New("3 "+ER_SQL_WHERE_FILED_CNT_MISMATCH)
		}

		return fields_s, sgns_s, vals_s, ic_s, nil
	}
	return nil, nil, nil, nil,nil
} 

//returns:
//	sql_s query string
//	vals_s slice of validated, sanatized parameters
//	error
func GetSQLWhereFromArgs(rfltArgs reflect.Value, fieldSep string, modelMD *model.ModelMD, extraConds sql.FilterCondCollection) (string, []interface{}, error) {
	//fields.FieldCollection
	fields_s, sgns_s, vals_s, ic_s, err := parseSQLWhereFromArgs(rfltArgs, fieldSep, modelMD.GetFields())
	if err != nil {
		return "", nil, err
	}
	
	if (fields_s == nil || len(fields_s) == 0) && (extraConds == nil || len(extraConds) == 0) {
		return "", nil, nil
	}
	sql_s := "WHERE "
	cond_cnt := 0
	if fields_s != nil {
		for i, fld := range fields_s {						
			sql.AddCondExpr(fld, sgns_s[i], ic_s[i], i, "AND", &sql_s)
			cond_cnt++
		}
	}
	if extraConds != nil && len(extraConds) > 0 {
		expr_conds := "" //pure expressions
		for _, cond := range extraConds {
			if cond.Expression != "" {
				if expr_conds != "" {
					expr_conds += " AND "
				}
				expr_conds += cond.Expression
				
			}else if cond.FieldID != "" {
				sgn := cond.Sign
				if cond.Sign == "" {
					sgn = sql.SGN_SQL_E
				}
				sql.AddCondExpr(cond.FieldID, sgn, cond.InsCase, cond_cnt, "AND", &sql_s)
				if vals_s == nil {
					vals_s = make([]interface{},0)
				}
				vals_s = append(vals_s, cond.Value)
				cond_cnt++
			}
		}
		if expr_conds != "" {
			if cond_cnt > 0 {
				sql_s += " AND "
			}
			sql_s += expr_conds
		}
	}
	
	return sql_s, vals_s, nil
} 


