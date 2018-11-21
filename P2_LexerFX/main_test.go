package main

import (
	"testing"
)

func TestNewLexer(t *testing.T) {
	lexer, _ := NewLexer("lang.fx")
	if lexer.file != "lang.fx" {
		t.Error("Archivo incorrecto")
	}
	if lexer.line != 1 {
		t.Error("Linea incorrecta")
	}
	if lexer.lastrune != 0 {
		t.Error("Lastrune incorrecta")
	}
}

func TestGet(t *testing.T) {
	lexer, _ := NewLexer("lang.fx")
	if lexer.get() != 47 {
		t.Error("La runa debería ser 47 (/)")
	}
}

func TestLex(t *testing.T) {
	lexer, _ := NewLexer("lang.fx")
	token, _ := lexer.Lex()
	if token.lexema != "type" {
		t.Error("token.lexema deberia ser type")
	}
	if token.tokType != 11 {
		t.Error("token.tokType deberia ser 11")
	}
	if token.valor != "" {
		t.Error("token.valor deberia estar vacio")
	}

}
