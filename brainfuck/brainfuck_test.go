package brainfuck

import (
	"bytes"
	"os"
	"testing"
)

func TestAddition(t *testing.T) {
	m := RunProgram("+")

	if m.memory[0] != 1 {
		t.Errorf("+ should increment the first value in memory")
	}
}

func TestSubtraction(t *testing.T) {
	m := RunProgram("-")

	if m.memory[0] != -1 {
		t.Errorf("- should decrement the first value in memory")
	}
}

func TestPointerRight(t *testing.T) {
	m := RunProgram(">")

	if m.pointer != 1 {
		t.Errorf("> should move the pointer on one")
	}
}

func TestPointerLeft(t *testing.T) {
	// We can't go left right away (as this would put the pointer out of bounds)
	// So advance two and move back one
	m := RunProgram(">><")

	if m.pointer != 1 {
		t.Errorf("< should move the pointer back one")
	}
}

func TestLoop(t *testing.T) {
	// Test program sets the first cell to 3 and then decrements it down
	m := RunProgram("+++[-]+")

	if m.memory[0] != 1 {
		t.Errorf("[] should loop")
	}

	// Test nested
	m = RunProgram("+++[>++[>+<-]<-]")

	if m.memory[0] != 0 || m.memory[2] != 6 {
		t.Errorf("[] should loop")
	}
}

func TestInput(t *testing.T) {
	in := bytes.NewBufferString("a")
	out := new(bytes.Buffer)
	m := NewMachine(",", in, out)
	m.Run()

	if m.memory[0] != 97 {
		t.Errorf(", should read in one byte")
	}
}

func TestOutput(t *testing.T) {
	in := new(bytes.Buffer)
	out := new(bytes.Buffer)
	m := NewMachine(".", in, out)
	// futz with the memory so the current cell has an "a" in
	m.memory[0] = 97
	m.Run()

	if result, _, err := out.ReadRune(); result != 'a' || err != nil {
		t.Errorf(". should print one byte")
	}
}

func TestComplex(t *testing.T) {
	in := new(bytes.Buffer)
	out := new(bytes.Buffer)
	m := NewMachine("++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++.", in, out)
	m.Run()

	out.WriteTo(os.Stdout)
}
