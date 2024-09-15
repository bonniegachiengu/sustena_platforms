package interpreter

import (
    "fmt"
    "strings"
)

type Interpreter struct {
    variables map[string]interface{}
}

func NewInterpreter() *Interpreter {
    return &Interpreter{
        variables: make(map[string]interface{}),
    }
}

func (i *Interpreter) Interpret(code string) error {
    lines := strings.Split(code, "\n")
    for _, line := range lines {
        if err := i.interpretLine(strings.TrimSpace(line)); err != nil {
            return err
        }
    }
    return nil
}

func (i *Interpreter) interpretLine(line string) error {
    // Basic variable assignment
    if strings.Contains(line, "=") {
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            return fmt.Errorf("invalid assignment: %s", line)
        }
        varName := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])
        i.variables[varName] = value
    }
    // Add more interpretation logic here

    return nil
}

func (i *Interpreter) GetVariable(name string) (interface{}, bool) {
    value, exists := i.variables[name]
    return value, exists
}

// Add more methods as needed
