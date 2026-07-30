entity Payment {
  field worker   ref(Worker) required
  field amount   int         required
  field method   string      default("bank_transfer")
  field paid_at  date        default(now())

  permission read  role(manager)
  permission write role(manager)
}

page "Payments" {
  icon: "credit-card"
  route: "/payments"
  view: table(Payment) {
    columns: [worker, amount, method, paid_at]
    action "Manage" -> page("Manage Payment", payment: $row)
    action "Add Payment" -> form(Payment)
  }
}

page "Manage Payment" {
  route: "/payments/:id"
  hidden: true
  view: detail(Payment, id: $params.id) {
    tab "Details" -> form(Payment)
  }
}

hook on Payment.create {
  emit "payment.created"
}
