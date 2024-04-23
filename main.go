package main

import (
	"fmt"
	"os"

	"github.com/Devansh3712/gg/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gg",
	Short: "gg - git rewrite in go",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.AddCommand(cmd.InitCmd)

	cmd.CatFileCmd.Flags().BoolVarP(&cmd.Type, "type", "t", false, "Show the object type")
	cmd.CatFileCmd.Flags().BoolVarP(&cmd.Size, "size", "s", false, "Show the object size")
	cmd.CatFileCmd.Flags().BoolVarP(&cmd.Pretty, "pretty", "p", false, "Pretty-print the contents of <object> based on its type")
	rootCmd.AddCommand(cmd.CatFileCmd)

	cmd.HashObjectCmd.Flags().BoolVarP(&cmd.Write, "write", "w", false, "Write object into object database")
	rootCmd.AddCommand(cmd.HashObjectCmd)

	cmd.LsTreeCmd.Flags().BoolVar(&cmd.NameOnly, "name-only", false, "List only filenames")
	rootCmd.AddCommand(cmd.LsTreeCmd)

	rootCmd.AddCommand(cmd.WriteTreeCmd)

	cmd.CommitTreeCmd.Flags().StringVarP(&cmd.Parent, "parent", "p", "", "ID of parent commit object")
	cmd.CommitTreeCmd.Flags().StringVarP(&cmd.Message, "message", "m", "", "A paragraph in the commit log message")
	rootCmd.AddCommand(cmd.CommitTreeCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing gg: %s\n", err)
		os.Exit(1)
	}
}
