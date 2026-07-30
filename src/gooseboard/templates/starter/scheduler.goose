entity Shift {
  field worker     ref(Worker) required
  field group      ref(Group)
  field date       date        required
  field start_time string      required
  field end_time   string      required

  permission read  role(coach, manager)
  permission write role(coach, manager)
}

page "Scheduler" {
  icon: "calendar"
  route: "/scheduler"
  view: table(Shift) {
    columns: [date, worker, group, start_time, end_time]
    action "Manage" -> page("Manage Shift", shift: $row)
    action "Add Shift" -> form(Shift)
  }
}

page "Manage Shift" {
  route: "/scheduler/:id"
  hidden: true
  view: detail(Shift, id: $params.id) {
    tab "Details" -> form(Shift)
    tab "Conflicts" -> plugin("schedule-conflicts", shift: $params.id)
  }
}

hook on Shift.create {
  emit "shift.created"
}
