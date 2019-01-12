package parser

const (
	R0 Instruction = iota
	R1
	R2
	R3
	R4
	R5
	R6
	R7
	R8
	R9
	R10
	R11
	R12
	R13
	R14
	R15

	SCREEN   Instruction = 16384
	KEYBOARD Instruction = 24576
)

const (
	SP Instruction = iota
	LCL
	ARG
	THIS
	THAT
)

// jump constants
const (
	_ = iota
	jmpJGT
	jmpJEQ
	jmpJGE
	jmpJLT
	jmpJNE
	jmpJLE
	jmpJMP
)

// destination constants
const (
	_ = iota
	dstM
	dstD
	dstMD
	dstA
	dstAM
	dstAD
	dstAMD
)

// comp constants
const (
	compZero      = 42  // 0101010
	compOne       = 63  // 0111111
	compMinusOne  = 58  // 0111010
	compD         = 12  // 0001100
	compA         = 48  // 0110000
	compM         = 112 // 1110000 a=1
	compNotD      = 13  // 0001101
	compNotA      = 49  // 0110001
	compNotM      = 113 // 1110001 a=1
	compMinusD    = 15  // 0001111
	compMinusA    = 51  // 0110011
	compMinusM    = 115 // 1110011 a=1
	compDPlusOne  = 31  // 0011111
	compAPlusOne  = 55  // 0110111
	compMPlusOne  = 119 // 1110111 a=1
	compDMinusOne = 14  // 0001110
	compAMinusOne = 50  // 0110010
	compMMinusOne = 114 // 1110010 a=1
	compDPlusA    = 2   // 0000010
	compDPlusM    = 66  // 1000010 a=1
	compDMinusA   = 19  // 0010011
	compDMinusM   = 83  // 1010011 a=1
	compAMinusD   = 7   // 0000111
	compMMinusD   = 71  // 1000111 a=1
	compDAndA     = 0   // 0000000
	compDAndM     = 64  // 1000000 a=1
	compDOrA      = 21  // 0010101
	compDorM      = 85  // 1010101 a=1
)

var (
	SymbolTable  map[string]Instruction
	compTable    map[string]int
	destTable    map[string]int
	jumpTable    map[string]int
	refTable     map[string][]int
	refTableKeys []string // used to keep the indentifiers order
)

func init() {
	SymbolTable = make(map[string]Instruction)

	SymbolTable["R0"] = R0
	SymbolTable["R1"] = R1
	SymbolTable["R2"] = R2
	SymbolTable["R3"] = R3
	SymbolTable["R4"] = R4
	SymbolTable["R5"] = R5
	SymbolTable["R6"] = R6
	SymbolTable["R7"] = R7
	SymbolTable["R8"] = R8
	SymbolTable["R9"] = R9
	SymbolTable["R10"] = R10
	SymbolTable["R11"] = R11
	SymbolTable["R12"] = R12
	SymbolTable["R13"] = R13
	SymbolTable["R14"] = R14
	SymbolTable["R15"] = R15

	SymbolTable["SCREEN"] = SCREEN
	SymbolTable["KBD"] = KEYBOARD

	SymbolTable["SP"] = SP
	SymbolTable["LCL"] = LCL
	SymbolTable["ARG"] = ARG
	SymbolTable["THIS"] = THIS
	SymbolTable["THAT"] = THAT

	compTable = make(map[string]int)
	compTable["0"] = compZero
	compTable["1"] = compOne
	compTable["-1"] = compMinusOne
	compTable["D"] = compD
	compTable["A"] = compA
	compTable["M"] = compM
	compTable["!D"] = compNotD
	compTable["!A"] = compNotA
	compTable["!M"] = compNotM
	compTable["-D"] = compMinusD
	compTable["-A"] = compMinusA
	compTable["-M"] = compMinusM
	compTable["D+1"] = compDPlusOne
	compTable["A+1"] = compAPlusOne
	compTable["M+1"] = compMPlusOne
	compTable["D-1"] = compDMinusOne
	compTable["A-1"] = compAMinusOne
	compTable["M-1"] = compMMinusOne
	compTable["D+A"] = compDPlusA
	compTable["D+M"] = compDPlusM
	compTable["D-A"] = compDMinusA
	compTable["D-M"] = compDMinusM
	compTable["A-D"] = compAMinusD
	compTable["M-D"] = compMMinusD
	compTable["D&A"] = compDAndA
	compTable["D&M"] = compDAndM
	compTable["D|A"] = compDOrA
	compTable["D|M"] = compDorM

	destTable = make(map[string]int)
	destTable["M"] = dstM
	destTable["D"] = dstD
	destTable["MD"] = dstMD
	destTable["A"] = dstA
	destTable["AM"] = dstAM
	destTable["AD"] = dstAD
	destTable["AMD"] = dstAMD

	jumpTable = make(map[string]int)
	jumpTable["JGT"] = jmpJGT
	jumpTable["JEQ"] = jmpJEQ
	jumpTable["JGE"] = jmpJGE
	jumpTable["JLT"] = jmpJLT
	jumpTable["JNE"] = jmpJNE
	jumpTable["JLE"] = jmpJLE
	jumpTable["JMP"] = jmpJMP

	refTable = make(map[string][]int)
	refTableKeys = make([]string, 0)
}
