package cmd

import (
	"fmt"

	"logsnap/collector/factory"

	"github.com/urfave/cli/v2"
)

func supportedProgramsAction(c *cli.Context) error {
	fmt.Println("支持的程序列表：")
	for _, program := range factory.GetSupportedProcessorTypes() {
		fmt.Println(program)
	}
	return nil
}
