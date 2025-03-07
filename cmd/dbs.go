package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/st3w4r/santa-cruz/manager"
)

func listDbs(cmd *cobra.Command, args []string) error {

	ds, err := manager.InitDBManger()
	if err != nil {
		return err
	}

	dbs, err := ds.ListDbs()
	if err != nil {
		return err
	}

	for _, db := range dbs {
		fmt.Println(db)
	}

	return nil
}

func listDbsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List databases",
		RunE:  listDbs,
	}
	return cmd
}

func addDb(cmd *cobra.Command, args []string) error {
	fmt.Println("Add database")
	ds, err := manager.InitDBManger()
	if err != nil {
		return err
	}
	// get absolute path
	pathAbs, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}

	// check if file exists
	if fi, err := os.Stat(pathAbs); os.IsNotExist(err) || fi.IsDir() {
		fmt.Println("File does not exist")
		return err
	}

	// check if file is a sqlite db
	if !strings.HasSuffix(pathAbs, ".db") {
		return errors.New("File is not a sqlite db")
	}

	// get the name
	var name string
	name = cmd.Flag("name").Value.String()
	if name == "" {
		_, name = filepath.Split(pathAbs)
		if name == "" {
			return errors.New("Invalid path")
		}
		name = strings.TrimSuffix(name, ".db")
	}

	desc := cmd.Flag("description").Value.String()

	// remove the extension
	fmt.Println("Path: "+pathAbs, "Name: "+name)

	db, err := ds.AddDb(name, pathAbs, desc)
	if err != nil {
		return err
	}

	fmt.Println(db)

	return nil
}

func addDbCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add sqlite database, usage: cruz add <path>",
		Args:  cobra.ExactArgs(1),
		RunE:  addDb,
	}
	cmd.Flags().StringP("name", "n", "", "Name of the database")
	cmd.Flags().StringP("description", "d", "", "Description of the database")
	return cmd
}
