package dsl

import (
	"reflect"
	"testing"
)

func assertPanics(t *testing.T, want string, fn func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic %q, got no panic", want)
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("expected panic value of type string, got %T: %v", r, r)
		}
		if msg != want {
			t.Fatalf("expected panic %q, got %q", want, msg)
		}
	}()
	fn()
}

func TestParsePanelHeaderOnly(t *testing.T) {
	panel := NewParser(`panel "my-panel" { title: "My Panel" }`).ParsePanel()
	if panel.ID != "my-panel" {
		t.Fatalf("expected id %q, got %q", "my-panel", panel.ID)
	}
	if panel.Title != "My Panel" {
		t.Fatalf("expected title %q, got %q", "My Panel", panel.Title)
	}
	if len(panel.Entities) != 0 || len(panel.Pages) != 0 || panel.Nav != nil {
		t.Fatalf("expected empty panel body, got %+v", panel)
	}
}

func TestParsePanelHeaderWithTheme(t *testing.T) {
	panel := NewParser(`panel "my-panel" { theme: import("./theme.goose") }`).ParsePanel()
	if panel.Theme != "./theme.goose" {
		t.Fatalf("expected theme %q, got %q", "./theme.goose", panel.Theme)
	}
}

func TestParsePanelWithoutWrapper(t *testing.T) {
	src := `entity User { field name string }`
	panel := NewParser(src).ParsePanel()
	if len(panel.Entities) != 1 {
		t.Fatalf("expected 1 entity, got %d", len(panel.Entities))
	}
	if panel.Entities[0].Name != "User" {
		t.Fatalf("expected entity name User, got %q", panel.Entities[0].Name)
	}
}

func TestParseEntityWithFields(t *testing.T) {
	src := `entity User {
		field Name string required
		field Age number
	}`
	panel := NewParser(src).ParsePanel()
	entity := panel.Entities[0]
	if entity.Name != "User" {
		t.Fatalf("expected entity name User, got %q", entity.Name)
	}
	if len(entity.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(entity.Fields))
	}
	want := FieldNode{Name: "Name", Type: "string", Required: true}
	if !reflect.DeepEqual(entity.Fields[0], want) {
		t.Fatalf("field 0 mismatch: want %+v, got %+v", want, entity.Fields[0])
	}
	want = FieldNode{Name: "Age", Type: "number"}
	if !reflect.DeepEqual(entity.Fields[1], want) {
		t.Fatalf("field 1 mismatch: want %+v, got %+v", want, entity.Fields[1])
	}
}

func TestParseFieldDefaultString(t *testing.T) {
	src := `entity User { field name string default("Anon") }`
	field := NewParser(src).ParsePanel().Entities[0].Fields[0]
	if field.Default != "Anon" {
		t.Fatalf("expected default %q, got %v", "Anon", field.Default)
	}
}

func TestParseFieldDefaultNumber(t *testing.T) {
	src := `entity User { field age number default(42) }`
	field := NewParser(src).ParsePanel().Entities[0].Fields[0]
	if field.Default != "42" {
		t.Fatalf("expected default %q, got %v", "42", field.Default)
	}
}

func TestParseFieldDefaultCall(t *testing.T) {
	src := `entity User { field hireDate date default(now()) }`
	field := NewParser(src).ParsePanel().Entities[0].Fields[0]
	want := CallExpr{Name: "now"}
	if !reflect.DeepEqual(field.Default, want) {
		t.Fatalf("expected default %+v, got %+v", want, field.Default)
	}
}

func TestParseFieldDefaultBoolAndNull(t *testing.T) {
	cases := []struct {
		src  string
		want any
	}{
		{`entity User { field active bool default(true) }`, true},
		{`entity User { field active bool default(false) }`, false},
		{`entity User { field owner string default(null) }`, nil},
	}
	for _, c := range cases {
		field := NewParser(c.src).ParsePanel().Entities[0].Fields[0]
		if field.Default != c.want {
			t.Fatalf("for %q: expected default %v, got %v", c.src, c.want, field.Default)
		}
	}
}

func TestParseFieldRefType(t *testing.T) {
	src := `entity Post { field author ref(User) }`
	field := NewParser(src).ParsePanel().Entities[0].Fields[0]
	if field.Type != "ref" {
		t.Fatalf("expected type %q, got %q", "ref", field.Type)
	}
	if field.RefTo != "User" {
		t.Fatalf("expected RefTo %q, got %q", "User", field.RefTo)
	}
}

