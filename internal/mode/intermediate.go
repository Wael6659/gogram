// Copyright (c) 2024 RoseLoverX

package mode

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/amarnathcjd/gogram/internal/encoding/tl"
)

type intermediate struct {
	conn io.ReadWriter
}

var _ Mode = (*intermediate)(nil)

var transportModeIntermediate = [...]byte{0xee, 0xee, 0xee, 0xee} // meta:immutable

func (*intermediate) getModeAnnouncement() []byte {
	return transportModeIntermediate[:]
}

func (m *intermediate) WriteMsg(msg []byte) error {
	size := make([]byte, tl.WordLen)
	binary.LittleEndian.PutUint32(size, uint32(len(msg)))
	if _, err := m.conn.Write(size); err != nil {
		return err
	}
	if _, err := m.conn.Write(msg); err != nil {
		return err
	}

	return nil
}

func (m *intermediate) ReadMsg() ([]byte, error) {
	sizeBuf := make([]byte, tl.WordLen)
	n, err := io.ReadFull(m.conn, sizeBuf)
	if err != nil {
		return nil, err
	}
	if n != tl.WordLen {
		return nil, fmt.Errorf("size is not length of int32, expected 4 bytes, got %d", n)
	}

	size := binary.LittleEndian.Uint32(sizeBuf)

	// incase the data byte is too large, its error in first 4 bytes
	// we read the size of the message, without checking the size

	if size > 1<<30 { // 1GB, arbitrary limit
		var msg = make([]byte, 1<<18) // 256kb
		n, err = io.ReadFull(m.conn, msg)
		if err != nil {
			return nil, err
		}

		var msgWithSize = make([]byte, n)
		copy(msgWithSize, msg[:n])

		return msgWithSize, nil
	}

	msg := make([]byte, int(size))
	n, err = io.ReadFull(m.conn, msg)
	if err != nil {
		return nil, err
	}
	if n != int(size) {
		return nil, fmt.Errorf("expected to read %d bytes, got %d", size, n)
	}

	return msg, nil
}
