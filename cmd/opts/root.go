package opts

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "connect-emp",
	Short:   "Connect Emp",
	Long:    "Connect Emp - An employee service for Connect",
	Version: "0.1.0",
}
