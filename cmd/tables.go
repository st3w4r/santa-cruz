package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/st3w4r/santa-cruz/manager"
)

func listTables(cmd *cobra.Command, args []string) error {
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

	mdb, err := manager.ConnectToManagedDb(db)
	if err != nil {
		return err
	}

	tables, err := manager.ListTablesDb(mdb)
	if err != nil {
		return err
	}

	data := [][]string{}

	if len(tables) == 0 {
		fmt.Println("No tables found")
		return nil
	}

	for _, table := range tables {
		data = append(data, []string{
			table.Name,
			table.TblName,
			table.Sql,
		})

	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Table Name", "SQL"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	return nil
}

func listTablesCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "tables",
		Short:   "List tables",
		Aliases: []string{"tbl"},
		Args:    cobra.ExactArgs(1),
		RunE:    listTables,
	}
}
