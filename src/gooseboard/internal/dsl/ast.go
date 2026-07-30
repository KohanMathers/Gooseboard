package dsl

type Node interface{}

type PanelNode struct {
	ID       string
	Title    string
	Theme    string
	Entities []EntityNode
	Pages    []PageNode
	Nav      []NavSectionNode
	Hooks    []HookNode
}

type FieldNode struct {
	Name     string
	Type     string
	RefTo    string
	Required bool
	Default  any
}

type PermissionNode struct {
	Action string
	Roles  []string
}

type EntityNode struct {
	Name        string
	Fields      []FieldNode
	Permissions []PermissionNode
}

type Arg struct {
	Name  string
	Value any
}

type CallExpr struct {
	Name string
	Args []Arg
}

type SigilExpr struct {
	Path []string
}

type ActionNode struct {
	Label  string
	Target CallExpr
}

type TabNode struct {
	Label  string
	Target CallExpr
}

type ViewNode struct {
	Kind    string
	Entity  string
	Args    []Arg
	Columns []string
	Actions []ActionNode
	Tabs    []TabNode
}

type PageNode struct {
	Title  string
	Route  string
	Icon   string
	Hidden bool
	View   ViewNode
}

type NavSectionNode struct {
	Title string
	Pages []string
	Roles []string
}

type HookNode struct {
	Entity string
	Event  string
	Emits  []string
}
