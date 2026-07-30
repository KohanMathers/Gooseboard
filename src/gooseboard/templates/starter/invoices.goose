entity Invoice {
  field number    string required
  field customer  string required
  field amount    int    required
  field status    string default("draft")
  field issued_at date   default(now())

  permission read  role(manager)
  permission write role(manager)
}

page "Invoices" {
  icon: "file-text"
  route: "/invoices"
  view: table(Invoice) {
    columns: [number, customer, amount, status, issued_at]
    action "Manage" -> page("Manage Invoice", invoice: $row)
    action "Create Invoice" -> form(Invoice)
  }
}

page "Manage Invoice" {
  route: "/invoices/:id"
  hidden: true
  view: detail(Invoice, id: $params.id) {
    tab "Details" -> form(Invoice)
    tab "Payment History" -> plugin("invoice-ledger", invoice: $params.id)
  }
}

hook on Invoice.create {
  emit "invoice.created"
}
