package main

// API surface for external app to call via FFI (using a shared .so binary)

/*
#include <stdlib.h>
#include <stdbool.h>

typedef struct {
    bool successful;
    char* result;
    char* std_output;
} RunResult;
*/
import "C"
import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/ruegerj/monkey/compiler"
	"github.com/ruegerj/monkey/lexer"
	"github.com/ruegerj/monkey/parser"
	"github.com/ruegerj/monkey/vm"
)

//export CompileAndRun
func CompileAndRun(input *C.char) C.RunResult {
	// monkey patch stdout to capture console output
	oldOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	inputStr := C.GoString(input)
	l := lexer.New(inputStr)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		var out bytes.Buffer
		out.WriteString("PARSING ERROR: \n")

		for _, e := range p.Errors() {
			out.WriteString("\n\t")
			out.WriteString(e)
		}

		stdout := captureStdoutAndRestore(r, w, oldOut)
		return createRunResult(false, out.String(), stdout)
	}

	comp := compiler.New()
	err := comp.Compile(program)
	if err != nil {
		compilationErr := fmt.Sprintf("COMPILATION ERROR: %s", err.Error())
		stdout := captureStdoutAndRestore(r, w, oldOut)
		return createRunResult(false, compilationErr, stdout)
	}

	byteCode := comp.Bytecode()
	machine := vm.New(byteCode)
	err = machine.Run()
	if err != nil {
		bytecodeErr := fmt.Sprintf("BYTECODE ERROR: %s", err.Error())
		stdout := captureStdoutAndRestore(r, w, oldOut)
		return createRunResult(false, bytecodeErr, stdout)
	}

	result := machine.LastPoppedStackElem()
	stdout := captureStdoutAndRestore(r, w, oldOut)

	return createRunResult(true, result.Inspect(), stdout)
}

func main() {}

func createRunResult(success bool, result string, stdout string) C.RunResult {
	// cRunResult := (*C.struct_RunResult)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_RunResult{}))))
	cRunResult := C.RunResult{}

	cRunResult.successful = C.bool(success)
	cRunResult.result = C.CString(result)
	cRunResult.std_output = C.CString(stdout)

	return cRunResult
}

func captureStdoutAndRestore(read *os.File, out *os.File, originOut *os.File) string {
	// copy the output in a separate goroutine so printing can't block indefinitely
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, read)
		outC <- buf.String()
	}()

	out.Close()
	os.Stdout = originOut
	captured := <-outC

	return captured
}
