package main

import (
	"fmt"
	"myapp/cmd"
	"os"

	"github.com/spf13/cobra"
	"github.com/zuiwuchang/gosdk"
)

func main() {
	var root = &cobra.Command{
		Use:   "myapp",
		Short: "myapp example",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(`args:`, gosdk.Args)
			fmt.Println(`dir:`, gosdk.Dir)
			fmt.Println(`yaegi:`, gosdk.Yaegi)
			fmt.Println(`app version:`, gosdk.AppVersion)
			fmt.Println(`app date:`, gosdk.AppDate)
			fmt.Println(`app commit:`, gosdk.AppCommit)
			fmt.Println(`env:`, os.Environ())
			fmt.Println(`os.Args`, os.Args)
		},
	}
	root.AddCommand(cmd.CreateCUI())

	root.SetArgs(os.Args[1:])
	root.Execute()
}
