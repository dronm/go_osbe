package osbe

import(
	"reflect"
	"strconv"
	"fmt"	
	"context"
	"math/rand"
	"crypto/md5"
	"encoding/hex"		
	"math"
	
	"ds/pgds"	
	"osbe/fields"	
	"osbe/response"
	"osbe/model"
	"osbe/socket"
	"osbe/sql"
	
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgtype"
)

const (
	KEY_FLD_PREF = "old_"
	KEY_FLD_PREF_LEN = 4
	
	LSN_FIELD = "lsn"
	
	METH_COMPLETE_DEF_COUNT = 50
	DEF_GET_LIST_LIMIT = 500

	RESP_ER_DELETE_CONSTR_VIOL = 500
	RESP_ER_WRITE_CONSTR_VIOL = 600
)

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
func ArgsToInsertParams(rfltArgs reflect.Value, presetConds sql.FilterCondCollection, encryptKey string) (fieldIds string, fieldArgs string, fieldValues []interface{}) {	
	rfltArgs_o := reflect.Indirect(rfltArgs)
	arg_tp := rfltArgs_o.Type()
	fieldValues = make([]interface{}, 0)
	field_ind := 0
	for i := 0; i < rfltArgs_o.NumField(); i++ {						
		if fld_v, ok := rfltArgs_o.Field(i).Interface().(fields.ValExt); ok {
			if field_id, ok := arg_tp.Field(i).Tag.Lookup("json"); ok {
				var fld_add bool
				var fld_val interface{}
				
				//check preset value
				if presetConds != nil {
					for _, pr_f := range presetConds {
						if pr_f.FieldID == field_id {
							fld_val = pr_f.Value
							fld_add = true
							break
						}
					}
				}
				
				if !fld_add && fld_v.GetIsSet() {
					fld_val = fld_v
					fld_add = true
					
				}
				/*else if is_autoInc, ok := arg_tp.Field(i).Tag.Lookup("autoInc"); ok && is_autoInc=="true" {
					//add anyway with NULL
					//fld_add = true
					//fld_val = "DEFAULT"
					//does not work this way 
				}*/
				
				if fld_add {
					if fieldIds != "" {
						fieldIds += ","
						fieldArgs += ","
					}
					fieldIds += field_id
					
					fld_encr := false
					field_arg_param := "$"+strconv.Itoa(field_ind + 1)
					
					if encryptKey != "" {
						if is_encrypted, ok := arg_tp.Field(i).Tag.Lookup("encrypted"); ok && is_encrypted == "true" {
							fieldArgs += fmt.Sprintf(`PGP_SYM_ENCRYPT(%s, "%s")`, field_arg_param, encryptKey)
							fld_encr = true
						}
					}
					
					if !fld_encr {
						fieldArgs += field_arg_param
					}
					fieldValues = append(fieldValues, fld_val)
					field_ind++			
				}
			}
		}
	}
	return
}

