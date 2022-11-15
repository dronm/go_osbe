package osbe

import(
	"reflect"
	"errors"	
	"strconv"
	"fmt"	
	"context"
	"math/rand"
	"crypto/md5"
	"encoding/hex"	
	
	"osbe/fields"	
	"osbe/response"
	"osbe/model"
	"osbe/socket"
	
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgconn"
)

const (
	KEY_FLD_PREF = "old_"
	KEY_FLD_PREF_LEN = 4
	
	METH_COMPLETE_DEF_COUNT = 50

	RESP_ER_DELETE_CONSTR_VIOL = 500
)

//External argument validation
func ValidateExtArgs(app Applicationer, pm PublicMethod, contr Controller, argv reflect.Value) error {

	md_model := pm.GetFields()
	if md_model == nil {
		return nil
	}
	
	//combines all errors in one string	
	valid_err := ""
	
	var arg_fld reflect.Value
	var arg_fld_v reflect.Value
	var argv_empty = argv.IsZero()
	
	for fid, fld := range md_model {
		
		//fmt.Println("fid=", fid, "GetRequired=", fld.GetRequired(), "argv_empty=", argv_empty, "IsValid=", arg_fld.IsValid())		
		//,"IsSet=", arg_fld.FieldByName("IsSet").Bool(),"IsNull=", arg_fld.FieldByName("IsNull").Bool())
		
		if !argv_empty {
			//Indirect always returns object!
			arg_fld = reflect.Indirect(argv).FieldByName(fid)
		}
		
		//GetRequired is implemented by all fields
		if fld.GetRequired() && (argv_empty || (arg_fld.IsValid() && arg_fld.Kind() == reflect.Struct && (!arg_fld.FieldByName("IsSet").Bool() || arg_fld.FieldByName("IsNull").Bool()) ) ) {
			//required field has no value
			//fmt.Println("required field has no value")
			appendError(&valid_err, fmt.Sprintf(ER_PARSE_NOT_VALID_EMPTY, fld.GetDescr()) ) 
			
		}else if !argv_empty && arg_fld.IsValid() && arg_fld.Kind() == reflect.Struct {
			//fmt.Println("!argv_empty && arg_fld.IsValid()")
			
			//check if metadata field implements certain interfaces
			//if it does, call methods of these interfaces
			//fmt.Printf("fid=%s, arg_fld=%v\n",fid, arg_fld)	
			
			var err error
			arg_fld_v = arg_fld.FieldByName("TypedValue")
			switch fld.GetDataType() {
			case fields.FIELD_TYPE_FLOAT:
				err = fields.ValidateFloat(fld.(fields.FielderFloat), arg_fld_v.Float())				
				
			case fields.FIELD_TYPE_INT:
				err = fields.ValidateInt(fld.(fields.FielderInt), arg_fld_v.Int())				
				
			case fields.FIELD_TYPE_TEXT:
				err = fields.ValidateText(fld.(fields.FielderText), arg_fld_v.String())				

			case fields.FIELD_TYPE_JSON:
				err = fields.ValidateJSON(fld.(fields.FielderJSON), []byte(arg_fld_v.String()))


			case fields.FIELD_TYPE_TIME:
				err = fields.ValidateTime(fld.(fields.Fielder), arg_fld_v.String())				

			case fields.FIELD_TYPE_DATE:
				err = fields.ValidateDate(fld.(fields.Fielder), arg_fld_v.String())				

			case fields.FIELD_TYPE_DATETIME:
				err = fields.ValidateDateTime(fld.(fields.Fielder), arg_fld_v.String())				

			case fields.FIELD_TYPE_DATETIMETZ:
				err = fields.ValidateDateTimeTZ(fld.(fields.Fielder), arg_fld_v.String())				

			case fields.FIELD_TYPE_ENUM:
				err = fields.ValidateEnum(fld.(fields.FielderEnum), arg_fld_v.String())				
				
			/*default:
				appendError(&valid_err, "osbe.ValidateExtArgs: unsupported field type" ) 
			*/
			}
			if err != nil {
				appendError(&valid_err, err.Error() ) 
			}
		//}else if !argv_empty {
			//field is present in ext argg but is not in metadata
		//	app.GetLogger().Warnf("External argument %s is not present in metadata of %s.%s", fid, contr.GetID(), pm.GetID())
			//fmt.Println("Field",fid, "arg_fld=",arg_fld)
		//}else{
			//fmt.Println("Otherwise")
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

//Separates public method arguments into  fieldIds, fieldArgs, retFieldIds,  fieldValues
//fieldIds is a string containing all ids
//fieldArgs is a string with parameters ($1,$2...) to be used in query
//fieldValues interface values
//function is used for insert PublicMethod
func ArgsToInsertParams(rfltArgs reflect.Value) (fieldIds string, fieldArgs string, fieldValues []interface{}) {	
	rfltArgs_o := reflect.Indirect(rfltArgs)
	arg_tp := rfltArgs_o.Type()
	fieldValues = make([]interface{}, 0)
	field_ind := 0
	for i := 0; i < rfltArgs_o.NumField(); i++ {						
		if fld_v, ok := rfltArgs_o.Field(i).Interface().(fields.ValExt); ok {
			if field_id, ok := arg_tp.Field(i).Tag.Lookup("json"); ok {
				var fld_add bool
				var fld_val interface{}
				
				if fld_v.GetIsSet() {
					fld_val = fld_v
					fld_add = true
					
				}else if is_autoInc, ok := arg_tp.Field(i).Tag.Lookup("autoInc"); ok && is_autoInc=="true" {
					//add anyway with NULL
					//fld_add = true
					//fld_val = "DEFAULT"
					//does not work this way 
				}
				
				if fld_add {
					if fieldIds != "" {
						fieldIds += ","
						fieldArgs += ","
					}
					fieldIds += field_id
					fieldArgs += "$"+strconv.Itoa(field_ind+1)
					fieldValues = append(fieldValues, fld_val)
					field_ind++			
				}
			}
		}
	}
	return
}

//puts old_key to where query
func ArgsToUpdateParams(rfltArgs reflect.Value) (fieldQuery string, whereQuery string, fieldValues []interface{}, keys map[string]interface{}) {
	fieldValues = make([]interface{}, 0)
	keys = make(map[string]interface{})
	
	rfltArgs_o := reflect.Indirect(rfltArgs)
	arg_tp := rfltArgs_o.Type()
		
	field_ind := 0
	//field_id := arg_tp.Field(i).Name
	for i := 0; i < rfltArgs_o.NumField(); i++ {						
		if fld_v, ok := rfltArgs_o.Field(i).Interface().(fields.ValExt); ok && fld_v.GetIsSet() {
			if field_id, ok := arg_tp.Field(i).Tag.Lookup("json"); ok && (len(field_id)<=KEY_FLD_PREF_LEN || field_id[:KEY_FLD_PREF_LEN]!=KEY_FLD_PREF) {
				if fieldQuery != "" {
					fieldQuery += ","
				}
				fieldQuery += field_id + "=$"+strconv.Itoa(field_ind+1)
				fieldValues = append(fieldValues, fld_v)			
				field_ind++			
			}
		}
	}
	
	for i := 0; i < rfltArgs_o.NumField(); i++ {						
		if fld_v, ok := rfltArgs_o.Field(i).Interface().(fields.ValExt); ok && fld_v.GetIsSet() {
			if field_id, ok := arg_tp.Field(i).Tag.Lookup("json"); ok && len(field_id)>KEY_FLD_PREF_LEN && field_id[:KEY_FLD_PREF_LEN]==KEY_FLD_PREF {
				if whereQuery != "" {
					whereQuery += " AND "
				}
				whereQuery += field_id[KEY_FLD_PREF_LEN:] + "=$"+strconv.Itoa(field_ind+1)
				fieldValues = append(fieldValues, fld_v)			
				keys[field_id[KEY_FLD_PREF_LEN:]],_ = fld_v.Value()
				field_ind++			
			}
		}
	}
	return
}

//Implements controller insert method
//internally calls UpdateOnArgsWithConn
func UpdateOnArgs(app Applicationer, pm PublicMethod, sock socket.ClientSocketer, rfltArgs reflect.Value, relation string) error {
	pool_conn, pm_err := app.GetPrimaryPoolConn()
	if pm_err != nil {
		return pm_err
	}
	defer pool_conn.Release()
	conn := pool_conn.Conn()
	
	return UpdateOnArgsWithConn(conn, app, pm, sock, rfltArgs, relation)
}

//Implements controller insert method
func UpdateOnArgsWithConn(conn *pgx.Conn, app Applicationer, pm PublicMethod, sock socket.ClientSocketer, rfltArgs reflect.Value, relation string) error {
	f_query, w_query, f_values, keys := ArgsToUpdateParams(rfltArgs)		
	if f_query == "" || w_query == "" {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, ER_UPDATE_EMPTY)
	}
	q := fmt.Sprintf("UPDATE %s SET %s WHERE %s", relation, f_query, w_query)
//fmt.Println("Update query=", q, "params:", f_values)	
	_, err := conn.Exec(context.Background(), q, f_values...)
	if err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("UpdateOnArgsWithConn pgx.Conn.Exec(): %v",err))
	}
	
	//events
//fmt.Printf("Update keys=%v+\n", keys)		
	PublishEventsWithKeys(sock.GetID(), keys, app, pm)
	
	return nil
}

