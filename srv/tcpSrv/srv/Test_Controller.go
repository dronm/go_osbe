package main

import (
	"fmt"
	"errors"
	"encoding/json"
	"reflect"
	"context"
	
	"osbe"
	"osbe/fields"
	"osbe/srv"
	"osbe/socket"
	"osbe/model"
	"osbe/response"
	
	"github.com/jackc/pgx/v4"
	///pgxpool
)


//Test_Controller->insert
//Test_Controller->delete

//Controller
type Test_Controller struct {
	PublicMethods osbe.PublicMethodCollection	
}

func (c *Test_Controller) GetID() osbe.ControllerID {
	return osbe.ControllerID("Test")
}

func (c *Test_Controller) InitPublicMethods() {
	c.PublicMethods = make(osbe.PublicMethodCollection)
	
	//************************** method insert **********************************
	c.PublicMethods["insert"] = &Test_Controller_insert{
		ModelMetadata: Get_Test_obj_md(),
		EventList: make(osbe.PublicMethodEventList,1),
	}
	c.PublicMethods["insert"].AddEvent("Test.insert")
	
	//************************** method delete *************************************
	c.PublicMethods["delete"] = &Test_Controller_delete{
		ModelMetadata: Get_Test_obj_md(),
		EventList: make(osbe.PublicMethodEventList,1),
	}
	c.PublicMethods["delete"].AddEvent("Test.delete")

	//************************** method update *************************************
	c.PublicMethods["update"] = &Test_Controller_update{
		ModelMetadata: Get_Test_old_keys_md(),
		EventList: make(osbe.PublicMethodEventList,1),
	}
	c.PublicMethods["update"].AddEvent("Test.update")
	
	
	//************************** method get_object *************************************
	c.PublicMethods["get_object"] = &Test_Controller_get_object{
		ModelMetadata: Get_Test_keys_md(),
	}

	//************************** method get_list *************************************
	c.PublicMethods["get_list"] = &Test_Controller_get_list{
		ModelMetadata: Get_Test_cond_md(),
	}
	
}

func (c *Test_Controller) GetPublicMethod(publicMethodID osbe.PublicMethodID) (pm osbe.PublicMethod, err error) {
	pm, ok := c.PublicMethods[publicMethodID]
	if !ok {
		err = errors.New(fmt.Sprintf(osbe.ER_CONTOLLER_METH_NOT_DEFINED, string(publicMethodID), string(c.GetID())))
	}
	
	return
}

//**************************************************************************************
//Public method: insert
type Test_Controller_insert_argv struct {
	Argv Test_obj `json:"argv"`	
}

type Test_Controller_insert struct {
	ModelMetadata fields.FieldCollection
	EventList osbe.PublicMethodEventList
}

func (pm *Test_Controller_insert) GetEventList() osbe.PublicMethodEventList {
	return pm.EventList
}

func (pm *Test_Controller_insert) AddEvent(evId string) {
	pm.EventList[len(pm.EventList)-1] = evId
}

func (pm *Test_Controller_insert) GetModelMetadata() fields.FieldCollection {
	return pm.ModelMetadata
}

func (c *Test_Controller_insert) GetID() osbe.PublicMethodID {
	return osbe.PublicMethodID("insert")
}

//Public method Unmarshal to structure
func (pm *Test_Controller_insert) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &Test_Controller_insert_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}

	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}