func TestParseFieldRequiredAndDefaultTogether(t *testing.T) {
	src := `entity User { field name string required default("x") }`
	field := NewParser(src).ParsePanel().Entities[0].Fields[0]
	want := FieldNode{Name: "name", Type: "string", Required: true, Default: "x"}
	if !reflect.DeepEqual(field, want) {
		t.Fatalf("expected %+v, got %+v", want, field)
	}
}

func TestParseMultipleEntities(t *testing.T) {
	src := `entity User { field name string }
	entity Post { field title string }`
	panel := NewParser(src).ParsePanel()
	if len(panel.Entities) != 2 {
		t.Fatalf("expected 2 entities, got %d", len(panel.Entities))
	}
	if panel.Entities[0].Name != "User" || panel.Entities[1].Name != "Post" {
		t.Fatalf("unexpected entity names: %+v", panel.Entities)
	}
}

func TestParseEntityWithPermissions(t *testing.T) {
	src := `entity Worker {
		field name string
		permission read role(coach, manager)
		permission write role(manager)
	}`
	entity := NewParser(src).ParsePanel().Entities[0]
	if len(entity.Permissions) != 2 {
		t.Fatalf("expected 2 permissions, got %d", len(entity.Permissions))
	}
	want0 := PermissionNode{Action: "read", Roles: []string{"coach", "manager"}}
	if !reflect.DeepEqual(entity.Permissions[0], want0) {
		t.Fatalf("permission 0 mismatch: want %+v, got %+v", want0, entity.Permissions[0])
	}
	want1 := PermissionNode{Action: "write", Roles: []string{"manager"}}
	if !reflect.DeepEqual(entity.Permissions[1], want1) {
		t.Fatalf("permission 1 mismatch: want %+v, got %+v", want1, entity.Permissions[1])
	}
}

func TestParsePageWithRouteAndIcon(t *testing.T) {
	src := `page "Home" {
		route: "/home"
		icon: "home-icon"
		view: table(User) { columns: [name, age] }
	}`
	page := NewParser(src).ParsePanel().Pages[0]
	if page.Title != "Home" {
		t.Fatalf("expected title Home, got %q", page.Title)
	}
	if page.Route != "/home" {
		t.Fatalf("expected route /home, got %q", page.Route)
	}
	if page.Icon != "home-icon" {
		t.Fatalf("expected icon home-icon, got %q", page.Icon)
	}
}

func TestParsePageHidden(t *testing.T) {
	src := `page "Manage" { hidden: true view: list { columns: [] } }`
	page := NewParser(src).ParsePanel().Pages[0]
	if !page.Hidden {
		t.Fatalf("expected page to be hidden")
	}
}

func TestParseViewWithEntityAndColumns(t *testing.T) {
	src := `page "Home" { view: table(User) { columns: [name, age, email] } }`
	view := NewParser(src).ParsePanel().Pages[0].View
	if view.Kind != "table" {
		t.Fatalf("expected kind table, got %q", view.Kind)
	}
	if view.Entity != "User" {
		t.Fatalf("expected entity User, got %q", view.Entity)
	}
	wantCols := []string{"name", "age", "email"}
	if !reflect.DeepEqual(view.Columns, wantCols) {
		t.Fatalf("expected columns %v, got %v", wantCols, view.Columns)
	}
}

func TestParseViewWithoutEntity(t *testing.T) {
	src := `page "Home" { view: list { columns: [name] } }`
	view := NewParser(src).ParsePanel().Pages[0].View
	if view.Kind != "list" {
		t.Fatalf("expected kind list, got %q", view.Kind)
	}
	if view.Entity != "" {
		t.Fatalf("expected no entity, got %q", view.Entity)
	}
}

func TestParseViewEmptyColumns(t *testing.T) {
	src := `page "Home" { view: list { columns: [] } }`
	view := NewParser(src).ParsePanel().Pages[0].View
	if len(view.Columns) != 0 {
		t.Fatalf("expected 0 columns, got %v", view.Columns)
	}
}

