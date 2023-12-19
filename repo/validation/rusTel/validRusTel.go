package validRusTel

import "strings"

var valid_prefixes []string

func isNotDigit(s string) bool {
	is_not_digit := func(c rune) bool { return c < '0' || c > '9' }
	return strings.IndexFunc(s, is_not_digit) == -1
}

func initValidPrefixes() {
	if valid_prefixes == nil {
		valid_prefixes = []string{"900","901","902","903","904","905","906","908","909","910","911","912","913","914","915","916","917","918","919","920","921","922","923","924","925","926","927","928","929","930","931","932","933","934","936","937","938","939","941","950","951","952","953","954","955","956","958","960","961","962","963","964","965","966","967","968","969","970","971","980","981","982","983","984","985","987","988","989","991","992","993","994","995","996","997","999"}
	}
}

func GetValidPrefixes() []string {
	initValidPrefixes()
	return valid_prefixes
}

func Check(tel string) bool {
	initValidPrefixes()
	//all digits
	//len=11, if len=10, add first 7
	if len(tel) == 10 {
		tel = "7" + tel
	}
	if len(tel) != 11 {
		return false
	}
	
	if !isNotDigit(tel) {
		return false
	}
	
	for _, t := range valid_prefixes {
		if t == tel[1:4] {
			return true
		}
	}
	return false
}