//puts old_key to where query
//pg specific function
func ArgsToUpdateParams(rfltArgs reflect.Value, presetConds sql.FilterCondCollection) (fieldQuery string, whereQuery string, fieldValues []interface{}, keys map[string]interface{}) {
	fieldValues = make([]interface{}, 0)
	keys = make(map[string]interface{})
	
	rfltArgs_o := reflect.Indirect(rfltArgs)
	arg_tp := rfltArgs_o.Type()
		
	field_ind := 0
	for i := 0; i < rfltArgs_o.NumField(); i++ {						
		if fld_v, ok := rfltArgs_o.Field(i).Interface().(fields.ValExt); ok {
			if field_id, ok := arg_tp.Field(i).Tag.Lookup("json"); ok && (len(field_id)<=KEY_FLD_PREF_LEN || field_id[:KEY_FLD_PREF_LEN]!=KEY_FLD_PREF) {
				
				//check preset value				
				if presetConds != nil {
					fld_found := false
					for _, pr_f := range presetConds {
						if pr_f.FieldID == field_id {
							//if f_ext_v, ok := pr_f.Value.(fields.ValExt); ok {
							//	fld_v = f_ext_v
							//	break									
							//}
							//simple values not fields.ValExt!
							fld_found = true
							fieldValues = append(fieldValues, pr_f.Value)			
							if fld_v.GetIsSet() {
								if fieldQuery != "" {
									fieldQuery += ","
								}
								fieldQuery += field_id + "=$"+strconv.Itoa(field_ind+1)								
							}
							if whereQuery != "" {
								whereQuery += " AND "
							}
							whereQuery += field_id + "=$"+strconv.Itoa(field_ind+1)
							field_ind++			
							
							break			
						}
					}
					if fld_found {
						continue
					}					
				}
				
				if fld_v.GetIsSet() {
					if fieldQuery != "" {
						fieldQuery += ","
					}
					fieldQuery += field_id + "=$"+strconv.Itoa(field_ind+1)
					fieldValues = append(fieldValues, fld_v)			
					field_ind++			
				}
			}
		}
	}
	
	for i := 0; i < rfltArgs_o.NumField(); i++ {						
		if fld_v, ok := rfltArgs_o.Field(i).Interface().(fields.ValExt); ok && fld_v.GetIsSet() {
			if field_id, ok := arg_tp.Field(i).Tag.Lookup("json"); ok && len(field_id)>KEY_FLD_PREF_LEN && field_id[:KEY_FLD_PREF_LEN]==KEY_FLD_PREF {
				//check preset value				
				if presetConds != nil {
					fld_found := false
					for _, pr_f := range presetConds {
						if pr_f.FieldID == field_id {
							//simple values not fields.ValExt!
							//ALWAYS!
							fld_found = true
							
							break
						}
					}
					if fld_found {
						continue
					}					
				}
				
				if fld_v.GetIsSet() {
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
	}
	return
}

//Implements controller insert method
//internally calls UpdateOnArgsWithConn
func UpdateOnArgs(app Applicationer, pm PublicMethod, resp *response.Response, sock socket.ClientSocketer, rfltArgs reflect.Value, modelMD *model.ModelMD, presetConds sql.FilterCondCollection) error {
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err := d_store.GetPrimary()
	if err != nil {
		return err
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
	
	return UpdateOnArgsWithConn(conn, app, pm, resp, sock, rfltArgs, modelMD, presetConds)
}

//Implements controller insert method
func UpdateOnArgsWithConn(conn *pgx.Conn, app Applicationer, pm PublicMethod, resp *response.Response, sock socket.ClientSocketer, rfltArgs reflect.Value, modelMD *model.ModelMD, presetConds sql.FilterCondCollection) error {
	f_query, w_query, f_values, keys := ArgsToUpdateParams(rfltArgs, presetConds)		
	if f_query == "" || w_query == "" {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, ER_UPDATE_EMPTY)
	}
	q := fmt.Sprintf("UPDATE %s SET %s WHERE %s", modelMD.Relation, f_query, w_query)
//fmt.Println("Update query=", q, "params:", f_values)	
	par, err := conn.Exec(context.Background(), q, f_values...)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23514" {
			//custom error
			return NewPublicMethodError(RESP_ER_WRITE_CONSTR_VIOL, ER_WRITE_CONSTR_VIOL)
		}else{
			return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("UpdateOnArgsWithConn pgx.Conn.Exec(): %v",err))
		}
	}
	
	lsn := GetDbLsn(conn)	
	resp.AddModel(model.New_MethodResult_Model(par.RowsAffected(), lsn))	
	//events
	PublishEventsWithKeys(sock.GetToken(), keys, app, pm, lsn)
	
	return nil
}

