package schema

import "testing"

func TestHandBuiltPanel(t *testing.T) {
	p := Panel{
		Title: "Acme Panel",
		Entities: map[string]Entity{
			"Customer": {
				Name: "Customer",
				Fields: []Field{
					{Name: "name", Type: FieldString, Required: true},
					{Name: "credits", Type: FieldInt, Default: 0},
				},
			},
		},
		Pages: map[string]Page{
			"Customers": {
				Title: "Customers",
				Route: "/customers",
				View:  View{Kind: ViewTable, Entity: "Customer", Columns: []string{"name", "credits"}},
			},
		},
	}

	if len(p.Entities) != 1 {
		t.Fatalf("expected 1 entity, got %d", len(p.Entities))
	}

	_ = p
}
