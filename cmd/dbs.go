package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
		fmt.Println("ID:          ", db.ID)
		fmt.Println("Name:        ", db.Name)
		fmt.Println("Path:        ", db.Path)
		fmt.Println("Description: ", db.Description.String)
		fmt.Println("Created At:  ", db.CreatedAt)
		fmt.Println("Updated At:  ", db.UpdatedAt)
		fmt.Println("----")
	}

	return nil
}

func listDbsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List databases",
		RunE:    listDbs,
	}
	return cmd
}

func addDb(cmd *cobra.Command, args []string) error {
	fmt.Println("Add database")
	ds, err := manager.InitDBManger()
	if err != nil {
		return err
	}

	pathAbs, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}

	if fi, err := os.Stat(pathAbs); os.IsNotExist(err) || fi.IsDir() {
		fmt.Println("File does not exist")
		return err
	}

	if !strings.HasSuffix(pathAbs, ".db") {
		return errors.New("File is not a sqlite db")
	}

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

	db, err := ds.AddDb(name, pathAbs, desc)
	if err != nil {
		return err
	}

	fmt.Println("Database added successfully")
	fmt.Println("ID:          ", db.ID)
	fmt.Println("Name:        ", db.Name)
	fmt.Println("Path:        ", db.Path)
	fmt.Println("Description: ", db.Description.String)
	fmt.Println("Created At:  ", db.CreatedAt)
	fmt.Println("Updated At:  ", db.UpdatedAt)

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

func removeDb(cmd *cobra.Command, args []string) error {
	fmt.Println("Remove database from tracking:")

	ds, err := manager.InitDBManger()
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	db, err := ds.GetDb(int64(id))
	if err != nil {
		return err
	}

	fmt.Println("----")
	fmt.Println("ID:          ", db.ID)
	fmt.Println("Name:        ", db.Name)
	fmt.Println("Path:        ", db.Path)
	fmt.Println("Description: ", db.Description.String)
	fmt.Println("Created At:  ", db.CreatedAt)
	fmt.Println("Updated At:  ", db.UpdatedAt)
	fmt.Println("----")

	fmt.Println("Are you sure you want to stop tracking this database? (y/n)")
	var confirm string
	fmt.Scanln(&confirm)

	if confirm == "y" {
		err := ds.RemoveDb(int64(id))
		if err != nil {
			return err
		}
		fmt.Println("Tracking removed successfully")
	}

	return nil
}

func removeDbCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Short:   "Remove database tracking, usage: cruz remove <id>",
		Args:    cobra.ExactArgs(1),
		RunE:    removeDb,
	}
	return cmd
}

func showDb(cmd *cobra.Command, args []string) error {

	ds, err := manager.InitDBManger()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	db, err := ds.GetDb(int64(id))
	if err != nil {
		return err
	}

	fmt.Println("ID:          ", db.ID)
	fmt.Println("Name:        ", db.Name)
	fmt.Println("Path:        ", db.Path)
	fmt.Println("Description: ", db.Description.String)
	fmt.Println("Created At:  ", db.CreatedAt)
	fmt.Println("Updated At:  ", db.UpdatedAt)

	return nil
}

func showDbCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show database, usage: cruz show <id>",
		Args:  cobra.ExactArgs(1),
		RunE:  showDb,
	}
	return cmd
}
