package stat

import (
	"sync"
	"time"
)

type SrvStat struct {
	mx sync.RWMutex
	startTime time.Time
	maxClientCount uint
	clientCount uint
	downloadedBytes uint64
	uploadedBytes uint64
	handshakes uint64
	//requests uint64
}

func (s SrvStat) GetMaxClientCount() uint {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.maxClientCount
}

func (s SrvStat) GetClientCount() uint {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.clientCount
}

func (s *SrvStat) OnClientDisconnceted() {
	s.mx.Lock()
	s.clientCount--
	s.mx.Unlock()
}

func (s SrvStat) GetDownloadedBytes() uint64 {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.downloadedBytes
}

func (s *SrvStat) IncDownloadedBytes(bt uint64) {
	s.mx.Lock()
	s.downloadedBytes =+ bt
	s.mx.Unlock()
}
	
func (s SrvStat) GetUploadedBytes() uint64 {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.uploadedBytes
}
	
func (s *SrvStat) IncUploadedBytes(bt uint64) {
	s.mx.Lock()
	s.uploadedBytes =+ bt
	s.mx.Unlock()
}
	
func (s SrvStat) GetHandshakes() uint64 {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.handshakes
}
	
func (s *SrvStat) IncHandshakes() {
	s.mx.Lock()
	s.handshakes++
	s.clientCount++
	if s.maxClientCount < s.clientCount {
		s.maxClientCount = s.clientCount
	}
	s.mx.Unlock()
}


func (s SrvStat) GetRunSeconds() uint64{
	s.mx.Lock()
	defer s.mx.Unlock()
	return uint64(time.Now().Sub(s.startTime).Seconds())

}

func NewSrvStat() *SrvStat{
	return &SrvStat{startTime: time.Now()}
}