//https://stackoverflow.com/questions/21986780/is-it-possible-to-retrieve-a-column-value-by-name-using-golang-database-sql
//Method implemenation
//http://localhost:59000/?c=Test_Controller&f=insert&v=xml&f1=579&f2=%D0%9D%D1%83%20%D0%BE%D0%BE%D1%87%D0%B5%D0%BD%D1%8C%20%D0%B4%D0%BB%D0%B8%D0%BD%D0%BD%D1%8B%D0%B9%20%D1%82%D0%B5%D0%BA%D1%81%D1%82&f3=385&f4=true
func (pm *Test_Controller_insert) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	args := rfltArgs.Interface().(Test_obj)
	
	pool,err := app.GetServerPool().GetPrimary()
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("db.ServerPool.GetPrimary(): %v",err))
	}
	
	pool_conn, err := pool.Pool.Acquire(context.Background())
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgxpool.Pool.Acquire(): %v",err))
	}
	defer pool_conn.Release()
	
	conn := pool_conn.Conn()
	var rows pgx.Rows
	//Preparing query
	//If all columns are set, except autoinc
	if !args.Id.IsSet && args.F1.IsSet && args.F2.IsSet && args.F3.IsSet&& args.F4.IsSet {
		_, err = conn.Prepare(context.Background(), "Test_Controller_insert", "INSERT INTO test (f1,f2,f3,f4) VALUES ($1, $2, $3, $4) RETURNING id")
		if err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Prepare(): %v",err))
		}
					
		rows, err = conn.Query(context.Background(), "Test_Controller_insert", args.F1, args.F2, args.F3, args.F4)
		if err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Query(): %v",err))
		}
	}else{
		field_ids,field_args,field_values := osbe.ArgsToInsertParams(rfltArgs)		
		q := fmt.Sprintf("INSERT INTO test (%s) VALUES (%s) RETURNING id", field_ids, field_args)
		//fmt.Println("Q=",q)
		rows, err = conn.Query(context.Background(), q, field_values...)
		if err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Query(): %v",err))
		}		
	}
	
	if rows.Next() {
		row := &Test_keys{}
		rows.Scan(&row.Id)
		resp.AddModel(model.New_InsertedId_Model(row))
	}
	
	//events
	osbe.PublishPublicMethodEvents(app, pm, sock.GetToken())
	
	return nil
}


//**************************************************************************************
//Public method: delete
type Test_Controller_delete_argv struct {
	Argv Test_keys `json:"argv"`	
}

type Test_Controller_delete struct {
	ModelMetadata fields.FieldCollection
	EventList osbe.PublicMethodEventList
}

func (pm *Test_Controller_delete) GetEventList() osbe.PublicMethodEventList {
	return pm.EventList
}

func (pm *Test_Controller_delete) AddEvent(evId string) {
	pm.EventList[len(pm.EventList)-1] = evId
}

func (pm *Test_Controller_delete) GetModelMetadata() fields.FieldCollection {
	return pm.ModelMetadata
}

func (c *Test_Controller_delete) GetID() osbe.PublicMethodID {
	return osbe.PublicMethodID("delete")
}

//Public method Unmarshal to structure
func (pm *Test_Controller_delete) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &Test_Controller_delete_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//Method implemenation
//http://localhost:59000/?c=Test_Controller&f=delete&id=10&v=json
func (pm *Test_Controller_delete) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	args := rfltArgs.Interface().(Test_keys)
	
	pool,err := app.GetServerPool().GetPrimary()
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("db.ServerPool.GetPrimary(): %v",err))
	}
	
	pool_conn, err := pool.Pool.Acquire(context.Background())
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgxpool.Pool.Acquire(): %v",err))
	}
	defer pool_conn.Release()
	
	conn := pool_conn.Conn()
	
	_, err = conn.Prepare(context.Background(), "Test_Controller_delete", "DELETE from test WHERE id = $1")
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Prepare(): %v",err))
	}
	
	par, err := conn.Exec(context.Background(), "Test_Controller_delete", args.Id.GetValue())
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Exec(): %v",err))
	}
	
	del_m := model.New_Deleted_Model(par.RowsAffected())
	
	resp.AddModel(del_m)
	
	//events
	osbe.PublishPublicMethodEvents(app, pm, sock.GetToken())
		
	return nil
}


//**************************************************************************************
//Public method: get_object
type Test_Controller_get_object_argv struct {
	Argv Test_keys `json:"argv"`	
}

type Test_Controller_get_object struct {
	ModelMetadata fields.FieldCollection
	EventList osbe.PublicMethodEventList
}

func (pm *Test_Controller_get_object) AddEvent(evId string) {
	pm.EventList[len(pm.EventList)-1] = evId
}

