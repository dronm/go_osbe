package xml

import (
	"reflect"
	"fmt"
	"strings"

	"osbe/model"
)

const (
	XML_HEADER = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
)

type Nullable interface {
	GetIsNull() bool
}

func getField(fld_id string, fld_i interface{}, omit_if_empty bool) string {
	if fld_v, ok := fld_i.(Nullable); ok && fld_v.GetIsNull() {
		if !omit_if_empty {
			return fmt.Sprintf(`<%s xsi:nil="true"/>`, fld_id)
		}
	}else{
		fld_val_s := ""				
		if fld_v, ok := fld_i.(fmt.Stringer); ok {
			fld_val_s = EscapeForXML(fld_v.String())
		}else{	
		
			switch fld_i.(type) {
			case int:
				fld_val_s = fmt.Sprintf("%d", fld_i.(int))
				
			case int32:
				fld_val_s = fmt.Sprintf("%d", fld_i.(int32))

			case int64:
				fld_val_s = fmt.Sprintf("%d", fld_i.(int64))

			case float32:
				fld_val_s = fmt.Sprintf("%f", fld_i.(float32))

			case float64:
				fld_val_s = fmt.Sprintf("%f", fld_i.(float64))

			case bool:
				v_bool := fld_i.(bool)
				if v_bool {
					fld_val_s = "true"
				}else{
					fld_val_s = "false"
				}

			case string:
				fld_val_s = EscapeForXML(fld_i.(string))
				
			default:
				fld_val_s = EscapeForXML(fmt.Sprintf("%s",fld_i))
			}					
		}
		if fld_val_s != "" || !omit_if_empty {
			return fmt.Sprintf(`<%s>%s</%s>`, fld_id, fld_val_s, fld_id)	
		}					
	}
	return ""
}

//returns xml string or empty string if it is not a struct/map[string]
//only fields with json tag are included
//if xml tag is present and omitempty=true and XML field value is empty,field is not included
//
//func rowToXML(row interface{}) string {
func rowToXML(v reflect.Value) string {		
	//v := reflect.ValueOf(row)
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		if v.IsNil() {
			break
		}
		v = v.Elem()
	}
	xml_s := ""
	if v.Kind() == reflect.Struct {		
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			if t.Field(i).Anonymous {
				//xml_s += rowToXML(v.Field(i).Interface())
				xml_s += rowToXML(v.Field(i))
				continue
			}
			fld_id, ok := t.Field(i).Tag.Lookup("json")
			if !ok {
				continue
			}
			omit_if_empty := false
			if xml_tag, ok := t.Field(i).Tag.Lookup("xml"); ok {
				xml_tag_vals := strings.Split(xml_tag,",")
				for _,xml_tag_v := range xml_tag_vals {
					if xml_tag_v == "omitempty" {
						omit_if_empty = true	
					}
				}
			}
			xml_s += getField(fld_id, v.Field(i).Interface(), omit_if_empty)
			//fld_i := v.Field(i).Interface()
			/*
			if fld_v, ok := fld_i.(Nullable); ok && fld_v.GetIsNull() {
				if !omit_if_empty {
					xml_s += fmt.Sprintf(`<%s xsi:nil="true"/>`, fld_id)
				}
				
			}else{
				fld_val_s := ""				
				if fld_v, ok := fld_i.(fmt.Stringer); ok {
					fld_val_s = EscapeForXML(fld_v.String())
				}else{	
				
					switch fld_i.(type) {
					case int:
						fld_val_s = fmt.Sprintf("%d", fld_i.(int))
						
					case int32:
						fld_val_s = fmt.Sprintf("%d", fld_i.(int32))

					case int64:
						fld_val_s = fmt.Sprintf("%d", fld_i.(int64))

					case float32:
						fld_val_s = fmt.Sprintf("%f", fld_i.(float32))

					case float64:
						fld_val_s = fmt.Sprintf("%f", fld_i.(float64))

					case bool:
						v_bool := fld_i.(bool)
						if v_bool {
							fld_val_s = "true"
						}else{
							fld_val_s = "false"
						}

					case string:
						fld_val_s = EscapeForXML(fld_i.(string))
						
					default:
						fld_val_s = EscapeForXML(fmt.Sprintf("%s",fld_i))
					}					
				}
				if fld_val_s != "" || !omit_if_empty {
					xml_s += fmt.Sprintf(`<%s>%s</%s>`, fld_id, fld_val_s, fld_id)	
				}					
			}
			*/
		}
	}else if v.Kind() == reflect.Map {
		//accept map[string]value
		for _, e := range v.MapKeys() {
			if fld_id, ok := e.Interface().(string); ok {
				xml_s += getField(fld_id, v.MapIndex(e).Interface(), false)
			}else{
				break
			}
		}
//fmt.Println("MapToXML=",xml_s)		
		
	}else{
		fmt.Println("rowToXML skeeping reflect type=", v.Kind())	
	}
	return xml_s
}

func ModelsToXML(models model.ModelCollection) ([]byte, error){
	xml_s := XML_HEADER
	xml_s += "<document>"
	
	for _, m := range models {
		
		raw_d := m.GetRawData()
		if len(raw_d) > 0 {
			//xml_s += fmt.Sprintf(`<model id="%s">`, m.GetID()) + string(raw_d) + `</model>`
			xml_s += string(raw_d)
			continue
		}
		
		is_sys := 0
		if m.GetSysModel() {
			is_sys = 1
		}
		
		//agg functions
		agg_funcs_s := ""
		agg_vals := m.GetAggFunctionValues()
		if len(agg_vals) > 0 {
			for _,agg_v := range agg_vals {
//fmt.Println(agg_v.Val)
				agg_funcs_s += fmt.Sprintf(` %s="%s"`, agg_v.Alias, agg_v.ValStr)
			}
		}
		
		xml_s += fmt.Sprintf(`<model id="%s" sysModel="%d" rowsPerPage="%d" listFrom="%d"%s>`, m.GetID(), is_sys, m.GetRowsPerPage(),
					m.GetListFrom(), agg_funcs_s)
		for _, row := range m.GetRows() {		
			if xml_row := rowToXML(reflect.ValueOf(row)); xml_row != "" {
			//if xml_row := rowToXML(row); xml_row != "" {
				xml_s += `<row xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">`+xml_row+`</row>`
			}
		}
		xml_s += `</model>`
	}
	
	xml_s += "</document>"
	
	return []byte(xml_s), nil		
}

func EscapeForXML(s string) string {
	res := strings.ReplaceAll(s, "&", "&amp;")
	res = strings.ReplaceAll(res, "<", "&lt;")
	res = strings.ReplaceAll(res, ">", "&gt;")
	res = strings.ReplaceAll(res, `"`, "&quot;")
	res = strings.ReplaceAll(res, "'", "&apos;")
	return res
}

