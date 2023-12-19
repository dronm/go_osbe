package contact

//Controller method model
import (
	"reflect"
	
	"osbe/model"
	"osbe/fields"
)

type Post_complete_argv struct {
	Argv *Post_complete `json:"argv"`	
}

//Exported model metadata
var Post_complete_md fields.FieldCollection

func Post_complete_Model_init() {	
	Post_complete_md = fields.GenModelMD(reflect.ValueOf(Post_complete{}))
}

//
type Post_complete struct {
	/*Count fields.ValInt `json:"count" default:10`	
	Ic fields.ValInt `json:"ic" default:1 minValue:0 maxValue:1`
	Mid fields.ValInt `json:"mid" default:1 minValue:0 maxValue:1`
	Ord_directs fields.ValText `json:"ord_directs" length:500`
	Field_sep fields.ValText `json:"field_sep" length:2`
	*/
	model.Complete_Model
	Name fields.ValText `json:"name" matchField:"true" required:"true"`
}
