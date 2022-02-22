package rpc

import (
	"encoding/binary"
	"io"
	"net"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-21

// Transport struct
type Transport struct {
	conn net.Conn
}

// NewTransport creates a transport
func NewTransport(conn net.Conn) *Transport {
	return &Transport{conn}
}

// Send data
func (t *Transport) Send(req Data) error {
	// Encode request into bytes
	b, err := Encode(req)
	if err != nil {
		return err
	}
	buf := make([]byte, 4+len(b))

	// set header field
	binary.BigEndian.PutUint32(buf[:4], uint32(len(b)))
	// set content field
	copy(buf[4:], b)

	_, err = t.conn.Write(buf)
	return err
}

// Receive data
func (t *Transport) Receive() (Data, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(t.conn, header)
	if err != nil {
		return Data{}, err
	}

	// read header field
	dataLen := binary.BigEndian.Uint32(header)
	// read content field
	content := make([]byte, dataLen)

	_, err = io.ReadFull(t.conn, content)
	if err != nil {
		return Data{}, err
	}
	// Decode response from bytes
	rsp, err := Decode(content)
	return rsp, err
}
