package web_client

import(
	"fmt"
	"net/http"
	"io"
	"os"
	"errors"
)

const (
	HTTP_PREF = "http://"
	HTTPS_PREF = "https://"
	DEF_VIEW = "ViewJSON"
)

type WebClient struct {
	Server string
	Token string
}

func GetQueryString(contr, meth, view, tmpl, token string, params map[string]interface{}) string {

	if params == nil && (contr != "" || meth != "" || view != "" || tmpl != "") {
		params = make(map[string]interface{})
	}
	if contr != "" {
		params["c"] = contr
	}
	if meth != "" {
		params["f"] = meth
	}
	if view != "" {
		params["v"] = view
	}
	if tmpl != "" {
		params["t"] = tmpl
	}
	if token != "" {
		params["token"] = token
	}
	
	cmd := ""
	for key, val := range params {
		val_s := ""
		switch v := val.(type) {
			case float64, float32:
				val_s = fmt.Sprintf("%f", v)
				
			case int64, int, int32:
				val_s = fmt.Sprintf("%d", v)
			case string:	
				val_s = v		
		}
		
		if val_s != "" {
			if cmd != "" {
				cmd += "&"
			}
			cmd += key + "="+val_s
		}
	}
	
	return cmd
}

func (c *WebClient) SendGet(contr, meth, view, tmpl string, params map[string]interface{}) error {
	if c.Server == "" {
		return errors.New("server param not set")
	}
	if c.Token == "" {
		return errors.New("token param not set")
	}
	
	if c.Server[:len(HTTP_PREF)] != HTTP_PREF && c.Server[:len(HTTPS_PREF)] != HTTPS_PREF {
		c.Server = HTTP_PREF + c.Server
	}
	
	if c.Server[len(c.Server)-1:] != "/" {
		c.Server += "/"
	}
	
	if view == "" {
		view = DEF_VIEW
	}
	url := c.Server + "?" + GetQueryString(contr, meth, view, tmpl, c.Token, params)

	resp, err := http.Get(url) 
	if err != nil { 
		return err
	} 
	
	defer resp.Body.Close()
	fmt.Println("Query: "+url)
	fmt.Println("Reponse:")
	io.Copy(os.Stdout, resp.Body)
	fmt.Println("")
	
	return nil	
}