//Implements controller insert method
//internally calls InsertOnArgsWithConn
func InsertOnArgs(app Applicationer, pm PublicMethod, resp *response.Response, sock socket.ClientSocketer, rfltArgs reflect.Value, relation string, retModel interface{}) error {
	pool_conn, pm_err := app.GetPrimaryPoolConn()
	if pm_err != nil {
		return pm_err
	}
	defer pool_conn.Release()
	conn := pool_conn.Conn()
	
	return InsertOnArgsWithConn(conn, app, pm, resp, sock, rfltArgs, relation, retModel)
}

//Implements controller insert method
func InsertOnArgsWithConn(conn *pgx.Conn, app Applicationer, pm PublicMethod, resp *response.Response, sock socket.ClientSocketer, rfltArgs reflect.Value, relation string, retModel interface{}) error {
	field_ids, field_args, f_values := ArgsToInsertParams(rfltArgs)		
	
	ret_field_ids:= "" //return all key fields
	keys := make(map[string]interface{})
	row_val := reflect.ValueOf(retModel).Elem()		
	row_fields := make([]interface{}, 0) //row_val.NumField()
	row_t := row_val.Type()
	for i := 0; i < row_val.NumField(); i++ {
		if field_id, ok := row_t.Field(i).Tag.Lookup("json"); ok {
			if ret_field_ids != "" {
				ret_field_ids += ", "
			}
			ret_field_ids += field_id
			keys[field_id] = nil
			value_field := row_val.Field(i)
			row_fields = append(row_fields, value_field.Addr().Interface())
		}else{
			return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("Field: %s, no json tag!", row_t.Field(i).Name))
		}
	}
	q := ""
	if field_ids == "" {
		q += fmt.Sprintf("INSERT INTO %s DEFAULT VALUES RETURNING %s", relation, ret_field_ids)
	}else{
		q += fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s", relation, field_ids, field_args, ret_field_ids)
	}
