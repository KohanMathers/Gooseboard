entity Announcement {
  field title     string required
  field body      string required
  field pinned    bool   default(false)
  field posted_at date   default(now())

  permission read  role(coach, manager)
  permission write role(manager)
}

page "Company Announcements" {
  icon: "megaphone"
  route: "/announcements"
  view: table(Announcement) {
    columns: [title, pinned, posted_at]
    action "Manage" -> page("Manage Announcement", announcement: $row)
    action "Add Announcement" -> form(Announcement)
  }
}

page "Manage Announcement" {
  route: "/announcements/:id"
  hidden: true
  view: detail(Announcement, id: $params.id) {
    tab "Content" -> form(Announcement)
  }
}

hook on Announcement.create {
  emit "announcement.created"
}
