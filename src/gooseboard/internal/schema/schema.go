package schema

type FieldType string

const (
	FieldString FieldType = "string"
	FieldInt    FieldType = "int"
	FieldBool   FieldType = "bool"
	FieldDate   FieldType = "date"
	FieldRef    FieldType = "ref"
)

type Field struct {
	Name     string
	Type     FieldType
	RefTo    string
	Required bool
	Default  any
}

type Permission struct {
	Action string
	Roles  []string
}

type Entity struct {
	Name        string
	Fields      []Field
	Permissions []Permission
}

type ViewKind string

const (
	ViewTable  ViewKind = "table"
	ViewForm   ViewKind = "form"
	ViewDetail ViewKind = "detail"
	ViewPlugin ViewKind = "plugin"
)

type View struct {
	Kind       ViewKind
	Entity     string
	Columns    []string
	PluginName string
	Children   []View
}

type Page struct {
	Title  string
	Route  string
	Icon   string
	Hidden bool
	Roles  []string
	View   View
}

type NavSection struct {
	Title string
	Pages []string
	Roles []string
}

type Panel struct {
	ID       string
	Title    string
	Theme    string
	Entities map[string]Entity
	Pages    map[string]Page
	Nav      []NavSection
}
