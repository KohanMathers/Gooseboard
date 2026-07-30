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