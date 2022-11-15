package model

import (
	"osbe/fields"
	"io/ioutil"
)

const FILE_MODEL_ID ModelID = "File"

type TFile struct {
	Name string `json:"name"`
	Size int `json:"size"`
}
type FileExport struct {
	File *TFile
	Content *fields.ValBytea `json:"content"`
}


func NewFileModel(f *TFile, content *fields.ValBytea) Modeler{
	m := &Model{ID: FILE_MODEL_ID, Rows: make([]ModelRow, 1)}
	m.Rows[0] = &FileExport{File: f, Content: content}
	return m
}

func NewFileModelFromFile(f *TFile, fileName string) (Modeler, error) {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	f.Size = len(dat)
	return NewFileModel(f, &fields.ValBytea{fields.Val{true, false}, dat}), nil
}
