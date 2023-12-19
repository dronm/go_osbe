package captcha

import (
		
	"osbe/fields"
)

type Captcha_get struct {
	Id fields.ValText `json:"id"`
	Width fields.ValInt `json:"width"`
	Height fields.ValInt `json:"height"`
	Count fields.ValInt `json:"count"`
}
type Captcha_get_argv struct {
	Argv *Captcha_get `json:"argv"`	
}

