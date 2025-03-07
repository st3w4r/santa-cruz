package cmd

import (
	"log"

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
		log.Println(db)
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
	log.Println("Add database")

	ds, err := manager.InitDBManger()
	if err != nil {
		return err
	}

	db, err := ds.AddDb("test", "test.db")
	if err != nil {
		return err
	}

	log.Println(db)

	return nil
}

func addDbCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add database",
		RunE:  addDb,
	}
	return cmd
}
