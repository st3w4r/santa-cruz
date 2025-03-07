package main

import (
	"context"
	_ "embed"
	"log"

	_ "modernc.org/sqlite"

	"github.com/st3w4r/santa-cruz/cmd"
	"github.com/st3w4r/santa-cruz/config"
)

// //go:embed schema.sql
// var ddl string

// func run() error {
// 	fmt.Println("Hello Santa Cruz")

// 	ctx := context.Background()

// 	db, err := sql.Open("sqlite", ":memory:")
// 	if err != nil {
// 		return err
// 	}

// 	if _, err := db.ExecContext(ctx, ddl); err != nil {
// 		return err
// 	}

// 	queries := dbsqlc.New(db)

// 	dbs, err := queries.ListDatabases(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	log.Println(dbs)

// 	timeNow := time.Now().UTC().Format(time.RFC3339)

// 	createdDb, err := queries.CreateDatabase(ctx, dbsqlc.CreateDatabaseParams{
// 		Name:      "test",
// 		CreatedAt: timeNow,
// 		UpdatedAt: timeNow,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	log.Println(createdDb)

// 	fetchedDb, err := queries.GetDatabase(ctx, createdDb.ID)
// 	if err != nil {
// 		return err
// 	}
// 	log.Println(fetchedDb)

// 	return nil
// }

func main() {

	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	cmd.NewCLI().ExecuteContext(context.Background())
}
