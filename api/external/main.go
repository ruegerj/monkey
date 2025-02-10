package main

// API surface for external app to call via FFI (using a shared .so binary)

import "C"
import (
	"bytes"
	"fmt"

	"github.com/ruegerj/monkey/compiler"
	"github.com/ruegerj/monkey/lexer"
	"github.com/ruegerj/monkey/parser"
	"github.com/ruegerj/monkey/vm"
)

//export CompileAndRun
func CompileAndRun(input *C.char) *C.char {
	inputStr := C.GoString(input)
	l := lexer.New(inputStr)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		var out bytes.Buffer
		out.WriteString("Parsing failed: \n")

		for _, e := range p.Errors() {
			out.WriteString("\n\t")
			out.WriteString(e)
		}

		return C.CString(out.String())
	}

	comp := compiler.New()
	err := comp.Compile(program)
	if err != nil {
		return C.CString(fmt.Sprintf("Compilation failed: %s", err.Error()))
	}

	byteCode := comp.Bytecode()
	machine := vm.New(byteCode)
	err = machine.Run()
	if err != nil {
		return C.CString(fmt.Sprintf("Bytecode execution failed: %s", err.Error()))
	}

	lastPopped := machine.LastPoppedStackElem()
	return C.CString(lastPopped.Inspect())
}

func main() {}
