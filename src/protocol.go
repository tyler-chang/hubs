package protocol

import (
	"encoding/binary"
	"errors"
	"io"
	"net"

	"github.com/gansidui/gotcp"
)

// Packet 数据包
type Packet struct {
	buff []byte
}

// Serialize 获取数据包 buff
func (this *Packet) Serialize() []byte {
	return this.buff
}

// GetLength 获取数据包长度
func (this *Packet) GetLength() uint32 {
	return binary.BigEndian.Uint32(this.buff[0:4])
}

// GetBody 获取数据包体
func (this *Packet) GetBody() []byte {
	return this.buff[4:]
}

// NewPacket 创建一个新的数据包
func NewPacket(buff []byte, hasLengthField bool) *Packet {
	p := &Packet{}

	if hasLengthField {
		p.buff = buff
	} else {
		p.buff = make([]byte, 4+len(buff))
		binary.BigEndian.PutUint32(p.buff[0:4], uint32(len(buff)))
		copy(p.buff[4:], buff)
	}

	return p
}

// Protocol 协议
type Protocol struct {
}

// ReadPacket 读取数据
func (this *Protocol) ReadPacket(conn *net.TCPConn) (gotcp.Packet, error) {
	var (
		lengthBytes []byte = make([]byte, 4)
		length      uint32
	)

	// read length
	if _, err := io.ReadFull(conn, lengthBytes); err != nil {
		return nil, err
	}
	if length = binary.BigEndian.Uint32(lengthBytes); length > 1024 {
		return nil, errors.New("the size of packet is larger than the limit")
	}

	buff := make([]byte, 4+length)
	copy(buff[0:4], lengthBytes)

	// read body ( buff = lengthBytes + body )
	if _, err := io.ReadFull(conn, buff[4:]); err != nil {
		return nil, err
	}

	return NewPacket(buff, true), nil
}
