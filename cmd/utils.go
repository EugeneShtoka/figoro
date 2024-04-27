package cmd

import "fmt"

func showError(msg string, err error) {
	logger.Error().Err(err).Msg(msg)
	fmt.Printf("%s. %v\n", msg, err)
}