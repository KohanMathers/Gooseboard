package dsl

type Parser struct {
	lex *Lexer
	cur Token
}

func NewParser(src string) *Parser {
	p := &Parser{lex: NewLexer(src)}
	p.advance()
	return p
}

func (p *Parser) advance() { p.cur = p.lex.NextSignificant() }

func (p *Parser) peek() Token {
	savedPos, savedLine := p.lex.pos, p.lex.line
	tok := p.lex.NextSignificant()
	p.lex.pos, p.lex.line = savedPos, savedLine
	return tok
}

func (p *Parser) expect(k TokenKind) Token {
	if p.cur.Kind != k {
		// TODO: Proper error handling
		panic("unexpected token")
	}
	t := p.cur
	p.advance()
	return t
}

func (p *Parser) expectBool() bool {
	switch p.cur.Kind {
	case TokTrue:
		p.advance()
		return true
	case TokFalse:
		p.advance()
		return false
	default:
		panic("expected bool, got: " + p.cur.Value)
	}
}

func (p *Parser) isKeyword(word string) bool {
	return p.cur.Kind == TokKeyword && p.cur.Value == word
}

func (p *Parser) expectKeyword(word string) Token {
	if !p.isKeyword(word) {
		panic("expected keyword " + word + ", got: " + p.cur.Value)
	}
	t := p.cur
	p.advance()
	return t
}

func (p *Parser) ParsePanel() PanelNode {
	var panel PanelNode
	for p.cur.Kind != TokEOF {
		switch {
		case p.isKeyword("panel"):
			p.parsePanelHeader(&panel)
		case p.isKeyword("entity"):
			panel.Entities = append(panel.Entities, p.parseEntity())
		case p.isKeyword("page"):
			panel.Pages = append(panel.Pages, p.parsePage())
		case p.isKeyword("nav"):
			panel.Nav = p.parseNav()
		case p.isKeyword("hook"):
			panel.Hooks = append(panel.Hooks, p.parseHook())
		default:
			panic("unexpected top-level token: " + p.cur.Value)
		}
	}
	return panel
}

func (p *Parser) parsePanelHeader(panel *PanelNode) {
	p.expectKeyword("panel")
	idTok := p.expect(TokString)
	panel.ID = idTok.Value
	p.expect(TokLBrace)
	for p.cur.Kind != TokRBrace {
		switch {
		case p.cur.Kind == TokIdent && p.cur.Value == "title":
			p.advance()
			p.expect(TokColon)
			titleTok := p.expect(TokString)
			panel.Title = titleTok.Value
		case p.cur.Kind == TokIdent && p.cur.Value == "theme":
			p.advance()
			p.expect(TokColon)
			panel.Theme = p.parseImportPath()
		case p.isKeyword("entity"):
			panel.Entities = append(panel.Entities, p.parseEntity())
		case p.isKeyword("page"):
			panel.Pages = append(panel.Pages, p.parsePage())
		case p.isKeyword("nav"):
			panel.Nav = p.parseNav()
		case p.isKeyword("hook"):
			panel.Hooks = append(panel.Hooks, p.parseHook())
		default:
			panic("unexpected token in panel header: " + p.cur.Value)
		}
	}
	p.expect(TokRBrace)
}

func (p *Parser) parseImportPath() string {
	p.expectKeyword("import")
	p.expect(TokLParen)
	pathTok := p.expect(TokString)
	p.expect(TokRParen)
	return pathTok.Value
}

func (p *Parser) parseEntity() EntityNode {
	p.expectKeyword("entity")
	nameTok := p.expect(TokIdent)
	entity := EntityNode{Name: nameTok.Value}
	p.expect(TokLBrace)
	for p.cur.Kind != TokRBrace {
		switch {
		case p.isKeyword("field"):
			entity.Fields = append(entity.Fields, p.parseField())
		case p.isKeyword("permission"):
			entity.Permissions = append(entity.Permissions, p.parsePermission())
		default:
			panic("expected 'field' or 'permission' keyword in entity, got: " + p.cur.Value)
		}
	}
	p.expect(TokRBrace)
	return entity
}

