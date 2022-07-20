package utils

import (
	"errors"
	"fmt"
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

func ParseAddr() {

}
