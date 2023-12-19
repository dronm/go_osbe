package notif

import(
	"net/http"
	"encoding/json"
	"context"
	"bytes"
	"os"
	"io"
	"io/ioutil"
	"fmt"
	"mime/multipart"
	b64 "encoding/base64"
	"errors"
	"strings"
	
	"github.com/jackc/pgx/v4"
)

type NotifProvider string
const (
	PROV_EMAIL NotifProvider = "email"
	PROV_SMS NotifProvider = "sms"
	PROV_TM NotifProvider = "tm"
	PROV_WA NotifProvider = "wa"
	PROV_VB NotifProvider = "vb"
)

//TMMessage is a Telegram provider message
type TMMessage struct {
	ChatID string `json:"chat_id"`
	Text string `json:"text"`
}
//NewNotif returns new notif message with one provider
func (msg *TMMessage) NewNotif(messageType string) *NotifMessage {
	new_m := NotifMessage{MessageType: messageType,
		Providers: []NotifProvider{PROV_TM},
		Message: make(map[NotifProvider]map[string]interface{}),
	}
	msg.SetParams(&new_m)
	return &new_m
}
//SetParams sets NotifMessage parameters from TMMessage structure
func (msg *TMMessage) SetParams(notifMsg *NotifMessage) {
	notifMsg.Message[PROV_TM] = make(map[string]interface{})
	notifMsg.Message[PROV_TM]["chat_id"] = msg.ChatID
	notifMsg.Message[PROV_TM]["text"] = msg.Text
}

//---------------------------------------------------------------------------------------------------
//WAMessage is a WhatsUp provider message
type WAMessage struct {
	Tel string `json:"tel"`
	Text string `json:"text"`
}
//NewNotif returns new notif message with one provider
func (msg *WAMessage) NewNotif(messageType string) *NotifMessage {
	new_m := NotifMessage{MessageType: messageType,
		Providers: []NotifProvider{PROV_WA},
		Message: make(map[NotifProvider]map[string]interface{}),
	}
	msg.SetParams(&new_m)
	return &new_m
}

//SetParams sets NotifMessage parameters from WAMessage structure
func (msg *WAMessage) SetParams(notifMsg *NotifMessage) {
	notifMsg.Message[PROV_WA] = make(map[string]interface{})
	notifMsg.Message[PROV_WA]["tel"] = msg.Tel
	notifMsg.Message[PROV_WA]["text"] = msg.Text
}

//--------------------------------------------------------------------------------------------------
//SMSMessage is a SMS provider message
type SMSMessage struct {
	Tel string `json:"tel"`
	Text string `json:"text"`
}
//NewNotif returns new notif message with one provider
func (msg *SMSMessage) NewNotif(messageType string) *NotifMessage {
	new_m := NotifMessage{MessageType: messageType,
		Providers: []NotifProvider{PROV_SMS},
		Message: make(map[NotifProvider]map[string]interface{}),
	}
	msg.SetParams(&new_m)
	return &new_m
}
//SetParams sets NotifMessage parameters from SMSMessage structure
func (msg *SMSMessage) SetParams(notifMsg *NotifMessage) {
	notifMsg.Message[PROV_SMS] = make(map[string]interface{})
	notifMsg.Message[PROV_SMS]["tel"] = msg.Tel
	notifMsg.Message[PROV_SMS]["text"] = msg.Text
}

//--------------------------------------------------------------------------------------------------
//SMSMessage is an EmailMessage provider message
type EmailMessage struct {
	FromAddr string `json:"from_addr"`
	FromName string `json:"from_name"`
	ToAddr string `json:"to_addr"`
	ToName string `json:"to_name"`
	ReplyName string `json:"reply_name"`
	Body string `json:"body"`
	SenderAddr string `json:"sender_addr"`
	Subject string `json:"subject"`
	Attachments []string `json:"attachments"`	//pathes to attachments
	AttachmentAlias []string 			//Aliases for file names, same Len as Attachments, no paths
							//Will not be sent to server
}
//NewNotif returns new notif message with email provider
func (msg *EmailMessage) NewNotif(messageType string) *NotifMessage {
	new_m := NotifMessage{MessageType: messageType,
		Providers: []NotifProvider{PROV_EMAIL},
		Message: make(map[NotifProvider]map[string]interface{}),
	}
	msg.SetParams(&new_m)
	return &new_m
}

//SetParams sets NotifMessage parameters from EmailMessage structure
func (msg *EmailMessage) SetParams(notifMsg *NotifMessage) {
	notifMsg.Message[PROV_EMAIL] = make(map[string]interface{})
	notifMsg.Message[PROV_EMAIL]["from_addr"] = msg.FromAddr
	notifMsg.Message[PROV_EMAIL]["from_name"] = msg.FromName
	notifMsg.Message[PROV_EMAIL]["to_addr"] = msg.ToAddr
	notifMsg.Message[PROV_EMAIL]["to_name"] = msg.ToName
	notifMsg.Message[PROV_EMAIL]["reply_name"] = msg.ReplyName
	notifMsg.Message[PROV_EMAIL]["body"] = msg.Body
	notifMsg.Message[PROV_EMAIL]["sender_addr"] = msg.SenderAddr
	notifMsg.Message[PROV_EMAIL]["subject"] = msg.Subject
	notifMsg.Message[PROV_EMAIL]["attachments"] = msg.Attachments
	notifMsg.Message[PROV_EMAIL]["attachment_alias"] = msg.AttachmentAlias
}

