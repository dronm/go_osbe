package fields

import (
	"reflect"
	"time"
	"strconv"
	"strings"
)

type ValEnumer interface {
	GetValues() []string
}
type ValInteger interface {
	GetValue() int64
}
type ValTexter interface {
	GetValue() string
}
type ValBytera interface {
	GetValue() []byte
}
type ValBooler interface {
	GetValue() bool
}
type ValTimer interface {
	GetValue() time.Time
}
type ValFloater interface {
	GetValue() float64
}
type ValArrayer interface {
	GetValue() []interface{}
}
type ValAssocArrayer interface {
	GetValue() map[string]interface{}
}

//function is executed on start, so we can panic
func GenModelMD(v reflect.Value) FieldCollection{	
	t := v.Type()	
	if t.Kind() != reflect.Struct {
		return nil
	}
	res := make(FieldCollection, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		//id := field.Tag.Get("json")
		id, id_exists := field.Tag.Lookup("json")
		if !id_exists && t.Field(i).Anonymous {
			//skeep embeded structures
			continue
		}else if !id_exists {
			panic("GenMDDescr field does not have json tag: "+field.Name)
		}
		
		if val, ok := v.Field(i).Interface().(ValEnumer); ok {
			res[field.Name] = &FieldEnum{Values: val.GetValues()}
		
		}else if _, ok := v.Field(i).Interface().(ValInteger); ok {
			fld := &FieldInt{}
			if tag_val, ok := field.Tag.Lookup("maxValue"); ok {
				v_i,err := StrToInt(tag_val)
				if err != nil {
					panic("GenModelMD ValInteger lookup maxValue convert error on field: "+field.Name)
				}
				fld.SetMaxValue(NewParamInt64(v_i))
			}
			if tag_val, ok := field.Tag.Lookup("minValue"); ok {
				v_i,err := StrToInt(tag_val)
				if err != nil {
					panic("GenModelMD ValInteger lookup minValue convert error on field: "+field.Name)
				}
				fld.SetMinValue(NewParamInt64(v_i))
			}
			if tag_val, ok := field.Tag.Lookup("notZero"); ok {
				v_b,_ := StrToBool(tag_val)
				fld.SetNotZero(NewParamBool(v_b))
			}			
			res[field.Name] = fld
		
		}else if _, ok := v.Field(i).Interface().(ValTexter); ok {
			fld := &FieldText{}
			if tag_val, ok := field.Tag.Lookup("length"); ok {
				v_i, err := strconv.Atoi(tag_val) //simple int
				if err != nil {
					panic("GenModelMD ValTexter lookup length convert error on field: "+field.Name)
				}
				fld.SetLength(NewParamInt(v_i))
			}
			res[field.Name] = fld
		
		}else if _, ok := v.Field(i).Interface().(ValBytera); ok {
			switch v.Field(i).Interface().(type) {
			case ValJSON:
				res[field.Name] = &FieldJSON{}
			default:
				res[field.Name] = &FieldBytea{}
			}
		
		}else if _, ok := v.Field(i).Interface().(ValBooler); ok {
			res[field.Name] = &FieldBool{}
		
		}else if _, ok := v.Field(i).Interface().(ValTimer); ok {
			switch v.Field(i).Interface().(type) {
			case ValTime:
				res[field.Name] = &FieldTime{}

			case ValDate:
				res[field.Name] = &FieldDate{}

			case ValDateTime:
				res[field.Name] = &FieldDateTime{}

			default:
				//all ValDateTimeTZ:
				res[field.Name] = &FieldDateTimeTZ{}
			}
		
		}else if _, ok := v.Field(i).Interface().(ValFloater); ok {
			fld := &FieldFloat{}
			if tag_val, ok := field.Tag.Lookup("maxValue"); ok {
				v_f,err := StrToFloat(tag_val)
				if err != nil {
					panic("GenModelMD ValFloater lookup maxValue convert error on field: "+field.Name)
				}
				fld.SetMaxValue(NewParamFloat(v_f))
			}
			if tag_val, ok := field.Tag.Lookup("minValue"); ok {
				v_f,err := StrToFloat(tag_val)
				if err != nil {
					panic("GenMDDescr lookup minValue convert error on field: "+field.Name)
				}
				fld.SetMinValue(NewParamFloat(v_f))
			}
			if v, ok := field.Tag.Lookup("notZero"); ok {
				v_b,_ := StrToBool(v)
				fld.SetNotZero(NewParamBool(v_b))
			}
			if tag_val, ok := field.Tag.Lookup("precision"); ok {
				v_i,err := StrToInt(tag_val)
				if err != nil {
					panic("GenMDDescr ValFloater lookup precision convert error on field: "+field.Name)
				}
				fld.SetPrecision(NewParamInt(int(v_i)))
				//res[field.Name].SetPrecision(v_i)
			}
			if tag_val, ok := field.Tag.Lookup("length"); ok {
				v_i,err := StrToInt(tag_val)
				if err != nil {
					panic("GenMDDescr ValFloater lookup length float convert error on field: "+field.Name)
				}
				fld.SetLength(NewParamInt(int(v_i)))
			}
			
			res[field.Name] = fld
		
		
		}else if _, ok := v.Field(i).Interface().(ValArrayer); ok {
			res[field.Name] = &FieldArray{}
		
		}else if _, ok := v.Field(i).Interface().(ValAssocArrayer); ok {
			res[field.Name] = &FieldAssocArray{}
		
		}else{
			//panic(fmt.Sprintf("GenMDDescr unsupported field type in struct %s, field num:%d", id, i))
			//v.Field(i).Kind()==reflect.Slice ||reflect.Struct
			fld := &FieldBytea{}
			res[field.Name] = fld
			
		}
		
		//common attributes
		res[field.Name].SetId(id)
		res[field.Name].SetAlias(field.Tag.Get("alias"))		
		res[field.Name].SetDescr(field.Tag.Get("descr"))
		res[field.Name].SetOrderInList(byte(i))

		if encrypted, ok := field.Tag.Lookup("encrypted"); ok && encrypted == "true" {
			res[field.Name].SetEncrypted(true)
		}

		if no_val_on_copy, ok := field.Tag.Lookup("noValueOnCopy"); ok && no_val_on_copy == "true" {
			res[field.Name].SetNoValueOnCopy(true)
		}
		
		if tag_val, ok := field.Tag.Lookup("defOrder"); ok {
			res[field.Name].SetDefOrder(NewParamBool((strings.ToUpper(tag_val) == "ASC")))
			var order_ind int64
			if tag_ind, ok := field.Tag.Lookup("defOrderIndex"); ok {				
				var err error
				order_ind, err = StrToInt(tag_ind)
				if err != nil {
					panic("GenMDDescr defOrderIndex convert error on field: " + field.Name)
				}				
			}
			res[field.Name].SetDefOrderIndex(byte(order_ind))
		}
		
		prim,_ := StrToBool(field.Tag.Get("primaryKey"))
		res[field.Name].SetPrimaryKey(prim)
		
		req,_ := StrToBool(field.Tag.Get("required"))
		res[field.Name].SetRequired(req)				
		
		s_col,_ := StrToBool(field.Tag.Get("sysCol"))
		res[field.Name].SetSysCol(s_col)
		//for md
		res[field.Name].SetFieldIndex(len(res)-1)
		
	}
	return res
}
