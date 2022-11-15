package osbe

import (
	"reflect"
	"strings"
	
	"osbe/model"
)

const (
	SQL_STATEMENT = "ORDER BY ";
	SQL_DIR_ASC = "ASC"
	SQL_DIR_DESC = "DESC"
	
	PAR_DIR_ASC = "asc"
	PAR_DIR_DESC = "desc"	
)

const (
	DIRECT_ASC SQLDirectType = iota
	DIRECT_DESC SQLDirectType = iota
)

type SQLDirectType byte

type SQLOrder struct {
	Expr string
	Direct SQLDirectType
}

type SQLOrders []SQLOrder

func addSQLOrderByExpr(expr string, direct SQLDirectType, sql *string) {
	if *sql != "" {
		*sql += ", "
	}else{
		*sql = SQL_STATEMENT
	}
	*sql += expr + " "
	if direct == DIRECT_DESC {
		*sql += SQL_DIR_DESC
	}else{
		*sql += SQL_DIR_ASC
	}	
}

//returns fields, directs as slices
func parseSQLOrderByFromArgs(rfltArgs reflect.Value, fieldSep string) ([]string, []SQLDirectType) {		
	if ids := GetTextArgValByName(rfltArgs, "Ord_fields", ""); ids != "" {		
		fields_s := strings.Split(ids, fieldSep)
		f_cnt := len(fields_s)
		if f_cnt == 0 {
			return nil, nil
		}
		
		var directs_s []SQLDirectType
		if directs := GetTextArgValByName(rfltArgs, "Ord_directs", ""); directs != "" {
			directs_str := strings.Split(directs, fieldSep)
			directs_s = make([]SQLDirectType, len(directs_str))
			ind := 0
			for _, dir := range directs_str {
				if dir == PAR_DIR_ASC || strings.ToUpper(dir) == strings.ToUpper(PAR_DIR_ASC) {
					directs_s[ind] = DIRECT_ASC
				}else{
					directs_s[ind] = DIRECT_DESC
				}
				ind++
			}
		}
		/*
		f_dir := reflect.Indirect(rfltArgs).FieldByName("Ord_directs")
		f_dir_i := f_dir.Interface()
		if f_dir_i != nil {
			if directs, ok := f_dir_i.(fields.ValText); ok && directs.IsSet {
				directs_str := strings.Split(directs.GetValue(), fieldSep)
				directs_s = make([]SQLDirectType, len(directs_str))
				ind := 0
				for _, dir := range directs_str {
					if dir == PAR_DIR_ASC || strings.ToUpper(dir) == strings.ToUpper(PAR_DIR_ASC) {
						directs_s[ind] = DIRECT_ASC
					}else{
						directs_s[ind] = DIRECT_DESC
					}
					ind++
				}
			}
		}
		*/
		if directs_s == nil {
			//missing directs
			directs_s = make([]SQLDirectType, len(fields_s))
			for i:= 0; i<len(directs_s); i++ {
				directs_s[i] = DIRECT_ASC
			}
		}
		
		if  len(directs_s) < f_cnt  {
			//missing some directs					
			for i:= len(directs_s); i<f_cnt; i++ {
				directs_s = append(directs_s, DIRECT_ASC)
			}
		}
			
		return fields_s, directs_s
	}
	return nil, nil
} 

func NewSQLOrderByFromArgs(rfltArgs reflect.Value, fieldSep string) *SQLOrders {		
	fields_s, directs_s := parseSQLOrderByFromArgs(rfltArgs, fieldSep)
	if fields_s == nil || directs_s == nil {
		return nil
	}
	
	//o := &SQLOrderByCont{orders: make([]SQLOrder, len(fields_s))}
	o := make(SQLOrders, len(fields_s))
	for ind, fld := range fields_s {			
		//o.orders[ind] = SQLOrder{Expr: fld, Direct: directs_s[ind]}
		o[ind] = SQLOrder{Expr: fld, Direct: directs_s[ind]}
	}
	return &o

}

func GetSQLOrderByFromArgsOrDefault(rfltArgs reflect.Value, fieldSep string, modelMD *model.ModelMD, encryptKey string) string {		
	sql := GetSQLOrderByFromArgs(rfltArgs, fieldSep)
	if sql == "" && modelMD != nil {
		sql = GetSQLDefaultOrderBy(modelMD, encryptKey)
	}
	return sql
}

func GetSQLOrderByFromArgs(rfltArgs reflect.Value, fieldSep string) string {		
	fields_s, directs_s := parseSQLOrderByFromArgs(rfltArgs, fieldSep)
	if fields_s == nil || directs_s == nil {
		return ""
	}
	
	sql := ""
	for ind, fld := range fields_s {						
		addSQLOrderByExpr(fld, directs_s[ind], &sql)
	}
	return sql
} 

func GetSQLDefaultOrderBy(modelMD *model.ModelMD, encryptKey string) string {
	/*
	sql := ""
	for _, fld := range modelMD.GetFields() {
		if  o := fld.GetDefOrder(); o.IsSet {
			var direct SQLDirectType
			if o.Value {
				direct = DIRECT_ASC
			}else{
				direct = DIRECT_DESC
			}
			addSQLOrderByExpr(fld.GetId(), direct, &sql)
		}
	}
	*/
	sql := modelMD.GetFieldDefOrder(encryptKey)
	if sql != "" {
		return SQL_STATEMENT + " "+sql
	}
	return ""
}
