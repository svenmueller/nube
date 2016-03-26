package common

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/fatih/color"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/spf13/viper"
)

//MissingArgumentsError is an error returned if there are too few arguments for a command.
type MissingArgumentsError struct {
	Command string
}

type outputErrors struct {
	Errors []outputError `json:"errors"`
}

type outputError struct {
	Detail string `json:"detail"`
}

var (
	_ error = &MissingArgumentsError{}

	colorErr = color.New(color.FgRed).SprintFunc()("Error")

	// defines what should happen if an error occurs
	errorAction = func() {
		os.Exit(1)
	}
)

// NewMissingArgumentsError creates a MissingArgumentsError instance.
func NewMissingArgumentsError(command *cobra.Command) *MissingArgumentsError {
	return &MissingArgumentsError{Command: command.Name()}
}

func (e *MissingArgumentsError) Error() string {
	return fmt.Sprintf("Command '%s' is missing required arguments", e.Command)
}

func HandleError(err error, cmd ...*cobra.Command) {
	if err == nil {
		return
	}

	output := viper.GetString("output")

	switch output {
	default:
		if len(cmd) > 0 {
			cmd[0].Help()
		}
		fmt.Fprintf(color.Output, "\n%s: %v\n", colorErr, err)
	case "json":
		es := outputErrors{
			Errors: []outputError{
				{Detail: err.Error()},
			},
		}

		b, _ := json.Marshal(&es)
		fmt.Println(string(b))
	}

	errorAction()
}