//fmt.Println("InsertOnArgs q=",q, "field_values=%v", f_values)	
	rows, err := conn.Query(context.Background(), q, f_values...)
	if err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Query(): %v",err))
	}	

	if rows.Next() {		
		if err := rows.Scan(row_fields...); err != nil {		
			return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Rows.Scan(): %v",err))	
		}
		m := model.New_InsertedId_Model(relation)
		resp.AddModel(m)
	}
	if err := rows.Err(); err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, err.Error())
	}
	
	//events
	i:= 0
	for key,_ := range keys {
		keys[key] = row_fields[i]
		i++
	}
	PublishEventsWithKeys(sock.GetID(), keys, app, pm)
	
	return nil
}

/*
Нужна функция для передачи запроса SELECT * FROM <table> WHERE с установкой в структуру тех полей, которые есть в структуре
func QueryRowModel(conn *pgx.Conn, rowModel ObjectRow, query string, condVals []interface{}) error {
	row, err := conn.Query(context.Background(), query, condVals...)
	if err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Query(): %v",err))
	}
	var row_fields []interface{}
	if rows.Next() {
		//дескрипторы!!!		
		if err := rows.Scan(row_fields...); err != nil {		
			return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Rows.Scan(): %v",err))	
		}
		m := model.New_InsertedId_Model(scanStruct)
		resp.AddModel(m)
	}
	if err := rows.Err(); err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, err.Error())
	}
}
func RowModelToStruct(scanStruct ObjectRow, condQuery string, condVals []interface{}, conn *pgx.Conn) error {	
	fields := ""
	scan_fields := make([]interface{}, 0)
	row_val := reflect.ValueOf(scanStruct).Elem()
	row_t := row_val.Type()
	for i := 0; i < row_val.NumField(); i++ {
		if field_id, ok := row_t.Field(i).Tag.Lookup("json"); ok {
			if fields != "" {
				fields += ", "
			}
			fields += field_id
			scan_fields = append(scan_fields, row_val.Field(i).Addr().Interface())
		}
	}
	if condQuery != "" && condQuery[0:1] != " " {
		condQuery = " "+condQuery
	}
	query := fmt.Sprintf("SELECT %s FROM %s%s", fields, scanStruct.GetDataTable(), condQuery)
	if err := conn.QueryRow(context.Background(), query, condVals...).Scan(scan_fields...); err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.QueryRow(): %v",err))
	}
	return nil
}
*/