//Implements controller insert method
//internally calls InsertOnArgsWithConn
func InsertOnArgs(app Applicationer, pm PublicMethod, resp *response.Response, sock socket.ClientSocketer, rfltArgs reflect.Value, modelMD *model.ModelMD, retModel interface{}, presetConds sql.FilterCondCollection) error {
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err := d_store.GetPrimary()
	if err != nil {
		return err
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
	
	return InsertOnArgsWithConn(conn, app, pm, resp, sock, rfltArgs, modelMD, retModel, presetConds)
}

//Implements controller insert method
func InsertOnArgsWithConn(conn *pgx.Conn, app Applicationer, pm PublicMethod, resp *response.Response, sock socket.ClientSocketer, rfltArgs reflect.Value, modelMD *model.ModelMD, retModel interface{}, presetConds sql.FilterCondCollection) error {
	field_ids, field_args, f_values := ArgsToInsertParams(rfltArgs, presetConds, app.GetEncryptKey())		
	
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
		q += fmt.Sprintf("INSERT INTO %s DEFAULT VALUES RETURNING %s", modelMD.Relation, ret_field_ids)
	}else{
		q += fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s", modelMD.Relation, field_ids, field_args, ret_field_ids)
	}
//fmt.Println("InsertOnArgs q=",q, "field_values=%v", f_values)	

	if err := conn.QueryRow(context.Background(), q, f_values...).Scan(row_fields...); err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23514" {
			//custom error
			return NewPublicMethodError(RESP_ER_WRITE_CONSTR_VIOL, ER_WRITE_CONSTR_VIOL)
		}else{
			return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("InsertOnArgsWithConn pgx.Conn.QueryRow(): %v",err))
		}		
	}
	
	m := model.New_InsertedKey_Model(retModel)
	resp.AddModel(m)
	/*
	rows, err := conn.Query(context.Background(), q, f_values...)
	if err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Query(): %v",err))
	}	

	if rows.Next() {		
		if err := rows.Scan(row_fields...); err != nil {		
			return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Rows.Scan(): %v",err))	
		}
		m := model.New_InsertedKey_Model(retModel)
		resp.AddModel(m)
	}
	if err := rows.Err(); err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, err.Error())
	}
	*/
	
	//events
	i:= 0
	for key,_ := range keys {
		keys[key] = row_fields[i]
		i++
	}
	lsn := GetDbLsn(conn)
	resp.AddModel(model.New_MethodResult_Model(1, lsn))
	PublishEventsWithKeys(sock.GetToken(), keys, app, pm, lsn)
	
	return nil
}

//Common function for deleting object from DB based on argument keys
func DeleteOnArgKeys(app Applicationer, pm PublicMethod, resp *response.Response, sock socket.ClientSocketer, rfltArgs reflect.Value, modelMD *model.ModelMD, presetConds sql.FilterCondCollection) error {
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err := d_store.GetPrimary()
	if err != nil {
		return err
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()

	return DeleteOnArgKeysWithConn(conn, app, pm, resp, sock, rfltArgs, modelMD, presetConds)
}

//Implements controller delete method
func DeleteOnArgKeysWithConn(conn *pgx.Conn, app Applicationer, pm PublicMethod, resp *response.Response, sock socket.ClientSocketer, rfltArgs reflect.Value, modelMD *model.ModelMD, presetConds sql.FilterCondCollection) error {
	
	rfltArgs_o := reflect.Indirect(rfltArgs)
	arg_tp := rfltArgs_o.Type()
	
	f_values := make([]interface{}, 0) //arg_tp.NumField()
	keys := make(map[string]interface{})
	
	ids_key := ""
	where_sql := ""
	field_ind := 0
	
	//add all preset values to delete condition
	var added_fields map[string]bool
	if presetConds != nil {
		added_fields := make(map[string]bool, len(presetConds))
		for _, pr_f := range presetConds {
			if where_sql != "" {
				where_sql += " AND "
			}
			if pr_f.FieldID != "" {
				where_sql += pr_f.FieldID + " = $"+strconv.Itoa(field_ind+1)		
				f_values = append(f_values, pr_f.Value)
				field_ind++
			}else if  pr_f.Expression !="" {
				//expression
				where_sql += pr_f.Expression
			}
			added_fields[pr_f.FieldID] = true			
		}
	}
	
	for i := 0; i < rfltArgs_o.NumField(); i++ {						
		if fld_v, ok := rfltArgs_o.Field(i).Interface().(fields.ValExt); ok && fld_v.GetIsSet() {
			if field_id, ok := arg_tp.Field(i).Tag.Lookup("json"); ok {
				if _, ok := added_fields[field_id]; ok {
					//added already
					continue
				}
				if where_sql != "" {
					where_sql += " AND "
				}
				where_sql += field_id + " = $"+strconv.Itoa(field_ind+1)
				ids_key += "_"+field_id
			
				f_values = append(f_values,fld_v)
				keys[field_id],_ = fld_v.Value()				
				field_ind++			
			}
		}
	}
	if where_sql == "" {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, ER_NO_KEYS)
	}
	
	//might exists different where clauses due to preset filters
	//query can be prepared in case of presetConds = nil
	q_id := ""
	q := fmt.Sprintf(`DELETE FROM %s WHERE %s`, modelMD.Relation, where_sql)
	if presetConds == nil {
		q_id = modelMD.Relation + ids_key + "_delete"
		_, err := conn.Prepare(context.Background(), q_id, q)
		if err != nil {
			return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Prepare(): %v",err))
		}
	}else{
		//got preset canditions, cannot preapre
		q_id = q
	}
//fmt.Println("DeleteOnArgKeys q=", q, "f_values=", f_values)				
	par, err := conn.Exec(context.Background(), q_id, f_values...)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23503" {
			//custom error
			return NewPublicMethodError(RESP_ER_DELETE_CONSTR_VIOL, ER_DELETE_CONSTR_VIOL)
		}else{
			return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Exec(): %v",err))
		}
	}
	
	lsn := GetDbLsn(conn)
	resp.AddModel(model.New_MethodResult_Model(par.RowsAffected(), lsn))
	//events
	PublishEventsWithKeys(sock.GetToken(), keys, app, pm, lsn)
		
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

