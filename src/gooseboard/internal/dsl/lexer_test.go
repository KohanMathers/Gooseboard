package dsl

import "testing"

var allTokenKinds = []TokenKind{
	TokEOF, TokIdent, TokString, TokNumber,
	TokLBrace, TokRBrace, TokLParen, TokRParen,
	TokColon, TokComma, TokKeyword,
	TokLBracket, TokRBracket, TokArrow, TokDot, TokSigil,
	TokNull, TokTrue, TokFalse, TokComment,
}

func kindName(k TokenKind) string {
	switch k {
	case TokEOF:
		return "EOF"
	case TokIdent:
		return "Ident"
	case TokString:
		return "String"
	case TokNumber:
		return "Number"
	case TokLBrace:
		return "LBrace"
	case TokRBrace:
		return "RBrace"
	case TokLParen:
		return "LParen"
	case TokRParen:
		return "RParen"
	case TokColon:
		return "Colon"
	case TokComma:
		return "Comma"
	case TokKeyword:
		return "Keyword"
	case TokLBracket:
		return "LBracket"
	case TokRBracket:
		return "RBracket"
	case TokArrow:
		return "Arrow"
	case TokDot:
		return "Dot"
	case TokSigil:
		return "Sigil"
	case TokNull:
		return "Null"
	case TokTrue:
		return "True"
	case TokFalse:
		return "False"
	case TokComment:
		return "Comment"
	default:
		return "Unknown"
	}
}

func assertTokens(t *testing.T, src string, expected []Token) {
	t.Helper()
	lexer := NewLexer(src)
	for i, want := range expected {
		got := lexer.Next()
		if got.Kind != want.Kind || got.Value != want.Value || got.Line != want.Line {
			t.Errorf("token %d mismatch for %q:\n  want %s(%q) line %d\n  got  %s(%q) line %d",
				i, src, kindName(want.Kind), want.Value, want.Line, kindName(got.Kind), got.Value, got.Line)
		}
	}
}

func TestBasicLexerHappyPath(t *testing.T) {
	assertTokens(t, `panel "My Panel" { field "Name" }`, []Token{
		{Kind: TokKeyword, Value: "panel", Line: 1},
		{Kind: TokString, Value: "My Panel", Line: 1},
		{Kind: TokLBrace, Value: "{", Line: 1},
		{Kind: TokKeyword, Value: "field", Line: 1},
		{Kind: TokString, Value: "Name", Line: 1},
		{Kind: TokRBrace, Value: "}", Line: 1},
		{Kind: TokEOF, Value: "", Line: 1},
	})
}

func TestBasicLexerWithNoWhitespace(t *testing.T) {
	assertTokens(t, `panel"My Panel"{field"Name"}`, []Token{
		{Kind: TokKeyword, Value: "panel", Line: 1},
		{Kind: TokString, Value: "My Panel", Line: 1},
		{Kind: TokLBrace, Value: "{", Line: 1},
		{Kind: TokKeyword, Value: "field", Line: 1},
		{Kind: TokString, Value: "Name", Line: 1},
		{Kind: TokRBrace, Value: "}", Line: 1},
		{Kind: TokEOF, Value: "", Line: 1},
	})
}

func TestBasicLexerWithComments(t *testing.T) {
	s := `panel "My Panel" { // This is a comment
	field "Name" // Another comment
	}`
	assertTokens(t, s, []Token{
		{Kind: TokKeyword, Value: "panel", Line: 1},
		{Kind: TokString, Value: "My Panel", Line: 1},
		{Kind: TokLBrace, Value: "{", Line: 1},
		{Kind: TokComment, Value: " This is a comment", Line: 1},
		{Kind: TokKeyword, Value: "field", Line: 2},
		{Kind: TokString, Value: "Name", Line: 2},
		{Kind: TokComment, Value: " Another comment", Line: 2},
		{Kind: TokRBrace, Value: "}", Line: 3},
		{Kind: TokEOF, Value: "", Line: 3},
	})
}

func TestLexerCommentToken(t *testing.T) {
	assertTokens(t, "// hello world\npanel", []Token{
		{Kind: TokComment, Value: " hello world", Line: 1},
		{Kind: TokKeyword, Value: "panel", Line: 2},
		{Kind: TokEOF, Value: "", Line: 2},
	})
}

func TestLexerEmptyComment(t *testing.T) {
	assertTokens(t, "//\npanel", []Token{
		{Kind: TokComment, Value: "", Line: 1},
		{Kind: TokKeyword, Value: "panel", Line: 2},
		{Kind: TokEOF, Value: "", Line: 2},
	})
}

