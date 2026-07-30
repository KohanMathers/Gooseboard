package main

import (
	"fmt"
	"net/http"
	"os"

	"gooseboard/internal/engine"
	"gooseboard/internal/schema"
	"gooseboard/internal/store"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: gooseboard <run|init|version> [args]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "version":
		fmt.Println("gooseboard v0.0.1")
	case "run":
		temporaryRun()
	case "init":
		fmt.Println("not implemented yet")
	default:
		fmt.Printf("unknown command %q\n", os.Args[1])
		os.Exit(1)
	}
}

func temporaryRun() {
	p := schema.Panel{
		Title: "Acme Panel",
		Entities: map[string]schema.Entity{
			"Customer": {
				Name: "Customer",
				Fields: []schema.Field{
					{Name: "name", Type: schema.FieldString, Required: true},
					{Name: "credits", Type: schema.FieldInt, Default: 0},
				},
			},
		},
		Pages: map[string]schema.Page{
			"Customers": {
				Title: "Customers",
				Route: "/customers",
				View:  schema.View{Kind: schema.ViewTable, Entity: "Customer", Columns: []string{"name", "credits"}},
			},
		},
	}

	s, err := store.NewSQLiteStore("data.db")
	if err != nil {
		panic(err)
	}
	if err := s.Migrate(p.Entities); err != nil {
		panic(err)
	}

	e := engine.New(p, s)
	http.ListenAndServe(":8080", e.Router())
}
