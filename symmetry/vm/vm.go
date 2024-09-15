package vm

import (
    "fmt"
    "strconv"
)

type VM struct {
    stack   []interface{}
    memory  map[string]interface{}
    program []Instruction
    pc      int // Program counter
}

type Instruction struct {
    OpCode  string
    Operand interface{}
}

func NewVM() *VM {
    return &VM{
        stack:   make([]interface{}, 0),
        memory:  make(map[string]interface{}),
        program: make([]Instruction, 0),
        pc:      0,
    }
}

func (vm *VM) LoadProgram(program []Instruction) {
    vm.program = program
    vm.pc = 0
}

func (vm *VM) Run() error {
    for vm.pc < len(vm.program) {
        instruction := vm.program[vm.pc]
        if err := vm.executeInstruction(instruction); err != nil {
            return err
        }
        vm.pc++
    }
    return nil
}

func (vm *VM) executeInstruction(instruction Instruction) error {
    switch instruction.OpCode {
    case "PUSH":
        vm.stack = append(vm.stack, instruction.Operand)
    case "POP":
        if len(vm.stack) == 0 {
            return fmt.Errorf("stack underflow")
        }
        vm.stack = vm.stack[:len(vm.stack)-1]
    case "ADD":
        if len(vm.stack) < 2 {
            return fmt.Errorf("not enough operands for ADD")
        }
        b := vm.pop()
        a := vm.pop()
        result, err := vm.add(a, b)
        if err != nil {
            return err
        }
        vm.stack = append(vm.stack, result)
    case "STORE":
        if len(vm.stack) < 1 {
            return fmt.Errorf("not enough operands for STORE")
        }
        value := vm.pop()
        key, ok := instruction.Operand.(string)
        if !ok {
            return fmt.Errorf("STORE operand must be a string")
        }
        vm.memory[key] = value
    case "LOAD":
        key, ok := instruction.Operand.(string)
        if !ok {
            return fmt.Errorf("LOAD operand must be a string")
        }
        value, exists := vm.memory[key]
        if !exists {
            return fmt.Errorf("variable %s not found", key)
        }
        vm.stack = append(vm.stack, value)
    default:
        return fmt.Errorf("unknown opcode: %s", instruction.OpCode)
    }
    return nil
}

func (vm *VM) pop() interface{} {
    if len(vm.stack) == 0 {
        return nil
    }
    value := vm.stack[len(vm.stack)-1]
    vm.stack = vm.stack[:len(vm.stack)-1]
    return value
}

func (vm *VM) add(a, b interface{}) (interface{}, error) {
    // Try to convert operands to float64
    aFloat, aErr := vm.toFloat64(a)
    bFloat, bErr := vm.toFloat64(b)

    if aErr == nil && bErr == nil {
        return aFloat + bFloat, nil
    }

    // If not both numbers, try string concatenation
    aStr, aOk := a.(string)
    bStr, bOk := b.(string)
    if aOk && bOk {
        return aStr + bStr, nil
    }

    return nil, fmt.Errorf("cannot add %v and %v", a, b)
}

func (vm *VM) toFloat64(v interface{}) (float64, error) {
    switch v := v.(type) {
    case float64:
        return v, nil
    case int:
        return float64(v), nil
    case string:
        return strconv.ParseFloat(v, 64)
    default:
        return 0, fmt.Errorf("cannot convert %v to float64", v)
    }
}

func (vm *VM) GetMemory() map[string]interface{} {
    return vm.memory
}