func TestLexerNextSignificantSkipsComments(t *testing.T) {
	s := `panel "My Panel" { // This is a comment
	field "Name" // Another comment
	}`
	lexer := NewLexer(s)
	expected := []Token{
		{Kind: TokKeyword, Value: "panel", Line: 1},
		{Kind: TokString, Value: "My Panel", Line: 1},
		{Kind: TokLBrace, Value: "{", Line: 1},
		{Kind: TokKeyword, Value: "field", Line: 2},
		{Kind: TokString, Value: "Name", Line: 2},
		{Kind: TokRBrace, Value: "}", Line: 3},
		{Kind: TokEOF, Value: "", Line: 3},
	}
	for i, want := range expected {
		got := lexer.NextSignificant()
		if got.Kind != want.Kind || got.Value != want.Value || got.Line != want.Line {
			t.Errorf("token %d mismatch for %q:\n  want %s(%q) line %d\n  got  %s(%q) line %d",
				i, s, kindName(want.Kind), want.Value, want.Line, kindName(got.Kind), got.Value, got.Line)
		}
	}
}

func TestLexerNextSignificantSkipsMultipleConsecutiveComments(t *testing.T) {
	s := "// first\n// second\n// third\npanel"
	lexer := NewLexer(s)

	got := lexer.NextSignificant()
	if got.Kind != TokKeyword || got.Value != "panel" || got.Line != 4 {
		t.Fatalf("expected panel keyword on line 4, got %+v", got)
	}

	got = lexer.NextSignificant()
	if got.Kind != TokEOF {
		t.Fatalf("expected EOF, got %+v", got)
	}
}

func TestLexerNextSignificantWithNoComments(t *testing.T) {
	lexer := NewLexer(`panel "X"`)

	got := lexer.NextSignificant()
	if got.Kind != TokKeyword || got.Value != "panel" {
		t.Fatalf("expected panel keyword, got %+v", got)
	}

	got = lexer.NextSignificant()
	if got.Kind != TokString || got.Value != "X" {
		t.Fatalf("expected string X, got %+v", got)
	}

	got = lexer.NextSignificant()
	if got.Kind != TokEOF {
		t.Fatalf("expected EOF, got %+v", got)
	}
}

func TestLexerNextSignificantTrailingCommentReachesEOF(t *testing.T) {
	lexer := NewLexer(`panel // trailing`)

	got := lexer.NextSignificant()
	if got.Kind != TokKeyword || got.Value != "panel" {
		t.Fatalf("expected panel keyword, got %+v", got)
	}

	got = lexer.NextSignificant()
	if got.Kind != TokEOF {
		t.Fatalf("expected EOF, got %+v", got)
	}
}

func TestLexerDelimiters(t *testing.T) {
	assertTokens(t, `{}()[]`, []Token{
		{Kind: TokLBrace, Value: "{", Line: 1},
		{Kind: TokRBrace, Value: "}", Line: 1},
		{Kind: TokLParen, Value: "(", Line: 1},
		{Kind: TokRParen, Value: ")", Line: 1},
		{Kind: TokLBracket, Value: "[", Line: 1},
		{Kind: TokRBracket, Value: "]", Line: 1},
		{Kind: TokEOF, Value: "", Line: 1},
	})
}

func TestLexerOperators(t *testing.T) {
	assertTokens(t, `: , -> .`, []Token{
		{Kind: TokColon, Value: ":", Line: 1},
		{Kind: TokComma, Value: ",", Line: 1},
		{Kind: TokArrow, Value: "->", Line: 1},
		{Kind: TokDot, Value: ".", Line: 1},
		{Kind: TokEOF, Value: "", Line: 1},
	})
}

func TestLexerLiterals(t *testing.T) {
	assertTokens(t, `true false null`, []Token{
		{Kind: TokTrue, Value: "true", Line: 1},
		{Kind: TokFalse, Value: "false", Line: 1},
		{Kind: TokNull, Value: "null", Line: 1},
		{Kind: TokEOF, Value: "", Line: 1},
	})
}

func TestLexerNumbers(t *testing.T) {
	assertTokens(t, `0 42 12345`, []Token{
		{Kind: TokNumber, Value: "0", Line: 1},
		{Kind: TokNumber, Value: "42", Line: 1},
		{Kind: TokNumber, Value: "12345", Line: 1},
		{Kind: TokEOF, Value: "", Line: 1},
	})
}

