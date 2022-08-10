package socketBuff

import (
	"bytes"
	"errors"
	"math"
	"net"
	"strconv"
)

const (
	Int32Size = 4 // use for kind and size
)

type SocketBuff struct {
	Kind    int32 // if change this, you must change the kind check in the Write()
	Size    int32 // if change this, you must change the size check in the Write()
	Message []byte
}

func Read(conn net.Conn) (*SocketBuff, error) {
	kind, err := readKind(conn)
	if err != nil {
		return nil, err
	}

	size, err := readSize(conn)
	if err != nil {
		return nil, err
	}

	message, err := readMessage(conn, size)
	if err != nil {
			return nil, err
		}
	}
	for i, v := range message {
		println(i, v)
	}

	return &SocketBuff{
		Kind:    int32(kind),
		Size:    int32(len(message)),
		Message: message,
	}, nil
}

func Write(conn net.Conn, kind int, bytes []byte) error {
	if kind > math.MaxInt32 {
		return errors.New("kind overflow")
	}
	kindByte := make([]byte, Int32Size)
	copy(kindByte, strconv.Itoa(kind))

	size := len(bytes)
	if size > math.MaxInt32 {
		return errors.New("size overflow")
	}
	sizeByte := make([]byte, Int32Size)
	copy(sizeByte, strconv.Itoa(size))

	joinByte := append(kindByte, sizeByte...)
	joinByte = append(joinByte, bytes...)

	_, err := conn.Write(joinByte)
	if err != nil {
		return err
	}

	for _, v := range joinByte {
		println(v)
	}
	return nil
}

func readKind(conn net.Conn) (int, error) {
	buf := make([]byte, Int32Size)
	_, err := conn.Read(buf)
	if err != nil {
		return 0, err
	}
	for i, v := range buf {
		println(i, v)
	}

	buf = bytes.Trim(buf, "\x00")
	bufString := string(buf)
	if bufString == "" {
		return 0, nil
	}
	toInt, err := strconv.ParseInt(bufString, 10, 32)
	if err != nil {
		return 0, err
	}

	return int(toInt), nil
}

func readSize(conn net.Conn) (int, error) {
	buf := make([]byte, Int32Size)
	_, err := conn.Read(buf)
	if err != nil {
		return 0, err
	}
	for i, v := range buf {
		println(i, v)
	}

	buf = bytes.Trim(buf, "\x00")
	bufString := string(buf)
	if bufString == "" {
		return 0, nil
	}
	toInt, err := strconv.ParseInt(bufString, 10, 32)
	if err != nil {
		return 0, err
	}

	return int(toInt), nil
}

func readMessage(conn net.Conn, size int) ([]byte, error) {
	buf := make([]byte, size)
	_, err := conn.Read(buf)
	if err != nil {
		return []byte{}, err
	}

	return buf, nil
}
