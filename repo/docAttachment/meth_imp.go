package docAttachment

import (
	"fmt"
	"os"
	"os/exec"	
	"strings"
	"errors"
	"crypto/md5"
	"encoding/hex"
	b64 "encoding/base64"
	"io/ioutil"
	"io"
	"bytes"
	"context"
	"mime/multipart"
	"path/filepath"
	
	"osbe"
	"osbe/view"
	"osbe/model"
	"osbe/response"	
	"osbe/srv/httpSrv"
	
	"github.com/jackc/pgx/v4"
)

const (
	CACHE_DIR = "CACHE"
	
	ER_UNSUPPORTED_MIME_CODE = 2000
	ER_UNSUPPORTED_MIME_DESCR = "Неподдерживаемый тип файла"
)

type previewRow struct {
	Id string `json:"id"`
	Cont string `json:"cont"`
}

//file ID is unique inside ref
func GetAttachmentCacheFileName(baseDir string, refDataType string, refID HttpInt, fileID string) string {
	return filepath.Join(baseDir, CACHE_DIR, GetMd5(fmt.Sprintf("att_%s%d_%s", refDataType, refID, fileID)))
}

func GetPreviewCacheFileName(baseDir string, refDataType string, refID HttpInt, fileID string) string {
	return filepath.Join(baseDir, CACHE_DIR, GetMd5(fmt.Sprintf("prev_%s%d_%s", refDataType, refID, fileID)))
}

func runCMD(progName, commands, previewName string, toPDF bool) error {
	cmd_args := strings.Split(commands, " ")
	cmd := exec.Command(progName, cmd_args...)
	err := cmd.Run()
	if err != nil { 
		return errors.New(fmt.Sprintf("Error converting document to image: %v, params:%s %s", err, progName, commands)) 
	}
	
	if toPDF {
		var thbn_n string
		if view.FileExists(previewName + "-1.jpg") {
			thbn_n = previewName + "-1.jpg"
			
		}else if view.FileExists(previewName + "-01.jpg") {
			thbn_n = previewName + "-01.jpg"			
			
		}else if view.FileExists(previewName + "-001.jpg") {
			thbn_n = previewName + "-001.jpg"
		}
		//thbn_n -->> previewName
		os.Rename(thbn_n, previewName)		
	}
		
	return nil
}

//realName for mime type!!!
//attName - attachment name
//pName - preview name
//realName
func GenThumbnail(attName, pName, realName string) error {
	var fExt string
	f_parts := strings.Split(realName, ".")
	if len(f_parts) > 0 {
		fExt = strings.ToLower(f_parts[len(f_parts)-1])
	}
	
	pdftoppm_fmt := "-l 1 -scale-to 300 -jpeg %s %s" //-q no comment or errors
	
	var cmd_name string
	var cmd_s string
	var pdf bool
	if fExt == "doc" || fExt == "docx" || fExt == "xls" || fExt == "xlsx" ||  fExt == "odt" ||  fExt == "ods"{
		//openoffice first to pdf

		//export HOME=CACHE && /usr/lib/libreoffice/program/./soffice --headless --convert-to pdf --outdir CACHE CACHE %s
		if err := runCMD("soffice",
				fmt.Sprintf("--headless --convert-to pdf:writer_pdf_Export --outdir CACHE %s", attName),
				pName, true); err != nil {
			return err
		}
		//got full pdf attName.pdf
		//pdf to image
		if err := runCMD("pdftoppm", fmt.Sprintf(pdftoppm_fmt, attName + ".pdf", pName), pName, true); err != nil {
			return err
		}
		os.Remove(attName + ".pdf") //remove temp full pdf file
		return nil
		
	}else if fExt == "pdf" {
		pdf = true
		cmd_name = "pdftoppm"
		cmd_s = fmt.Sprintf(pdftoppm_fmt, attName, pName)		
		
	}else {
		cmd_name = "convert"
		cmd_s = fmt.Sprintf("-define jpeg:size=500x180 %s -auto-orient -thumbnail 250x100 -unsharp 0x.5 %s", attName, pName)		
	}
	
	return runCMD(cmd_name, cmd_s, pName, pdf)

}

//adds to response new model with  Base64 content
func AddPreviewModel(resp *response.Response, fileID string, pCont []byte) {
	ret_model := &model.Model{ID: model.ModelID("Preview_Model"), Rows: make([]model.ModelRow, 0)}
	ret_model.Rows = append(ret_model.Rows, &previewRow{Id: fileID, Cont: b64.StdEncoding.EncodeToString(pCont)})
	resp.AddModel(ret_model)
}

