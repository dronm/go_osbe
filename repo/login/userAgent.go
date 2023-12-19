package login

import (
	"encoding/json"
	
	"github.com/mssola/user_agent" //User agent parser
)

type userAgent struct {
	Platform     	string `json:"platform"`
	OSName      	string `json:"osName"`
	OSVersion   	string `json:"osVersion"`
	Mozilla     	string `json:"mozilla"`
	Localization	string `json:"localization"`
	EngineName	string `json:"engineName"`
	EngineVersion	string `json:"engineVersion"`
	BrowserName  	string `json:"browserName"`
	BrowserVersion  string `json:"browserVersion"`
	Bot          	bool `json:"bot"`
	Mobile       	bool `json:"mobile"`
}

func GetUserAgentFieldValue(userAgentHeader string) ([]byte, error) {
	ua := user_agent.New(userAgentHeader)
	os_inf := ua.OSInfo()	
	ua_s := userAgent{Platform: ua.Platform(),
		OSName: os_inf.Name,
		OSVersion: os_inf.Version,
		Mozilla: ua.Mozilla(),
		Mobile: ua.Mobile(),
		Localization: ua.Localization(),
		Bot: ua.Bot(),
		}
	ua_s.EngineName, ua_s.EngineVersion = ua.Engine()
	ua_s.BrowserName, ua_s.BrowserVersion = ua.Browser()
	
	return json.Marshal(&ua_s)
}


