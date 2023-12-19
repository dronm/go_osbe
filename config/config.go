// Config package manages an application configuration.
// It contains the minimal configuration parameters.
// Database, session, ws server, log level are supported. 
package config

import (
	"encoding/json"
	"io/ioutil"
	"bytes"	
)

// DbStorage describes common db connections: one primary and several secondary servers.
type DbStorage struct {
	Primary string `json:"primary"`
	Secondaries map[string]string `json:"secondaries"`
}

// Session holds configuration for sessions.
type Session struct {
	MaxLifeTime int64 `json:"maxLifeTime"`
	MaxIdleTime int64 `json:"maxIdleTime"`
	EncKey string `json:"encKey"`
	DestroyAllTime string `json:"destroy_all_time"`
}

// AppConfig is the main application configuration structure.
type AppConfig struct {
	LogLevel string `json:"logLevel"`	//debug|warn|note|error
	Db DbStorage `json:"db"`		
	WSServer string `json:"wsServer"`	//Web socket server host:port
	TLSCert string `json:"TLSCert"`		//path to a TLS sertificate file
	TLSKey string `json:"TLSKey"`		//path to a TLS key file
	TLSWSServer string `json:"TLSwsServer"`	//TLS server host:port
	AppID string `json:"appId"`		// Application ID
	TemplateDir string `json:"templateDir"`	//Server template directory
	Session Session `json:"session"`
	ReportErrors bool `json:"reportErrors"` //If set to true public method error will be send to client,
						//otherwise error will be logged, short text will be sent to client
	
	DebugQueries bool `json:"debugQueries"`	
						
	XSLTDir string `json:"XSLTDir"`	
	DefaultLocale string `json:"defaultLocale"`					
	
	TechMail string `json:"techMail"`	//Author name && email
	Author string `json:"author"`
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
func (c *AppConfig) GetDestroyAllTime() string {
	return c.Session.DestroyAllTime
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

func (a *AppConfig) GetDebugQueries() bool {
	return a.DebugQueries
}

// ReadConf reads configiration from json file
func ReadConf(fileName string, c interface{}) error{
	file, err := ioutil.ReadFile(fileName)
	if err == nil {
		file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
		err = json.Unmarshal([]byte(file), c)		
	}
	return err
}

// WriteConf writes configiration from arbitary struct to json file
func WriteConf(fileName string, c interface{}) error{
	cont_b, err := json.Marshal(c)
	if err == nil {
		err = ioutil.WriteFile(fileName, cont_b, 0644)
	}
	return err
}


