package ws_parser

import (
	"fmt"
	"testing"
)

func TestParseData1(t *testing.T) {
	dataFrameBytes := []byte{
		byte(0b1_0_0_0_0001), // FIN 1 OpenCode 1
		byte(0b1_1111111),    // MASK 1 PayloadLen 127
		// the next 8 bytes is payload length (bytes)
		byte(0b00000000),
		byte(0b00000000),
		byte(0b00000000),
		byte(0b00000000),
		byte(0b00000000),
		byte(0b00000000),
		byte(0b00000000),
		byte(0b11111111), // payload length is 255 bytes
		// mask key
		byte(0b00000000),
		byte(0b00000000),
		byte(0b00000000),
		byte(0b11111111), // 1
		// PayloadData ...
	}
	// 255 bytes
	dataFrameBytes = append(dataFrameBytes, make([]byte, 255)...)

	data, err := ParseData(dataFrameBytes)
	fmt.Println(data, err)
	fmt.Println(data.Payload.GetData())
}

func TestParseData2(t *testing.T) {
	dataFrameBytes := []byte{
		byte(0b1_0_0_0_0010), // FIN 1 OpenCode 2
		byte(0b0_0000001),    // MASK 0 PayloadLen 1
		// PayloadData
		byte(0b01111111),
	}
	data, err := ParseData(dataFrameBytes)
	fmt.Println(data, err)
}

func TestParseData3(t *testing.T) {
	dataFrameBytes := []byte{
		byte(0b1_0_0_0_0001), // FIN 1 OpenCode 2
		byte(0b0_1111110),    // MASK 0 PayloadLen 126
		// the next 8 bytes is payload length (bytes)
		byte(0b00000000),
		byte(0b01111111), // payload length is 127 bytes
		// PayloadData...
	}
	dataFrameBytes = append(dataFrameBytes, make([]byte, 127)...)
	data, err := ParseData(dataFrameBytes)
	fmt.Println(data, err)
}

func TestWrapperPayload1(t *testing.T) {
	payloads := make([]byte, 125)
	frameBytes := WrapperPayload(payloads, TextDataType)
	for i, frameByte := range frameBytes {
		if i > 1 {
			continue
		}
		fmt.Printf("NO: %d, Binary: %b \n", i, frameByte)
	}
}

func TestWrapperPayload2(t *testing.T) {
	payloads := make([]byte, 456)
	frameBytes := WrapperPayload(payloads, TextDataType)
	for i, frameByte := range frameBytes {
		if i > 3 {
			continue
		}
		fmt.Printf("NO: %d, Binary: %b \n", i, frameByte)
	}
}

func TestWrapperPayload3(t *testing.T) {
	payloads := make([]byte, (1<<16)-1)
	frameBytes := WrapperPayload(payloads, TextDataType)
	for i, frameByte := range frameBytes {
		if i > 3 {
			continue
		}
		fmt.Printf("NO: %d, Binary: %b \n", i, frameByte)
	}
}

func TestWrapperPayload4(t *testing.T) {
	payloads := make([]byte, 1<<16+(1<<15-1))
	frameBytes := WrapperPayload(payloads, TextDataType)
	for i, frameByte := range frameBytes {
		if i > 9 {
			continue
		}
		fmt.Printf("NO: %d, Binary: %b \n", i, frameByte)
	}
	// fmt.Println(ParseData(frameBytes))
}

func TestWrapperPayloads(t *testing.T) {
	payloads := make([]byte, 1<<16)
	framesBytes := WrapperPayloads(payloads, 5)
	for no, framesByte := range framesBytes {
		for i, b := range framesByte {
			if i > 3 {
				continue
			}
			fmt.Printf("NO: %d/%d Binary %b \n", no, i, b)
		}
		fmt.Println()
		frame, err := ParseData(framesByte)
		fmt.Println(frame, err)
		fmt.Println()
	}

}
