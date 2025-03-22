package cmd

import (
	"context"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/st3w4r/santa-cruz/manager"
)

func query(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	ds, err := manager.InitDBManger()
	if err != nil {
		return err
	}

	var dbName string
	id, err := strconv.Atoi(args[0])
	if err != nil {
		dbName = args[0]
	}

	db, err := ds.GetDb(int64(id), dbName)
	if err != nil {
		return err
	}

	mdb, err := manager.ConnectToManagedDb(db)
	if err != nil {
		return err
	}
	defer mdb.Db.Close()

	ctx := context.Background()
	rows, err := mdb.Db.QueryContext(ctx, args[1])
	if err != nil {
		return err
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	row := make([][]byte, len(cols))
	rowPtr := make([]any, len(cols))
	for i := range row {
		rowPtr[i] = &row[i]
	}
	data := [][]string{}

	for rows.Next() {
		_ = rows.Scan(rowPtr...)

		rowData := []string{}
		for _, r := range row {
			rowData = append(rowData, string(r))
		}
		data = append(data, rowData)

	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(cols)
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	return nil

}

func queryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query <id/name> '<sql_query>'",
		Aliases: []string{"q"},
		Short:   "Execute a query on a database",
		Args:    cobra.ExactArgs(2),
		Example: "cruz query 1 'SELECT * FROM table'",
		RunE:    query,
	}

	return cmd
}
