package test_app

import (
	"errors"
	"encoding/json"	
	"reflect"
	"fmt"

	"osbe"
	"osbe/fields"
	"osbe/evnt"
	"osbe/srv"
	"osbe/socket"
	"osbe/model"
	"osbe/response"
)

//Controller
type Event_Controller struct {
	PublicMethods osbe.PublicMethodCollection
}

func (c *Event_Controller) GetID() osbe.ControllerID {
	return osbe.ControllerID("Event")
}

func (c *Event_Controller) InitPublicMethods() {
	c.PublicMethods = make(osbe.PublicMethodCollection)
	
	//*********************************** method subscribe
	c.PublicMethods["subscribe"] = &Event_Controller_subscribe{
		ModelMetadata: evnt.Get_Event_subscr_md(),
	}	
	
	//********************************** method unsubscribe
	c.PublicMethods["unsubscribe"] = &Event_Controller_unsubscribe{
		ModelMetadata: evnt.Get_Event_subscr_md(),
	}

	//********************************** method publish
	c.PublicMethods["publish"] = &Event_Controller_publish{
		ModelMetadata: evnt.Get_Event_publish_md(),
	}
}

func (c *Event_Controller) GetPublicMethod(publicMethodID osbe.PublicMethodID) (pm osbe.PublicMethod, err error) {
	pm, ok := c.PublicMethods[publicMethodID]
	if !ok {
		err = errors.New(fmt.Sprintf(osbe.ER_CONTOLLER_METH_NOT_DEFINED, string(publicMethodID), string(c.GetID())))
	}
	
	return
}

type Event_Controller_event_argv struct {
	Argv evnt.Event_subscr `json:"argv"`	
}

type Event_Controller_publish_argv struct {
	Argv evnt.Event_publish `json:"argv"`	
}

//******************* subscribe ***************************************
type Event_Controller_subscribe struct {
	ModelMetadata fields.FieldCollection
	EventList osbe.PublicMethodEventList
}

func (pm *Event_Controller_subscribe) GetEventList() osbe.PublicMethodEventList {
	return pm.EventList
}

func (pm *Event_Controller_subscribe) AddEvent(evId string) {
	pm.EventList[len(pm.EventList)-1] = evId
}

func (pm *Event_Controller_subscribe) GetModelMetadata() fields.FieldCollection {
	return pm.ModelMetadata
}

func (pm *Event_Controller_subscribe) GetFields() fields.FieldCollection{
	return pm.ModelMetadata
}

func (c *Event_Controller_subscribe) GetID() osbe.PublicMethodID {
	return osbe.PublicMethodID("subscribe")
}


//Public method Unmarshal to structure
func (pm *Event_Controller_subscribe) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &Event_Controller_event_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//Method implemenation
func (pm *Event_Controller_subscribe) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {

	args := rfltArgs.Interface().(evnt.Event_subscr)
	ev_sock, ok := sock.(*evnt.EvntSocket)
	if !ok {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, ("Event_Controller_subscribe.Run(): cannot cast socket to *evnt.EvntSocket"))
	}	
	app.GetEvntServer().SubsrAddDbListener(&args, ev_sock)
	return nil
}


//******************* subscribe ***************************************


//******************* unsubscribe ***************************************
type Event_Controller_unsubscribe struct {
	ModelMetadata fields.FieldCollection
	EventList osbe.PublicMethodEventList
}

func (pm *Event_Controller_unsubscribe) GetEventList() osbe.PublicMethodEventList {
	return pm.EventList
}

func (pm *Event_Controller_unsubscribe) AddEvent(evId string) {
	pm.EventList[len(pm.EventList)-1] = evId
}

func (pm *Event_Controller_unsubscribe) GetModelMetadata() fields.FieldCollection {
	return pm.ModelMetadata
}

func (pm *Event_Controller_unsubscribe) GetFields() fields.FieldCollection{
	return pm.ModelMetadata
}

