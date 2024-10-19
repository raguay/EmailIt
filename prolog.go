package main

import (
	"github.com/ichiban/prolog"
)

func (b *App) runProlog(program string, input string) string {
	//
	// Create the output string.
	//
	result := ""

	//
	// Create an interpreter instance.
	//
	p := prolog.New(nil, nil)

	//
	// Add the input as a clause for the program.
	//
	full := `Input(` + input + `)` + input

	//
	// Load the program.
	//
	if err := p.Exec(full); err != nil {
		//
		// There was an error with the program.
		//
		b.err = err.Error()
	} else {
		//
		// The program was loaded fine. Get the results.
		//
	}
	return result
}
