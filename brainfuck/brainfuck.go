package brainfuck

import (
	"bytes"
	"fmt"
	"io"
)

type machine struct {
	memory  []int8
	pointer int

	program        string
	programCounter int

	input  io.Reader
	output io.Writer

	debug bool
}

func NewMachine(s string, in io.Reader, out io.Writer) *machine {
	// Make 2k's worth of "RAM"
	mem := make([]int8, 2048)
	return &machine{
		program: s,
		memory:  mem,
		input:   in,
		output:  out,
	}
}

func RunProgram(s string) *machine {
	in := new(bytes.Buffer)
	out := new(bytes.Buffer)
	m := NewMachine(s, in, out)
	m.Run()
	return m
}

func (m *machine) String() string {
	ins := "EOL"
	if m.programCounter < len(m.program) {
		ins = fmt.Sprintf("%c", m.program[m.programCounter])
	}
	return fmt.Sprintf("pc: %d, ins: %s, ir: %d, irv: %d\n",
		m.programCounter, ins, m.pointer, m.memory[m.pointer])
}

func (m *machine) consumeUntilBalanced(deep uint8, undeep uint8, incr int) {
	depth := 1
	m.programCounter += incr
	for depth > 0 && m.programCounter < len(m.program) && m.programCounter >= 0 {
		if m.debug {
			fmt.Printf("pc: %d, ins: %c, depth: %d\n",
				m.programCounter, m.program[m.programCounter], depth)
		}
		switch m.program[m.programCounter] {
		case deep:
			depth++
		case undeep:
			depth--
		}
		m.programCounter += incr
	}
}

func (m *machine) Run() {
	for m.programCounter < len(m.program) {
		if m.debug {
			fmt.Printf("%s", m)
		}
		switch m.program[m.programCounter] {
		case '+':
			m.memory[m.pointer]++
		case '-':
			m.memory[m.pointer]--
		case '>':
			m.pointer++
			if m.pointer > len(m.memory) {
				panic("Incremented pointer beyond range of memory")
			}
		case '<':
			m.pointer--
			if m.pointer < 0 {
				panic("Decremented pointer beyond range of memory")
			}
		case '[':
			if m.memory[m.pointer] == 0 {
				m.consumeUntilBalanced('[', ']', 1)
				// Program counter must go back as we increment it later
				m.programCounter -= 2
			}
		case ']':
			if m.memory[m.pointer] != 0 {
				m.consumeUntilBalanced(']', '[', -1)
			}
		case ',':
			in := make([]byte, 1)
			_, err := m.input.Read(in)
			if err != nil && err != io.EOF {
				panic("Failed to read from input: " + err.Error())
			}
			m.memory[m.pointer] = int8(in[0])
		case '.':
			m.output.Write([]byte{byte(m.memory[m.pointer])})
		}

		m.programCounter++
	}
	if m.debug {
		fmt.Printf("%s", m)
	}
}
