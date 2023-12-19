package menu

//Controller method model
import (
	"osbe/model"
	"osbe/fields"
)

type View_complete_argv struct {
	Argv *View_complete `json:"argv"`	
}

//
type View_complete struct {
	model.Complete_Model
	User_descr fields.ValText `json:"user_descr" matchField:"true" required:"true"`
}
