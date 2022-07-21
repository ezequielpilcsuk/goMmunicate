package utils

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"goMunication/member"
	"goMunication/message"
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

func WrapMessage(sender member.Member, data []byte) (message message.Message) {
	message.Sender = sender
	message.Data = data
	message.ID = uuid.New()
	return message
}

func getFinalAddr() {

}
