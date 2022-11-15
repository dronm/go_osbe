package evnt

import (
)

//subscribe
//fields.ValArray
type Event_subscr struct {
	Events []Event `json:"events" required:"true"`
}

type Event_subscr_argv struct {
	Argv *Event_subscr `json:"argv"`	
}

//Common structure for subscribe/unsubscribe
//fields.ValAssocArray
type Event struct {
	Id string `json:"id" required:"true" length:100`
	Params map[string]interface{} `json:"params"`
}

type Event_argv struct {
	Argv *Event `json:"argv"`	
}