//Implements controller get_object method
func GetObjectOnArgs(app Applicationer, resp *response.Response, rfltArgs reflect.Value, scanModelMD *model.ModelMD, scanStruct interface{}) error {
	
	pool_conn, pm_err := app.GetSecondaryPoolConn()
	if pm_err != nil {
		return pm_err
	}
	defer pool_conn.Release()
	conn := pool_conn.Conn()
	
	//fields with key values
	field_vals := make([]interface{}, 0)	
	where := ""
	rfltArgs_o := reflect.Indirect(rfltArgs)
	arg_tp := rfltArgs_o.Type()		
	where_ind := 0
	for i := 0; i < rfltArgs_o.NumField(); i++ {						
		if fld_v, ok := rfltArgs_o.Field(i).Interface().(fields.ValExt); ok && fld_v.GetIsSet() {
			if field_id, ok := arg_tp.Field(i).Tag.Lookup("json"); ok {
				if where != "" {
					where += " AND "
				}
				where += field_id + "=$"+strconv.Itoa(where_ind+1)
				field_vals = append(field_vals, fld_v)			
				where_ind++			
			}
		}
	}
	if where_ind == 0 {
		//should not happen if keys are marked as required in get object model
		//return NewPublicMethodError(response.RESP_ER_INTERNAL, ER_NO_KEYS)
		//happens when http requests insert with get_object without key
		return nil
	}
	
	relation := scanModelMD.DbRelation.Relation
	query_id := relation + "_get_object"
	_, err := conn.Prepare(context.Background(), query_id, fmt.Sprintf("SELECT %s FROM %s WHERE %s", scanModelMD.GetFieldList(), relation, where))
	if err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Prepare(): %v",err))
	}
	
	rows, err := conn.Query(context.Background(), query_id, field_vals...)	
	if err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Query(): %v",err))
	}

	m := &model.Model{ID: model.ModelID(scanModelMD.DbRelation.ID), Rows: make([]model.ModelRow, 0)}	
	for rows.Next() {
		row := scanStruct
		row_val := reflect.ValueOf(row).Elem()
		row_fields := make([]interface{}, 0) //row_val.NumField()
		row_t := row_val.Type()
		for i := 0; i < row_val.NumField(); i++ {
			if _, ok := row_t.Field(i).Tag.Lookup("json"); ok {
				value_field := row_val.Field(i)
				//row_fields[i] = value_field.Addr().Interface()
				row_fields = append(row_fields, value_field.Addr().Interface())
			}
		}
	
		if err := rows.Scan(row_fields...); err != nil {		
			return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Rows.Scan(): %v",err))	
		}
		m.Rows = append(m.Rows, row)
	}
	if err := rows.Err(); err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, err.Error())
	}
	
	resp.AddModel(m)
	return nil
}

