entity Group {
  field name        string   required
  field description string
  field lead        ref(Worker)

  permission read  role(coach, manager)
  permission write role(manager)
}

page "Groups" {
  icon: "layout-grid"
  route: "/groups"
  view: table(Group) {
    columns: [name, lead, description]
    action "Manage" -> page("Manage Group", group: $row)
    action "Add Group" -> form(Group)
  }
}

page "Manage Group" {
  route: "/groups/:id"
  hidden: true
  view: detail(Group, id: $params.id) {
    tab "Profile" -> form(Group)
    tab "Members" -> plugin("group-roster", group: $params.id)
  }
}

hook on Group.create {
  emit "group.created"
}
