package osbe

import (
	"fmt"
	"reflect"
	"context"
	"database/sql/driver"
	"encoding/xml"
	"os"
	"io/ioutil"
	"os/exec"
	
	"osbe/response"
	"osbe/model"
	"ds/pgds"
	
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type XMLer interface {
	Value() (driver.Value, error)
	GetIsNull() bool
}

func GetExcelFileOnArgs(app Applicationer, resp *response.Response, rfltArgs reflect.Value, scanModelMD *model.ModelMD, scanModel interface{}, clientFilName string) error {
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	defer d_store.Release(pool_conn, conn_id)
	pool_conn, conn_id, err := d_store.GetSecondary("")
	if err != nil {
		return err
	}
	conn := pool_conn.Conn()

	query, query_tot, where_params, _, _, err := GetListQuery(conn, rfltArgs, scanModelMD, nil, app.GetEncryptKey())
	if err != nil {
		return err
	}
	
	model_xml, err := QueryResultToXML(scanModelMD.ID, scanModel, query, query_tot, where_params, conn)
	if err != nil {
		return err
	}
	xml_s := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><document>%s</document>`, model_xml)

	xml_file, err := ioutil.TempFile("", "excel")
	if err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("GetExcelFile ioutil.TempFile: %v", err))
	}
    	if _, err := xml_file.WriteString(xml_s); err != nil {
    		return nil
    	}	
	xml_file.Close()
	defer os.Remove(xml_file.Name())
fmt.Println("xml_file.Name()=",xml_file.Name())
	out_file := xml_file.Name() + ".out"
	cmd := exec.Command("xalan", "-in", xml_file.Name(), "-xsl", app.GetConfig().GetXSLTDir()+"/ModelsToExcel.xls.xsl", "-out", out_file)
	err = cmd.Run()
	if err != nil { 
		return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("GetExcelFile exec.Command: %v", err))
	}
	_, err = os.Stat(out_file)
	if err != nil && os.IsNotExist(err) {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, "GetExcelFile exec.Command: file is missing")
		
	}else if err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL,  fmt.Sprintf("GetExcelFile os.Stat: %v", err))
	}
	defer os.Remove(out_file)	
	
	/*
	if _, err := app.XSLTransform(nil, xml_file.Name(), app.GetConfig().GetXSLTDir()+"/ModelsToExcel.xls.xsl", out_file); err != nil {
		return err
	}
	defer os.Remove(out_file)
	*/
	
fmt.Println("out_file=",out_file)	
	file_m, err := model.NewFileModelFromFile(&model.TFile{Name: clientFilName}, out_file)
	if err != nil {
		return NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("GetExcelFile NewFileModelFromFile: %v", err))
	}
	resp.AddModel(file_m)
	
	return nil
}

func QueryResultToXML(modelID string, scanModel interface{}, query string, queryTotal string, whereParams []interface{}, conn *pgx.Conn) (string, error) {
	//tot
	tot_cnt := 0 
	if queryTotal != "" {
		row_tot := conn.QueryRow(context.Background(), queryTotal, whereParams...)					
		if err := row_tot.Scan(&tot_cnt); err != nil && err != pgx.ErrNoRows {
			return "", NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("QueryResultToXML total pgx.Rows.Scan(): %v",err))	
		}
	}	
	//var rows pgx.Rows
	rows, err := conn.Query(context.Background(), query, whereParams...)	
	if err != nil {
		return "", NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("QueryResultToXML pgx.Conn.Query(): %v",err))
	}
		
	row_val := reflect.ValueOf(scanModel).Elem()
	row_t := row_val.Type()
	fld_ids := make([]string, row_val.NumField())
	for i := 0; i < row_val.NumField(); i++ {
		if field_id, ok := row_t.Field(i).Tag.Lookup("json"); ok {
			fld_ids[i] = field_id
		}
	}
	
	xml_rows := ""
	//`<rows>`
	r_cnt := 0
	for rows.Next() {
		row := reflect.New(reflect.ValueOf(scanModel).Elem().Type()).Interface().(model.ModelRow)
		row_val := reflect.ValueOf(row).Elem()
		row_fields := make([]interface{}, row_val.NumField())	
		for i := 0; i < row_val.NumField(); i++ {
			value_field := row_val.Field(i)
			row_fields[i] = value_field.Addr().Interface()			
		}
	
		if err := rows.Scan(row_fields...); err != nil {		
			return "", NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("QueryResultToXML pgx.Rows.Scan(): %v",err))	
		}
		xml_rows += `<row xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">`
		for i, f_val := range row_fields {
			if f_val_i,ok := f_val.(XMLer); ok {
				if f_val_i.GetIsNull() {
					xml_rows += fmt.Sprintf(`<%s  xsi:nil="true"/>`, fld_ids[i])
					
				}else{
					fld_val, err := f_val_i.Value()
					if err != nil {
						return "", NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("QueryResultToXML Value(): %v",err))	
					}
					xml_v, err := xml.Marshal(fld_val)
					if err != nil {
						//return "", NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("QueryResultToXML xml.Marshal(): %v",err))	
						xml_v = []byte(err.Error())
					}
					xml_rows += fmt.Sprintf(`<%s>%s</%s>`, fld_ids[i], string(xml_v), fld_ids[i])
				}				
			}
		}
		xml_rows += `</row>`
		r_cnt++		
	}
	//xml_rows += `</rows>`	
	
	if err := rows.Err(); err != nil {
		return "", NewPublicMethodError(response.RESP_ER_INTERNAL, err.Error())
	}
	
	if queryTotal == "" {
		tot_cnt = r_cnt
	}
	
	return fmt.Sprintf(`<model id="%s" totalCount="%d">%s</model>`, modelID, tot_cnt, xml_rows),
		nil
}