//Returns: query, total query and where_params
//string - query, string - total query, []interface{} - where params
func GetListQuery(rfltArgs reflect.Value, scanModelMD *model.ModelMD, conn *pgx.Conn) (string, string, []interface{}, error) {
	if scanModelMD.DbRelation == nil {
		return "", "", nil, NewPublicMethodError(response.RESP_ER_INTERNAL, "Model DbRelation not set")
	}

	f_sep := ArgsFieldSep(rfltArgs)
	orderby_sql := GetSQLOrderByFromArgs(rfltArgs, f_sep)	
	if orderby_sql == "" {
		//order is not set, trying set default column order
		for _, fld := range scanModelMD.GetFields() {
			if  o := fld.GetDefOrder(); o.IsSet {
				var direct SQLDirectType
				if o.Value {
					direct = DIRECT_ASC
				}else{
					direct = DIRECT_ASC
				}
				addSQLOrderByExpr(fld.GetId(), direct, &orderby_sql)
			}
		}
	}
	limit_sql := GetSQLLimitFromArgs(rfltArgs)
	
	where_sql, where_params, err := GetSQLWhereFromArgs(rfltArgs, f_sep, scanModelMD.GetFields(), nil)
	if err != nil {
		return "", "", nil, NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("%v",err))
	}

	relation := scanModelMD.DbRelation.Relation
	query_tmpl := fmt.Sprintf("SELECT %s FROM %s", scanModelMD.GetFieldList(), relation)
	query_tot_tmpl := fmt.Sprintf("SELECT count(*) FROM %s", relation)

	query := ""
	query_tot := ""
	if orderby_sql == "" && limit_sql == "" && where_sql == "" {
		query = relation+"_get_list"
		_, err = conn.Prepare(context.Background(), query, query_tmpl)
		if err != nil {
			return "", "", nil, NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Prepare(): %v",err))
		}		
		
		query_tot = relation+"_get_list_tot"
		_, err = conn.Prepare(context.Background(), query_tot, query_tot_tmpl)
		if err != nil {
			return "", "", nil, NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Prepare(): %v",err))
		}		
		
		
	}else{
		//custom query
		query = query_tmpl
		query_tot = query_tot_tmpl
		if where_sql != "" {
			query += " "+where_sql
			query_tot += " "+where_sql
		}
		if orderby_sql != "" {
			query += " "+orderby_sql
		}		
		if limit_sql != "" {
			query += " "+limit_sql
		}		
	}
//fmt.Println("GetListQuery", query)	
	return query, query_tot, where_params, nil
}

//Executes query and returns it result as model
func QueryResultToModel(modelID string, scanModel interface{}, query string, queryTotal string, condValues []interface{}, conn *pgx.Conn, sysModel bool) (model.Modeler, error) {
	//tot
	tot_cnt := 0 
	if queryTotal != "" {
		row_tot := conn.QueryRow(context.Background(), queryTotal, condValues...)					
		if err := row_tot.Scan(&tot_cnt); err != nil && err != pgx.ErrNoRows {
			return nil, NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("QueryResultToModel total pgx.Rows.Scan(): %v",err))	
		}
	}	
	//var rows pgx.Rows
	rows, err := conn.Query(context.Background(), query, condValues...)	
	if err != nil {
		return nil, NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("QueryResultToModel pgx.Conn.Query(): %v",err))
	}
	
	m := &model.Model{ID: model.ModelID(modelID), TotalCount: tot_cnt, SysModel: sysModel, Rows: make([]model.ModelRow, 0)}	
	for rows.Next() {
		row := reflect.New(reflect.ValueOf(scanModel).Elem().Type()).Interface().(model.ModelRow)
		row_val := reflect.ValueOf(row).Elem()
		row_fields := make([]interface{}, 0) //row_val.NumField()
		row_t := row_val.Type()
		for i := 0; i < row_val.NumField(); i++ {
			if _, ok := row_t.Field(i).Tag.Lookup("json"); ok {
				value_field := row_val.Field(i)
				row_fields = append(row_fields, value_field.Addr().Interface())
			}
		}
	
		if err := rows.Scan(row_fields...); err != nil {		
			return nil, NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("QueryResultToModel pgx.Rows.Scan(): %v",err))	
		}
		m.Rows = append(m.Rows, &row)		
	}
	if err := rows.Err(); err != nil {
		return nil, NewPublicMethodError(response.RESP_ER_INTERNAL, err.Error())
	}
	if queryTotal == "" {
		m.TotalCount = len(m.Rows)
	}
	
	return m, nil
}

