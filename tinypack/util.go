package tinypack

import (
	"fmt"
	"io"
)

func ensureWrite(w Writer, data []byte) error {
	nWritten, err := w.Write(data)
	if err == nil && nWritten != len(data) {
		return fmt.Errorf("nWritten (%v) != len(data) (%v)", nWritten, len(data))
	}
	return err
}

func ensureRead(r Reader, data []byte) error {
	nRead, err := io.ReadFull(r, data)
	if err == nil && nRead != len(data) {
		return fmt.Errorf("nRead (%v) != len(data) (%v)", nRead, len(data))
	}
	return err
}
