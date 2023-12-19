package model

import (
	"encoding/json"
	"fmt"
)

type AggFunctionValue struct {	
	Alias string
	Val interface{}
	ValStr string
}

/*const(
	MODEL_SORT_ORDER_ASC ModelSortOrder = iota
	MODEL_SORT_ORDER_DESC
)*/

//type ModelSortOrder int

type ModelID string

type ModelCollection map[ModelID]Modeler
//*Model

type ModelRow interface{}

type Modeler interface {
	MarshalJSON() ([]byte, error)
	GetID() ModelID
	GetRowsPerPage() int
	SetRowsPerPage(int)
	SetListFrom(int)
	GetListFrom() int
	GetSysModel() bool
	GetRow(int) ModelRow
	GetRows() []ModelRow
	GetRowCount() int
	GetRawData() []byte
	GetAggFunctionValues() []*AggFunctionValue
	GetMetadata() *ModelMD
}

type Model struct {
	ID ModelID `json:"id"`
	//DataTable string `json:"dataTable"`
	//DefaultOrder map[string]ModelSortOrder `json:"defaultOrder"`
	SysModel bool `json:"sysModel"`
	RowsPerPage int `json:"rowsPerPage"`
	ListFrom int `json:"listFrom"`
	//TotalCount int `json:"totalCount"`
	RawData []byte `json:"-"`
	Rows []ModelRow `json:"rows"`
	AggFunctionValues []*AggFunctionValue
	Metadata *ModelMD `json:"-"`
}

func (m *Model) GetID() ModelID {
	return m.ID
}

func (m *Model) GetRowsPerPage() int {
	return m.RowsPerPage
}

func (m *Model) SetRowsPerPage(cnt int) {
	m.RowsPerPage = cnt
}

func (m *Model) GetListFrom() int {
	return m.ListFrom
}

func (m *Model) SetListFrom(cnt int) {
	m.ListFrom = cnt
}

func (m *Model) GetSysModel() bool {
	return m.SysModel
}

func (m *Model) GetRow(ind int) ModelRow {
	return m.Rows[ind]
}

func (m *Model) AddRow(row ModelRow) {
	m.Rows = append(m.Rows, row)
}

func (m *Model) GetRowCount() int {
	return len(m.Rows)
}

func (m *Model) GetRows() []ModelRow {
	return m.Rows
}

func (m *Model) GetRawData() []byte {
	return m.RawData
}

func (m *Model) GetMetadata() *ModelMD {
	return m.Metadata
}

func (m *Model) GetAggFunctionValues() []*AggFunctionValue {
	return m.AggFunctionValues
}

func NewModel(id string) *Model{
	return &Model{ID: ModelID(id), Rows:make([]ModelRow, 0)}
}
/*func (m *Model) MarshalJSON() ([]byte, error) {
	return json.Marshal(m)
}*/


//Short variant to be used with data
type ModelData struct {
	ID ModelID `json:"id"`
	Rows []ModelRow `json:"rows"`
}

type ModelDataTotalCount struct {
	ID ModelID `json:"id"`
	TotalCount int `json:"totalCount"`
	Rows []ModelRow `json:"rows"`
}

//custom marshal, skeeping certain fields

func (m *Model) MarshalJSON() ([]byte, error) {
	return ModelMarshalJSON(m)
}
/*
 * @ToDo
 * AggFunctions []*AggFunction{Alias, Value interface{}??}
 * util.QueryResultToModel &model.Model{}
 * m.AddAggFunction("totCount", val)
 * m.AddAggFunction("totTotal", val)
 * ModelMarshalJSON custom, add all agg functions
 * {"id": "m.GetID()", ALL_AGG_FUNCTIONS, "rows": json.Marshal(m.GetRows()) }
 * ну или создать map[string]interface{id, ALL_AGG_FUNCTIONS, rows:m.GetRows()}
 *		json.Marshal(map)
 * ModelMarshalXML - обавить сюда же???
 */
func ModelMarshalJSON(m Modeler) ([]byte, error) {
	raw := m.GetRawData()
	if len(raw) > 0 {
		return raw, nil
		
	}else {
		agg_funcs_s := ""
		agg_vals := m.GetAggFunctionValues()
		if len(agg_vals) >0 {
			for _,agg_v := range agg_vals {
				agg_funcs_s += fmt.Sprintf(`, "%s": %v`, agg_v.Alias, agg_v.Val)
			}
		}
		rows_b, err := json.Marshal(m.GetRows())
		if err != nil {
			return nil, err
		}
		return []byte(fmt.Sprintf(`{"id": "%s"%s, "rows": %s}`, string(m.GetID()), agg_funcs_s, string(rows_b))), nil
	}	
	
	/*
	cnt := m.GetTotalCount()
	if len(raw) > 0 {
		return raw, nil
		
	}else if cnt > 0{
		return json.Marshal(&ModelDataTotalCount{
			ID: m.GetID(),
			TotalCount: cnt,
			Rows: m.GetRows(),
		})
		
	}else{
		return json.Marshal(&ModelData{
			ID: m.GetID(),
			Rows: m.GetRows(),
		})
	}
	*/
}

