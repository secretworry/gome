package protocol

import (
	"io"
	"unsafe"
	"testing"
	"bytes"
	"reflect"
)

type TestMessage struct {}

func (t *TestMessage) WriteTo(w io.Writer) (n int64, err error) {
	n = 0
	return
}

func (t *TestMessage) ReadFrom(r io.Reader) (n int64, err error) {
	n = 0
	return
}

type TestInt8Message struct {
	v int8
}

func (t *TestInt8Message) WriteTo(w io.Writer) (n int64, err error) {
	n = int64(unsafe.Sizeof(t.v))
	w.Write([]byte{byte(t.v)})
	return
}

func (t *TestInt8Message) ReadFrom(r io.Reader) (n int64, err error) {
	n = 0
	buf := make([]byte, 1)
	_, err = r.Read(buf)
	if err != nil {
		return;
	}
	t.v = int8(buf[0])
	return
}

var testProtocol *Protocol = New(1).
	RegisterMessageType(&TestMessage{}, 1).
	RegisterMessageType(&TestInt8Message{}, 2)

func NewExample() {
	New(1)
}

func RegisterMessageExample() {
	New(1).RegisterMessageType(&TestMessage{}, 1)
}

func TestProtocol_WriteTo(t *testing.T) {
	message := &TestInt8Message{8}
	expect := []byte{
		0x06,
		0x22,
		0x01,
		0x02,
		0x08,
	}
	buf := new(bytes.Buffer)
	_, err := testProtocol.WriteTo(buf, message)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if (!reflect.DeepEqual(buf.Bytes(), expect)) {
		t.Logf("Illegal data %v", buf.Bytes())
		t.FailNow()
	}
}

func TestProtocol_ReadFrom(t *testing.T) {
	data := []byte{
		0x06,
		0x22,
		0x01,
		0x02,
		0x08,

	}
	buf := bytes.NewBuffer(data)
	message, err := testProtocol.ReadFrom(buf)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	testMessage, ok := message.(*TestInt8Message)
	if !ok {
		t.Logf("Illegal message type %v", reflect.TypeOf(testMessage))
		t.FailNow()
	}

	if testMessage.v != 8 {
		t.Logf("Illegal message value %v", testMessage.v)
		t.FailNow()
	}

}

