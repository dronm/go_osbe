package osbe

import (
	"reflect"
	"fmt"
	"context"
	
	"osbe/model"
	
	"github.com/jackc/pgx/v4"
)

//Document per page count can be set:
//		1) for model in DocPerPageCount property
//		2) With constant DocPerPageCountConstantID
//		3) Otherwise DEF_DOC_PER_PAGE_COUNT will be used
// Model may have LimitCount. If set no more then this value will be served. Otherwise no more then 'Document per page count' value will be served.
// Document default count value constant ID must be specified at startup. It should return int value.
//

//Default count value
//@ToDo: subscribe to DocPerPageCount update local event!!!
var DocPerPageCount int
var DocPerPageCountConstantID string

const (
	SQL_STATEMENT_LIMIT_2 = "OFFSET %d LIMIT %d";
	SQL_STATEMENT_LIMIT_1 = "LIMIT %d";
	DEF_DOC_PER_PAGE_COUNT int = 50
	DEF_LIMIT_COUNT int = 5000
)

func SetDocPerPageCount(conn *pgx.Conn) error {
	return conn.QueryRow(context.Background(), fmt.Sprintf("SELECT const_%s_val()", DocPerPageCountConstantID)).Scan(&DocPerPageCount)
}

func parseSQLLimitFromArgs(rfltArgs reflect.Value) (int, int) {		
	return int(GetIntArgValByName(rfltArgs, "From", 0)), int(GetIntArgValByName(rfltArgs, "Count", 0))
} 

func GetSQLLimitFromArgs(rfltArgs reflect.Value, scanModelMD *model.ModelMD, conn *pgx.Conn, docPerPageCount int) (string, int, int, error) {		
	from_v, count_v := parseSQLLimitFromArgs(rfltArgs)
		
	if docPerPageCount == 0 && count_v == 0 && scanModelMD != nil && scanModelMD.DocPerPageCount > 0  {
		count_v = scanModelMD.DocPerPageCount
		
	}else if docPerPageCount == 0 && count_v == 0 && DocPerPageCountConstantID != "" {
		if DocPerPageCount == 0 {
			if err := SetDocPerPageCount(conn); err != nil {
				return "", 0, 0, err
			}
		}
		count_v = DocPerPageCount
	}
	/*
	if scanModelMD != nil {
		if scanModelMD.LimitCount > 0 && (count_v == 0 || count_v > scanModelMD.LimitCount) {
			count_v = scanModelMD.LimitCount
			
		}else if scanModelMD.LimitConstant != "" && DocPerPageCount > 0 && (count_v == 0 || count_v > DocPerPageCount) {
			count_v = DocPerPageCount
			
		}else if scanModelMD.LimitConstant != "" && DocPerPageCount == 0 {			
			if err := SetDocPerPageCount(conn, scanModelMD.LimitConstant); err != nil {
				return "", 0, 0,err
			}
			if DocPerPageCount > 0 && (count_v == 0 || count_v > DocPerPageCount) {
				count_v = DocPerPageCount
			}
		}
	}
	*/
	
	if count_v == 0 {
		count_v = DEF_DOC_PER_PAGE_COUNT
		
	}else if count_v > DEF_LIMIT_COUNT {
		count_v = DEF_LIMIT_COUNT //global limit
	}
	//model limit - the highest priority
	if scanModelMD != nil && scanModelMD.LimitCount > count_v {
		count_v = scanModelMD.LimitCount
	}
		
//fmt.Println("GetSQLLimitFromArgs count_v=", count_v)	
	if from_v ==0 && count_v ==0 {
		return "", 0, 0, nil
		
	}else if from_v ==0 {
		return fmt.Sprintf(SQL_STATEMENT_LIMIT_1, count_v), 0, count_v, nil
	}
	return fmt.Sprintf(SQL_STATEMENT_LIMIT_2, from_v, count_v), from_v, count_v, nil
}
