package config

import (
	"encoding/json"
	"io/ioutil"
	"bytes"	
)

//json configuration file storage

type DbStorage struct {
	Primary string `json:"primary"`
	Secondaries map[string]string `json:"secondaries"`
}

type Session struct {
	MaxLifeTime int64 `json:"maxLifeTime"`
	MaxIdleTime int64 `json:"maxIdleTime"`
	EncKey string `json:"encKey"`
}

type AppConfig struct {
	LogLevel string `json:"logLevel"`
	Db DbStorage `json:"db"`
	WSServer string `json:"wsServer"`
	TLSCert string `json:"TLSCert"`
	TLSKey string `json:"TLSKey"`
	TLSWSServer string `json:"TLSwsServer"`
	AppID string `json:"appId"`
	TemplateDir string `json:"templateDir"`
	Session Session `json:"session"`
	ReportErrors bool `json:"reportErrors"` //If set to true public method error will be send to client,
						//otherwise error will be logged, short text will be sent to client
	XSLTDir string `json:"XSLTDir"`	
	DefaultLocale string `json:"defaultLocale"`					
	
	TechMail string `json:"techMail"`
	Author string `json:"author"`
}

func (c *AppConfig) ReadConf(fileName string) error{
	file, err := ioutil.ReadFile(fileName)
	if err == nil {
		file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
		err = json.Unmarshal([]byte(file), c)		
	}
	return err
}

func (c *AppConfig) WriteConf(fileName string) error{
	cont_b, err := json.Marshal(c)
	if err == nil {
		err = ioutil.WriteFile(fileName, cont_b, 0644)
	}
	return err
}

func (c *AppConfig) GetDb() DbStorage {
	return c.Db
}

func (c *AppConfig) GetWSServer() string {
	return c.WSServer
}

func (c *AppConfig) GetTLSWSServer() string {
	return c.TLSWSServer
}

func (c *AppConfig) GetTLSKey() string {
	return c.TLSKey
}

func (c *AppConfig) GetTLSCert() string {
	return c.TLSCert
}

func (c *AppConfig) GetAppID() string {
	return c.AppID
}

func (c *AppConfig) GetLogLevel() string {
	return c.LogLevel
}

func (c *AppConfig) GetSessMaxLifeTime() int64 {
	return c.Session.MaxLifeTime
}

func (c *AppConfig) GetSessMaxIdleTime() int64 {
	return c.Session.MaxIdleTime
}

func (c *AppConfig) GetSessEncKey() string {
	return c.Session.EncKey
}

func (c *AppConfig) GetTemplateDir() string {
	return c.TemplateDir
}

func (a *AppConfig) GetSession() Session{
	return a.Session
}

func (a *AppConfig) GetReportErrors() bool{
	return a.ReportErrors
}

func (a *AppConfig) GetXSLTDir() string {
	return a.XSLTDir
}

func (a *AppConfig) GetDefaultLocale() string {
	return a.DefaultLocale
}

func (a *AppConfig) GetTechMail() string {
	return a.TechMail
}
func (a *AppConfig) GetAuthor() string {
	return a.Author
}

