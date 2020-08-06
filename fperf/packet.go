package fperf

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/fregie/gotool/freconn"
)

type CtrlType uint8

func (t *CtrlType) Uint8() uint8 { return uint8(*t) }
func ParseType(t uint8) CtrlType { return CtrlType(t) }

const (
	CtrlType_Unknown CtrlType = iota
	START
	DATA
	FIN
	STAT
)

const CtrlHeaderLen = 12 //Byte
type CtrlHeader struct {
	CtrlType      CtrlType
	Reserved      [7]byte
	PayloadLength uint32
	Data          []byte
}

func (h *CtrlHeader) Pack() []byte {
	packed := make([]byte, CtrlHeaderLen+len(h.Data))
	packed[0] = h.CtrlType.Uint8()
	h.PayloadLength = uint32(len(h.Data))
	binary.BigEndian.PutUint32(packed[8:], h.PayloadLength)
	copy(packed[CtrlHeaderLen:], h.Data)
	return packed
}

func PacketFromReader(r io.Reader) (*CtrlHeader, error) {
	headerBuffer := make([]byte, CtrlHeaderLen)
	_, err := io.ReadFull(r, headerBuffer)
	if err != nil {
		return nil, fmt.Errorf("read ctrl header failed:%w", err)
	}
	h := &CtrlHeader{
		CtrlType:      ParseType(headerBuffer[0]),
		PayloadLength: binary.BigEndian.Uint32(headerBuffer[8:]),
	}
	h.Data = make([]byte, h.PayloadLength)
	if h.PayloadLength > 0 {
		_, err := io.ReadFull(r, h.Data)
		if err != nil {
			return nil, err
		}
	}
	return h, nil
}

type initData struct {
	Start        int64
	TestDuration int64
}

func (f *Fperf) sendStart() error {
	start := time.Now().Unix()
	jsonByte, err := json.Marshal(&initData{Start: start, TestDuration: int64(f.TestDuration.Seconds())})
	if err != nil {
		return err
	}
	header := &CtrlHeader{
		CtrlType: START,
		Data:     jsonByte,
	}
	_, err = f.CtrlConn.Write(header.Pack())
	if err != nil {
		return err
	}
	return nil
}

func (f *Fperf) recvStart() (*initData, error) {
	header, err := PacketFromReader(f.CtrlConn)
	if err != nil {
		return nil, err
	}
	if header.CtrlType != START {
		return nil, errors.New("first ctrl type not start")
	}
	var data initData
	err = json.Unmarshal(header.Data, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *Fperf) sendFIN() error {
	header := &CtrlHeader{
		CtrlType: FIN,
		Data:     make([]byte, 0),
	}
	_, err := f.CtrlConn.Write(header.Pack())
	if err != nil {
		return err
	}
	return nil
}

func (f *Fperf) recvFIN() error {
	header, err := PacketFromReader(f.CtrlConn)
	if err != nil {
		return nil
	}
	if header.CtrlType != FIN {
		return errors.New("first ctrl type not start")
	}
	return nil
}

func (f *Fperf) sendStat() error {
	jsonByte, err := json.Marshal(f.Stat)
	if err != nil {
		return err
	}
	header := &CtrlHeader{
		CtrlType:      STAT,
		PayloadLength: uint32(len(jsonByte)),
		Data:          jsonByte,
	}
	_, err = f.CtrlConn.Write(header.Pack())
	if err != nil {
		return err
	}
	return nil
}

func (f *Fperf) recvStat() error {
	header, err := PacketFromReader(f.CtrlConn)
	if err != nil {
		return err
	}
	if header.CtrlType != STAT {
		return errors.New("wrong ctrl type,need STAT")
	}
	f.PeerStat = &freconn.Stat{}
	err = json.Unmarshal(header.Data, f.PeerStat)
	if err != nil {
		return err
	}
	return nil
}
