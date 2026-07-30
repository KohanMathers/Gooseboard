package main

import (
	"fmt"
	"net/http"
	"os"

	"gooseboard/internal/compiler"
	"gooseboard/internal/dsl"
	"gooseboard/internal/engine"
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
	ps := `
panel "example-business" {
  title: "Very Cool Business"
  theme: import("./theme.goose")
}

entity Worker {
  field name        string   required
  field dob         date
  field group       ref(Group)
  field hours_worked int      default(0)
  field hire_date   date     default(now())

  permission read  role(coach, manager)
  permission write role(manager)
}

page "Workers" {
  icon: "users"
  route: "/workers"
  view: table(Worker) {
    columns: [name, group, hours_worked, hire_date]
    action "Manage" -> page("Manage Worker", worker: $row)
    action "Add Worker" -> form(Worker)
  }
}

page "Manage Worker" {
  route: "/workers/:id"
  hidden: true
  view: detail(Worker, id: $params.id) {
    tab "Profile" -> form(Worker)
    tab "Meetings" -> plugin("meeting-history", worker: $params.id)
    tab "Payment" -> plugin("coin-ledger", worker: $params.id)
  }
}

nav {
  section "Core" [Workers, Groups, Scheduler]
  section "Finance" [Payments, Invoices, Orders]
  section "Content" [News, "Company Announcements"]
  section "Admin" [Accounts, "Company Settings"] role(manager)
}

hook on Worker.create {
  emit "worker.created"
}`

	p := dsl.NewParser(ps)
	ast := p.ParsePanel()
	c, err := compiler.Compile(ast)
	if err != nil {
		fmt.Printf("error compiling DSL: %v\n", err)
		os.Exit(1)
	}
	s, err := store.NewSQLiteStore("data.db")
	if err != nil {
		fmt.Printf("error creating store: %v\n", err)
		os.Exit(1)
	}
	if err := s.Migrate(c.Entities); err != nil {
		fmt.Printf("error migrating store: %v\n", err)
		os.Exit(1)
	}
	e := engine.New(c, s)
	http.ListenAndServe(":8080", e.Router())
}
