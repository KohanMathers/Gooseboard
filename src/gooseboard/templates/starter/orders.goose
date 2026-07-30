entity Order {
  field customer    string required
  field items       string required
  field total       int    required
  field status       string default("pending")
  field ordered_at  date   default(now())

  permission read  role(coach, manager)
  permission write role(manager)
}

page "Orders" {
  icon: "shopping-cart"
  route: "/orders"
  view: table(Order) {
    columns: [customer, total, status, ordered_at]
    action "Manage" -> page("Manage Order", order: $row)
    action "Add Order" -> form(Order)
  }
}

page "Manage Order" {
  route: "/orders/:id"
  hidden: true
  view: detail(Order, id: $params.id) {
    tab "Details" -> form(Order)
    tab "Fulfillment" -> plugin("order-fulfillment", order: $params.id)
  }
}

hook on Order.create {
  emit "order.created"
}
