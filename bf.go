/*
 Brainfuck-Go ( http://github.com/kgabis/brainfuck-go )
 Copyright (c) 2013 Krzysztof Gabis

 Permission is hereby granted, free of charge, to any person obtaining a copy
 of this software and associated documentation files (the "Software"), to deal
 in the Software without restriction, including without limitation the rights
 to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 copies of the Software, and to permit persons to whom the Software is
 furnished to do so, subject to the following conditions:

 The above copyright notice and this permission notice shall be included in
 all copies or substantial portions of the Software.

 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 THE SOFTWARE.
*/
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Instruction struct {
	operator uint16
	operand  uint16
}

const (
	op_inc_dp = iota
	op_dec_dp
	op_inc_val
	op_dec_val
	op_out
	op_in
	op_jmp_fwd
	op_jmp_bck
)

const data_size int = 65535

func compile_bf(input string) (program []Instruction, err error) {
	var pc, jmp_pc uint16 = 0, 0
	jmp_stack := make([]uint16, 0)
	for _, c := range input {
		switch c {
		case '>':
			program = append(program, Instruction{op_inc_dp, 0})
		case '<':
			program = append(program, Instruction{op_dec_dp, 0})
		case '+':
			program = append(program, Instruction{op_inc_val, 0})
		case '-':
			program = append(program, Instruction{op_dec_val, 0})
		case '.':
			program = append(program, Instruction{op_out, 0})
		case ',':
			program = append(program, Instruction{op_in, 0})
		case '[':
			program = append(program, Instruction{op_jmp_fwd, 0})
			jmp_stack = append(jmp_stack, pc)
		case ']':
			if len(jmp_stack) == 0 {
				return nil, errors.New("Compilation error.")
			}
			jmp_pc = jmp_stack[len(jmp_stack)-1]
			jmp_stack = jmp_stack[:len(jmp_stack)-1]
			program = append(program, Instruction{op_jmp_bck, jmp_pc})
			program[jmp_pc].operand = pc
		default:
			pc--
		}
		pc++
	}
	if len(jmp_stack) != 0 {
		return nil, errors.New("Compilation error.")
	}
	return
}

func execute_bf(program []Instruction) {
	data := make([]int16, data_size)
	var data_ptr uint16 = 0
	reader := bufio.NewReader(os.Stdin)
	for pc := 0; pc < len(program); pc++ {
		switch program[pc].operator {
		case op_inc_dp:
			data_ptr++
		case op_dec_dp:
			data_ptr--
		case op_inc_val:
			data[data_ptr]++
		case op_dec_val:
			data[data_ptr]--
		case op_out:
			fmt.Printf("%c", data[data_ptr])
		case op_in:
			read_val, _ := reader.ReadByte()
			data[data_ptr] = int16(read_val)
		case op_jmp_fwd:
			if data[data_ptr] == 0 {
				pc = int(program[pc].operand)
			}
		case op_jmp_bck:
			if data[data_ptr] > 0 {
				pc = int(program[pc].operand)
			}
		default:
			panic("Unknown operator.")
		}
	}
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: %s filename\n", args[0])
		return
	}
	filename := args[1]
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading %s\n", filename)
		return
	}
	program, err := compile_bf(string(fileContents))
	if err != nil {
		fmt.Println(err)
		return
	}
	execute_bf(program)
}
