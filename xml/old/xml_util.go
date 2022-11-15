package xml_util

import (
	"fmt"
	"reflect"
	//"errors"
	"encoding/xml"

	"osbe/model"
	"osbe/fields"	
)

const (
	XML_HEADER = `<?xml version="1.0" encoding="UTF-8"?>`
	FLD_NUL_VAL = ` xsi:nil="true"`
)

func ModelsToXML(models model.ModelCollection) ([]byte, error){
	xml_s := XML_HEADER
	xml_s += "<document>"
	
	//var is_sys int
	for _, m := range models {
		is_sys := 0
		if m.GetSysModel() {
			is_sys = 1
		}
		xml_s += fmt.Sprintf(`<model id="%s" sysModel="%d" rowsPerPage="%d" totalCount="%d">`, m.GetID(), is_sys, m.GetRowsPerPage(), m.GetTotalCount())
		xml_s += `<rows>`
		for _, row := range m.GetRows() {
			v := reflect.ValueOf(row)
			if v.Kind() == reflect.Ptr {
				v = v.Elem();
			}			
			if v.Kind() != reflect.Struct {
data, err := xml.MarshalIndent(row, "", "  ")			
if err != nil {
	return nil, err 
}
xml_s+=string(data)
continue
/*		
fmt.Println("v=",v)
fmt.Println("v.Type()=",v.Type())
fmt.Println("v.Kind()=",v.Kind())
v := reflect.ValueOf(v)
if v.Kind() == reflect.Ptr {
	v = v.Elem();
}
fmt.Println("v2=",v)
fmt.Println("v2.Type()=",v.Type())
fmt.Println("v2.Kind()=",v.Kind())
if v.Kind() == reflect.Struct {
	fmt.Println("STRUCT!!! v2.NumField=",v.NumField())
	//v.Field(1).Addr().Interface()
}

v = reflect.ValueOf(v)
if v.Kind() == reflect.Ptr {
	v = v.Elem();
}
fmt.Println("v3=",v)
fmt.Println("v3.Type()=",v.Type())
fmt.Println("v3.Kind()=",v.Kind())
if v.Kind() == reflect.Ptr {
	v = v.Elem();
}
if v.Kind() == reflect.Struct {
	fmt.Println("STRUCT!!! v3.NumField=",v.NumField())
}

fmt.Println("v4=",v)
fmt.Println("v4.Type()=",v.Type())
fmt.Println("v4.Kind()=",v.Kind())
fmt.Println("v4Num=",v.NumField())
v = reflect.ValueOf(v)
if v.Kind() == reflect.Ptr {
	v = v.Elem();
}
fmt.Println("v5=",v)
fmt.Println("v5.Type()=",v.Type())
fmt.Println("v5.Kind()=",v.Kind())
fmt.Println("v5Num=",v.NumField())
v = reflect.ValueOf(v)
if v.Kind() == reflect.Ptr {
	v = v.Elem();
}
fmt.Println("v6=",v)
fmt.Println("v6.Type()=",v.Type())
fmt.Println("v6.Kind()=",v.Kind())
fmt.Println("v6Num=",v.NumField())


return nil, errors.New("ModelsToXML row is not a struct")*/
			}else{
//fmt.Println("v=",v)
//fmt.Println("v.Type()=",v.Type())
//fmt.Println("v.Kind()=",v.Kind())
			
}

			t := v.Type()
									
			xml_s += `<row xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">`			
			for i := 0; i < v.NumField(); i++ {
				fld_id := string(t.Field(i).Tag.Get("json"))
				
				fld_val := ""
				fld_null :=""
				fld_i := v.Field(i).Interface()
				switch fld_i.(type) {
					case fields.ValBool:
						v_bool := fld_i.(fields.ValBool)
						if v_bool.IsNull {
							fld_null = FLD_NUL_VAL
						}else{						
							v_bool_v,_ := v_bool.MarshalJSON()
							fld_val = string(v_bool_v)
						}
					case fields.ValInt:
						v_int := fld_i.(fields.ValInt)
						if v_int.IsNull {
							fld_null = FLD_NUL_VAL
						}else{													
							v_int_v,_ := v_int.MarshalJSON()
							fld_val = string(v_int_v)
						}
					case fields.ValFloat:
						v_float := fld_i.(fields.ValFloat)
						if v_float.IsNull {
							fld_null = FLD_NUL_VAL
						}else{																				
							v_float_v,_ := v_float.MarshalJSON()
							fld_val = string(v_float_v)
						}
					case fields.ValText:
						v_txt := fld_i.(fields.ValText)
						if v_txt.IsNull {
							fld_null = FLD_NUL_VAL
						}else{																				
							fld_val = v_txt.GetValue()
						}
					case int:
						fld_val = fmt.Sprintf("%d", fld_i.(int))
						
					case int32:
						fld_val = fmt.Sprintf("%d", fld_i.(int32))

					case int64:
						fld_val = fmt.Sprintf("%d", fld_i.(int64))

					case float32:
						fld_val = fmt.Sprintf("%f", fld_i.(float32))

					case float64:
						fld_val = fmt.Sprintf("%f", fld_i.(float64))

					case bool:
						v_bool := fld_i.(bool)
						var v_bool_v byte
						if v_bool {
							v_bool_v = 1
						}else{
							v_bool_v = 0
						}
						fld_val = fmt.Sprintf("%d", v_bool_v)


					case string:
						fld_val = fmt.Sprintf("%s", fld_i.(string))
						
					default:
						fld_val = fmt.Sprintf("%s",fld_i)
						
				}
				
				xml_s += fmt.Sprintf(`<%s%s>%s</%s>`, fld_id, fld_null, fld_val, fld_id)
			}			
			
			xml_s += `</row>`
		}	
		xml_s += `</rows>`
		xml_s += `</model>`
	}
	
	xml_s += "</document>"
	return []byte(xml_s), nil
}

