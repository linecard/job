package run

import (
	"fmt"
	"strings"

	"github.com/linecard/run/catalog"
	"github.com/linecard/run/internal/docker"
	"github.com/linecard/run/internal/equip"
	"github.com/linecard/run/internal/kube"
	"github.com/linecard/run/internal/parse"
	"github.com/spf13/cobra"
)

func NewScriptCmd(name, desc string, scriptArgs map[int]parse.Arg) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: desc,
		RunE: func(cmd *cobra.Command, args []string) error {
			script, err := catalog.Read(name)
			if err != nil {
				return err
			}

			var flagArgs = []string{}
			for pos := 0; pos < len(scriptArgs); pos++ {
				arg := scriptArgs[pos]
				flagArgs = append(flagArgs, cmd.Flag(arg.Name).Value.String())
			}

			switch {
			case rootCmdFlags.local:
				ctr, err := docker.Run(
					cmd.Context(),
					rootCmdFlags.attach,
					equip.Pack(script, strings.Join(flagArgs, " ")),
					name,
				)
				if err != nil {
					return err
				}

				fmt.Println(ctr)

				return nil
			default:
				pod, err := kube.Run(
					cmd.Context(),
					rootCmdFlags.attach,
					rootCmdFlags.bind,
					rootCmdFlags.secrets,
					equip.Pack(script, strings.Join(flagArgs, " ")),
					name,
				)
				if err != nil {
					return err
				}

				fmt.Println(pod)
			}

			return nil
		},
	}

	for _, arg := range scriptArgs {
		cmd.Flags().String(arg.Name, arg.Value, arg.Desc)

		if arg.Value == "" {
			cmd.MarkFlagRequired(arg.Name)
		}
	}

	cmd.AddCommand(NewInspectCmd(name))

	return cmd
}

func NewInspectCmd(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspect script",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.NoArgs(cmd, args); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			script, err := catalog.Read(name)
			if err != nil {
				return err
			}

			fmt.Println(string(script))

			return nil
		},
	}

	return cmd
}