func GetModelLsnValue(modelStruct interface{}) string {
	return GetTextArgValByName(reflect.ValueOf(modelStruct), LSN_FIELD, "")
}

//Implements controller get_object method
//rfltArgs holds condition fields
func GetObjectOnArgs(app Applicationer, resp *response.Response, rfltArgs reflect.Value, modelMD *model.ModelMD, modelStruct interface{}, presetConds sql.FilterCondCollection) error {		
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err := d_store.GetSecondary(GetModelLsnValue(modelStruct))
	if err != nil {
		return err
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
	
	//fields with key values
	field_vals := make([]interface{}, 0)	
	cond_sql := ""
	rfltArgs_o := reflect.Indirect(rfltArgs)
	arg_tp := rfltArgs_o.Type()		
	cond_ind := 0
	for i := 0; i < rfltArgs_o.NumField(); i++ {						
		if fld_v, ok := rfltArgs_o.Field(i).Interface().(fields.ValExt); ok && fld_v.GetIsSet() {
			if field_id, ok := arg_tp.Field(i).Tag.Lookup("json"); ok {
				if cond_sql != "" {
					cond_sql += " AND "
				}
				cond_sql += field_id + "=$"+strconv.Itoa(cond_ind+1)
				field_vals = append(field_vals, fld_v)			
				cond_ind++			
			}
		}
	}
	if len(field_vals) == 0 {
		//should not happen if keys are marked as required in get object model
		//return NewPublicMethodError(response.RESP_ER_INTERNAL, ER_NO_KEYS)
		//happens when http requests insert with get_object without key
		return nil
	}
	
	relation := modelMD.Relation
	query_id := relation + "_get_object"
	sql_s := fmt.Sprintf("SELECT %s FROM %s WHERE %s", modelMD.GetFieldList(app.GetEncryptKey()), relation, cond_sql)
//fmt.Println("conn.Prepare", sql_s)
	_, err = conn.Prepare(context.Background(), query_id, sql_s)
	if err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Prepare(): %v",err))
	}
//fmt.Println("conn.Query", query_id)	
	rows, err := conn.Query(context.Background(), query_id, field_vals...)	
	if err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Query(): %v",err))
	}

	m := &model.Model{ID: model.ModelID(modelMD.ID), Rows: make([]model.ModelRow, 0)}	
	for rows.Next() {
		row := modelStruct
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
		return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Rows.Next(): %v",err))
	}
	
	resp.AddModel(m)
	return nil
}

