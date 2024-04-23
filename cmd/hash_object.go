package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Write bool

var HashObjectCmd = &cobra.Command{
	Use:   "hash-object",
	Short: "Compute object ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		if Write {
			StoreBlob(path)
		}
		fmt.Fprintf(os.Stdout, "%x\n", HashBlob(path))
	},
}

func HashBlob(path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}
	hash := HashObject("blob", content)
	return hash
}

func StoreBlob(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}

	data := fmt.Sprintf("blob %d\x00%s", len(content), content)
	compressed, err := ZLibCompress([]byte(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error compressing object: %s\n", err)
		os.Exit(1)
	}

	hash := HashBlob(path)
	if err := WriteObject(fmt.Sprintf("%x", hash), compressed); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing compressed object to file: %s\n", err)
	}
}
