package main

import (
	"encoding/json"
	"io/ioutil"
	"bytes"	
)

//json configuration file storage

type AppConfig struct {
	Server string `json:"server"`
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

