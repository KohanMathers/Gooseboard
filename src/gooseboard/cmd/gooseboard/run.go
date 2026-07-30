package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gooseboard/internal/compiler"
	"gooseboard/internal/dsl"
	"gooseboard/internal/engine"
	"gooseboard/internal/store"
)

func runCommand(templateDir string, port int) error {
	var ast dsl.PanelNode

	err := filepath.WalkDir(templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".goose") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		mergeASTs(&ast, dsl.NewParser(string(src)).ParsePanel())
		return nil
	})
	if err != nil {
		return err
	}

	panel, err := compiler.Compile(ast)
	if err != nil {
		return err
	}

	st, err := store.NewSQLiteStore(templateDir + "/data.db")
	if err != nil {
		return err
	}
	if err := st.Migrate(panel.Entities); err != nil {
		return err
	}

	e := engine.New(panel, st)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), e.Router())
}

func mergeASTs(dst *dsl.PanelNode, src dsl.PanelNode) {
	if src.ID != "" {
		dst.ID = src.ID
	}
	if src.Title != "" {
		dst.Title = src.Title
	}
	if src.Theme != "" {
		dst.Theme = src.Theme
	}
	dst.Entities = append(dst.Entities, src.Entities...)
	dst.Pages = append(dst.Pages, src.Pages...)
	dst.Nav = append(dst.Nav, src.Nav...)
	dst.Hooks = append(dst.Hooks, src.Hooks...)
}
