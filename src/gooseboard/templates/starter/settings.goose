entity CompanySettings {
  field company_name  string required
  field support_email string
  field timezone      string default("UTC")
  field currency       string default("USD")

  permission read  role(manager)
  permission write role(manager)
}

page "Company Settings" {
  icon: "settings"
  route: "/settings"
  view: form(CompanySettings) {}
}
