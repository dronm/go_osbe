package validRusBanckAcc

import (
	"strings"
	"context"
	
	"github.com/jackc/pgx/v4"	
)

const (
	BIK_LEN = 9
	ACC_NUM_LEN = 20
	
	CHECK_BIK_QUERY = `SELECT TRUE FROM banks.banks WHERE bik = $1`
)

func isNotDigit(s string) bool {
	is_not_digit := func(c rune) bool { return c < '0' || c > '9' }
	return strings.IndexFunc(s, is_not_digit) == -1
}

func CheckBik(bik string, conn *pgx.Conn, ctx context.Context) bool {
	if len(bik) != BIK_LEN || !isNotDigit(bik) {
		return false
	}
	
	found := false
	if err := conn.QueryRow(ctx, CHECK_BIK_QUERY, bik).Scan(&found); err != nil || !found {
		return false
	} 	
	return true
}

func CheckAcc(bik string, accNum string) bool {
	if len(bik) != BIK_LEN || !isNotDigit(bik) {
		return false
	}
	if len(accNum) != ACC_NUM_LEN || !isNotDigit(accNum) {
		return false
	}
	
	s := bik[6 : 9] + accNum
	w := [23]int{7, 1, 3, 7, 1, 3, 7, 1, 3, 7, 1, 3, 7, 1, 3, 7, 1, 3, 7, 1, 3, 7, 1};
	sm := 0
	for i := 0; i <= 22; i++ {
		ch := int(s[i : i + 1][0] - '0')
		sm += ( ch * w[i] ) % 10;
	}
	
	if sm % 10 == 0 {
		return true
	}
	return false
}