func TestParseViewWithActionsAndTargetKeywordName(t *testing.T) {
	src := `page "Workers" {
		view: table(Worker) {
			columns: [name]
			action "Manage" -> page("Manage Worker", worker: $row)
			action "Add Worker" -> form(Worker)
		}
	}`
	view := NewParser(src).ParsePanel().Pages[0].View
	if len(view.Actions) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(view.Actions))
	}
	if view.Actions[0].Label != "Manage" || view.Actions[0].Target.Name != "page" {
		t.Fatalf("unexpected action 0: %+v", view.Actions[0])
	}
	wantArgs := []Arg{{Value: "Manage Worker"}, {Name: "worker", Value: SigilExpr{Path: []string{"row"}}}}
	if !reflect.DeepEqual(view.Actions[0].Target.Args, wantArgs) {
		t.Fatalf("unexpected action 0 args: %+v", view.Actions[0].Target.Args)
	}
	if view.Actions[1].Label != "Add Worker" || view.Actions[1].Target.Name != "form" {
		t.Fatalf("unexpected action 1: %+v", view.Actions[1])
	}
	wantArgs1 := []Arg{{Value: "Worker"}}
	if !reflect.DeepEqual(view.Actions[1].Target.Args, wantArgs1) {
		t.Fatalf("unexpected action 1 args: %+v", view.Actions[1].Target.Args)
	}
}

func TestParseDetailViewWithTabsAndSigilArgs(t *testing.T) {
	src := `page "Manage Worker" {
		view: detail(Worker, id: $params.id) {
			tab "Profile" -> form(Worker)
			tab "Meetings" -> plugin("meeting-history", worker: $params.id)
		}
	}`
	view := NewParser(src).ParsePanel().Pages[0].View
	if view.Kind != "detail" || view.Entity != "Worker" {
		t.Fatalf("unexpected view: %+v", view)
	}
	wantArgs := []Arg{{Name: "id", Value: SigilExpr{Path: []string{"params", "id"}}}}
	if !reflect.DeepEqual(view.Args, wantArgs) {
		t.Fatalf("unexpected view args: %+v", view.Args)
	}
	if len(view.Tabs) != 2 {
		t.Fatalf("expected 2 tabs, got %d", len(view.Tabs))
	}
	if view.Tabs[0].Label != "Profile" || view.Tabs[0].Target.Name != "form" {
		t.Fatalf("unexpected tab 0: %+v", view.Tabs[0])
	}
	if view.Tabs[1].Label != "Meetings" || view.Tabs[1].Target.Name != "plugin" {
		t.Fatalf("unexpected tab 1: %+v", view.Tabs[1])
	}
	wantTab1Args := []Arg{{Value: "meeting-history"}, {Name: "worker", Value: SigilExpr{Path: []string{"params", "id"}}}}
	if !reflect.DeepEqual(view.Tabs[1].Target.Args, wantTab1Args) {
		t.Fatalf("unexpected tab 1 args: %+v", view.Tabs[1].Target.Args)
	}
}

func TestParsePageMissingViewPanics(t *testing.T) {
	assertPanics(t, "page Home has no view", func() {
		NewParser(`page "Home" { route: "/home" }`).ParsePanel()
	})
}

func TestParseMultiplePages(t *testing.T) {
	src := `page "Home" { view: list { columns: [a] } }
	page "About" { view: list { columns: [b] } }`
	panel := NewParser(src).ParsePanel()
	if len(panel.Pages) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(panel.Pages))
	}
	if panel.Pages[0].Title != "Home" || panel.Pages[1].Title != "About" {
		t.Fatalf("unexpected page titles: %+v", panel.Pages)
	}
}

func TestParseNavSections(t *testing.T) {
	src := `nav {
		section "Main" [Home, "Settings"]
		section "Other" ["About"]
	}`
	nav := NewParser(src).ParsePanel().Nav
	if len(nav) != 2 {
		t.Fatalf("expected 2 nav sections, got %d", len(nav))
	}
	if nav[0].Title != "Main" || !reflect.DeepEqual(nav[0].Pages, []string{"Home", "Settings"}) {
		t.Fatalf("unexpected nav section 0: %+v", nav[0])
	}
	if nav[1].Title != "Other" || !reflect.DeepEqual(nav[1].Pages, []string{"About"}) {
		t.Fatalf("unexpected nav section 1: %+v", nav[1])
	}
}

