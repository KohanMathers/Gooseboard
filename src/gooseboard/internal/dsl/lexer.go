package dsl

import "strings"

type TokenKind int

const (
	TokEOF TokenKind = iota
	TokIdent
	TokString
	TokNumber
	TokLBrace
	TokRBrace
	TokLParen
	TokRParen
	TokColon
	TokComma
	TokKeyword
	TokLBracket
	TokRBracket
	TokArrow
	TokDot
	TokSigil
	TokNull
	TokTrue
	TokFalse
	TokComment
)

type Token struct {
	Kind  TokenKind
	Value string
	Line  int
}

var keywords = map[string]bool{
	"panel": true, "entity": true, "field": true, "page": true,
	"view": true, "nav": true, "section": true, "hook": true,
	"on": true, "permission": true, "role": true, "action": true,
	"tab": true, "import": true, "required": true, "default": true,
	"true": true, "false": true, "null": true, "emit": true, "ref": true,
}

type Lexer struct {
	src  []rune
	pos  int
	line int
}

func NewLexer(src string) *Lexer {
	return &Lexer{src: []rune(src), line: 1}
}

func (l *Lexer) Next() Token {
	l.skipWhitespace()
	if l.pos >= len(l.src) {
		return Token{Kind: TokEOF, Line: l.line}
	}

	ch := l.src[l.pos]

	switch {
	case ch == '/' && l.pos+1 < len(l.src) && l.src[l.pos+1] == '/':
		return l.readComment()
	case ch == '{':
		l.pos++
		return Token{Kind: TokLBrace, Value: "{", Line: l.line}
	case ch == '}':
		l.pos++
		return Token{Kind: TokRBrace, Value: "}", Line: l.line}
	case ch == '(':
		l.pos++
		return Token{Kind: TokLParen, Value: "(", Line: l.line}
	case ch == ')':
		l.pos++
		return Token{Kind: TokRParen, Value: ")", Line: l.line}
	case ch == ':':
		l.pos++
		return Token{Kind: TokColon, Value: ":", Line: l.line}
	case ch == ',':
		l.pos++
		return Token{Kind: TokComma, Value: ",", Line: l.line}
	case ch == '[':
		l.pos++
		return Token{Kind: TokLBracket, Value: "[", Line: l.line}
	case ch == ']':
		l.pos++
		return Token{Kind: TokRBracket, Value: "]", Line: l.line}
	case ch == '-' && l.pos+1 < len(l.src) && l.src[l.pos+1] == '>':
		l.pos += 2
		return Token{Kind: TokArrow, Value: "->", Line: l.line}
	case ch == '.':
		l.pos++
		return Token{Kind: TokDot, Value: ".", Line: l.line}
	case ch == '$':
		return l.readSigil()
	case ch == '"':
		return l.readString()
	case isDigit(ch):
		return l.readNumber()
	case isIdentStart(ch):
		return l.readIdentOrKeyword()
	default:
		l.pos++
		// TODO: Proper error handling
		return l.Next()
	}
}

func (l *Lexer) NextSignificant() Token {
	tok := l.Next()
	for tok.Kind == TokComment {
		tok = l.Next()
	}
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.src) {
		ch := l.src[l.pos]
		switch ch {
		case '\n':
			l.line++
			l.pos++
		case ' ', '\t', '\r':
			l.pos++
		default:
			return
		}
	}
}

func (l *Lexer) readComment() Token {
	line := l.line
	l.pos += 2
	start := l.pos
	for l.pos < len(l.src) && l.src[l.pos] != '\n' {
		l.pos++
	}
	return Token{Kind: TokComment, Value: string(l.src[start:l.pos]), Line: line}
}

func isDigit(ch rune) bool { return ch >= '0' && ch <= '9' }
func isIdentStart(ch rune) bool {
	return ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}
func isIdentPart(ch rune) bool { return isIdentStart(ch) || isDigit(ch) }

func (l *Lexer) readString() Token {
	l.pos++
	line := l.line
	var sb strings.Builder
	for l.pos < len(l.src) && l.src[l.pos] != '"' {
		ch := l.src[l.pos]
		if ch == '\\' && l.pos+1 < len(l.src) {
			l.pos++
			switch l.src[l.pos] {
			case '"':
				sb.WriteRune('"')
			case '\\':
				sb.WriteRune('\\')
			case 'n':
				sb.WriteRune('\n')
			case 't':
				sb.WriteRune('\t')
			case 'r':
				sb.WriteRune('\r')
			default:
				sb.WriteRune(l.src[l.pos])
			}
			l.pos++
			continue
		}
		if ch == '\n' {
			l.line++
		}
		sb.WriteRune(ch)
		l.pos++
	}
	l.pos++
	return Token{Kind: TokString, Value: sb.String(), Line: line}
}

func (l *Lexer) readNumber() Token {
	start := l.pos
	for l.pos < len(l.src) && isDigit(l.src[l.pos]) {
		l.pos++
	}
	return Token{Kind: TokNumber, Value: string(l.src[start:l.pos]), Line: l.line}
}

func (l *Lexer) readSigil() Token {
	l.pos++
	start := l.pos
	for l.pos < len(l.src) && isIdentPart(l.src[l.pos]) {
		l.pos++
	}
	return Token{Kind: TokSigil, Value: string(l.src[start:l.pos]), Line: l.line}
}

func (l *Lexer) readIdentOrKeyword() Token {
	start := l.pos
	for l.pos < len(l.src) && isIdentPart(l.src[l.pos]) {
		l.pos++
	}
	word := string(l.src[start:l.pos])
	if keywords[word] {
		switch word {
		case "true":
			return Token{Kind: TokTrue, Value: word, Line: l.line}
		case "false":
			return Token{Kind: TokFalse, Value: word, Line: l.line}
		case "null":
			return Token{Kind: TokNull, Value: word, Line: l.line}
		}
		return Token{Kind: TokKeyword, Value: word, Line: l.line}
	}
	return Token{Kind: TokIdent, Value: word, Line: l.line}
}