func AddQueryResult(resp *response.Response, modelID string, scanModel interface{}, query string, queryTotal string, condValues []interface{}, conn *pgx.Conn, sysModel bool) error {
	model, err := QueryResultToModel(modelID, scanModel, query, queryTotal, condValues, conn, sysModel)
	if err != nil {
		return err
	}
	resp.AddModel(model)
	return nil
}

func GetListOnArgs(app Applicationer, resp *response.Response, rfltArgs reflect.Value, scanModelMD *model.ModelMD, scanModel interface{}) error {
	if scanModelMD == nil {
		app.GetLogger().Error("osbe.GetListOnArgs (util.go) scanModelMD not defined. Potentially error prone code!")
	}
	pool_conn, pm_err := app.GetSecondaryPoolConn()
	if pm_err != nil {
		return pm_err
	}
	defer pool_conn.Release()
	conn := pool_conn.Conn()

	query, query_tot, where_params, err := GetListQuery(rfltArgs, scanModelMD, conn)
	if err != nil {
		return err
	}
	
	return AddQueryResult(resp, scanModelMD.DbRelation.ID, scanModel, query, query_tot, where_params, conn, false)		
}


//Common function for deleting object from DB based on argument keys
func DeleteOnArgKeys(app Applicationer, pm PublicMethod, resp *response.Response, sock socket.ClientSocketer, rfltArgs reflect.Value, relation string) error {
	pool_conn, pm_err := app.GetPrimaryPoolConn()
	if pm_err != nil {
		return pm_err
	}
	defer pool_conn.Release()
	conn := pool_conn.Conn()

	return DeleteOnArgKeysWithConn(conn, app, pm, resp, sock, rfltArgs, relation)
}

//Implements controller delete method
func DeleteOnArgKeysWithConn(conn *pgx.Conn, app Applicationer, pm PublicMethod, resp *response.Response, sock socket.ClientSocketer, rfltArgs reflect.Value, relation string) error {
	
	rfltArgs_o := reflect.Indirect(rfltArgs)
	arg_tp := rfltArgs_o.Type()
	
	f_values := make([]interface{}, arg_tp.NumField())
	keys := make(map[string]interface{})
	
	ids_key := ""
	where_sql := ""
	field_ind := 0
	for i := 0; i < rfltArgs_o.NumField(); i++ {						
		if fld_v, ok := rfltArgs_o.Field(i).Interface().(fields.ValExt); ok && fld_v.GetIsSet() {
			if field_id, ok := arg_tp.Field(i).Tag.Lookup("json"); ok {
				if where_sql != "" {
					where_sql += " AND "
				}
				where_sql += field_id + " = $"+strconv.Itoa(field_ind+1)
				f_values[field_ind] = fld_v
				keys[field_id],_ = fld_v.Value()
				ids_key += "_"+field_id
				field_ind++			
			}
		}
	}
	if where_sql == "" {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, ER_NO_KEYS)
	}

	q := fmt.Sprintf(`DELETE FROM %s WHERE %s`, relation, where_sql)
//fmt.Println("DeleteOnArgKeys q=", q, "f_values=", f_values)			
	_, err := conn.Prepare(context.Background(), relation+ids_key+"_delete", q)
	if err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Prepare(): %v",err))
	}
	
	par, err := conn.Exec(context.Background(), relation+ids_key+"_delete", f_values...)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23503" {
			//custom error
			return NewPublicMethodError(RESP_ER_DELETE_CONSTR_VIOL, ER_DELETE_CONSTR_VIOL)
		}else{
			return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Exec(): %v",err))
		}
	}
	del_m := model.New_Deleted_Model(par.RowsAffected())	
	resp.AddModel(del_m)
	
	//events
	PublishEventsWithKeys(sock.GetID(), keys, app, pm)
		
	return nil

}

//Implements controller complete method
//Internally calls CompleteOnArgsWithConn
func CompleteOnArgs(app Applicationer, resp *response.Response, rfltArgs reflect.Value, scanModelMD *model.ModelMD, scanModel interface{}) error {	
	pool_conn, pm_err := app.GetSecondaryPoolConn()
	if pm_err != nil {
		return pm_err
	}
	defer pool_conn.Release()
	conn := pool_conn.Conn()
	return CompleteOnArgsWithConn(conn, app, resp, rfltArgs, scanModelMD, scanModel)
}

