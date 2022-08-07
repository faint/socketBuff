package socketbuf

import (
	"net"
	"strconv"
)

const (
	KindSize = 4
	SizeSize = 4
	SizeMax  = 512
)

type SocketBuff struct {
	Kind  int32
	Size  int32
	Bytes []byte
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

	bytes := make([]byte, size)
	for len(bytes) < size { // 实际长度大于缓存长度时，需多次读取
		splitBytes, err := readBytes(conn, size)
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, splitBytes...)
	}

	return &SocketBuff{
		Kind:  int32(kind),
		Size:  int32(len(bytes)),
		Bytes: bytes,
	}, nil
}

func Write(conn net.Conn, kind int, bytes []byte) error {
	kindByte := make([]byte, KindSize)
	kindByte = []byte(strconv.Itoa(kind))

	size := len(bytes)
	sizeByte := make([]byte, SizeSize)
	sizeByte = []byte(strconv.Itoa(size))

	joinByte := append(kindByte, sizeByte...)
	joinByte = append(joinByte, bytes...)

	_, err := conn.Write(joinByte)
	if err != nil {
		return err
	}

	return nil
}

func readKind(conn net.Conn) (int, error) {
	buf := make([]byte, KindSize)
	_, err := conn.Read(buf)
	if err != nil {
		return 0, err
	}

	toInt, err := strconv.Atoi(string(buf))
	if err != nil {
		return 0, err
	}

	return toInt, nil
}

func readSize(conn net.Conn) (int, error) {
	buf := make([]byte, SizeMax)
	_, err := conn.Read(buf)
	if err != nil {
		return 0, err
	}

	toInt, err := strconv.Atoi(string(buf))
	if err != nil {
		return 0, err
	}

	return toInt, nil
}

func readBytes(conn net.Conn, size int) ([]byte, error) {
	buf := make([]byte, size)
	_, err := conn.Read(buf)
	if err != nil {
		return []byte{}, err
	}

	return buf, nil
}
