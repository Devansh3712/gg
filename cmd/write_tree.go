package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	ModeBlob = "100644"
	ModeDir  = "40000"
)

var WriteTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "Create tree object from the current index",
	Run: func(cmd *cobra.Command, args []string) {
		tree := WriteTree(".")
		fmt.Fprintf(os.Stdout, "%x\n", tree.Hash)
	},
}

func WriteTree(path string) *Tree {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading files and directories: %s\n", err)
		os.Exit(1)
	}

	var tree Tree
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".git") || file.Name() == "." {
			continue
		}
		newPath := filepath.Join(path, file.Name())
		if file.IsDir() {
			dirTree := WriteTree(newPath)
			tree.Entries = append(tree.Entries, TreeObject{
				Mode: []byte(ModeDir),
				Name: []byte(file.Name()),
				Hash: dirTree.Hash,
			})
		} else {
			StoreBlob(newPath)
			tree.Entries = append(tree.Entries, TreeObject{
				Mode: []byte(ModeBlob),
				Name: []byte(file.Name()),
				Hash: []byte(HashBlob(newPath)),
			})
		}
	}
	content := tree.Bytes()
	tree.Hash = HashObject("tree", content)

	data := fmt.Sprintf("tree %d\x00%s", len(content), content)
	compressed, err := ZLibCompress([]byte(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error compressing object: %s\n", err)
		os.Exit(1)
	}
	if err := WriteObject(fmt.Sprintf("%x", tree.Hash), compressed); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing compressed object: %s\n", err)
		os.Exit(1)
	}
	return &tree
}
