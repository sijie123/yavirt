package util

const (
	Byte int64 = 1 << (10 * iota)
	KB
	MB
	GB
)

func ConvToMB(b int64) int64 {
	return b / MB
}
