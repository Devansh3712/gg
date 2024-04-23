package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a git directory",
	Run: func(cmd *cobra.Command, args []string) {
		GitInit()
	},
}

func GitInit() {
	dirs := []string{".git", ".git/objects", ".git/refs"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", dir)
			os.Exit(1)
		}
	}

	headContent := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(".git/HEAD", headContent, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Initialized git directory")
}
