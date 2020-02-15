package steam

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"time"

)

var (
	steamGuardCodeTranslations = []byte{50, 51, 52, 53, 54, 55, 56, 57, 66, 67, 68, 70, 71, 72, 74, 75, 77, 78, 80, 81, 82, 84, 86, 87, 88, 89}
)

type SteamGuard struct {
	SharedSecret	string	`json:"shared_secret,omitempty"`
	SerialNumber	string	`json:"serial_number,omitempty"`
	RevocationCode	string	`json:"revocation_code,omitempty"`
	URI		string	`json:"uri,omitempty"`
	ServerTime	string	`json:"server_time,omitempty"`
	AccountName	string	`json:"account_name,omitempty"`
	TokenGID	string	`json:"token_gid,omitempty"`
	IdentitySecret	string	`json:"identity_secret,omitempty"`
	Secret		string	`json:"secret_1,omitempty"`
	Status		int	`json:"status,omitempty"`
	DeviceID	string	`json:"device_id,omitempty"`
	FullyEnrolled	bool	`json:"fully_enrolled"`
}

func (s *Session) GenerateSteamGuardCode(secret string) (code string, err error) {
	var (
		t		time.Time
		sharedSecret	[]byte
		buf		= new(bytes.Buffer)
		hashedData	[]byte
		codeBytes	= make([]byte, 5)
	)

	t = s.GetSteamTime()
	if sharedSecret, err = base64.StdEncoding.DecodeString(secret); err != nil {
		return
	}

	binary.Write(buf, binary.BigEndian, t.Unix()/30)
	var mac = hmac.New(sha1.New, sharedSecret)
	mac.Write(buf.Bytes())
	hashedData = mac.Sum(nil)

	var b = byte(hashedData[19] & 0xF)
	var codePoint = int(hashedData[b]&0x7F)<<24 | int(hashedData[b+1]&0xFF)<<16 | int(hashedData[b+2]&0xFF)<<8 | int(hashedData[b+3]&0xFF)
	var translationCount = len(steamGuardCodeTranslations)

	for i := 0; i < 5; i++ {
		codeBytes[i] = steamGuardCodeTranslations[codePoint%translationCount]
		codePoint /= translationCount
	}

	code = string(codeBytes)
	return
}
