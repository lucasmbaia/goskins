package steamland

import (
	"encoding/binary"
	"strconv"
	"math"
	"bytes"
)

type Packet struct {
	EMsg	    EMsg
	IsProto	    bool
	TargetJobId JobId
	SourceJobId JobId
	Data	    []byte
}

func NewPacket(data []byte) (p *Packet, err error) {
	var (
		t	uint32
		emsg	EMsg
		buffer	*bytes.Reader
	)

	if err = binary.Read(bytes.NewReader(data), binary.LittleEndian, &t); err != nil {
		return
	}

	emsg = NewEMsg(t)
	buffer = bytes.NewReader(data)

	if emsg == EMsg_ChannelEncryptRequest || emsg == EMsg_ChannelEncryptResult {
		var header *MsgHdr = NewMsgHdr()
		header.Msg = emsg

		if err = header.Deserialize(buffer); err != nil {
			return
		}

		p = &Packet{
			EMsg:		emsg,
			IsProto:	true,
			TargetJobId:	JobId(header.Proto.GetJobidTarget()),
			SourceJobId:	JobId(header.Proto.GetJobidSource()),
			Data:		data,
		}
	} else {
		fmt.Println("OUTRO TIPO")
	}

	return
}

func (p *Packet) ReadMsg(b MessageBody) *Msg {

}