func (pm *Test_Controller_get_object) GetEventList() osbe.PublicMethodEventList {
	return pm.EventList
}

func (pm *Test_Controller_get_object) GetModelMetadata() fields.FieldCollection {
	return pm.ModelMetadata
}

func (c *Test_Controller_get_object) GetID() osbe.PublicMethodID {
	return osbe.PublicMethodID("get_object")
}

//Public method Unmarshal to structure
func (pm *Test_Controller_get_object) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &Test_Controller_get_object_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//Method implemenation
//http://localhost:59000/?c=Test_Controller&f=get_object&id=10&v=json
func (pm *Test_Controller_get_object) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	args := rfltArgs.Interface().(Test_keys)
	
	pool,err := app.GetServerPool().GetPrimary()
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("db.ServerPool.GetPrimary(): %v",err))
	}
	
	pool_conn, err := pool.Pool.Acquire(context.Background())
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgxpool.Pool.Acquire(): %v",err))
	}
	defer pool_conn.Release()
	
	conn := pool_conn.Conn()
	
	_, err = conn.Prepare(context.Background(), "Test_Controller_get_object", "SELECT id,f1,f2,f3,f4 FROM test WHERE id = $1")
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Prepare(): %v",err))
	}
	
	row := &Test_obj{}
	rows, err := conn.Query(context.Background(), "Test_Controller_get_object", args.Id)
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Query(): %v",err))
	}
	
	m := &model.Model{ID: "TestDialog"}
	if rows.Next() {
		if err := rows.Scan(&row.Id, &row.F1, &row.F2, &row.F3, &row.F4); err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Rows.Scan(): %v",err))	
		}
		m.Rows = make([]model.ModelRow, 1)		
		m.Rows[0] = row
	}else{
		m.Rows = make([]model.ModelRow, 0)
	}
	
	resp.AddModel(m)
	
	return nil
}

//**************************************************************************************
//Public method: get_list
type Test_Controller_get_list_argv struct {
	Argv Test_cond `json:"argv"`	
}

type Test_Controller_get_list struct {
	ModelMetadata fields.FieldCollection
	EventList osbe.PublicMethodEventList
}

func (pm *Test_Controller_get_list) AddEvent(evId string) {
	pm.EventList[len(pm.EventList)-1] = evId
}

func (pm *Test_Controller_get_list) GetEventList() osbe.PublicMethodEventList {
	return pm.EventList
}

func (pm *Test_Controller_get_list) GetModelMetadata() fields.FieldCollection {
	return pm.ModelMetadata
}

func (c *Test_Controller_get_list) GetID() osbe.PublicMethodID {
	return osbe.PublicMethodID("get_list")
}

//Public method Unmarshal to structure
func (pm *Test_Controller_get_list) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &Test_Controller_get_list_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//Method implemenation
//http://localhost:59000/?c=Test_Controller&f=get_list&v=json
func (pm *Test_Controller_get_list) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	//args := rfltArgs.Interface().(Test_cond)
	
	f_sep := osbe.ArgsFieldSep(rfltArgs)
	oredrby_sql := osbe.GetSQLOrderByFromArgs(rfltArgs, f_sep)
	
	if oredrby_sql != "" {
		fmt.Println("OrderBy=",oredrby_sql)
	}
	
	limit_sql := osbe.GetSQLLimitFromArgs(rfltArgs)
	if limit_sql != "" {
		fmt.Println("LIMIT=",limit_sql)
	}
	
	where_sql, where_params, err := osbe.GetSQLWhereFromArgs(rfltArgs, f_sep, Get_Test_obj_md())
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("%v",err))
	}
	if where_sql != "" {
		fmt.Println("where_sql=",where_sql)
	}
	
	fmt.Println("f_sep="+f_sep)
	
	pool,err := app.GetServerPool().GetPrimary()
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("db.ServerPool.GetPrimary(): %v",err))
	}
	
	pool_conn, err := pool.Pool.Acquire(context.Background())
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgxpool.Pool.Acquire(): %v",err))
	}
	defer pool_conn.Release()
	
	conn := pool_conn.Conn()
	
	query := ""
	//no params
	if oredrby_sql == "" && limit_sql == "" && where_sql == "" {
		query = "Test_Controller_get_list"
		_, err = conn.Prepare(context.Background(), query, "SELECT * FROM test ORDER BY id")
		if err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Prepare(): %v",err))
		}		
		
	}else{
		//custom query
		query = "SELECT * FROM test"
		if where_sql != "" {
			query += " "+where_sql
		}
		if oredrby_sql != "" {
			query += " "+oredrby_sql
		}		
		if limit_sql != "" {
			query += " "+limit_sql
		}		
	}
	fmt.Println("query=",query)
	var rows pgx.Rows
	//if where_params != nil && len(where_params) {
		rows, err = conn.Query(context.Background(), query, where_params...)
	//}else{
	//	rows, err := conn.Query(context.Background(), query)
	//}
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Query(): %v",err))
	}		
		
	m := &model.Model{ID: "TestList", Rows: make([]model.ModelRow, 0)}
	for rows.Next() {
		row := &Test_obj{}
		if err := rows.Scan(&row.Id, &row.F1, &row.F2, &row.F3, &row.F4); err != nil {		
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Rows.Scan(): %v",err))	
		}		
		m.Rows = append(m.Rows, row)
	}
	resp.AddModel(m)
	
	return nil
}

