entity NewsPost {
  field title        string required
  field body         string required
  field author       ref(Worker)
  field published_at date   default(now())

  permission read  role(coach, manager)
  permission write role(manager)
}

page "News" {
  icon: "newspaper"
  route: "/news"
  view: table(NewsPost) {
    columns: [title, author, published_at]
    action "Manage" -> page("Manage News Post", post: $row)
    action "Add Post" -> form(NewsPost)
  }
}

page "Manage News Post" {
  route: "/news/:id"
  hidden: true
  view: detail(NewsPost, id: $params.id) {
    tab "Content" -> form(NewsPost)
  }
}

hook on NewsPost.create {
  emit "news_post.created"
}