func (c *Event_Controller_unsubscribe) GetID() osbe.PublicMethodID {
	return osbe.PublicMethodID("unsubscribe")
}

//Public method Unmarshal to structure
func (pm *Event_Controller_unsubscribe) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &Event_Controller_event_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}
//Method implementations
func (pm *Event_Controller_unsubscribe) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	args := rfltArgs.Interface().(evnt.Event_subscr)
	ev_sock, ok := sock.(*evnt.EvntSocket)
	if !ok {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, ("Event_Controller_unsubscribe.Run(): cannot cast socket to *evnt.EvntSocket"))
	}		
	app.GetEvntServer().SubsrRemoveDbListener(&args, ev_sock)
	
	return nil
}
//******************* unsubscribe ***************************************


//******************* publish ***************************************
//Public method Event
type Event_Controller_publish struct {
	ModelMetadata fields.FieldCollection
	EventList osbe.PublicMethodEventList
}

func (pm *Event_Controller_publish) GetEventList() osbe.PublicMethodEventList {
	return pm.EventList
}

func (pm *Event_Controller_publish) AddEvent(evId string) {
	pm.EventList[len(pm.EventList)-1] = evId
}

func (pm *Event_Controller_publish) GetModelMetadata() fields.FieldCollection {
	return pm.ModelMetadata
}

func (pm *Event_Controller_publish) GetFields() fields.FieldCollection{
	return pm.ModelMetadata
}

func (c *Event_Controller_publish) GetID() osbe.PublicMethodID {
	return osbe.PublicMethodID("publish")
}

//Public method Unmarshal to structure
func (pm *Event_Controller_publish) Unmarshal(payload []byte) (res reflect.Value, err error) {
	//argument structrure
	argv := &Event_Controller_publish_argv{}

	err = json.Unmarshal(payload, argv)
	if err != nil {
		return
	}
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}
func (pm *Event_Controller_publish) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {

	args := rfltArgs.Interface().(evnt.Event_publish)
	
	if _, ok := app.GetEvntServer().Events.GetEvent(args.Id.GetValue()); ok {
		evnt_emitter_id := ""
		assoc_params := args.Params.GetValue()
		if v, ok := assoc_params[evnt.EVNT_PARAM_EMITTER_ID]; ok  {
			if v_str, ok := v.(string); ok {
				evnt_emitter_id = v_str
			}
		}
		
		logger := app.GetLogger()
		//event exists
		logger.Debugf("publishEvent emitter_id: %s", evnt_emitter_id)
		
		/*loc_ev := app.GetLocalEvents()
		if loc_ev != nil {
			if loc_ev_pm, ok := loc_ev[args.Id.GetValue()]; ok {
				//Здесь надо вызвать локальный метод с передачей всех параметров
			}
		}*/		
		for _, app_serv := range app.GetServers() {			
			for sock_item := range app_serv.GetClientSockets().Iter() {			
				sock, ok := sock_item.Socket.(*evnt.EvntSocket)
				if !ok {
					logger.Error("Event_Controller_publish.Run(): cannot cast socket to *evnt.EvntSocket")
					continue
				}
				token := sock.GetToken()
				if token == "" || evnt_emitter_id != token {
					//all events of other emitters
					if ok := sock.Events.HasEvent(args.Id.GetValue()); ok {						
						sendEventToClient(app, serv, sock, args)
					}
				}
			}
		}
	}
	
	return nil
}

//Sends data to a specified socket
func sendEventToClient(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, args evnt.Event_publish) {
	resp := response.NewResponse("", app.GetMD().Version.Value)				
	m := &model.Model{ID: evnt.EVNT_MODEL_ID, Rows: make([]model.ModelRow,1)}
	m.Rows[0] = &evnt.Event_publish{
		Id: args.Id,
		Params: args.Params,
	}
	resp.AddModel(m)
	app.SendToClient(serv, sock, resp, "json")

}
