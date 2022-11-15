package stat

import (
	"reflect"
	
	"osbe"
	"osbe/model"
	"osbe/fields"
	"osbe/srv"
	"osbe/socket"
	"osbe/response"
	
)

//Controller
type SrvStatistics_Controller struct {
	osbe.Base_Controller
}

func NewSrvStatistics_Controller() *SrvStatistics_Controller{
	c := &SrvStatistics_Controller{osbe.Base_Controller{ID: "SrvStatistics", PublicMethods: make(osbe.PublicMethodCollection)}}
	
	//************************** method get_statistics **********************************
	c.PublicMethods["get_statistics"] = &SrvStatistics_Controller_get_statistics{
		osbe.Base_PublicMethod{
			ID: "get_statistics",
		},
	}
	
	return c
}

//**************************************************************************************
//Public method: get_statistics
type SrvStatistics_Controller_get_statistics struct {
	osbe.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *SrvStatistics_Controller_get_statistics) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	return res, nil
}

type statServer interface{
	GetStatistics() SrvStater
}

//custom method
func (pm *SrvStatistics_Controller_get_statistics) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	rows := make([]model.ModelRow, 0)
	for srv_name, srv := range app.GetServers() {
		if stat_srv, ok := srv.(statServer); ok {
			stat := stat_srv.GetStatistics()
			m_row := &SrvStatistics{Name: fields.ValText{TypedValue: srv_name},
				Max_client_count: fields.ValInt{TypedValue: int64(stat.GetMaxClientCount())},
				Client_count: fields.ValInt{TypedValue: int64(stat.GetClientCount())},
				Downloaded_bytes: fields.ValUint{TypedValue: stat.GetDownloadedBytes()},
				Uploaded_bytes: fields.ValUint{TypedValue: stat.GetUploadedBytes()},
				Handshakes: fields.ValUint{TypedValue: stat.GetHandshakes()},
				Run_seconds: fields.ValUint{TypedValue: stat.GetRunSeconds()},				
			}
			rows = append(rows, m_row)
		}
	}
	resp.AddModel(&model.Model{ID: model.ModelID(app.GetMD().Models["SrvStatistics"].ID), Rows: rows})
	return nil
}

