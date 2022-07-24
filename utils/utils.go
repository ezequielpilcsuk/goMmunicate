package utils

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	gm "goMunication/src"
	"os"
)

// CheckErr checks for errors and displays and exits if positive
func CheckErr(err error) {
	if err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			fmt.Fprintf(os.Stderr, "Tempo de conexão excedida: %s", err.Error())
		} else {
			fmt.Fprintf(os.Stderr, "Erro fatal: %s", err.Error())
		}
		os.Exit(1)
	}
}

func WrapMessage(data []byte) gm.Message {
	tmpBuf := bytes.NewBuffer(data)
	tmpMsg := new(gm.Message)
	decoder := gob.NewDecoder(tmpBuf)
	CheckErr(decoder.Decode(tmpMsg))

	return *tmpMsg
}

func UnwrapMessage(msg gm.Message) []byte {
	tmpBuf := new(bytes.Buffer)
	encoder := gob.NewEncoder(tmpBuf)
	CheckErr(encoder.Encode(msg))
	return tmpBuf.Bytes()
}
