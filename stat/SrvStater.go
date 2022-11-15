package stat

type SrvStater interface {
	GetMaxClientCount() uint
	GetClientCount() uint
	GetDownloadedBytes() uint64
	IncDownloadedBytes(bt uint64)
	GetUploadedBytes() uint64
	IncUploadedBytes(bt uint64)
	GetHandshakes() uint64
	IncHandshakes()
	GetRunSeconds() uint64
	OnClientDisconnceted()
}
