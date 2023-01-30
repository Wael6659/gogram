package session

import (
	"encoding/base64"
	"strings"
)

const (
	Prefix = "1Bvq"
)

type (
	StringSession struct {
		AuthKey     []byte
		AuthKeyHash []byte
		DCID        int
		IpAddr      string
		AppID       int32
		AppHash     string
		Encoded     string
	}
)

func (s StringSession) Encode() []byte {
	return []byte(Prefix + base64.RawURLEncoding.EncodeToString(s.AuthKey))
}

func (s StringSession) AppendAllValues() string {
	return string(s.GetAuthKey()) + "::" + string(s.GetAuthKeyHash()) + "::" + string(s.IpAddr) + "::" + string(rune(s.DCID)) + "::" + string(s.AppHash) + "::" + string(rune(s.AppID))
}

func (s StringSession) GetAuthKey() string {
	return string(s.AuthKey)
}

func (s StringSession) GetAuthKeyHash() string {
	return string(s.AuthKeyHash)
}

func (s StringSession) EncodeToString() string {
	return Prefix + base64.RawURLEncoding.EncodeToString([]byte(s.AppendAllValues()))
}

func (s StringSession) GetDCID() int {
	return s.DCID
}

func (s StringSession) Decode() (AuthKey, AuthKeyHash []byte, DcID int, IpAddr string, AppHash string, AppID int32, err error) {
	Decoded, err := base64.RawURLEncoding.DecodeString(s.Encoded[len(Prefix):])
	if err != nil {
		return nil, nil, 0, "", "", 0, err
	}
	DecodedString := string(Decoded)
	AuthKey, AuthKeyHash, IpAddr, DcID, AppHash, AppID = s.SplitValues(DecodedString)
	return AuthKey, AuthKeyHash, DcID, IpAddr, AppHash, AppID, nil
}

func (s StringSession) SplitValues(DecodedString string) (AuthKey, AuthKeyHash []byte, IpAddr string, DcID int, AppHash string, AppID int32) {
	Sep := strings.Split(DecodedString, "::")
	if len(Sep) != 4 {
		return nil, nil, "", 0, "", 0
	}
	AuthKey = []byte(Sep[0])
	AuthKeyHash = []byte(Sep[1])
	IpAddr = Sep[2]
	DcID = int(Sep[3][0])
	AppHash = Sep[4]
	AppID = int32(Sep[5][0])
	return AuthKey, AuthKeyHash, IpAddr, DcID, AppHash, AppID
}
