package evnt

import (
	"encoding/json"
	"reflect"
	
	"osbe"
	"osbe/srv"
	"osbe/socket"
	"osbe/response"
	"osbe/model"
)

/*
select pg_notify('Vendor.delete',
json_build_object('params',json_build_object('emitterId','1111'))::text)
*/

//Controller
type Event_Controller struct {
	osbe.Base_Controller
	EvntServer *EvntSrv
}

func (c *Event_Controller) InitPublicMethods() {
	c.ID = "Event"
	
	Event_Model_init()

	c.PublicMethods = make(osbe.PublicMethodCollection)
	
	//************************** method subscribe **********************************
	c.PublicMethods["subscribe"] = &Event_Controller_subscribe{
		osbe.Base_PublicMethod{
			ID: "subscribe",
			ModelMetadata: Event_subscr_md,
		},
		c,
	}

	//************************** method unsubscribe **********************************
	c.PublicMethods["unsubscribe"] = &Event_Controller_unsubscribe{
		osbe.Base_PublicMethod{
			ID: "unsubscribe",
			ModelMetadata: Event_subscr_md,
		},
		c,
	}

	//************************** method publish **********************************
	c.PublicMethods["publish"] = &Event_Controller_publish{
		osbe.Base_PublicMethod{
			ID: "publish",
			ModelMetadata: Event_publish_md,
		},
	}
}

//**************************************************************************************
//Public method: subscribe
type Event_Controller_subscribe struct {	
	osbe.Base_PublicMethod
	Contr *Event_Controller
}
//Public method Unmarshal to structure
func (pm *Event_Controller_subscribe) Unmarshal(payload []byte) (reflect.Value, error) {

	//argument structrure
	var res reflect.Value
	argv := &Event_subscr_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	return res, nil
}

//custom method
func (pm *Event_Controller_subscribe) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return handleSubscription(pm.Contr, sock, rfltArgs, true)
}

//**************************************************************************************
//Public method: unsubscribe
type Event_Controller_unsubscribe struct {
	osbe.Base_PublicMethod
	Contr *Event_Controller
}
//Public method Unmarshal to structure
func (pm *Event_Controller_unsubscribe) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &Event_subscr_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}
	res = reflect.ValueOf(&argv.Argv).Elem()
	return res, nil
}

//custom method
func (pm *Event_Controller_unsubscribe) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return handleSubscription(pm.Contr, sock, rfltArgs, false)
}

//**************************************************************************************
//Public method: publish
type Event_Controller_publish struct {
	osbe.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *Event_Controller_publish) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &Event_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}
	res = reflect.ValueOf(&argv.Argv).Elem()
	return res, nil
}

//custom method
func (pm *Event_Controller_publish) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	
	args := rfltArgs.Interface().(*Event)	
	emitter_id := ""
	if v, ok := args.Params[EVNT_PARAM_EMITTER_ID]; ok {
		emitter_id, ok = v.(string)
		if !ok {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, "emitterId not a string")		
		}
		delete(args.Params, EVNT_PARAM_EMITTER_ID)
	}
	//response object for other clients
	sock_resp := response.NewResponse("", app.GetMD().Version.Value)
	m := &model.Model{ID: EVNT_MODEL_ID, Rows: make([]model.ModelRow, 1)}	
	m.Rows[0] = &Event{Id: args.Id, Params: args.Params}
	sock_resp.AddModel(m)
	
	//iterate all sockets of all (ws or tcp) servers and send to all interested in event
	for _, s := range app.GetServers() {
		sock_list := s.GetClientSockets()
		if sock_list == nil {
			continue
		}
		for sock_item := range sock_list.Iter() {
			 if sock_item.Socket.GetID() != emitter_id {
			 	if sock_item_s, ok := sock_item.Socket.(*EvntSocket); ok {
				 	if ok := sock_item_s.Events.HasEvent(args.Id); ok {
					 	app.SendToClient(s, sock_item.Socket, sock_resp, "ViewJSON")
					}
				}
			 }
		}
	}
	
	
	return nil
}

func handleSubscription(contr *Event_Controller, sock socket.ClientSocketer, rfltArgs reflect.Value, addEvent bool) error {
	args := rfltArgs.Interface().(*Event_subscr)
	for _,ev := range args.Events {
		s, ok := sock.(*EvntSocket)
		if !ok {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, "client socket does not implement EvntSocket")	
		}
		if addEvent {			
			contr.EvntServer.AddDbListener(ev.Id, s)
		}else{
			contr.EvntServer.RemoveDbListener(ev.Id, s)
		}
	}
	return nil

}

