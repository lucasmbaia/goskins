package steamland

type JobId uint64

type MessageBody interface {
	Serializable
	GetEMsg() EMsg
}

func (j JobId) String() string {
	if j == math.MaxUint64 {
		return "(none)"
	}
	return strconv.FormatUint(uint64(j), 10)
}

