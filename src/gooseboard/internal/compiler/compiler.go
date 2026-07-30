package compiler

import (
	"fmt"

	"gooseboard/internal/dsl"
	"gooseboard/internal/schema"
)

func Compile(ast dsl.PanelNode) (schema.Panel, error) {
	panel := schema.Panel{
		ID:       ast.ID,
		Title:    ast.Title,
		Theme:    ast.Theme,
		Entities: map[string]schema.Entity{},
		Pages:    map[string]schema.Page{},
	}

	for _, en := range ast.Entities {
		entity := schema.Entity{Name: en.Name}
		for _, f := range en.Fields {
			entity.Fields = append(entity.Fields, schema.Field{
				Name:     f.Name,
				Type:     schema.FieldType(f.Type),
				RefTo:    f.RefTo,
				Required: f.Required,
				Default:  f.Default,
			})
		}
		for _, perm := range en.Permissions {
			entity.Permissions = append(entity.Permissions, schema.Permission{
				Action: perm.Action,
				Roles:  perm.Roles,
			})
		}
		panel.Entities[en.Name] = entity
	}

	for _, pg := range ast.Pages {
		if pg.View.Entity != "" {
			if _, ok := panel.Entities[pg.View.Entity]; !ok {
				return panel, fmt.Errorf("page %q references unknown entity %q", pg.Title, pg.View.Entity)
			}
		}
		panel.Pages[pg.Title] = schema.Page{
			Title:  pg.Title,
			Route:  pg.Route,
			Icon:   pg.Icon,
			Hidden: pg.Hidden,
			View: schema.View{
				Kind:    schema.ViewKind(pg.View.Kind),
				Entity:  pg.View.Entity,
				Columns: pg.View.Columns,
			},
		}
	}

	for _, section := range ast.Nav {
		resolved := schema.NavSection{Title: section.Title, Roles: section.Roles}
		for _, pageTitle := range section.Pages {
			if _, ok := panel.Pages[pageTitle]; !ok {
				return panel, fmt.Errorf("nav section %q references unknown page %q", section.Title, pageTitle)
			}
			resolved.Pages = append(resolved.Pages, pageTitle)
		}
		panel.Nav = append(panel.Nav, resolved)
	}

	return panel, nil
}
