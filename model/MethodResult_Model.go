package model

const METHOD_RESULT_MODEL_ID ModelID = "MethodResult"

type MethodResult_Model_row struct {
	AffectedRows int64 `json:"affected_rows"`
	Lsn string `json:"lsn"`
}

func New_MethodResult_Model(affectedRows int64, lsn string) Modeler{
	m := &Model{ID: METHOD_RESULT_MODEL_ID, Rows: make([]ModelRow, 1)}
	m.Rows[0] = &MethodResult_Model_row{AffectedRows: affectedRows, Lsn: lsn}
	return m
}