//**************************************************************************************
//Public method: update
type Test_Controller_update_argv struct {
	Argv Test_old_keys `json:"argv"`	
}

type Test_Controller_update struct {
	ModelMetadata fields.FieldCollection
	EventList osbe.PublicMethodEventList
}

func (pm *Test_Controller_update) GetEventList() osbe.PublicMethodEventList {
	return pm.EventList
}

func (pm *Test_Controller_update) AddEvent(evId string) {
	pm.EventList[len(pm.EventList)-1] = evId
}

func (pm *Test_Controller_update) GetModelMetadata() fields.FieldCollection {
	return pm.ModelMetadata
}

func (c *Test_Controller_update) GetID() osbe.PublicMethodID {
	return osbe.PublicMethodID("update")
}

//Public method Unmarshal to structure
func (pm *Test_Controller_update) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &Test_Controller_update_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//Method implemenation
//http://localhost:59000/?c=Test_Controller&f=update&old_id=10&v=json
func (pm *Test_Controller_update) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	args := rfltArgs.Interface().(Test_old_keys)
	
	pool,err := app.GetServerPool().GetPrimary()
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("db.ServerPool.GetPrimary(): %v",err))
	}
	
	pool_conn, err := pool.Pool.Acquire(context.Background())
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgxpool.Pool.Acquire(): %v",err))
	}
	defer pool_conn.Release()
	
	conn := pool_conn.Conn()
	
	var rows pgx.Rows	
	if !args.Id.IsSet && args.F1.IsSet && args.F2.IsSet && args.F3.IsSet&& args.F4.IsSet {
		_, err = conn.Prepare(context.Background(), "Test_Controller_update", "UPDATE test SET f1=$1, f2=$2, f3=$3, f4=$4 WHERE id = $5")
		if err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Prepare(): %v",err))
		}
					
		rows, err = conn.Query(context.Background(), "Test_Controller_update", args.F1, args.F2, args.F3, args.F4, args.Old_id)
		if err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Query(): %v",err))
		}
	}else{
		//custom column order
		f_query, w_query, field_values := osbe.ArgsToUpdateParams(rfltArgs)		
		q := fmt.Sprintf("UPDATE test SET %s WHERE %s", f_query, w_query)
		fmt.Println("UpdateQuery=",q)
		rows, err = conn.Query(context.Background(), q, field_values...)
		if err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Query(): %v",err))
		}		
	}
	
	if rows.Next() {
		row := &Test_keys{}
		rows.Scan(&row.Id)
		resp.AddModel(model.New_InsertedId_Model(row))
	}
	
	//events
	osbe.PublishPublicMethodEvents(app, pm, sock.GetToken())
	
	return nil
}