func (p *Parser) parseField() FieldNode {
	p.expectKeyword("field")
	nameTok := p.expect(TokIdent)
	field := FieldNode{Name: nameTok.Value}
	field.Type, field.RefTo = p.parseFieldType()

	for {
		switch {
		case p.isKeyword("required"):
			p.advance()
			field.Required = true
		case p.isKeyword("default"):
			p.advance()
			p.expect(TokLParen)
			field.Default = p.parseValue()
			p.expect(TokRParen)
		default:
			return field
		}
	}
}

func (p *Parser) parseFieldType() (typ string, refTo string) {
	if p.isKeyword("ref") {
		p.advance()
		p.expect(TokLParen)
		entTok := p.expect(TokIdent)
		p.expect(TokRParen)
		return "ref", entTok.Value
	}
	typTok := p.expect(TokIdent)
	return typTok.Value, ""
}

func (p *Parser) parsePermission() PermissionNode {
	p.expectKeyword("permission")
	actionTok := p.expect(TokIdent)
	p.expectKeyword("role")
	roles := p.parseIdentListParen()
	return PermissionNode{Action: actionTok.Value, Roles: roles}
}

func (p *Parser) parseIdentListParen() []string {
	p.expect(TokLParen)
	var items []string
	for p.cur.Kind != TokRParen {
		idTok := p.expect(TokIdent)
		items = append(items, idTok.Value)
		if p.cur.Kind == TokComma {
			p.advance()
		}
	}
	p.expect(TokRParen)
	return items
}

func (p *Parser) parseValue() any {
	switch {
	case p.cur.Kind == TokSigil:
		return p.parseSigil()
	case p.cur.Kind == TokIdent:
		nameTok := p.cur
		p.advance()
		if p.cur.Kind == TokLParen {
			return CallExpr{Name: nameTok.Value, Args: p.parseArgList()}
		}
		return nameTok.Value
	default:
		return p.parseLiteralValue()
	}
}

func (p *Parser) parseSigil() SigilExpr {
	tok := p.expect(TokSigil)
	path := []string{tok.Value}
	for p.cur.Kind == TokDot {
		p.advance()
		idTok := p.expect(TokIdent)
		path = append(path, idTok.Value)
	}
	return SigilExpr{Path: path}
}

func (p *Parser) parseArgList() []Arg {
	p.expect(TokLParen)
	var args []Arg
	for p.cur.Kind != TokRParen {
		args = append(args, p.parseArg())
		if p.cur.Kind == TokComma {
			p.advance()
		}
	}
	p.expect(TokRParen)
	return args
}

func (p *Parser) parseArg() Arg {
	if p.cur.Kind == TokIdent && p.peek().Kind == TokColon {
		name := p.cur.Value
		p.advance()
		p.advance()
		return Arg{Name: name, Value: p.parseValue()}
	}
	return Arg{Value: p.parseValue()}
}

func (p *Parser) parseLiteralValue() any {
	switch p.cur.Kind {
	case TokString:
		t := p.expect(TokString)
		return t.Value
	case TokNumber:
		t := p.expect(TokNumber)
		return t.Value
	case TokTrue:
		p.advance()
		return true
	case TokFalse:
		p.advance()
		return false
	case TokNull:
		p.advance()
		return nil
	default:
		panic("expected literal value, got: " + p.cur.Value)
	}
}

func (p *Parser) parsePage() PageNode {
	p.expectKeyword("page")
	titleTok := p.expect(TokString)
	page := PageNode{Title: titleTok.Value}

	p.expect(TokLBrace)
	sawView := false
	for p.cur.Kind != TokRBrace {
		switch {
		case p.cur.Kind == TokIdent && p.cur.Value == "route":
			p.advance()
			p.expect(TokColon)
			routeTok := p.expect(TokString)
			page.Route = routeTok.Value
		case p.cur.Kind == TokIdent && p.cur.Value == "icon":
			p.advance()
			p.expect(TokColon)
			iconTok := p.expect(TokString)
			page.Icon = iconTok.Value
		case p.cur.Kind == TokIdent && p.cur.Value == "hidden":
			p.advance()
			p.expect(TokColon)
			page.Hidden = p.expectBool()
		case p.isKeyword("view"):
			page.View = p.parseView()
			sawView = true
		default:
			panic("unexpected token in page: " + p.cur.Value)
		}
	}
	p.expect(TokRBrace)

	if !sawView {
		panic("page " + page.Title + " has no view")
	}
	return page
}

