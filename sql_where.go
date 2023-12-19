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
	
	JOIN_PAR_AND = "a"
	JOIN_PAR_OR = "o"
)

type ConditionJoin int
func (c ConditionJoin) Sql() string {
	switch c {
	case CONDITION_JOIN_AND:
		return "AND"
	case CONDITION_JOIN_OR:
		return "OR"
	}
	return "UNDEFIND_JOIN"
}

const (
	CONDITION_JOIN_AND ConditionJoin = iota
	CONDITION_JOIN_OR
)

type ArgConditions struct {
	Fields []string
	Signs []sql.SQLCondition
	Values []interface{}
	InsCases []bool
	Joins []ConditionJoin
}

//parses reflect.Value, extracts data from cond_fields, cond_sgns, cond_ic, cond_vals, cond_joins
//returns
//cond_fields - slice of string
//cond_sgns - slice of sql.SQLCondition
//cond_vals - slice of interface{}
//cond_ic - slice of []bool
func ParseSQLWhereFromArgs(rfltArgs reflect.Value, fieldSep string, modelMetadata fields.FieldCollection) (*ArgConditions, error) {
	if ids := GetTextArgValByName(rfltArgs, "Cond_fields", ""); ids != "" {
		arg_conds := ArgConditions{}
		//fields
		arg_conds.Fields = strings.Split(ids, fieldSep) //fld_t.GetValue()
		f_cnt := len(arg_conds.Fields)
		if f_cnt == 0 {
			return nil, nil
		}
		
		//signs		
		if sgns := GetTextArgValByName(rfltArgs, "Cond_sgns", ""); sgns != "" {					
			sgns_str := strings.Split(sgns, fieldSep)
			if f_cnt != len(sgns_str) {
				//field count mismatch
				return nil, errors.New("1 "+ER_SQL_WHERE_FILED_CNT_MISMATCH)
			}
			arg_conds.Signs = make([]sql.SQLCondition, f_cnt)
			for ind, sgn := range sgns_str {
				switch sgn {
					case SGN_PAR_E:
						arg_conds.Signs[ind] = sql.SGN_SQL_E
					case SGN_PAR_L:
						arg_conds.Signs[ind] = sql.SGN_SQL_L
					case SGN_PAR_G:
						arg_conds.Signs[ind] = sql.SGN_SQL_G
					case SGN_PAR_LE:
						arg_conds.Signs[ind] = sql.SGN_SQL_LE
					case SGN_PAR_GE:
						arg_conds.Signs[ind] = sql.SGN_SQL_GE
					case SGN_PAR_LK:
						arg_conds.Signs[ind] = sql.SGN_SQL_LK
					case SGN_PAR_NE:
						arg_conds.Signs[ind] = sql.SGN_SQL_NE
					case SGN_PAR_I:
						arg_conds.Signs[ind] = sql.SGN_SQL_I
					case SGN_PAR_IN:
						arg_conds.Signs[ind] = sql.SGN_SQL_IN
					case SGN_PAR_INCL:
						arg_conds.Signs[ind] = sql.SGN_SQL_INCL
					case SGN_PAR_ANY:
						arg_conds.Signs[ind] = sql.SGN_SQL_ANY
					case SGN_PAR_OVERLAP:
						arg_conds.Signs[ind] = sql.SGN_SQL_OVERLAP							
					default:
						return nil, errors.New(fmt.Sprintf(ER_SQL_WHERE_UNKNOWN_COND, sgn))
				}
			}
		}
				
		//ics
		arg_conds.InsCases = make([]bool, f_cnt) //defaults false
		if ics := GetTextArgValByName(rfltArgs, "Cond_ic", ""); ics != "" {		
			ics_str := strings.Split(ics, fieldSep)			
			for i, ic := range ics_str {
				if i == f_cnt {
					break
				}
				arg_conds.InsCases[i], _ = fields.StrToBool(ic)					
			}		
		}

		//joins
		arg_conds.Joins = make([]ConditionJoin, f_cnt) //defaults AND
		if joins := GetTextArgValByName(rfltArgs, "Cond_joins", ""); joins != "" {		
			join_str := strings.Split(joins, fieldSep)			
			for i := 0; i < f_cnt; i++ {
				arg_conds.Joins[i] = CONDITION_JOIN_AND
				if i < len(join_str) && join_str[i] == JOIN_PAR_OR {
					arg_conds.Joins[i] = CONDITION_JOIN_OR
				}
			}		
		}
		//values				
		vals := GetTextArgValByName(rfltArgs, "Cond_vals", "")
		if vals == "" && f_cnt == 1 {
			//one field with null value
			vals = "null"
		}
		if vals != "" {		
			vals_str := strings.Split(vals, fieldSep)
			if f_cnt != len(vals_str) {
				//field count mismatch
				return nil, errors.New("2 "+ER_SQL_WHERE_FILED_CNT_MISMATCH)
			}
			arg_conds.Values = make([]interface{}, f_cnt)
			//cast string value to real field type value
			var valid_err strings.Builder
			var md_field_ids map[string]fields.Fielder //case insensetive field ids
			for ind, val_str := range vals_str {
				if len(val_str) == 0 {
					appendError(&valid_err, "field value not set")
					continue
				}
				
				//in most cases first letter is capitalized
				id := strings.ToUpper(string(arg_conds.Fields[ind][:1])) + string(arg_conds.Fields[ind][1:])										
				model_f, ok := modelMetadata[id]
				if !ok {
					//case insensetive check!!!	
					if md_field_ids == nil {
						md_field_ids = make(map[string]fields.Fielder, len(modelMetadata))
						for _, m_f := range modelMetadata {
							m_f_id := m_f.GetId()
							if !ok && m_f_id == arg_conds.Fields[ind] &&
							len(arg_conds.Fields) == 1 {
								model_f = m_f
								ok = true
								break
								
							}else if !ok && m_f_id == arg_conds.Fields[ind] {
								model_f = m_f
								ok = true
							}
							md_field_ids[m_f_id] = m_f
						}
					}
					if !ok {
						if model_f_lc, ok_lc := md_field_ids[arg_conds.Fields[ind]]; ok_lc {
							model_f = model_f_lc
							ok = true
						}
					}						
				}
				if ok {						
					var err error
					var val_i interface{}
					
					if strings.ToUpper(val_str) == "NULL" {
						val_i = nil //for all types
						
					//might be wild char signs % -at the begining and at the end of the val_str!!!
					}else if val_str[0:1]=="%" || val_str[len(val_str)-1:] == "%" {
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
							if arg_conds.InsCases[ind] {
								arg_conds.InsCases[ind] = false
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
							if arg_conds.InsCases[ind] {
								arg_conds.InsCases[ind] = false
							}
							
						case fields.FIELD_TYPE_BOOL:
							tp_v,_ := fields.StrToBool(val_str)
							val_i = tp_v
							if arg_conds.InsCases[ind] {
								arg_conds.InsCases[ind] = false
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
							if len(val_str) == 10 {
								//cast to date
								arg_conds.Fields[ind] += "::date"
							}
							
						case fields.FIELD_TYPE_DATETIMETZ:
							err = fields.ValidateDateTimeTZ(model_f.(fields.Fielder), val_str)
							if err == nil {
								val_i = val_str
							}
							if len(val_str) == 10 {
								//cast to date
								arg_conds.Fields[ind] += "::date"
							}
							
						default:
							err = errors.New(fmt.Sprintf("'%s' unsupported condition field type",arg_conds.Fields[ind])) 
						}
					}
					if err != nil {
						appendError(&valid_err, err.Error() ) 
					}else{
						arg_conds.Values[ind] = val_i
					}						
				}else{
					return nil, errors.New(fmt.Sprintf("parseSQLWhereFromArgs(): field %s not found in model", id))
				}	
			}
			if valid_err.Len() > 0 {
				return nil, errors.New(valid_err.String())
			}
		}
		
//fmt.Println("vals_s=",arg_conds.Values, "Len=", len(arg_conds.Values))
//fmt.Println("f_cnt=",f_cnt)
//fmt.Println("fields_s=",fields_s)
//can be nil if cystom is set
		if arg_conds.Values == nil || arg_conds.Signs == nil {
			return nil, errors.New("3 "+ER_SQL_WHERE_FILED_CNT_MISMATCH)
		}

		return &arg_conds, nil
	}
	return nil, nil
} 

//returns:
//	sql_s query string
//	vals_s slice of validated, sanatized parameters
//	error
func GetSQLWhereFromArgs(rfltArgs reflect.Value, fieldSep string, modelMD *model.ModelMD, extraConds sql.FilterCondCollection) (string, []interface{}, error) {
	arg_conds, err := ParseSQLWhereFromArgs(rfltArgs, fieldSep, modelMD.GetFields())
	if err != nil {
		return "", nil, err
	}
	if (arg_conds == nil || arg_conds.Fields == nil || len(arg_conds.Fields) == 0) && (extraConds == nil || len(extraConds) == 0) {
		return "", nil, nil
	}
	var arg_cond_values []interface{}
	if arg_conds != nil && arg_conds.Values != nil {
		arg_cond_values = arg_conds.Values
	}
	sql_s := "WHERE "
	cond_cnt := 0
	
	if arg_conds!= nil && arg_conds.Fields != nil {
		or_join_exists := false
		cond_sql := ""
		for i, fld := range arg_conds.Fields {						
			if arg_conds.Values[i] == nil  && (arg_conds.Signs[i] == sql.SGN_SQL_I || arg_conds.Signs[i] == sql.SGN_SQL_IN) {
				//null cases
				if i > 0 {
					cond_sql += " "
				}
				cond_sql += fmt.Sprintf("(%s %s NULL)", fld, arg_conds.Signs[i])
				//remove nil
				arg_cond_values[i] = arg_cond_values[len(arg_cond_values) - 1]
				arg_cond_values = arg_cond_values[:len(arg_cond_values) - 1]
//fmt.Println("len(arg_cond_values)=", len(arg_cond_values))				
			}else{
				sql.AddCondExpr(fld, arg_conds.Signs[i], arg_conds.InsCases[i], i, arg_conds.Joins[i].Sql(), &cond_sql)
				cond_cnt++
			}
			if arg_conds.Joins[i] == CONDITION_JOIN_OR && !or_join_exists {
				or_join_exists = true
			}
		}
		//OR always in paranthesis
		if cond_cnt > 1 && or_join_exists {
			cond_sql = "(" + cond_sql + ")"
		}
		sql_s+= cond_sql
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
				if arg_cond_values == nil {
					arg_cond_values = make([]interface{}, 0)
				}
				arg_cond_values = append(arg_cond_values, cond.Value)
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
	
	return sql_s, arg_cond_values, nil
} 


