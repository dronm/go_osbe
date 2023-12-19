package xml

import (
	"reflect"
	"fmt"
	"strings"

	"osbe/model"
	"osbe/fields"
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
	var xml_s strings.Builder
	if v.Kind() == reflect.Struct {		
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			if t.Field(i).Anonymous {
				//xml_s += rowToXML(v.Field(i))
				xml_s.WriteString(rowToXML(v.Field(i)))
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
			xml_s.WriteString(getField(fld_id, v.Field(i).Interface(), omit_if_empty))
		}
	}else if v.Kind() == reflect.Map {
		//accept map[string]value
		for _, e := range v.MapKeys() {
			if fld_id, ok := e.Interface().(string); ok {
				xml_s.WriteString(getField(fld_id, v.MapIndex(e).Interface(), false))
			}else{
				break
			}
		}
	}else{
		fmt.Println("rowToXML skeeping reflect type=", v.Kind())	
	}
	return xml_s.String()
}

func ModelToXML(m model.Modeler) string {
	var xml_s strings.Builder
	raw_d := m.GetRawData()
	if len(raw_d) > 0 {
		xml_s.WriteString(string(raw_d))
		return xml_s.String()
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
			agg_funcs_s += fmt.Sprintf(` %s="%s"`, agg_v.Alias, agg_v.ValStr)
		}
	}
	
	xml_s.WriteString(fmt.Sprintf(`<model id="%s" sysModel="%d" rowsPerPage="%d" listFrom="%d"%s>`, m.GetID(), is_sys, m.GetRowsPerPage(),
				m.GetListFrom(), agg_funcs_s))
	for _, row := range m.GetRows() {		
		if xml_row := rowToXML(reflect.ValueOf(row)); xml_row != "" {
			xml_s.WriteString(`<row xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">`)
			xml_s.WriteString(xml_row)
			xml_s.WriteString(`</row>`)
		}
	}
	xml_s.WriteString(`</model>`)
	return xml_s.String()
}

func MetadataToXML(md *model.ModelMD) string {
	var xml_s strings.Builder
	xml_s.WriteString(fmt.Sprintf(`<metadata modelId="%s">`, md.ID))
	//correct order
	f_list := make([]fields.Fielder, len(md.Fields))
	for i := 0; i<len(md.Fields); i ++ {
		for _, f := range md.Fields {		
			if i == f.GetFieldIndex() {
				f_list[i] = f
				break
			}
		}
	}	
	for _, f := range f_list {		
		attrs := ""
		alias := f.GetAlias()
		if alias != "" {
			attrs+= fmt.Sprintf(` alias="%s"`, alias)
		}
		sys_col := f.GetSysCol()
		if sys_col {
			attrs+= ` sysCol="TRUE"`
		}
		xml_s.WriteString(fmt.Sprintf(`<field id="%s" dataType="%d"%s/>`, f.GetId(), f.GetDataType(), attrs))
	}
	xml_s.WriteString(`</metadata>`)
	return xml_s.String()
}

func ModelsToXML(models model.ModelCollection, includeMD bool) string {
	var xml_s strings.Builder
	for _, m := range models {		
		if includeMD && m.GetMetadata() != nil{
			//add md
			xml_s.WriteString(MetadataToXML(m.GetMetadata()))
		}
		xml_s.WriteString(ModelToXML(m))
	}
	return xml_s.String()
}

func Marshal(models model.ModelCollection, includeMD bool) ([]byte, error){
	/*xml_s := XML_HEADER +
		"<document>" +
		ModelsToXML(models, includeMD) +
		"</document>"	
	return []byte(xml_s), nil		
	*/
	res := append([]byte(XML_HEADER), []byte("<document>")...)
	res = append(res, []byte(ModelsToXML(models, includeMD))...)
	res = append(res, []byte("</document>")...)
	return res, nil
}

func EscapeForXML(s string) string {
	res := strings.ReplaceAll(s, "&", "&amp;") //#38
	res = strings.ReplaceAll(res, "<", "&lt;") //#60
	res = strings.ReplaceAll(res, ">", "&gt;") //#62
	res = strings.ReplaceAll(res, `"`, "&quot;") //#34
	res = strings.ReplaceAll(res, "'", "&apos;") //#39
	return res
}