func StoreAttachment(conn *pgx.Conn, ref *Ref_Type, fileInfo *Content_info_Type, fileData []byte, previewData []byte) error {
	fileInfo.Size = int64(len(fileData))
	if _, err := conn.Exec(context.Background(), `BEGIN` ); err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("StoreAttachment BEGIN: %v",err))
	}	

	if _, err := conn.Exec(context.Background(),
		`DELETE FROM attachments
		WHERE ref->>'dataType' = $1 AND (ref->'keys'->>'id')::int = $2 AND content_info->>'id' = $3`,
			ref.DataType, ref.Keys.Id, fileInfo.Id,
		); err != nil {
		
		conn.Exec(context.Background(), `ROLLBACK` )
		
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("StoreAttachment DELETE: %v",err))
	}	

	if _, err := conn.Exec(context.Background(),
		`INSERT INTO attachments
		(ref, content_info, content_data, content_preview)
		VALUES ($1, $2, $3, $4)`,
			ref,
			fileInfo,	
			fileData,
			previewData,
		); err != nil {
		
		conn.Exec(context.Background(), `ROLLBACK` )
		
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("StoreAttachment conn.QueryRow(): %v",err))
	}

	if _, err := conn.Exec(context.Background(), `COMMIT` ); err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("StoreAttachment COMMIT: %v",err))
	}	
	return nil
}

func GenAttachmentThumbnail(baseDir string, refDataType string, refID HttpInt, fileInfo *Content_info_Type, attBuf io.Reader) ([]byte, error) {
	att_n := GetAttachmentCacheFileName(baseDir, refDataType, refID, fileInfo.Id)
	file_att, err := os.OpenFile(att_n, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return []byte{}, osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("GenAttachmentThumbnail os.OpenFile(): %v",err))
	}
	defer file_att.Close()
	_, err = io.Copy(file_att, attBuf)
	if err != nil {
		return []byte{}, osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("GenAttachmentThumbnail io.Copy(): %v",err))
	}
	
	preview_fn := GetPreviewCacheFileName(baseDir, refDataType, refID, fileInfo.Id)
	if err := GenThumbnail(att_n, preview_fn, fileInfo.Name); err != nil {
		return []byte{}, osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("GenAttachmentThumbnail GenThumbnail(): %v",err))
	}
	defer os.Remove(preview_fn)
	
	var preview_bt []byte
	preview_bt, err = ioutil.ReadFile(preview_fn)
	if err != nil {
		return []byte{}, osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("GenAttachmentThumbnail os.ReadFile(): %v", err))
	}
	return preview_bt, nil		
}

//
//Structs described in attachments.go
//file multipart.File io.Reader
func AddFileThumbnailToDb(conn *pgx.Conn, baseDir string, file io.Reader, fileInfo *Content_info_Type, ref *Ref_Type) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return []byte{}, osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("AddFileThumbnailToDb io.Copy(): %v",err))
	}		
	file_bt := buf.Bytes()	
	fileInfo.Size = int64(buf.Len())
	//thumbnail
	preview_bt, err := GenAttachmentThumbnail(baseDir, ref.DataType, ref.Keys.Id, fileInfo, buf)
	if err != nil {
		return []byte{}, osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("AddFileThumbnailToDb: %v",err))
	}

	if err := StoreAttachment(conn, ref, fileInfo, file_bt, preview_bt); err != nil {
		return []byte{}, err
	}
	
	return preview_bt, nil
}

func ClearCache(baseDir string, ref Ref_Type, contentID string) error {
	att_n := GetAttachmentCacheFileName(baseDir, ref.DataType, ref.Keys.Id, contentID)
	if view.FileExists(att_n) {
		if err := os.Remove(att_n); err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("ClearCache os.Remove(): %v",err))	
		}
	}
	preview_fn := GetPreviewCacheFileName(baseDir, ref.DataType, ref.Keys.Id, contentID)
	if view.FileExists(preview_fn) {
		if err := os.Remove(preview_fn); err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("ClearCache os.Remove(): %v",err))	
		}
	}
	return nil	
}

func RemoveAttachment(conn *pgx.Conn, baseDir string, ref Ref_Type, contentID string) error {
	if _, err := conn.Exec(context.Background(),
		`DELETE FROM attachments
		WHERE ref->>'dataType' = $1
			AND (ref->'keys'->>'id')::int = $2
			AND content_info->>'id' = $3`,
		ref.DataType,
		ref.Keys.Id,
		contentID,
	); err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("DeleteAttachment conn.Exec(): %v",err))	
	}
	
	return ClearCache(baseDir, ref, contentID)
}

func FileHeaderContainsMimes(fileH *multipart.FileHeader, mimes []httpSrv.MIME_TYPE) bool {
	if tp, ok := fileH.Header["Content-Type"]; ok {
		for _, tp_id := range tp {
			for _, m_id := range mimes {
				if m_id == httpSrv.MIME_TYPE(tp_id) {
					return true
				}
			}
		}
	}
	return false
}
func FileHeaderContainsMime(fileH *multipart.FileHeader, mimeId httpSrv.MIME_TYPE) bool {
	if tp, ok := fileH.Header["Content-Type"]; ok {
		for _, tp_id := range tp {
			if mimeId == httpSrv.MIME_TYPE(tp_id) {
				return true
			}
		}
	}
	return false
}

func GetMd5(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

