package docAttachment

/**
 * This is Attachment controller implimentation file.
 *
 */

import (
	"reflect"	
	"context"
	"fmt"
	"os"
	"time"
	
	"github.com/dronm/ds/pgds"
	
	"osbe"
	"osbe/srv"
	"osbe/view"
	"osbe/socket"
	"osbe/response"	
	"osbe/srv/httpSrv"
	"osbe/fields"	
	
	"github.com/jackc/pgx/v4/pgxpool"
	//"github.com/jackc/pgx/v4"
)



//Method implemenation get_object
func (pm *Attachment_Controller_get_object) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["AttachmentList"], &AttachmentList{}, sock.GetPresetFilter("AttachmentList"))	
}

//Method implemenation get_list
func (pm *Attachment_Controller_get_list) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["AttachmentList"], &AttachmentList{}, sock.GetPresetFilter("AttachmentList"))	
}


//Method implemenation
func (pm *Attachment_Controller_delete_file) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err_с := d_store.GetPrimary()
	if err_с != nil {
		return err_с
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
	
	args := rfltArgs.Interface().(*Attachment_delete_file)
	
	return RemoveAttachment(conn, app.GetBaseDir(), args.Ref, args.Content_id.GetValue())
}

//Method implemenation
func (pm *Attachment_Controller_clear_cache) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	args := rfltArgs.Interface().(*Attachment_clear_cache)
	return ClearCache(app.GetBaseDir(), args.Ref, args.Content_id.GetValue())
}

//Method implemenation
func (pm *Attachment_Controller_get_file) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	//returns file
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err_с := d_store.GetPrimary()
	if err_с != nil {
		return err_с
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
	
	args := rfltArgs.Interface().(*Attachment_get_file)
	
	var att_id int64
	var att_name fields.ValText
	if err := conn.QueryRow(context.Background(),
		`SELECT
			id,
			content_info->>'name'
		FROM attachments
		WHERE ref->>'dataType' = $1
			AND (ref->'keys'->>'id')::int = $2
			AND content_info->>'id' = $3`,
		args.Ref.DataType,
		args.Ref.Keys.Id,
		args.Content_id,
	).Scan(&att_id, &att_name); err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("Attachment_Controller_get_file conn.QueryRow(): %v",err))	
	}
	
	var cont_disp httpSrv.CONTENT_DISPOSITION
	if args.Inline.GetValue() == 1 {
		cont_disp = httpSrv.CONTENT_DISPOSITION_INLINE
	}else{
		cont_disp = httpSrv.CONTENT_DISPOSITION_ATTACHMENT
	}
	
	var cache_f *os.File
	var err error
	cache_fn := GetAttachmentCacheFileName(app.GetBaseDir(), args.Ref.DataType, args.Ref.Keys.Id, args.Content_id.GetValue())
	if view.FileExists(cache_fn) {
		//from cache		
		cache_f, err = os.Open(cache_fn)
		if err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("Attachment_Controller_get_file os.Open(): %v", err))
		}
		defer cache_f.Close()		
		if err := httpSrv.DownloadFile(resp, sock, cache_f, att_name.GetValue(), "", cont_disp); err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("Attachment_Controller_get_file DownloadFile(): %v", err))
		}
		
	}else{
		//no cache, read from db && save
		var f_cont []byte//&fields.ValBytea{}
		if err := conn.QueryRow(context.Background(),
			`SELECT
				content_data
			FROM attachments
			WHERE id = $1`,
			att_id,
		).Scan(&f_cont); err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("Attachment_Controller_get_file conn.QueryRow(): %v",err))	
		}
		
		cache_f, err = os.Create(cache_fn)
		if err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("Attachment_Controller_get_file os.Create(): %v", err))
		}
		defer cache_f.Close()
		if _, err := cache_f.Write(f_cont); err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("Attachment_Controller_get_file Write(): %v", err))
		}
		sock_http, ok := sock.(*httpSrv.HTTPSocket)
		if !ok {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, "Attachment_Controller_get_file no cache sock must be *HTTPSocket")
		}

		httpSrv.ServeContent(sock_http, &f_cont, att_name.GetValue(), "", time.Time{}, cont_disp)		
	}
	
	return nil
}

//Method implemenation
//Structs described in attachments.go
func (pm *Attachment_Controller_add_file) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	http_sock, ok := sock.(*httpSrv.HTTPSocket)
	if !ok {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, "Attachment_Controller_add_file Not HTTPSocket type")
	}	
	
	args := rfltArgs.Interface().(*Attachment_add_file)	

	file, file_h, err := http_sock.Request.FormFile("content_data[]")
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("Attachment_Controller_add_file Request.FormFile(): %v",err))
	}
	defer file.Close()

	//проверка по типу
	/*mimes := []httpSrv.MIME_TYPE{httpSrv.MIME_TYPE_pdf,
		httpSrv.MIME_TYPE_png,
		httpSrv.MIME_TYPE_jpg,
		httpSrv.MIME_TYPE_jpeg,
		httpSrv.MIME_TYPE_jpe,
		httpSrv.MIME_TYPE_xls,
		httpSrv.MIME_TYPE_xlsx,
		httpSrv.MIME_TYPE_docx,
		httpSrv.MIME_TYPE_doc,
		httpSrv.MIME_TYPE_odt,
		httpSrv.MIME_TYPE_zip,
		httpSrv.MIME_TYPE_txt,
		httpSrv.MIME_TYPE_bmp,
	}
	if !FileHeaderContainsMimes(file_h, mimes) {
		return osbe.NewPublicMethodError(ER_UNSUPPORTED_MIME_CODE, ER_UNSUPPORTED_MIME_DESCR)
	}*/
	/*
	if !FileHeaderContainsMime(file_h, httpSrv.MIME_TYPE_pdf) &&
		!FileHeaderContainsMime(file_h, httpSrv.MIME_TYPE_png) &&
		!FileHeaderContainsMime(file_h, httpSrv.MIME_TYPE_jpg) &&
		!FileHeaderContainsMime(file_h, httpSrv.MIME_TYPE_jpeg) &&
		!FileHeaderContainsMime(file_h, httpSrv.MIME_TYPE_jpe) {
		return osbe.NewPublicMethodError(ER_UNSUPPORTED_MIME_CODE, ER_UNSUPPORTED_MIME_DESCR)
	}
	*/

	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	pool_conn, conn_id, err_с := d_store.GetPrimary()
	if err_с != nil {
		return err_с
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
			
	args.Content_info.Name = file_h.Filename
	preview_bt, err := AddFileThumbnailToDb(conn, app.GetBaseDir(), file, &args.Content_info, &args.Ref);
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("Attachment_Controller_add_file AddFileThumbnailToDb(): %v",err))	
	}
	
	AddPreviewModel(resp, args.Content_info.Id, preview_bt)
		
	return nil
}

