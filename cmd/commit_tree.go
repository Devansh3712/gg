package cmd

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	Parent  string
	Message string
)

var CommitTreeCmd = &cobra.Command{
	Use:   "commit-tree",
	Short: "Create a new commit object",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hash := CommitTree(args[0], Parent, Message)
		fmt.Println(hash)
	},
}

func CommitTree(tree, parent, message string) string {
	path := fmt.Sprintf(".git/objects/%s/%s", tree[:2], tree[2:])
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error reading tree object: %s\n", err)
		os.Exit(1)
	}

	timestamp := time.Now()
	offset := timestamp.Format("+0000")
	// Todo: Add author name and email from a config file
	author := fmt.Sprintf(
		"Devansh Singh <devanshamity@gmail.com> %d %s",
		timestamp.Unix(), offset,
	)
	content := new(bytes.Buffer)
	content.Write([]byte("tree " + tree + "\n"))
	// Todo: Add multiple parents ([]string)
	if parent != "" {
		content.Write([]byte("parent " + parent + "\n"))
	}
	content.Write([]byte("author " + author + "\n"))
	content.Write([]byte("committer " + author + "\n"))
	content.Write([]byte("\n" + message + "\n"))

	data := fmt.Sprintf("commit %d\x00%s", content.Len(), content.Bytes())
	compressed, err := ZLibCompress([]byte(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error compressing object: %s\n", err)
		os.Exit(1)
	}

	hash := HashObject("commit", content.Bytes())
	if err := WriteObject(fmt.Sprintf("%x", hash), compressed); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing compressed object: %s\n", err)
		os.Exit(1)
	}
	return fmt.Sprintf("%x", hash)
}
