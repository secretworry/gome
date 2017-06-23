package protocol

import (
	"unsafe"
	"io"
	"reflect"
	"sync"
	"encoding/binary"
	"fmt"
	"log"
)

const (
	protocolVersionBits 	= 8 // bit
	protocolVersionOffset 	= 8 // bit
	protocolVersionMask 	= (0x1 << protocolVersionBits - 1) << protocolVersionOffset
	messageTypeBits 	= 8 // bit
	messageTypeOffset 	= 0 // bit
	messageTypeIdBytes 	= (protocolVersionBits + messageTypeBits) / 8 //byte
	messageTypeMask 	= (0x1 << messageTypeBits - 1) << messageTypeOffset
	signature uint16	= 0x0622
	signatureBytes 		= unsafe.Sizeof(signature)
	headerBytes 		= signatureBytes + messageTypeIdBytes
)

type ProtocolVersionType uint8

type MessageTypeType uint8

type MessageTypeIdType uint16

type Message interface {
	io.WriterTo
	io.ReaderFrom
}

type Protocol struct {
	version	ProtocolVersionType
	messageTypes map[MessageTypeType] reflect.Type
	reverseMessageTypes map[reflect.Type] MessageTypeType
}

type ErrMalformatedData struct {
	reason string
}

func (e *ErrMalformatedData) Error() string {
	return e.reason;
}


type ErrIllegalMessage struct {
	reason string
}

func (e *ErrIllegalMessage) Error() string {
	return e.reason;
}

func validateHeader(header []byte) (err error) {
	if inSignature := binary.BigEndian.Uint16(header); inSignature != signature {
		err = &ErrMalformatedData{fmt.Sprintf("Illegal header: illegal signature, expecting %v but got %v", signature, inSignature)}
	}
	return
}

func New(version ProtocolVersionType) *Protocol {
	return &Protocol{
		version: version,
		messageTypes: make(map[MessageTypeType] reflect.Type),
		reverseMessageTypes: make(map[reflect.Type] MessageTypeType),
	}
}

func (p *Protocol) RegisterMessageType(message Message, messageType MessageTypeType) *Protocol {
	if previousType := p.messageTypes[messageType]; previousType != nil {
		log.Panicf("Duplicate declaration of messageType %v", messageType)
	}
	typeOfMessage := reflect.TypeOf(message)
	if previousTypeOfMessage := p.reverseMessageTypes[typeOfMessage]; previousTypeOfMessage != 0 {
		log.Panicf("Conflict messageType %v for message %v", previousTypeOfMessage, typeOfMessage)
	}
	p.messageTypes[messageType] = typeOfMessage
	p.reverseMessageTypes[typeOfMessage] = messageType
	return p
}

var frameBytes = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, headerBytes)
		return &buf;
	},
}

func (p *Protocol) ReadFrom(r io.Reader) (message Message, err error) {
	bufp := frameBytes.Get().(*[]byte)
	defer frameBytes.Put(bufp)
	if err = p.readHeader(r, *bufp); err != nil {
		return
	}
	messageTypeId := MessageTypeIdType(binary.BigEndian.Uint16((*bufp)[signatureBytes:signatureBytes + messageTypeIdBytes]))
	message, err = p.readMessage(r, messageTypeId)
	return
}

func (p *Protocol) WriteTo(w io.Writer, message Message) (n int64, err error) {
	n = 0
	headerLength, err := p.writeHeader(w, message)
	n += headerLength
	if err != nil {
		return;
	}
	messageLength, err := message.WriteTo(w)
	n += messageLength
	if err != nil {
		return;
	}
	return;
}

func (p *Protocol) writeHeader(w io.Writer, message Message) (n int64, err error) {
	typeOfMessage := reflect.TypeOf(message)
	messageType := p.reverseMessageTypes[typeOfMessage]
	if messageType == 0 {
		err = &ErrIllegalMessage{fmt.Sprintf("Illegal message: message of type %v is unsupported by protocol %v", typeOfMessage, p)}
		return
	}
	messageTypeId := (uint16(p.version) << protocolVersionOffset) | (uint16(messageType) << messageTypeOffset)
	bufp := frameBytes.Get().(*[]byte)
	defer frameBytes.Put(bufp)
	binary.BigEndian.PutUint16((*bufp), signature)
	binary.BigEndian.PutUint16((*bufp)[signatureBytes:signatureBytes + messageTypeIdBytes], messageTypeId)
	writeLength, err := w.Write(*bufp)
	n = int64(writeLength)
	return;
}

func (p *Protocol) readHeader(r io.Reader, buf []byte) (err error) {
	_, err = io.ReadFull(r, buf)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			err = &ErrMalformatedData{"Illegal header: insuffcient bytes"}
		} else {
			err = &ErrMalformatedData{fmt.Sprintf("Reading header error: %v", err)}
		}
		return
	}
	err = validateHeader(buf)
	return
}

func (p *Protocol) readMessage(r io.Reader, messageTypeIdType MessageTypeIdType) (message Message, err error) {
	protocolVersion := uint16(messageTypeIdType) & protocolVersionMask >> protocolVersionOffset
	messageType	:= MessageTypeType(uint16(messageTypeIdType) & messageTypeMask >> messageTypeOffset)
	if ProtocolVersionType(protocolVersion) != p.version {
		err = &ErrMalformatedData{fmt.Sprintf("Illegal header: illegal protocol version, expecting %v but got %v", p.version, protocolVersion)}
		return
	}
	mType := p.messageTypes[messageType]
	if mType == nil {
		err = &ErrMalformatedData{fmt.Sprintf("Illegal message type: %v", messageType)}
		return
	}
	mValue := reflect.New(mType.Elem())
	message = mValue.Interface().(Message)
	_, err = message.ReadFrom(r)
	return
}

