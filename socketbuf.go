package socketBuff

import (
	"encoding/binary"
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
	println("kind", kind)

	size, err := readSize(conn)
	if err != nil {
		return nil, err
	}
	println("size", size)

	message, err := readMessage(conn, size)
	if err != nil {
		return nil, err
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
		println("kind buf", i, v)
	}

	toInt := binary.BigEndian.Uint32(buf)
	println("kind toInt:", toInt)

	return int(toInt), nil
}

func readSize(conn net.Conn) (int, error) {
	buf := make([]byte, Int32Size)
	_, err := conn.Read(buf)
	if err != nil {
		return 0, err
	}
	for i, v := range buf {
		println("size buf", i, v)
	}

	toInt := binary.BigEndian.Uint32(buf)
	println("size toInt:", toInt)

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
