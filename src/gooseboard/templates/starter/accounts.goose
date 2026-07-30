entity Account {
  field name   string required
  field email  string required
  field role   string required
  field active bool   default(true)

  permission read  role(manager)
  permission write role(manager)
}

page "Accounts" {
  icon: "shield"
  route: "/accounts"
  view: table(Account) {
    columns: [name, email, role, active]
    action "Manage" -> page("Manage Account", account: $row)
    action "Add Account" -> form(Account)
  }
}

page "Manage Account" {
  route: "/accounts/:id"
  hidden: true
  view: detail(Account, id: $params.id) {
    tab "Profile" -> form(Account)
    tab "Permissions" -> plugin("account-roles", account: $params.id)
  }
}

hook on Account.create {
  emit "account.created"
}