func TestLexerNumberDotIsNotDecimal(t *testing.T) {
	assertTokens(t, `3.14`, []Token{
		{Kind: TokNumber, Value: "3", Line: 1},
		{Kind: TokDot, Value: ".", Line: 1},
		{Kind: TokNumber, Value: "14", Line: 1},
		{Kind: TokEOF, Value: "", Line: 1},
	})
}

func TestLexerStringEscapes(t *testing.T) {
	assertTokens(t, `"a\"b\\c\nd\te\rf\zg"`, []Token{
		{Kind: TokString, Value: "a\"b\\c\nd\te\rfzg", Line: 1},
		{Kind: TokEOF, Value: "", Line: 1},
	})
}

func TestLexerUnterminatedStringReachesEOF(t *testing.T) {
	lexer := NewLexer(`"unterminated`)

	tok := lexer.Next()
	if tok.Kind != TokString || tok.Value != "unterminated" {
		t.Fatalf("expected unterminated string token, got %+v", tok)
	}

	tok = lexer.Next()
	if tok.Kind != TokEOF {
		t.Fatalf("expected EOF after unterminated string, got %+v", tok)
	}
}

func TestLexerSigil(t *testing.T) {
	assertTokens(t, `$user $item2 $`, []Token{
		{Kind: TokSigil, Value: "user", Line: 1},
		{Kind: TokSigil, Value: "item2", Line: 1},
		{Kind: TokSigil, Value: "", Line: 1},
		{Kind: TokEOF, Value: "", Line: 1},
	})
}

func TestLexerIdentifiersVsKeywords(t *testing.T) {
	assertTokens(t, `panelName my_var2 Panel`, []Token{
		{Kind: TokIdent, Value: "panelName", Line: 1},
		{Kind: TokIdent, Value: "my_var2", Line: 1},
		{Kind: TokIdent, Value: "Panel", Line: 1},
		{Kind: TokEOF, Value: "", Line: 1},
	})
}

func TestLexerAllKeywords(t *testing.T) {
	keywordList := []string{
		"panel", "entity", "field", "page", "view", "nav", "section", "hook",
		"on", "permission", "role", "action", "tab", "import", "required",
		"default", "emit",
	}
	for _, kw := range keywordList {
		assertTokens(t, kw, []Token{
			{Kind: TokKeyword, Value: kw, Line: 1},
			{Kind: TokEOF, Value: "", Line: 1},
		})
	}
}

func TestLexerUnknownCharacterIsSkipped(t *testing.T) {
	assertTokens(t, `panel @ # "X"`, []Token{
		{Kind: TokKeyword, Value: "panel", Line: 1},
		{Kind: TokString, Value: "X", Line: 1},
		{Kind: TokEOF, Value: "", Line: 1},
	})
}

func TestLexerMultilineString(t *testing.T) {
	s := "\"line one\nline two\" field"
	assertTokens(t, s, []Token{
		{Kind: TokString, Value: "line one\nline two", Line: 1},
		{Kind: TokKeyword, Value: "field", Line: 2},
		{Kind: TokEOF, Value: "", Line: 2},
	})
}