type NotifMessage struct {
	MessageType string `json:"messageType"`					// application defined message type, arbitary string
	Providers []NotifProvider `json:"providers"`				// list of providers for the message in order of priority	
	Message map[NotifProvider]map[string]interface{} `json:"message"`	// message, specific structure for every provider
}
func (nm *NotifMessage) AddMessage(prov NotifProvider) {
	if nm.Providers == nil {
		nm.Providers = make([]NotifProvider, 1)
		nm.Providers[0] = prov
	}else{
		nm.Providers = append(nm.Providers, prov)
	}
	if nm.Message == nil {
		nm.Message = make(map[NotifProvider]map[string]interface{})
	}
}
func (nm *NotifMessage) AddEmailMessage(msg *EmailMessage) {
	nm.AddMessage(PROV_EMAIL)
	msg.SetParams(nm)	
}
func NewNotifMessage(messageType string, providers ...NotifProvider) *NotifMessage {
	return &NotifMessage{MessageType: messageType,
		Providers: providers,
		Message: make(map[NotifProvider]map[string]interface{}),
	}
}

type Response struct {
	ID int `json:"id"`
	Error string `json:"error"`
}

type Notifier struct {
	AppName string `json:"appName"`		//login
	Pwd string `json:"pwd"`			//password
	Host string `json:"host"`
}

func NewNotifier(host, appName, pwd string) *Notifier {
	return &Notifier{AppName: appName,
		Pwd: pwd,	
		Host: host,
	}
}

type fileForSending struct {
	fileName string //path to file
	alias string
}
func (n *Notifier) Send(batch []*NotifMessage) ([]*Response, error) {
	var files []fileForSending //actual files with alias
	//check attachments: provider=email + param=attachments
	for i, m := range batch {
		if m.Message == nil {
			continue
		}
		prov_msg, ok := m.Message[PROV_EMAIL]
		if !ok {
			continue
		}
		att_i, ok := prov_msg["attachments"]
		if !ok {
			continue
		}
		att_list, ok := att_i.([]string)
		if !ok {
			continue
		}
		if files == nil {
			files = make([]fileForSending, 0)
		}
		var att_alias_list []string
		att_alias_i, ok := prov_msg["attachment_alias"]
		if ok {
			att_alias_list, _ = att_alias_i.([]string)
		}
		for att_ind, att_s := range att_list {
			att := fileForSending{fileName: att_s}
			if att_alias_list != nil && len(att_alias_list) > att_ind {
				att.alias = att_alias_list[att_ind]
			}
			files = append(files, att)
		}
		if m.Message[PROV_EMAIL]["attachment_alias"] != nil {
			batch[i].Message[PROV_EMAIL]["attachments"] = m.Message[PROV_EMAIL]["attachment_alias"]
			delete(batch[i].Message[PROV_EMAIL], "attachment_alias")
		}
	}
	
	json_b, err := json.Marshal(&batch)
	if err != nil {
		return nil, err
	}
    	var b bytes.Buffer
    	var cont_type string
    	if files == nil || len(files) == 0 {
    		//no attachments, content-type=application/json
    		b = *bytes.NewBuffer(json_b)
    		cont_type = "application/json"
    		
    	}else{
    		//attachments, multipart form
    		//messages=[]NotifMessage, fileN...
		w := multipart.NewWriter(&b)	
		fw, err := w.CreateFormField("messages")
		if err != nil {
			return nil, err
		}
		if _, err = io.Copy(fw, bytes.NewReader(json_b)); err != nil {
			return nil, err
		}		
		for i, fl := range files {
			f, err := os.Open(fl.fileName)
			if err != nil {
				return nil, err
			}
			var file_name string
			if fl.alias != "" {
				file_name = fl.alias
			}else{
				file_name = f.Name() //original name
			}
			if fw, err = w.CreateFormFile(fmt.Sprintf("file%d", i), file_name); err != nil {
				return nil, err
			}
			if _, err = io.Copy(fw, f); err != nil {
				return nil, err
			}
		}
		w.Close()
		cont_type = w.FormDataContentType() //this will contain the boundary
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", n.Host, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", n.AppName, n.Pwd)))) )
	req.Header.Set("Content-Type", cont_type)
	resp, err := client.Do(req)	
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("error http code: %d", resp.StatusCode))
	}
	
	resp_body := make([]*Response, 0)
	if err := json.Unmarshal(body, &resp_body); err != nil {
		return nil, err
	}
	
	return resp_body, nil
}

