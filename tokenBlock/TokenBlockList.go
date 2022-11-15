package tokenBlock

import (
	"sync"
	"time"
)

type TokenBlockList struct {
	mx sync.Mutex
	m map[string]time.Time
	cnt int
	expirTime time.Duration
	maxCount int
}

func (l *TokenBlockList) Contains(token string) bool{
	l.mx.Lock()
	_,ok := l.m[token]
	if ok {
		l.m[token] = time.Now()
	}
	l.mx.Unlock()
	return ok
}

func (l *TokenBlockList) Append(token string){
	l.mx.Lock()
	defer l.mx.Unlock()
		
	if l.cnt >= l.maxCount {
		//clearing expired...
		for tkn,tm := range l.m {
			if l.timeExpired(tm) {
				delete(l.m,tkn)
				l.cnt =- 1
			}
		}
	}
	
	if l.cnt < l.maxCount {
		l.cnt =+ 1
		l.m[token] = time.Now()
	}
}

func (l *TokenBlockList) timeExpired(tm time.Time) bool {
	n := time.Now()
	if n.Sub(tm) >= l.expirTime{
		return true
	}
	return false
}

func NewTokenBlockList() *TokenBlockList{
	return &TokenBlockList{m: make(map[string]time.Time)}
}
