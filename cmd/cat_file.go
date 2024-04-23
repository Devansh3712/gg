package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Type   bool
	Size   bool
	Pretty bool
)

var CatFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "Provide contents or details of repository objects",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		objectHash := args[0]
		object := CatFile(objectHash)
		switch {
		case Type:
			fmt.Println(object.Type)
		case Size:
			fmt.Println(object.Size)
		case Pretty:
			fmt.Print(object.Content)
		}
	},
}

type Object struct {
	Type    string
	Size    string
	Content string
}

func CatFile(hash string) Object {
	content, err := ReadObject(hash)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading object: %s\n", err)
		os.Exit(1)
	}
	// Create a reader for zlib using contents of the file read as buffer
	decompressed, err := ZLibDecompress(content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decompressing object: %s\n", err)
		os.Exit(1)
	}
	// Split type, size, file content from the decompressed file
	space := bytes.IndexByte(decompressed, ' ')
	null := bytes.IndexByte(decompressed, '\x00')
	return Object{
		Type:    string(decompressed[:space]),
		Size:    string(decompressed[space+1 : null]),
		Content: string(decompressed[null+1:]),
	}
}
