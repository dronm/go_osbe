package captcha

/**
 * Andrey Mikhalevich 16/12/22
 *
 * Controller implimentation file
 *
 */

import (
	"reflect"
	"fmt"
	b64 "encoding/base64"
	"bytes"	
	//"encoding/gob"
		
	"osbe"
	"osbe/srv"
	"osbe/socket"
	"osbe/model"
	"osbe/response"	
	"osbe/logger"
	
	"github.com/dronm/session"
	cpt "github.com/dchest/captcha"
)

const (
	DEF_WIDTH = 240
	DEF_HEIGHT = 80
	DEF_COUNT = 6
)

//Method implemenation
func (pm *Captcha_Controller_get) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	args := rfltArgs.Interface().(*Captcha_get)
	w := int(args.Width.GetValue())
	if w == 0 {
		w = DEF_WIDTH
	}
	h := int(args.Height.GetValue())
	if h == 0 {
		h = DEF_HEIGHT
	}
	cnt := int(args.Count.GetValue())
	if cnt == 0 {
		cnt = DEF_COUNT
	}
	return AddNewCaptcha(sock.GetSession(), app.GetLogger(), resp, args.Id.GetValue(), w, h, cnt)
}	

//*********************************************
type CaptchaSessVal struct {
	Id string `json:"id"`
	Digits []byte `json:"digits"`
}

//captchaStore implements cpt.Store interface
type captchaStore struct {
	Sess session.Session
	Log logger.Logger
	ID string
}

func (s captchaStore) Set(id string, digits []byte) {	
	if err := s.Sess.Set(s.ID, CaptchaSessVal{Id: id, Digits: digits}); err != nil {
		s.Log.Errorf("captchaStore Set() Session.Set(): %v", err)
	}
	if err := s.Sess.Flush(); err != nil {
		s.Log.Errorf("captchaStore Set() Session.Flush(): %v", err)
	}
}

func (s captchaStore) Get(id string, clear bool) []byte {
	captcha_data, err := getSessCaptchaData(s.Sess, s.ID)
	if err != nil {
		s.Log.Errorf("captchaStore getSessCaptchaData(): %v", err)
		return []byte{}
	}
	return captcha_data.Digits
}

//returns:
//	captcha id
//	image
//	error
func NewCaptcha(sess session.Session, log logger.Logger, captchaID string, w int, h int, charCount int) ([]byte, error) {	
	//gob.Register(CaptchaSessVal{})
	cap_store := captchaStore{Sess: sess, Log: log, ID: "captcha_"+captchaID}
	cpt.SetCustomStore(&cap_store)
	id := cpt.NewLen(charCount)
	var buf bytes.Buffer
	if err := cpt.WriteImage(&buf, id, w, h); err != nil {
		return []byte{}, err	
	}
	
	return buf.Bytes(), nil
}

func CaptchaVerify(sess session.Session, log logger.Logger, key []byte, captchaID string) (bool, error) {
	sess_val_id := "captcha_" + captchaID
	captcha_data, err := getSessCaptchaData(sess, sess_val_id)
	if err != nil {
		return false, err
	}
	sess.Delete(sess_val_id)
	if err := sess.Flush(); err != nil {
		return false, err
	}
	if len(captcha_data.Digits) != len(key) {
		return false, nil
	}
	//turn from ascii codes to numbers
	for i :=0; i<len(key); i++ {
		if (key[i]-48) != captcha_data.Digits[i] {
			return false, nil
		}
	}
	return true, nil
}

func getSessCaptchaData(sess session.Session, captchaID string) (*CaptchaSessVal, error) {
	captcha := CaptchaSessVal{}
	if err := sess.Get(captchaID, &captcha); err != nil {
		return nil, err
	}
	return &captcha, nil
}

func AddNewCaptcha(sess session.Session, log logger.Logger, resp *response.Response, captchaID string, w int, h int, charCount int) error {
	captcha_data, err := NewCaptcha(sess, log, captchaID, w, h, charCount)
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("AddNewCaptcha(): %v",err))
	}
	
	m := make([]model.ModelRow, 1)
	m[0] = struct{
		Img string `json:"img"`
	}{Img: b64.StdEncoding.EncodeToString(captcha_data)}
	resp.AddModel(&model.Model{ID: model.ModelID("Captcha_Model"), Rows: m})
	
	return nil
}


