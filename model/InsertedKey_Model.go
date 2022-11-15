package model

const INSERTED_MODEL_ID ModelID = "InsertedKey"	

func New_InsertedKey_Model(row ModelRow) Modeler{
	m := &Model{ID: INSERTED_MODEL_ID, Rows: make([]ModelRow, 1)}
	m.Rows[0] = row
	return m
}

