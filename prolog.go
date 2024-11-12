package main

import (
	"github.com/ichiban/prolog"
)

func (b *App) RunProlog(program string, input string) string {
	//
	// Create the output string.
	//
	b.err = ""
	result := ""

	//
	// Create an interpreter instance.
	//
	p := prolog.New(nil, nil)

	//
	// Add the input as a clause for the program.
	//
	full := `intext("` + input + `").` + "\n" + program

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
		sols, err := p.Query(`main(X).`)
		if err != nil {
			b.err = err.Error()
		} else {
			// Iterates over solutions.
			for sols.Next() {
				// Prepare a struct with fields which name corresponds with a variable in the query.
				var s struct {
					X string
				}
				if err := sols.Scan(&s); err != nil {
					b.err = err.Error()
				}
				result += s.X + "\n"
			}
		}
		defer sols.Close()

		// Check if an error occurred while querying.
		if err := sols.Err(); err != nil {
			b.err = err.Error()
		}
	}
	return result
}
