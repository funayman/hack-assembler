package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	// AtSymbol is the beginning of all A-Instructions in the Hack language
	AtSymbol = '@'

	// OpenParen is the beginning of LABEL
	OpenParen = '('

	// CloseParen is the ending marker for LABEL
	CloseParen = ')'

	// CommentStart is the beginning character of a comment in the hack language
	// If this is found at the begnning of a line, the line is ignored
	CommentStart = '/'
)

type Instruction uint16

func (inst Instruction) String() string {
	return fmt.Sprintf("%016b", inst)
}

type Parser struct {
	source *bufio.Reader // reader for the file
	ch     byte          // current character
	pc     uint16        // program counter
	memLoc uint16        // current position in memory
}

var (
	insts []Instruction
)

func New(source io.Reader) (*Parser, error) {
	p := &Parser{source: bufio.NewReader(source), memLoc: 16}
	p.advance()
	return p, nil
}

func isWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n'
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isNumeric(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (p *Parser) Parse() []Instruction {
	insts = make([]Instruction, 0)
	for p.ch != 0 {
		p.eatWhiteSpace()

		switch p.ch {
		case 0:
			break
		case CommentStart:
			p.eatComment()
		case OpenParen:
			p.parseLabel()
		case AtSymbol:
			insts = append(insts, p.parseAInstruction())
		default:
			insts = append(insts, p.parseCInstruction())
		}

		p.advance()
	}

	// fix indentifiers memory location
	for _, k := range refTableKeys {
		for _, index := range refTable[k] {
			insts[index] = Instruction(p.memLoc)
		}
		p.memLoc++
	}

	return insts
}

func (p *Parser) advance() {
	p.ch, _ = p.source.ReadByte()
}

func (p *Parser) eatWhiteSpace() {
	for ; isWhiteSpace(p.ch); p.advance() {
		// nom nom nom
	}
}

func (p *Parser) eatComment() {
	for ; p.ch != '\n'; p.advance() {
		// nom nom nom
	}
}

func (p *Parser) parseAInstruction() Instruction {
	// advance past the @ symbol
	p.advance()

	var buf strings.Builder
	for ; !isWhiteSpace(p.ch); p.advance() {
		// nom nom nom
		buf.WriteByte(p.ch)
	}

	symbol := buf.String()

	// if symbol in table, increase PC and return value
	if data, ok := SymbolTable[symbol]; ok {
		ref := refTable[symbol]
		refTable[symbol] = append(ref, len(insts))
		p.pc++
		return data
	}

	val, err := strconv.Atoi(symbol)
	if err != nil {
		// not a number, its an identifier
		SymbolTable[symbol] = Instruction(0)

		// add to reference list because the label (if its a label) has yet to be defined
		if ref, ok := refTable[symbol]; ok {
			refTable[symbol] = append(ref, len(insts))
		} else {
			refTable[symbol] = make([]int, 0)
			refTable[symbol] = append(ref, len(insts))
			refTableKeys = append(refTableKeys, symbol)
		}
	} else {
		// store the number for easier lookup
		SymbolTable[symbol] = Instruction(val)
	}

	p.pc++
	return SymbolTable[symbol]
}

func (p *Parser) parseCInstruction() Instruction {
	var parseJump, parseComp bool

	inst := 7 << 13 // 1110000000000000
	// dest = comp ; jump
	var buf strings.Builder
	for ; !isWhiteSpace(p.ch); p.advance() {
		switch p.ch {
		case '=':
			dest := destTable[buf.String()]
			inst |= (dest << 3)
			buf.Reset()

			parseComp = true
		case ';':
			comp := compTable[buf.String()]
			inst |= (comp << 6)
			buf.Reset()

			parseJump = true
			parseComp = false
		default:
			// nom nom nom
			buf.WriteByte(p.ch)
		}
	}

	// clean up buffer
	if parseComp {
		comp := compTable[buf.String()]
		inst |= (comp << 6)
	} else if parseJump {
		jump := jumpTable[buf.String()]
		inst |= jump
	}

	p.pc++
	return Instruction(inst)
}

func (p *Parser) parseLabel() {
	// advance past the open paren
	p.advance()

	var buf strings.Builder
	for ; p.ch != CloseParen; p.advance() {
		// get the label stuff
		buf.WriteByte(p.ch)
	}

	label := buf.String()
	SymbolTable[label] = Instruction(p.pc)

	// backtrack and update previous instructions
	if backInsts, ok := refTable[label]; ok {
		for _, index := range backInsts {
			insts[index] = Instruction(p.pc)
		}

		delete(refTable, label)

		// remove the "label" from the keys
		for i, k := range refTableKeys {
			if k == label {
				refTableKeys = append(refTableKeys[:i], refTableKeys[i+1:]...)
				break
			}
		}
	}
}