//Returns: query, total query and cond_params
//string - query, string - total query, []interface{} - condition params
//from int, count int
func GetListQuery(conn *pgx.Conn, rfltArgs reflect.Value, scanModelMD *model.ModelMD, extraConds sql.FilterCondCollection, encryptKey string) (string, string, []interface{}, int, int, error) {
	f_sep := ArgsFieldSep(rfltArgs)
	orderby_sql := GetSQLOrderByFromArgsOrDefault(rfltArgs, f_sep, scanModelMD, encryptKey)	
	limit_sql, from, cnt, err := GetSQLLimitFromArgs(rfltArgs, scanModelMD, conn, DEF_GET_LIST_LIMIT)
	if err != nil {
		return "", "", nil, 0, 0, NewPublicMethodError(response.RESP_ER_INTERNAL, err.Error())
	}
//fmt.Println("GetListQuery limit_sql=",limit_sql)		
//fmt.Println("GetListQuery GetSQLWhereFromArgs")	
	where_sql, cond_params, err := GetSQLWhereFromArgs(rfltArgs, f_sep, scanModelMD, extraConds)
	if err != nil {
		return "", "", nil, 0, 0, NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("%v",err))
	}

	relation := scanModelMD.Relation
	query_tmpl := fmt.Sprintf("SELECT %s FROM %s", scanModelMD.GetFieldList(encryptKey), relation)
	
	query_tot_tmpl := ""
	if scanModelMD.AggFunctions != nil {
		tot_expr := ""
		for _,agg_f := range scanModelMD.AggFunctions {
			if tot_expr != "" {
				tot_expr += ", "
			}
			tot_expr += fmt.Sprintf("%s AS %s", agg_f.Expr, agg_f.Alias)
		}
		query_tot_tmpl = fmt.Sprintf("SELECT %s FROM %s", tot_expr, relation)
	}
	
	query := ""
	query_tot := ""
	if orderby_sql == "" && limit_sql == "" && where_sql == "" {
		query = relation+"_get_list"
		_, err = conn.Prepare(context.Background(), query, query_tmpl)
		if err != nil {
			return "", "", nil, 0, 0, NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Prepare(): %v",err))
		}		
		
		if query_tot_tmpl != "" {
			query_tot = relation+"_get_list_tot"
			_, err = conn.Prepare(context.Background(), query_tot, query_tot_tmpl)
			if err != nil {
				return "", "", nil, 0, 0, NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Prepare(): %v",err))
			}		
		}		
		
	}else{
		//custom query
		query = query_tmpl		
		if where_sql != "" {
			query += " "+where_sql
			if query_tot_tmpl != "" {
				query_tot_tmpl += " "+where_sql
			}
		}
		if query_tot_tmpl != "" {
			query_tot = query_tot_tmpl
		}
				
		if orderby_sql != "" {
			query += " "+orderby_sql
		}		
		if limit_sql != "" {
			query += " "+limit_sql
		}		
	}
//fmt.Println("GetListQuery", scanModelMD.ID, "Q=", query, "query_tot=", query_tot)
//fmt.Println("cond_params=", cond_params)		
	return query, query_tot, cond_params, from , cnt, nil
}