func (p *Parser) parseView() ViewNode {
	p.expectKeyword("view")
	p.expect(TokColon)
	kindTok := p.expect(TokIdent)
	view := ViewNode{Kind: kindTok.Value}

	if p.cur.Kind == TokLParen {
		for _, arg := range p.parseArgList() {
			if arg.Name == "" {
				if ident, ok := arg.Value.(string); ok && view.Entity == "" {
					view.Entity = ident
					continue
				}
			}
			view.Args = append(view.Args, arg)
		}
	}

	p.expect(TokLBrace)
	for p.cur.Kind != TokRBrace {
		switch {
		case p.cur.Kind == TokIdent && p.cur.Value == "columns":
			p.advance()
			p.expect(TokColon)
			view.Columns = p.parseIdentList()
		case p.isKeyword("action"):
			view.Actions = append(view.Actions, p.parseAction())
		case p.isKeyword("tab"):
			view.Tabs = append(view.Tabs, p.parseTab())
		default:
			panic("unexpected token in view: " + p.cur.Value)
		}
	}
	p.expect(TokRBrace)
	return view
}

func (p *Parser) parseAction() ActionNode {
	p.expectKeyword("action")
	labelTok := p.expect(TokString)
	p.expect(TokArrow)
	return ActionNode{Label: labelTok.Value, Target: p.parseTarget()}
}

func (p *Parser) parseTab() TabNode {
	p.expectKeyword("tab")
	labelTok := p.expect(TokString)
	p.expect(TokArrow)
	return TabNode{Label: labelTok.Value, Target: p.parseTarget()}
}

func (p *Parser) parseTarget() CallExpr {
	if p.cur.Kind != TokIdent && p.cur.Kind != TokKeyword {
		panic("expected call target, got: " + p.cur.Value)
	}
	nameTok := p.cur
	p.advance()
	return CallExpr{Name: nameTok.Value, Args: p.parseArgList()}
}

func (p *Parser) parseIdentList() []string {
	p.expect(TokLBracket)
	var items []string
	for p.cur.Kind != TokRBracket {
		idTok := p.expect(TokIdent)
		items = append(items, idTok.Value)
		if p.cur.Kind == TokComma {
			p.advance()
		}
	}
	p.expect(TokRBracket)
	return items
}

func (p *Parser) parseNav() []NavSectionNode {
	p.expectKeyword("nav")
	var sections []NavSectionNode
	p.expect(TokLBrace)
	for p.cur.Kind != TokRBrace {
		if !p.isKeyword("section") {
			panic("expected 'section' keyword in nav, got: " + p.cur.Value)
		}
		p.advance()
		titleTok := p.expect(TokString)
		section := NavSectionNode{Title: titleTok.Value}

		p.expect(TokLBracket)
		for p.cur.Kind != TokRBracket {
			switch p.cur.Kind {
			case TokString:
				t := p.expect(TokString)
				section.Pages = append(section.Pages, t.Value)
			case TokIdent:
				t := p.expect(TokIdent)
				section.Pages = append(section.Pages, t.Value)
			default:
				panic("expected page name in nav section, got: " + p.cur.Value)
			}
			if p.cur.Kind == TokComma {
				p.advance()
			}
		}
		p.expect(TokRBracket)

		if p.isKeyword("role") {
			p.advance()
			section.Roles = p.parseIdentListParen()
		}

		sections = append(sections, section)
	}
	p.expect(TokRBrace)
	return sections
}

func (p *Parser) parseHook() HookNode {
	p.expectKeyword("hook")
	p.expectKeyword("on")
	entityTok := p.expect(TokIdent)
	p.expect(TokDot)
	eventTok := p.expect(TokIdent)
	hook := HookNode{Entity: entityTok.Value, Event: eventTok.Value}

	p.expect(TokLBrace)
	for p.cur.Kind != TokRBrace {
		if !p.isKeyword("emit") {
			panic("expected 'emit' keyword in hook, got: " + p.cur.Value)
		}
		p.advance()
		emitTok := p.expect(TokString)
		hook.Emits = append(hook.Emits, emitTok.Value)
	}
	p.expect(TokRBrace)
	return hook
}
