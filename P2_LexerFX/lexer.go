package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"unicode"
)

type tokType int

const (
	ComentTok   tokType = iota + 10 //10
	StrTok                          //11
	EOFTok                          //12
	EOLTok                          //13
	IntValTok                       //14
	FloatValTok                     //15
	OpParTok                        //16
	OpMulTok                        //17
	OpDivTok                        //18
	OpSumTok                        //19
	OpResTok                        //20
	CommaTok                        //21
	OpCorTok                        //22
	BoolTok                         //23
	ModuleTok                       //24
	MoreTok                         //25
	LessTok                         //26
	DotCommaTok                     //27
	DotsTok                         //28
	OrTok                           //29
	AndTok                          //30
	SquareTok                       //31
	ExclTok                         //32
	HexadecTok                      //33
	CorchTok                        //34
	AsigTok                         //35
	EqualTok                        //36
	RecordTok                       //37
	CommentTok                      //38
)

type RuneScanner interface {
	ReadRune() (r rune, size int, err error)
	UnreadRune() error
}
type Lexer struct {
	file     string
	line     int
	r        RuneScanner
	lastrune rune

	accepted []rune
}
type Token struct {
	lexema string
	tokType
	valor string
}

func NewLexer(file string) (l *Lexer, err error) {
	l = &Lexer{line: 1}
	// l.line = 1
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	} // Hay error luego devuelvo nada como lexema y el error
	// Si no hay error -> devuelvo el Lexema y nil de error
	l.file = file
	l.r = bufio.NewReader(f)
	return l, nil
}
func (l *Lexer) get() (r rune) {
	var err error
	r, _, err = l.r.ReadRune()
	if err == nil { // No hay error
		l.lastrune = r // La ultima runa es la actual
	} else if err == io.EOF { // Si encuentro un final de fichero debo devolver una runa vacia
		l.lastrune = 0
		return l.lastrune
	} else if err != nil { // Si hay error llamo a panic para controlarlo
		panic(err)
	}
	if r == '\n' { // Si la ultima runa es un fin de linea tengo una linea mas
		l.line++
	}
	l.accepted = append(l.accepted, r) // Si no se da ninguno de los casos anteriores se acepta la runa y se añade a las aceptadas
	return r
}
func (l *Lexer) unget() {
	var err error
	if l.lastrune == 0 { // Si la ultima linea es un EOF
		return
	}
	err = l.r.UnreadRune()
	if err == nil && l.lastrune == '\n' {
		l.line--
	}
	l.lastrune = unicode.ReplacementChar
	if len(l.accepted) != 0 {
		l.accepted = l.accepted[0:(len(l.accepted) - 1)] // Cojo la rodaja de las runas aceptadas menos la actual
	}
	if err != nil { // Si hay error llamo a panic para controlarlo
		panic(err)
	}
	return
}
func (l *Lexer) accept() (token string) { // Compruebo que mis runas son validas, no son espacios en blanco
	token = string(l.accepted)
	if token == "" && l.lastrune != 0 {
		panic(errors.New("empty token"))
	}
	l.accepted = nil
	return token
}
func (l *Lexer) Lex() (t Token, err error) {
	for r := l.get(); ; r = l.get() { // Voy cogiendo runas y leyendolas
		if unicode.IsSpace(r) { // Si es un espacio
			l.accept() // Se acepta como lexema
			continue
		}
		if unicode.IsLetter(r) { // Si es una letra
			l.unget()
			t, err = l.lexId()
			return t, err
		} else if unicode.IsDigit(r) {
			l.unget()
			t, err = l.lexNum()
			return t, err
		} else { // Si no es una letra será otra cosa
			switch r {
			case '(', ')', '*', '+', '-', ',', '[', ']', '%', '<', '>', ':', ';', '|', '&', '^', '!', '{', '}', '=', '.':
				t.tokType = TokType(r, l) // TokType sera una funcion que me reciba una runa y me diga que tipo de token es
				t.lexema = l.accept()
				return t, nil
			case '/':
				if l.get() != '/' {
					t.tokType = OpDivTok
					t.lexema = l.accept()
					return t, nil
				} else {
					for r = l.get(); r != 10; r = l.get() {
					}
					l.accept()
				}
			case 0: // Si es un final de fichero es 0
				t.tokType = EOFTok // El tipo de token es el de final de fichero
				l.accept()
				return t, nil
			default:
				errs := fmt.Sprintf("bad rune %c: %x", r, r)
				return t, errors.New(errs)
			}
		}
	}

	return t, err
}
func TokType(r rune, l *Lexer) tokType {
	var OperationType tokType
	switch r {
	case '(', ')':
		OperationType = OpParTok
	case '[', ']':
		OperationType = OpCorTok
	case '*':
		OperationType = OpMulTok
	case '+':
		OperationType = OpSumTok
	case '-':
		OperationType = OpResTok
	case ',':
		OperationType = CommaTok
	case '%':
		OperationType = ModuleTok
	case '<':
		OperationType = LessTok
	case '>':
		OperationType = MoreTok
	case ';':
		OperationType = DotCommaTok
	case ':':
		if l.get() == '=' {
			OperationType = AsigTok
		} else {
			OperationType = DotsTok
		}
	case '|':
		OperationType = OrTok
	case '&':
		OperationType = AndTok
	case '^':
		OperationType = SquareTok
	case '!':
		OperationType = ExclTok
	case '{', '}':
		OperationType = CorchTok
	case '=':
		OperationType = EqualTok
	case '.':
		OperationType = RecordTok
	}
	return OperationType
}
func (l *Lexer) lexId() (t Token, err error) {
	r := l.get()
	if !unicode.IsLetter(r) { // Si no es una letra reportamos error ya que deberia serlo
		return t, errors.New("bad Id, should not happen")
	}
	isAlpha := func(ar rune) bool { // devuelve un booleano si es un numero, letra o guion
		return unicode.IsDigit(ar) || unicode.IsLetter(ar) || ar == '-'
		// return unicode.IsLetter(ar) || ar == '-'
	}
	for r = l.get(); isAlpha(r); r = l.get() {
	}
	l.unget()

	t.lexema = l.accept()
	if t.lexema == "True" || t.lexema == "False" { // Si el lexema es True o False es un booleano y no un string
		t.tokType = BoolTok
	} else {
		t.tokType = StrTok
	}

	return t, err
}
func (l *Lexer) lexNum() (t Token, err error) {
	const (
		Es          = "Ee"
		Signs       = "+-"
		HexadeciSym = "x"
	)
	var hasDot bool = false
	var isHexadec bool = false
	r := l.get()
	if r == '.' { // Si la runa es un punto
		hasDot = true // Es un punto es true
		r = l.get()
	}
	for ; unicode.IsDigit(r); r = l.get() {
	}
	if strings.ContainsRune(HexadeciSym, r) { // Para procesar los tipos 0x2dfadfd, 0x46...
		r = l.get()
		for r = l.get(); unicode.IsDigit(r) || unicode.IsLetter(r); r = l.get() { // Mientras sea numero o letra voy cogiendo runas
		}
		isHexadec = true
	} // Si la runa contiene una x es un numero en hexadecimal
	if r == '.' {
		if hasDot { // Si es un punto de nuevo es que algo va mal no hay un float con dos puntos
			return t, errors.New("bad float [" + l.accept() + "]")
		}
		hasDot = true // Si no es un punto ponemos a true el booleano por si vuelve a pasar que sí este mal
		for r = l.get(); unicode.IsDigit(r); r = l.get() {
		}
	}
	switch {
	case strings.ContainsRune(Es, r): // Si la runa contiene un exponencial simbolo
		r = l.get()                         // Cojo la siguiente
		if strings.ContainsRune(Signs, r) { // Si contiene algun signo
			r = l.get() // Cojo la siguiente
		}
	case isHexadec: // Si es hexadecimal
		l.unget()
		t.lexema = l.accept()
		tokHexVal, err := strconv.ParseInt(t.lexema, 0, 64)
		if err != nil {
			return t, errors.New("bad Hex [" + t.lexema + "]")
		}
		t.tokType = HexadecTok
		t.valor = strconv.FormatInt(tokHexVal, 10)
	case hasDot: // Si es un punto lo suelto
		l.unget()
		break
	case !hasDot: // Si no es un punto debera ser un entero
		l.unget()
		t.lexema = l.accept()
		tokIntVal, err := strconv.ParseInt(t.lexema, 10, 64)
		if err != nil {
			return t, errors.New("bad Int [" + t.lexema + "]")
		}
		t.tokType = IntValTok
		t.valor = strconv.FormatInt(tokIntVal, 10)
		return t, nil
	default:
		return t, errors.New("bad float [" + l.accept() + "]")
	}

	// fmt.Println(hasDot) // Aqui debe valer true
	if hasDot {
		for r = l.get(); unicode.IsDigit(r); r = l.get() {
		}
		l.unget()
		t.lexema = l.accept()
		// tokFloatVal, err := strconv.ParseFloat(t.lexema, 64)
		tokFloatVal, err := strconv.ParseFloat(t.lexema, 64)
		if err != nil {
			return t, errors.New("bad float [" + t.lexema + "]")
		}
		t.tokType = FloatValTok
		t.valor = strconv.FormatFloat(tokFloatVal, 'E', -1, 64)
		return t, nil
	}

	return t, nil
}
func main() {
	const (
		BugMsg = "compiler error :"
		RunMsg = "runtime error :"
	)
	defer func() {
		if r := recover(); r != nil {
			errs := fmt.Sprint(r)
			if strings.HasPrefix(errs, "runtime error :") {
				errs = strings.Replace(errs, RunMsg, BugMsg, 1)
			}
			err := errors.New(errs)
			fmt.Fprintf(os.Stderr, "%s \n %s", err, debug.Stack())
		}
	}()
	fileName := os.Args[1]
	if len(os.Args) != 2 {
		log.Fatal("Introduce el numero de argumentos correcto")
	}
	lex, err := NewLexer(fileName)
	if err != nil {
		log.Fatal()
	}
	// fmt.Println(lex)
	// fmt.Println(lex.get())
	//
	// tok, _ := lex.Lex()
	// fmt.Printf("Type: %T Value: %v\n", tok, tok)
	// fmt.Println(tok)

	var notEnd bool = true
	for notEnd {
		tok, _ := lex.Lex()
		fmt.Println(tok)
		if tok.tokType == 12 { // Si el tipo de token es 12 es que es final de fichero
			notEnd = false
		}
	}

}
