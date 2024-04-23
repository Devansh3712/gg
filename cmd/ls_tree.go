package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var NameOnly bool

var LsTreeCmd = &cobra.Command{
	Use:   "ls-tree",
	Short: "List the contents of a tree object",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tree := LsTree(args[0])
		switch {
		case NameOnly:
			for _, object := range tree.Entries {
				fmt.Println(string(object.Name))
			}
		}
	},
}

type TreeObject struct {
	Mode []byte
	Name []byte
	Hash []byte
}

type Tree struct {
	Type    string
	Size    string
	Entries []TreeObject
	Hash    []byte
}

func (t *Tree) Bytes() []byte {
	entries := new(bytes.Buffer)
	for _, object := range t.Entries {
		data := fmt.Sprintf("%s %s\x00%s", object.Mode, object.Name, object.Hash)
		entries.Write([]byte(data))
	}
	return entries.Bytes()
}

func LsTree(hash string) Tree {
	content, err := ReadObject(hash)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading object: %s\n", err)
		os.Exit(1)
	}

	decompressed, err := ZLibDecompress(content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decompressing object: %s\n", err)
		os.Exit(1)
	}

	space := bytes.IndexByte(decompressed, ' ')
	null := bytes.IndexByte(decompressed, '\x00')
	tree := Tree{
		Type: string(decompressed[:space]),
		Size: string(decompressed[space+1 : null]),
	}
	// Parse all entries and store in tree object
	// Format of an entry
	// <mode> <size>\0<20_byte_sha>
	var objects []TreeObject
	entries := decompressed[null+1:]
	for len(entries) > 0 {
		space = bytes.IndexByte(entries, ' ')
		null = bytes.IndexByte(entries, '\x00')
		hashStart := null + 1
		hashEnd := hashStart + 20

		object := TreeObject{
			Mode: entries[:space],
			Name: entries[space+1 : null],
			Hash: entries[hashStart:hashEnd],
		}
		objects = append(objects, object)
		entries = entries[hashEnd:]
	}
	tree.Entries = objects
	return tree
}
