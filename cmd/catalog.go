// Copyright © 2016 Zellyn Hunter <zellyn@gmail.com>

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zellyn/diskii/lib/disk"
)

var shortnames bool // flag for whether to print short filenames
var debug bool

// catalogCmd represents the cat command, used to catalog a disk or
// directory.
var catalogCmd = &cobra.Command{
	Use:     "catalog",
	Aliases: []string{"cat", "ls"},
	Short:   "print a list of files",
	Long:    `Catalog a disk or subdirectory.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runCat(args); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(catalogCmd)
	catalogCmd.Flags().BoolVarP(&shortnames, "shortnames", "s", false, "whether to print short filenames (only makes a difference on Super-Mon disks)")
	catalogCmd.Flags().BoolVarP(&debug, "debug", "d", false, "pring debug information")
}

// runCat performs the actual catalog logic.
func runCat(args []string) error {
	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("cat expects a disk image filename, and an optional subdirectory")
	}
	op, err := disk.Open(args[0])
	if err != nil {
		return err
	}
	if debug {
		fmt.Printf("Got disk of type %q with underlying sector/block order %q.\n", op.Name(), op.Order())
	}
	subdir := ""
	if len(args) == 2 {
		if !op.HasSubdirs() {
			return fmt.Errorf("Disks of type %q cannot have subdirectories", op.Name())
		}
		subdir = args[1]
	}
	fds, err := op.Catalog(subdir)
	if err != nil {
		return err
	}
	for _, fd := range fds {
		if !shortnames && fd.Fullname != "" {
			fmt.Println(fd.Fullname)
		} else {
			fmt.Println(fd.Name)
		}
	}
	return nil
}