func TestParseNavEmptySection(t *testing.T) {
	src := `nav { section "Main" [] }`
	nav := NewParser(src).ParsePanel().Nav
	if len(nav) != 1 || nav[0].Title != "Main" || len(nav[0].Pages) != 0 {
		t.Fatalf("unexpected nav: %+v", nav)
	}
}

func TestParseNavSectionWithRole(t *testing.T) {
	src := `nav { section "Admin" [Accounts, "Company Settings"] role(manager) }`
	nav := NewParser(src).ParsePanel().Nav
	if !reflect.DeepEqual(nav[0].Roles, []string{"manager"}) {
		t.Fatalf("expected roles [manager], got %v", nav[0].Roles)
	}
}

func TestParseTopLevelHook(t *testing.T) {
	src := `hook on Worker.create {
		emit "worker.created"
	}`
	hooks := NewParser(src).ParsePanel().Hooks
	if len(hooks) != 1 {
		t.Fatalf("expected 1 hook, got %d", len(hooks))
	}
	want := HookNode{Entity: "Worker", Event: "create", Emits: []string{"worker.created"}}
	if !reflect.DeepEqual(hooks[0], want) {
		t.Fatalf("unexpected hook: %+v", hooks[0])
	}
}

func TestParseFullPanel(t *testing.T) {
	src := `panel "dashboard" {
		title: "Dashboard"
		entity User {
			field name string required
			field age number default(0)
		}
		page "Home" {
			route: "/home"
			icon: "home"
			view: table(User) { columns: [name, age] }
		}
		nav {
			section "Main" [Home]
		}
	}`
	panel := NewParser(src).ParsePanel()
	if panel.Title != "Dashboard" {
		t.Fatalf("expected title Dashboard, got %q", panel.Title)
	}
	if len(panel.Entities) != 1 || len(panel.Pages) != 1 || len(panel.Nav) != 1 {
		t.Fatalf("unexpected panel shape: %+v", panel)
	}
}

func TestParseUnexpectedTopLevelTokenPanics(t *testing.T) {
	assertPanics(t, "unexpected top-level token: field", func() {
		NewParser(`field`).ParsePanel()
	})
}

func TestParseUnexpectedTokenInPanelHeaderPanics(t *testing.T) {
	assertPanics(t, "unexpected token in panel header: field", func() {
		NewParser(`panel "P" { field }`).ParsePanel()
	})
}

func TestParseEntityMissingFieldKeywordPanics(t *testing.T) {
	assertPanics(t, "expected 'field' or 'permission' keyword in entity, got: page", func() {
		NewParser(`entity User { page }`).ParsePanel()
	})
}

func TestParseUnexpectedTokenInPagePanics(t *testing.T) {
	assertPanics(t, "unexpected token in page: entity", func() {
		NewParser(`page "Home" { entity }`).ParsePanel()
	})
}

func TestParseViewMissingColumnsPanics(t *testing.T) {
	assertPanics(t, "unexpected token in view: route", func() {
		NewParser(`page "Home" { view: list { route } }`).ParsePanel()
	})
}

func TestParseNavMissingSectionKeywordPanics(t *testing.T) {
	assertPanics(t, "expected 'section' keyword in nav, got: page", func() {
		NewParser(`nav { page }`).ParsePanel()
	})
}

func TestParseNavSectionExpectsPageNamePanics(t *testing.T) {
	assertPanics(t, "expected page name in nav section, got: 42", func() {
		NewParser(`nav { section "Main" [42] }`).ParsePanel()
	})
}

func TestParseExpectWrongTokenKindPanics(t *testing.T) {
	assertPanics(t, "unexpected token", func() {
		NewParser(`panel 42`).ParsePanel()
	})
}

func TestParseExpectKeywordWrongKeywordPanics(t *testing.T) {
	assertPanics(t, "expected keyword panel, got: entity", func() {
		NewParser(`entity`).parsePanelHeader(&PanelNode{})
	})
}

func TestParseLiteralValueUnexpectedTokenPanics(t *testing.T) {
	assertPanics(t, "expected literal value, got: {", func() {
		NewParser(`entity User { field name string default({}) }`).ParsePanel()
	})
}

func TestParseEntityMissingClosingBracePanics(t *testing.T) {
	assertPanics(t, "expected 'field' or 'permission' keyword in entity, got: ", func() {
		NewParser(`entity User { field name string`).ParsePanel()
	})
}