//Implements controller complete method
//args.Ic - insensetive case 1/0
//args.Mid 1 - %pattern%, 0 - pattern%
//scanModelMD
//scanModel
//pattern - pattern to match
//there is also another argument with the same name as model field marked as matchField=true in tag
func CompleteOnArgsWithConn(conn *pgx.Conn, app Applicationer, resp *response.Response, rfltArgs reflect.Value, scanModelMD *model.ModelMD, scanModel interface{}) error {
	if scanModelMD.DbRelation == nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, "Model DbRelation not set")
	}
	
	f_sep := ArgsFieldSep(rfltArgs)
	orderby_sql := GetSQLOrderByFromArgs(rfltArgs, f_sep)	
	limit_sql := GetSQLLimitFromArgs(rfltArgs)
	if limit_sql == "" {
		limit_sql = fmt.Sprintf(SQL_STATEMENT_LIMIT_1, METH_COMPLETE_DEF_COUNT)
	}
	
	v := reflect.Indirect(rfltArgs)
	
	v_ic := int(GetIntArgValByName(rfltArgs, "Ic", 1))
	v_mid := int(GetIntArgValByName(rfltArgs, "Mid", 0))	
	
	pattern := ""
	if v_mid == 1 {
		pattern = "'%'||"
	}
	if v_ic == 1 {
		pattern += "lower($1)||'%'"
	}else{
		pattern += "$1||'%'"
	}	
	
	cond_sql := ""	
	cond_vals := make([]interface{},1)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if match_f, ok := t.Field(i).Tag.Lookup("matchField"); ok && match_f=="true" {
			if field_id, ok := t.Field(i).Tag.Lookup("json"); ok && match_f=="true" {
				cond_sql = "coalesce("+field_id+",'') LIKE "+pattern
				cond_vals[0] = GetTextArgValByName(rfltArgs, t.Field(i).Name, "")
				break
			}
		}
	}
	if cond_sql == "" {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, ER_NO_WHERE)
	}
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s %s %s`, scanModelMD.GetFieldList(), scanModelMD.DbRelation.Relation, cond_sql, orderby_sql, limit_sql)
//fmt.Println("CompleteOnArgsWithConn query=",query, cond_vals[0])		
	return AddQueryResult(resp, scanModelMD.DbRelation.ID, scanModel, query, "", cond_vals, conn, false)
}

func PublishEventsWithKeys(sockID string, keys map[string]interface{}, app Applicationer, pm PublicMethod) {
	//events
	params := make(map[string]interface{})
	params["emitterId"] = sockID
	params["keys"] = keys
	PublishPublicMethodEvents(app, pm, params)		
}

//Generates MD5 hash
func GetMd5(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

//Generates unique identifier
func GenUniqID(maxLen int) string{
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, maxLen)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)	
}

//Helper function to get value as int64 from argument by name
func GetIntArgValByName(args reflect.Value, fieldName string, defVal int64) int64 {
	var v reflect.Value
	if (args.Kind() == reflect.Struct) {
		v = args
	}else{
		v = reflect.Indirect(args)	
	}
	val := defVal
	arg_fld := v.FieldByName(fieldName)
	if arg_fld.Kind() != reflect.Invalid && !arg_fld.IsZero(){
		if fld_v, ok := arg_fld.Interface().(fields.ValInt); ok && fld_v.IsSet {
			val = fld_v.GetValue()
		}	
	}
	return val
}

//Helper function to get value as int64 from argument by name
func GetTextArgValByName(args reflect.Value, fieldName string, defVal string) string {
	var v reflect.Value
	if (args.Kind() == reflect.Struct) {
		v = args
	}else{
		v = reflect.Indirect(args)	
	}
	val := defVal
	arg_fld := v.FieldByName(fieldName)
	if arg_fld.Kind() != reflect.Invalid && !arg_fld.IsZero(){
		if fld_v, ok := arg_fld.Interface().(fields.ValText); ok && fld_v.IsSet {
			val = fld_v.GetValue()
		}	
	}
	return val
}

//Helper function. Returns field separator of a condition query
func ArgsFieldSep(rfltArgs reflect.Value) string {
	return GetTextArgValByName(rfltArgs, "Field_sep", DEF_FIELD_SEP)
}
/*
func dump(data interface{}){
    b,_:=json.MarshalIndent(data, "", "  ")
    fmt.Print(string(b))
}
*/
