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

nav {
  section "Core" [Workers, Groups, Scheduler]
  section "Finance" [Payments, Invoices, Orders]
  section "Content" [News, "Company Announcements"]
  section "Admin" [Accounts, "Company Settings"] role(manager)
}

hook on Worker.create {
  emit "worker.created"
}