//----------------------------------------------------------------------------------
//SQL function can accept any namber of parameters and must have this signature:
//	tel text
//	body text
func NewWAMessageFromSQL(conn *pgx.Conn, sqlFunc string, sqlParamValues []interface{}) (*WAMessage, error) {
	var sql_params strings.Builder
	for i := 0; i < len(sqlParamValues); i++ {
		if i > 0 {
			sql_params.WriteString(",")
		}
		sql_params.WriteString(fmt.Sprintf("$%d", i+1))
	}
	msg := &WAMessage{}
	if err := conn.QueryRow(context.Background(),
		fmt.Sprintf(`SELECT
			tel,
			body
		FROM %s(%s)`, sqlFunc, sql_params.String()),
		sqlParamValues...).Scan(
				&msg.Tel,
				&msg.Text,
			); err != nil && err != pgx.ErrNoRows {
	
		return nil, err	
		
	}else if err == pgx.ErrNoRows {
		return nil, errors.New(ER_TEMPLATE_NOT_FOUND)
	} 
	
	if msg.Tel == "" {
		return nil, errors.New(ER_NO_TEL)
	}
	return msg, nil
}

func NewWAMessagesFromSQL(conn *pgx.Conn, sqlFunc string, sqlParamValues []interface{}) ([]*WAMessage, error) {
	var sql_params strings.Builder
	for i := 0; i < len(sqlParamValues); i++ {
		if i > 0 {
			sql_params.WriteString(",")
		}
		sql_params.WriteString(fmt.Sprintf("$%d", i+1))
	}
	msg_list := make([]*WAMessage, 0)
	rows, err := conn.Query(context.Background(),
		fmt.Sprintf(`SELECT
			tel,
			body
		FROM %s(%s)`, sqlFunc, sql_params.String()),
		sqlParamValues...)
	if err != nil {
		return nil, err
	} 
	for rows.Next() {
		msg := WAMessage{}
		if err := rows.Scan(&msg.Tel, &msg.Text); err != nil {		
			return nil, err
		}
		msg_list = append(msg_list, &msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()
	
	return msg_list, nil
}

//SQL function can accept any namber of parameters and must have this signature:
//	tel text
//	body text
func NewSMSMessageFromSQL(conn *pgx.Conn, sqlFunc string, sqlParamValues []interface{}) (*SMSMessage, error) {
	var sql_params strings.Builder
	for i := 0; i < len(sqlParamValues); i++ {
		if i > 0 {
			sql_params.WriteString(",")
		}
		sql_params.WriteString(fmt.Sprintf("$%d", i+1))
	}
	msg := &SMSMessage{}
	if err := conn.QueryRow(context.Background(),
		fmt.Sprintf(`SELECT
			tel,
			body
		FROM %s(%s)`, sqlFunc, sql_params.String()),
		sqlParamValues...).Scan(
				&msg.Tel,
				&msg.Text,
			); err != nil && err != pgx.ErrNoRows {
	
		return nil, err	
		
	}else if err == pgx.ErrNoRows {
		return nil, errors.New(ER_TEMPLATE_NOT_FOUND)
	} 
	
	if msg.Tel == "" {
		return nil, errors.New(ER_NO_TEL)
	}
	return msg, nil
}

//SQL function can accept any namber of parameters and must have this signature:
//	from_addr text
//	from_name text
//	reply_name text
//	sender_addr text
//	to_name text
//	body text
//	subject text
func NewEmailMessageFromSQL(conn *pgx.Conn, sqlFunc string, attachments []string, attachmentAlias []string, sqlParamValues []interface{}) (*EmailMessage, error) {
	var sql_params strings.Builder
	for i := 0; i < len(sqlParamValues); i++ {
		if i > 0 {
			sql_params.WriteString(",")
		}
		sql_params.WriteString(fmt.Sprintf("$%d", i+1))
	}
	msg := &EmailMessage{}
	if err := conn.QueryRow(context.Background(),
		fmt.Sprintf(`SELECT
			from_addr,
			from_name,
			reply_name,
			sender_addr,
			to_addr,
			to_name,		
			body,	
			subject		
		FROM %s(%s)`, sqlFunc, sql_params.String()),
		sqlParamValues...).Scan(
				&msg.FromAddr,
				&msg.FromName,
				&msg.ReplyName,
				&msg.SenderAddr,
				&msg.ToAddr,
				&msg.ToName,		
				&msg.Body,	
				&msg.Subject,
			); err != nil && err != pgx.ErrNoRows {
	
		return nil, err	
		
	}else if err == pgx.ErrNoRows {
		return nil, errors.New(ER_TEMPLATE_NOT_FOUND)
	} 
	
	if msg.ToAddr == "" {
		return nil, errors.New(ER_EMAIL_NO_TO_ADDR)
	}
	msg.Attachments = attachments
	msg.AttachmentAlias = attachmentAlias		
	return msg, nil
}

