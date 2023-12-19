package httpSrv

import "strings"

type MIME_TYPE string
const (
	MIME_TYPE_ai MIME_TYPE = "application/postscript"
	MIME_TYPE_aif MIME_TYPE = "audio/aiff"
	MIME_TYPE_ani MIME_TYPE = "application/x-navi-animation"
	MIME_TYPE_arc MIME_TYPE = "application/octet-stream"
	MIME_TYPE_arj MIME_TYPE = "application/arj"
	MIME_TYPE_asm MIME_TYPE = "text/x-asm"
	MIME_TYPE_asp MIME_TYPE = "text/asp"
	MIME_TYPE_asx MIME_TYPE = "application/x-mplayer2"
	MIME_TYPE_au MIME_TYPE = "audio/basic"
	MIME_TYPE_avi MIME_TYPE = "video/avi"
	MIME_TYPE_bmp MIME_TYPE = "image/bmp"
	MIME_TYPE_exe MIME_TYPE = "application/octet-stream"
	MIME_TYPE_gif MIME_TYPE = "image/gif"
	MIME_TYPE_png MIME_TYPE = "image/png"
	MIME_TYPE_gz MIME_TYPE = "application/x-gzip"
	MIME_TYPE_htm MIME_TYPE = "text/html"
	MIME_TYPE_html MIME_TYPE = "text/html"
	MIME_TYPE_ico MIME_TYPE = "image/x-icon"
	MIME_TYPE_inf MIME_TYPE = "application/inf"
	MIME_TYPE_jpe MIME_TYPE = "image/jpe"
	MIME_TYPE_jpg MIME_TYPE = "image/jpg"
	MIME_TYPE_jpeg MIME_TYPE = "image/jpeg"
	MIME_TYPE_js MIME_TYPE = "application/javascript"
	MIME_TYPE_json MIME_TYPE = "application/json"
	MIME_TYPE_lha MIME_TYPE = "application/lha"
	MIME_TYPE_list MIME_TYPE = "text/plain"
	MIME_TYPE_lst MIME_TYPE = "text/plain"
	MIME_TYPE_lzh MIME_TYPE = "application/octet-stream"
	MIME_TYPE_mov MIME_TYPE = "video/quicktime"
	MIME_TYPE_movie MIME_TYPE = "video/x-sgi-movie"
	MIME_TYPE_mp2 MIME_TYPE = "audio/mpeg"
	MIME_TYPE_mp3 MIME_TYPE = "audio/mpeg3"
	MIME_TYPE_mp4 MIME_TYPE = "video/mp4"
	MIME_TYPE_mpa MIME_TYPE = "video/mpeg"
	MIME_TYPE_mpeg MIME_TYPE = "video/mpeg"
	MIME_TYPE_mpg MIME_TYPE = "audio/mpeg"
	MIME_TYPE_mpga MIME_TYPE = "audio/mpeg"
	MIME_TYPE_pdf MIME_TYPE = "application/pdf"
	MIME_TYPE_pic MIME_TYPE = "image/pict"
	MIME_TYPE_ppt MIME_TYPE = "application/mspowerpoint"
	MIME_TYPE_qt MIME_TYPE = "video/quicktime"
	MIME_TYPE_qtif MIME_TYPE = "image/x-quicktime"
	MIME_TYPE_sgml MIME_TYPE = "text/sgml"
	MIME_TYPE_sh MIME_TYPE = "application/x-bsh"
	MIME_TYPE_shtml MIME_TYPE = "text/html"
	MIME_TYPE_tar MIME_TYPE = "application/x-tar"
	MIME_TYPE_tgz MIME_TYPE = "application/x-compressed"
	MIME_TYPE_tif MIME_TYPE = "image/tiff"
	MIME_TYPE_tiff MIME_TYPE = "image/tiff"
	MIME_TYPE_txt MIME_TYPE = "text/plain"
	MIME_TYPE_uri MIME_TYPE = "text/uri-list"
	MIME_TYPE_vcd MIME_TYPE = "application/x-cdlink"
	MIME_TYPE_vmd MIME_TYPE = "application/vocaltec-media-desc"
	MIME_TYPE_vsd MIME_TYPE = "application/x-visio"
	MIME_TYPE_vst MIME_TYPE = "application/x-visio"
	MIME_TYPE_vsw MIME_TYPE = "application/x-visio"
	MIME_TYPE_wav MIME_TYPE = "audio/wav"
	MIME_TYPE_wmf MIME_TYPE = "windows/metafile"
	MIME_TYPE_xls MIME_TYPE = "application/excel"
	MIME_TYPE_xm MIME_TYPE = "audio/xm"
	MIME_TYPE_xml MIME_TYPE = "text/xml"
	MIME_TYPE_z MIME_TYPE = "application/x-compressed"
	MIME_TYPE_zip MIME_TYPE = "application/zip"
	MIME_TYPE_gzip MIME_TYPE = "application/gzip"
	
	MIME_TYPE_odt MIME_TYPE = "application/vnd.oasis.opendocument.text"
	MIME_TYPE_ott MIME_TYPE = "application/vnd.oasis.opendocument.text-template"
	MIME_TYPE_sig MIME_TYPE = "application/pgp-signature"
	MIME_TYPE_doc MIME_TYPE = "application/msword"
	MIME_TYPE_docx MIME_TYPE = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	MIME_TYPE_xlsx MIME_TYPE = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	
)

func GetMimeTypeOnFileExt(fName string) MIME_TYPE {
	var fExt string
	f_parts := strings.Split(fName, ".")
	if len(f_parts) > 0 {
		fExt = f_parts[len(f_parts)-1]
	}
	switch fExt {
		case "xml":
			return MIME_TYPE_xml		
		case "pdf":
			return MIME_TYPE_pdf
		case "zip":
			return MIME_TYPE_zip
		case "gzip":
			return MIME_TYPE_gzip
		case "gif":
			return MIME_TYPE_gif
		case "png":
			return MIME_TYPE_png
		case "jpeg":
			return MIME_TYPE_jpeg
		case "jpg":
			return MIME_TYPE_jpg
		case "txt":
			return MIME_TYPE_txt
		case "html":
			return MIME_TYPE_html
		case "odt":
			return MIME_TYPE_odt
		case "ott":
			return MIME_TYPE_ott
		case "sig":
			return MIME_TYPE_sig
		case "doc":
			return MIME_TYPE_doc
		case "docx":
			return MIME_TYPE_docx
		case "arj":
			return MIME_TYPE_arj
		case "avi":
			return MIME_TYPE_avi
		case "bmp":
			return MIME_TYPE_bmp
			
		default:
			return MIME_TYPE_exe			
	}
}
