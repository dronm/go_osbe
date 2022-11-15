package osbe

import(
	"reflect"
	"errors"	
	"fmt"	
	
)

//External argument validation
func ValidateExtArgs(app Applicationer, pm PublicMethod, contr Controller, argv reflect.Value) error {

	//combines all errors in one string	
	valid_err := ""
	
	var arg_fld reflect.Value
	var arg_fld_v reflect.Value
	var argv_empty = argv.IsZero()
	for fid, fld := range pm.GetFields() {		
		if !argv_empty {
			arg_fld = argv.FieldByName(fid)
		}
		
		//GetRequired is implemented by all fields
		if fld.GetRequired() && (argv_empty || (arg_fld.IsValid() && (!arg_fld.FieldByName("IsSet").Bool() || arg_fld.FieldByName("IsNull").Bool()) ) ) {
			//required field has no value
			appendError(&valid_err, fmt.Sprintf(ER_PARSE_NOT_VALID_EMPTY, fld.GetDescr()) ) 
			
		}else if !argv_empty && arg_fld.IsValid(){
			//check if metadata field implements certain interfaces
			//if it does, call methods of these interfaces
			//fmt.Printf("fid=%s, arg_fld=%v\n",fid, arg_fld)	
			var err error
			arg_fld_v = arg_fld.FieldByName("TypedValue")
			switch fld.GetDataType() {
			case FIELD_TYPE_FLOAT:
				err = ValidateFloat(fld.(FielderFloat), arg_fld_v.Float())				
			case FIELD_TYPE_INT:
				err = ValidateInt(fld.(FielderInt), arg_fld_v.Int())				
			case FIELD_TYPE_TEXT:
				err = ValidateText(fld.(FielderText), arg_fld_v.String())				
			}
			if err != nil {
				appendError(&valid_err, err.Error() ) 
			}
		}else if !argv_empty {
			//field is present in ext argg but is not in metadata
			app.GetLogger().Warnf("External argument %s is not present in metadata of %s.%s", fid, contr.GetID(), pm.GetID())
			//fmt.Println("Field",fid, "arg_fld=",arg_fld)
		}
		
		//fmt.Println("Field",fid,"IsSet=",arg_fld.FieldByName("IsSet"),"IsNull=",arg_fld.FieldByName("IsNull"),"Value=",arg_fld.FieldByName("TypedValue"))
	}
	if valid_err != "" {
		return errors.New(valid_err)
	}
	
	return nil
}

func appendError(er *string, addStr string) {
	if *er !="" {
		*er+= ", "
	}
	*er+= addStr
}

