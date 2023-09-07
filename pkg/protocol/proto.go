package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/panjf2000/gnet/v2"
)

var ErrIncompletePacket = errors.New("incomplete packet")

const (
	magicNumber uint32 = 0x12346641

	sizeMagicNumber = 4
	sizeMsgType     = 1
	sizeBodyLen     = 4
)

type MsgType uint16

const (
	MsgTypeZero  MsgType = 0x00 // 错误
	MsgTypeHeart MsgType = 0x01 // 心跳
	MsgTypeData  MsgType = 0x02 // 数据
	MsgTypeKick  MsgType = 0x03 // kick
)

var magicNumberBytes []byte

func init() {
	magicNumberBytes = make([]byte, sizeMagicNumber)
	binary.BigEndian.PutUint32(magicNumberBytes, magicNumber)
}

// SimpleCodec Protocol format:
// *  |----------Len--------|------------------------------------MetaInfo------------------------------------|
// *  |---------4Byte-------|--------2Byte--------|----4Byte-------|----------------X------------------------|
// *  +------------------------------------------------------------------------------------------------------+
// *  |      Magic          |     MsgType         |   body len    |                Body(seqId)              |
// *  +------------------------------------------------------------------------------------------------------+

// PayloadLen - length of body in byte, 3 bytes big-endian integer.
// Magic - magic number
// MsgType - package type, 2 byte
//   0x01: heartbeat package
//   0x02: data package
//   0x03: disconnect message from server
// body - binary payload.

type SimpleCodec struct{}

func (codec SimpleCodec) Encode(msyType MsgType, buf []byte) ([]byte, error) {
	bodyOffset := sizeMagicNumber + sizeMsgType + sizeBodyLen
	msgLen := bodyOffset + len(buf)

	data := make([]byte, msgLen)
	// magic
	copy(data, magicNumberBytes)

	// msgType
	binary.BigEndian.PutUint16(data[sizeMagicNumber:sizeMagicNumber+sizeMsgType], uint16(msyType))

	// bodylen
	binary.BigEndian.PutUint32(data[sizeMagicNumber+sizeMsgType:bodyOffset], uint32(len(buf)))

	// body
	copy(data[bodyOffset:msgLen], buf)
	return data, nil
}

func (codec *SimpleCodec) Decode(c gnet.Conn) (MsgType, []byte, error) {
	bodyOffset := sizeMagicNumber + sizeMsgType + sizeBodyLen
	buf, _ := c.Peek(bodyOffset)
	if len(buf) < bodyOffset {
		return MsgTypeZero, nil, ErrIncompletePacket
	}

	if !bytes.Equal(magicNumberBytes, buf[:sizeMagicNumber]) {
		return MsgTypeZero, nil, errors.New("invalid magic number")
	}

	// msgType
	msg := binary.BigEndian.Uint16(buf[sizeMagicNumber : sizeMagicNumber+sizeMsgType])
	msgType := MsgType(msg)

	bodyLen := binary.BigEndian.Uint32(buf[sizeMagicNumber+sizeMsgType : bodyOffset])
	msgLen := bodyOffset + int(bodyLen)
	if c.InboundBuffered() < msgLen {
		return MsgTypeZero, nil, ErrIncompletePacket
	}
	buf, _ = c.Peek(msgLen)
	_, _ = c.Discard(msgLen)

	return msgType, buf[bodyOffset:msgLen], nil
}
