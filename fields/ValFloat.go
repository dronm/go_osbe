package fields

import (
	"errors"
	"strconv"
	"database/sql/driver"
	"strings"
	"math"
//	"fmt"
)

type ValFloat struct {
	Val
	TypedValue float64
}

func (v ValFloat) GetValue() float64{
	if v.IsNull {
		return 0
	}else{
		return v.TypedValue
	}	
}

func (v ValFloat) GetIsNull() bool{
	return v.IsNull
}

func (v ValFloat) GetIsSet() bool{
	return v.IsSet
}

func (v *ValFloat) SetValue(vF float64){
	v.TypedValue = vF
	v.IsSet = true
	v.IsNull = false
}

func (v ValFloat) SetNull(){
	v.TypedValue = 0
	v.IsSet = true
	v.IsNull = true
}

//Custom Float unmarshal
func (v *ValFloat) UnmarshalJSON(data []byte) error {
	v.IsSet = true
	v.TypedValue = 0 
	
	if ExtValIsNull(data){
		v.IsNull = true
		return nil
	}
	
	v_str := ExtRemoveQuotes(data)
	temp, err := StrToFloat(v_str)
	if err != nil {
		return err
	}
	/*v_str = strings.Replace(v_str, ",", ".", 1)	
	temp, err := strconv.ParseFloat(v_str, 64)
	if err != nil {
		return errors.New(ER_UNMARSHAL_FLOAT + err.Error())
	}
	*/
	v.TypedValue = temp
	
	return nil	
}

func (v ValFloat) String() string {
	return strconv.FormatFloat(v.TypedValue, 'f', -1, 64)
}

func (v *ValFloat) MarshalJSON() ([]byte, error) {
	if v.IsNull {
		return []byte(JSON_NULL), nil
		
	}else{
		return []byte(v.String()), nil
	}
}

//driver.Scanner, driver.Valuer interfaces
func (v *ValFloat) Scan(value interface{}) error {
	v.IsSet = true
	v.IsNull = false
	if value == nil {
		v.IsNull = true
		return nil
	}else{
//fmt.Println("ValFloat csan v=",value)	
		switch val := value.(type) {
			case float64:
				v.TypedValue = val
				return nil
			case float32:
				v.TypedValue = float64(val)
				return nil
			case int64:
				v.TypedValue = float64(val)
				return nil
			case string:	
				//0e0=0 1035e-2=10,35
				val_s := string(val)				
				if exp_p := strings.Index(val_s, "e"); exp_p == -1 {
					if i64, err := strconv.ParseInt(val_s, 10, 64); err == nil {
						v.TypedValue = float64(i64)
						return nil
					}					
				}else{					
					if num, err := strconv.ParseInt(val_s[:exp_p], 10, 64); err == nil {
						if exp, err := strconv.ParseInt(val_s[exp_p+1:], 10, 64); err == nil {
							v.TypedValue = float64(num) * math.Pow(10.0, float64(exp))							
							return nil
						}						
					}
				}
		}	
		return errors.New(ER_UNMARSHAL_FLOAT + "unsupported value")
		
	}
	return nil
}

func (v ValFloat) Value() (driver.Value, error) {
	if v.IsNull {
		return driver.Value(nil),nil
	}
	return driver.Value(v.TypedValue), nil
}

func StrToFloat(vStr string) (float64, error) {
	vStr = strings.Replace(vStr, ",", ".", 1)	
	temp, err := strconv.ParseFloat(vStr, 64)
	if err != nil {
		return 0, errors.New(ER_UNMARSHAL_FLOAT + err.Error())
	}
	return temp, nil
}

func NewValFloat(val float64, isNull bool) ValFloat{
	return ValFloat{Val{true, isNull}, val}
}