//Executes query and returns it result as model
func QueryResultToModel(modelMD *model.ModelMD, modelStruct interface{}, query string, queryTotal string, condValues []interface{}, conn *pgx.Conn, sysModel bool) (model.Modeler, error) {
	var agg_values []*model.AggFunctionValue
	if queryTotal != "" && modelMD.AggFunctions != nil {
		agg_values = make([]*model.AggFunctionValue, len(modelMD.AggFunctions))
		totals := make([]interface{}, len(modelMD.AggFunctions))
		for i, agg_f := range modelMD.AggFunctions {
			agg_v := &model.AggFunctionValue{Alias: agg_f.Alias}			
			totals[i] = &agg_v.Val
			agg_values[i] = agg_v
		}	
//fmt.Println("queryTotal=", queryTotal)			
		row_tot := conn.QueryRow(context.Background(), queryTotal, condValues...)					
		if err := row_tot.Scan(totals...); err != nil && err != pgx.ErrNoRows {
			return nil, NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("QueryResultToModel total pgx.Rows.Scan(): %v",err))	
		}
		for i,t_val := range agg_values {
			if t_num, ok := t_val.Val.(pgtype.Numeric); ok && !t_num.NaN {
				fl := float64(t_num.Int.Int64()) * math.Pow(10.0, float64(t_num.Exp))
				if float64(int64(fl)) == fl {
					agg_values[i].ValStr = fmt.Sprintf("%d", int64(fl))
				}else{
					agg_values[i].ValStr = fmt.Sprintf("%f", fl)
				}
				
			}else if t_num, ok := t_val.Val.(int64); ok {
				agg_values[i].ValStr = fmt.Sprintf("%d", t_num)
			}
		}
		
	}	

	rows, err := conn.Query(context.Background(), query, condValues...)	
	if err != nil {
		return nil, NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("QueryResultToModel pgx.Conn.Query(): %v, model: %s",err, modelMD.ID))
	}
	
	m := &model.Model{ID: model.ModelID(modelMD.ID), AggFunctionValues: agg_values, SysModel: sysModel, Rows: make([]model.ModelRow, 0)}	
	for rows.Next() {
		row := reflect.New(reflect.ValueOf(modelStruct).Elem().Type()).Interface().(model.ModelRow)
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
	//m.RowsPerPage???
	if err := rows.Err(); err != nil {
		return nil, NewPublicMethodError(response.RESP_ER_INTERNAL, err.Error())
	}
	/*if queryTotal == "" {
		m.TotalCount = len(m.Rows)
	}*/
	
	return m, nil
}

func AddQueryResult(resp *response.Response, modelMD *model.ModelMD, modelStruct interface{}, query string, queryTotal string, condValues []interface{}, conn *pgx.Conn, sysModel bool) error {
	m, err := QueryResultToModel(modelMD, modelStruct, query, queryTotal, condValues, conn, sysModel)
	if err != nil {
		return err
	}
	resp.AddModel(m)
	return nil
}

