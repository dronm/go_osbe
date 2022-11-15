package osbe

import (
	"reflect"
	"fmt"
	"context"
	
	"osbe/model"
	
	"github.com/jackc/pgx/v4"
)

//@ToDo: suscribe to limit_constant_value update local event!!!
var limit_constant_value int

const (
	SQL_STATEMENT_LIMIT_2 = "OFFSET %d LIMIT %d";
	SQL_STATEMENT_LIMIT_1 = "LIMIT %d";
)

func parseSQLLimitFromArgs(rfltArgs reflect.Value) (int, int) {		
	return int(GetIntArgValByName(rfltArgs, "From", 0)), int(GetIntArgValByName(rfltArgs, "Count", 0))
} 

func GetSQLLimitFromArgs(rfltArgs reflect.Value, scanModelMD *model.ModelMD, conn *pgx.Conn, defLimit int) (string, int, int, error) {		
	from_v, count_v := parseSQLLimitFromArgs(rfltArgs)
	
	if scanModelMD != nil {
		if scanModelMD.LimitCount > 0 && (count_v == 0 || count_v > scanModelMD.LimitCount) {
			count_v = scanModelMD.LimitCount
			
		}else if scanModelMD.LimitConstant != "" && limit_constant_value > 0 && (count_v == 0 || count_v > limit_constant_value) {
			count_v = limit_constant_value
			
		}else if scanModelMD.LimitConstant != "" && limit_constant_value == 0 {			
			if err := conn.QueryRow(context.Background(),
				fmt.Sprintf("SELECT const_%s_val()", scanModelMD.LimitConstant)).Scan(&limit_constant_value); err != nil {
				return "", 0, 0,err
			}
			if limit_constant_value > 0 && (count_v == 0 || count_v > limit_constant_value) {
				count_v = limit_constant_value
			}
		}
	}
	
	if count_v == 0 {
		count_v = defLimit
	}
	
	//Global limit???
	if from_v ==0 && count_v ==0 {
		return "", 0, 0, nil
		
	}else if from_v ==0 {
		return fmt.Sprintf(SQL_STATEMENT_LIMIT_1, count_v), 0, count_v, nil
	}
	return fmt.Sprintf(SQL_STATEMENT_LIMIT_2, from_v, count_v), from_v, count_v, nil
}