func TestFullLexerWithAllTokens(t *testing.T) {
	s := `import "theme" panel "Dashboard" { page "Home" { nav ["a", "b"] section "Main" { ` +
		`view UserView.Name entity "User" (role: "admin", permission: "read") { ` +
		`field "Name": required, field "Count": default, tab "Info" ` +
		`hook on "create" -> emit($user) action "Save" -> $user.save ` +
		`visible: true archived: false owner: null count: 42 } } } } // trailing comment`

	expected := []Token{
		{Kind: TokKeyword, Value: "import", Line: 1},
		{Kind: TokString, Value: "theme", Line: 1},
		{Kind: TokKeyword, Value: "panel", Line: 1},
		{Kind: TokString, Value: "Dashboard", Line: 1},
		{Kind: TokLBrace, Value: "{", Line: 1},
		{Kind: TokKeyword, Value: "page", Line: 1},
		{Kind: TokString, Value: "Home", Line: 1},
		{Kind: TokLBrace, Value: "{", Line: 1},
		{Kind: TokKeyword, Value: "nav", Line: 1},
		{Kind: TokLBracket, Value: "[", Line: 1},
		{Kind: TokString, Value: "a", Line: 1},
		{Kind: TokComma, Value: ",", Line: 1},
		{Kind: TokString, Value: "b", Line: 1},
		{Kind: TokRBracket, Value: "]", Line: 1},
		{Kind: TokKeyword, Value: "section", Line: 1},
		{Kind: TokString, Value: "Main", Line: 1},
		{Kind: TokLBrace, Value: "{", Line: 1},
		{Kind: TokKeyword, Value: "view", Line: 1},
		{Kind: TokIdent, Value: "UserView", Line: 1},
		{Kind: TokDot, Value: ".", Line: 1},
		{Kind: TokIdent, Value: "Name", Line: 1},
		{Kind: TokKeyword, Value: "entity", Line: 1},
		{Kind: TokString, Value: "User", Line: 1},
		{Kind: TokLParen, Value: "(", Line: 1},
		{Kind: TokKeyword, Value: "role", Line: 1},
		{Kind: TokColon, Value: ":", Line: 1},
		{Kind: TokString, Value: "admin", Line: 1},
		{Kind: TokComma, Value: ",", Line: 1},
		{Kind: TokKeyword, Value: "permission", Line: 1},
		{Kind: TokColon, Value: ":", Line: 1},
		{Kind: TokString, Value: "read", Line: 1},
		{Kind: TokRParen, Value: ")", Line: 1},
		{Kind: TokLBrace, Value: "{", Line: 1},
		{Kind: TokKeyword, Value: "field", Line: 1},
		{Kind: TokString, Value: "Name", Line: 1},
		{Kind: TokColon, Value: ":", Line: 1},
		{Kind: TokKeyword, Value: "required", Line: 1},
		{Kind: TokComma, Value: ",", Line: 1},
		{Kind: TokKeyword, Value: "field", Line: 1},
		{Kind: TokString, Value: "Count", Line: 1},
		{Kind: TokColon, Value: ":", Line: 1},
		{Kind: TokKeyword, Value: "default", Line: 1},
		{Kind: TokComma, Value: ",", Line: 1},
		{Kind: TokKeyword, Value: "tab", Line: 1},
		{Kind: TokString, Value: "Info", Line: 1},
		{Kind: TokKeyword, Value: "hook", Line: 1},
		{Kind: TokKeyword, Value: "on", Line: 1},
		{Kind: TokString, Value: "create", Line: 1},
		{Kind: TokArrow, Value: "->", Line: 1},
		{Kind: TokKeyword, Value: "emit", Line: 1},
		{Kind: TokLParen, Value: "(", Line: 1},
		{Kind: TokSigil, Value: "user", Line: 1},
		{Kind: TokRParen, Value: ")", Line: 1},
		{Kind: TokKeyword, Value: "action", Line: 1},
		{Kind: TokString, Value: "Save", Line: 1},
		{Kind: TokArrow, Value: "->", Line: 1},
		{Kind: TokSigil, Value: "user", Line: 1},
		{Kind: TokDot, Value: ".", Line: 1},
		{Kind: TokIdent, Value: "save", Line: 1},
		{Kind: TokIdent, Value: "visible", Line: 1},
		{Kind: TokColon, Value: ":", Line: 1},
		{Kind: TokTrue, Value: "true", Line: 1},
		{Kind: TokIdent, Value: "archived", Line: 1},
		{Kind: TokColon, Value: ":", Line: 1},
		{Kind: TokFalse, Value: "false", Line: 1},
		{Kind: TokIdent, Value: "owner", Line: 1},
		{Kind: TokColon, Value: ":", Line: 1},
		{Kind: TokNull, Value: "null", Line: 1},
		{Kind: TokIdent, Value: "count", Line: 1},
		{Kind: TokColon, Value: ":", Line: 1},
		{Kind: TokNumber, Value: "42", Line: 1},
		{Kind: TokRBrace, Value: "}", Line: 1},
		{Kind: TokRBrace, Value: "}", Line: 1},
		{Kind: TokRBrace, Value: "}", Line: 1},
		{Kind: TokRBrace, Value: "}", Line: 1},
		{Kind: TokComment, Value: " trailing comment", Line: 1},
		{Kind: TokEOF, Value: "", Line: 1},
	}

	assertTokens(t, s, expected)

	seen := make(map[TokenKind]bool)
	for _, tok := range expected {
		seen[tok.Kind] = true
	}
	for _, k := range allTokenKinds {
		if !seen[k] {
			t.Errorf("TestFullLexerWithAllTokens does not exercise token kind %s — add it to the source", kindName(k))
		}
	}
}
