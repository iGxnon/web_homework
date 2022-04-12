package ws_parser

import (
	"errors"
)

type OpenCodeType uint8

type Payload struct {
	rawPayload []byte

	// 2nd byte
	IsMasked   bool
	PayloadLen int64 // true length

	// MaskKey max length is 4 bytes, it is 0 if IsMask is false
	MaskKey [4]byte
}

var (
	CheckCodeAvailable = func(code OpenCodeType) bool {
		return code <= BinaryDataType || (code >= CloseConnType && code <= PongType)
	}
)

const (
	ContinuationDataType OpenCodeType = iota
	TextDataType
	BinaryDataType
	_
	_
	_
	_
	_
	CloseConnType
	PingType
	PongType
)

type WebsocketDataFrame struct {

	// raw raw bytes data
	raw []byte

	// 1st byte
	FIN      bool
	RSV      [3]byte
	OpenCode OpenCodeType

	// Payload
	Payload Payload
}

func ParseData(data []byte) (*WebsocketDataFrame, error) {
	ret := &WebsocketDataFrame{}
	err := ret.parse(data)
	return ret, err
}

// WrapperPayload only wrapper one frame!
// please make sure data length size no more than 8bytes
func WrapperPayload(payloadData []byte, code OpenCodeType) []byte {
	frame := make([]byte, 0)
	frame = append(frame, getFrameHeader(true, code, len(payloadData))...)
	frame = append(frame, payloadData...)
	return frame
}

// WrapperPayloads wrappers give splits frames and set OpenCode to 0x0
func WrapperPayloads(payloadData []byte, splits int) [][]byte {
	frames := make([][]byte, splits)
	for i := 0; i < splits; i++ {
		if i == splits-1 {
			length := len(payloadData)/splits + len(payloadData)%splits
			frames[i] = append(frames[i], getFrameHeader(true, ContinuationDataType, length)...)
			frames[i] = append(frames[i], payloadData[len(payloadData)-length:]...)
			break
		}
		length := len(payloadData) / splits
		frames[i] = append(frames[i], getFrameHeader(i == splits-1, ContinuationDataType, length)...)
		frames[i] = append(frames[i], payloadData[i*length:(i+1)*length]...)
	}
	return frames
}

func (p *Payload) GetData() []byte {
	truePayload := make([]byte, len(p.rawPayload))
	// unmask
	if p.IsMasked {
		for i := 0; i < len(p.rawPayload); i++ {
			j := i % 4
			truePayload[i] = p.rawPayload[i] ^ p.MaskKey[j]
		}
	} else {
		truePayload = p.rawPayload
	}
	return truePayload
}

// parse raw bits
func (w *WebsocketDataFrame) parse(data []byte) error {
	w.raw = data
	firstByte := data[0]
	w.FIN = (firstByte >> 7) == 1
	w.RSV = [3]byte{
		((1 << 6) & firstByte) >> 6,
		((1 << 5) & firstByte) >> 5,
		((1 << 4) & firstByte) >> 4,
	}
	w.OpenCode = OpenCodeType(15 & firstByte)
	if !CheckCodeAvailable(w.OpenCode) {
		return errors.New("invalid open code")
	}

	secondByte := data[1]
	w.Payload.IsMasked = (secondByte >> 7) == 1
	payloadLen := 127 & secondByte

	maskBegin := 2
	if payloadLen == 127 { // the next 8 bytes is length
		for i := 2; i <= 9; i++ {
			w.Payload.PayloadLen += int64(data[i]) << ((9 - i) * 8)
		}
		maskBegin = 10
	} else if payloadLen == 126 { // the next 2 bytes is length
		w.Payload.PayloadLen = int64(data[2])<<8 + int64(data[3])
		maskBegin = 4
	} else {
		w.Payload.PayloadLen = int64(payloadLen)
	}

	payLoadBegin := maskBegin
	if w.Payload.IsMasked { // if is masked, read the next 4 bytes next to maskBegin
		w.Payload.MaskKey = [4]byte{}
		for i := maskBegin; i < maskBegin+4; i++ {
			w.Payload.MaskKey[i-maskBegin] = data[i]
		}
		payLoadBegin = maskBegin + 4
	}
	// rawPayload
	w.Payload.rawPayload = data[payLoadBegin : payLoadBegin+int(w.Payload.PayloadLen)]
	w.raw = data[:payLoadBegin+int(w.Payload.PayloadLen)]

	// should not check
	//if int64(len(w.Payload.rawPayload)) != w.Payload.PayloadLen {
	//	return errors.New("real payload length cannot match the frame header gives")
	//}
	return nil
}

// getFrameHeader does not support mask,
// dataframe sent by server does not need to be masked
func getFrameHeader(fin bool, code OpenCodeType, length int) []byte {
	header := make([]byte, 2)
	// 1st byte
	if fin {
		header[0] = 1 << 7
	}
	header[0] |= byte(code)

	// 2nd byte
	maskLength := length
	if length > ((1 << 16) - 1) { // 2 bytes cannot contains
		maskLength = 127
	} else if length >= ((1 << 7) - 1) { // 7 bits cannot contains
		maskLength = 126
	}
	header[1] = byte(maskLength)

	if maskLength == 127 {
		header = append(header, make([]byte, 8)...)
		for i := 9; i >= 2; i-- {
			header[i] = byte(length >> (8 * (9 - i)))
		}
	} else if maskLength == 126 {
		header = append(header, make([]byte, 2)...)
		header[3] = byte(length)
		header[2] = byte(length >> 8)
	}

	return header
}