func GetListOnArgs(app Applicationer, resp *response.Response, rfltArgs reflect.Value, modelMD *model.ModelMD, modelStruct interface{}, presetConds sql.FilterCondCollection) error {
	if modelMD == nil {
		app.GetLogger().Error("osbe.GetListOnArgs (util.go) modelMD not defined. Potentially error prone code!")
	}
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err := d_store.GetSecondary(GetModelLsnValue(modelStruct))
	if err != nil {
		return err
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()

	query, query_tot, where_params, from,cnt, err := GetListQuery(conn, rfltArgs, modelMD, presetConds, app.GetEncryptKey())
	if err != nil {
		return err
	}
	
	//return AddQueryResult(resp, modelMD, modelStruct, query, query_tot, where_params, conn, false)		
	m, err := QueryResultToModel(modelMD, modelStruct, query, query_tot, where_params, conn, false)
	if err != nil {
		return err
	}
	m.SetRowsPerPage(cnt)
	m.SetListFrom(from)
	resp.AddModel(m)
	return nil
	
}


//Implements controller complete method
//Internally calls CompleteOnArgsWithConn
func CompleteOnArgs(app Applicationer, resp *response.Response, rfltArgs reflect.Value, scanModelMD *model.ModelMD, scanModel interface{}, presetConds sql.FilterCondCollection) error {	
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err := d_store.GetSecondary(GetModelLsnValue(scanModel))
	if err != nil {
		return err
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
	return CompleteOnArgsWithConn(conn, app, resp, rfltArgs, scanModelMD, scanModel, presetConds)
}

//Implements controller complete method
//args.Ic - insensetive case 1/0
//args.Mid 1 - %pattern%, 0 - pattern%
//scanModelMD
//scanModel
//pattern - pattern to match
//there is also another argument with the same name as model field marked as matchField=true in tag
func CompleteOnArgsWithConn(conn *pgx.Conn, app Applicationer, resp *response.Response, rfltArgs reflect.Value, scanModelMD *model.ModelMD, scanModel interface{}, presetConds sql.FilterCondCollection) error {
	f_sep := ArgsFieldSep(rfltArgs)
	orderby_sql := GetSQLOrderByFromArgsOrDefault(rfltArgs, f_sep, scanModelMD, app.GetEncryptKey())	
	limit_sql, _, _, err := GetSQLLimitFromArgs(rfltArgs, scanModelMD, conn, METH_COMPLETE_DEF_COUNT)
	if err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, err.Error())
	}
	/*if limit_sql == "" {
		limit_sql = fmt.Sprintf(SQL_STATEMENT_LIMIT_1, METH_COMPLETE_DEF_COUNT)
	}*/
	
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
			if field_id, ok := t.Field(i).Tag.Lookup("json"); ok {
				if v_ic == 1 {
					cond_sql = "coalesce(lower("+field_id+"),'') LIKE "+pattern
				}else{
					cond_sql = "coalesce("+field_id+",'') LIKE "+pattern
				}
				cond_vals[0] = GetTextArgValByName(rfltArgs, t.Field(i).Name, "")
				break
			}
		}
	}
	if cond_sql == "" {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, ER_NO_WHERE)
	}
	
	//add all preset values to condition
	if presetConds != nil {
		field_ind := 1
		for _, pr_f := range presetConds {
			if cond_sql != "" {
				cond_sql += " AND "
			}
			if pr_f.FieldID != "" {
				cond_sql += pr_f.FieldID + " = $"+strconv.Itoa(field_ind+1)		
				cond_vals = append(cond_vals, pr_f.Value)
				field_ind++
			}else if  pr_f.Expression !="" {
				//expression
				cond_sql += pr_f.Expression
			}
		}
	}
	
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s %s %s`, scanModelMD.GetFieldList(app.GetEncryptKey()), scanModelMD.Relation, cond_sql, orderby_sql, limit_sql)
//fmt.Println("CompleteOnArgsWithConn query=",query, cond_vals[0])		
	return AddQueryResult(resp, scanModelMD, scanModel, query, "", cond_vals, conn, false)
}

func PublishEventsWithKeys(sockID string, keys map[string]interface{}, app Applicationer, pm PublicMethod, lsn string) {
	//events
	params := make(map[string]interface{})
	params["emitterId"] = sockID
	params["keys"] = keys
	if lsn != "" {
		params[LSN_FIELD] = lsn
	}
	app.PublishPublicMethodEvents(pm, params)		
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

//Helper function to get value as sring from argument by name
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

func AddStructFieldsToList(v reflect.Value, fields *[]interface{}, fieldIDs *string) error {	
	//v := reflect.ValueOf(str)
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		if v.IsNil() {
			break
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct{
		return nil
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Anonymous {
			if err := AddStructFieldsToList(v.Field(i), fields, fieldIDs); err != nil {
				return err
			}
			
		}else if sql, ok := t.Field(i).Tag.Lookup("sql"); !ok || sql != "false" {			
			if field_id, ok := t.Field(i).Tag.Lookup("json"); ok {
				value_field := v.Field(i)
				*fields = append(*fields, value_field.Addr().Interface())
				if *fieldIDs != "" {
					*fieldIDs += ","
				}
				*fieldIDs += field_id
			}
		}
	}
	return nil
}

//Returns:
//	struct fields,
//	list of field IDs for select query
//	error if any
func MakeStructRowFields(resultStruct interface{}) ([]interface{}, string, error){	
	fields := make([]interface{}, 0)
	field_ids := ""
	AddStructFieldsToList(reflect.ValueOf(resultStruct), &fields, &field_ids)
	return fields, field_ids, nil
}

/*
func dump(data interface{}){
    b,_:=json.MarshalIndent(data, "", "  ")
    fmt.Print(string(b))
}
*/

func GetDbLsn(conn *pgx.Conn) string {
	lsn := ""
	conn.QueryRow(context.Background(), "SELECT pg_current_wal_lsn()").Scan(&lsn)
	return lsn
}


