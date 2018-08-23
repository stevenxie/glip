package portal

import (
	"fmt"
	"io"
)

// Write allows for the writing of data into a command's standard input.
func (p *Portal) Write(src []byte) (n int, err error) {
	// Prepare a fresh "Cmd" from "blueprint".
	p.Prepare()

	// Open a pipe to program Stdin.
	in, err := p.StdinPipe()
	if err != nil {
		return 0, stdinPipeErr(err)
	}

	// Start program; begin writing to it's Stdin from "src".
	if err = p.Start(); err != nil {
		return 0, startErr(err)
	}
	if n, err = in.Write(src); err != nil {
		return n, fmt.Errorf("portal: could not write to Stdin: %v", err)
	}

	// Close Stdin to signal to the program that we are done with it.
	if err = in.Close(); err != nil {
		return n, closeStdinErr(err)
	}

	// Wait for the program to exit.
	if err = p.Wait(); err != nil {
		return n, waitErr(err)
	}
	return n, err
}

// ReadFrom allows for the piping of data from a io.Writer into a command's
// standard output.
func (p *Portal) ReadFrom(r io.Reader) (n int64, err error) {
	// Prepare a fresh "Cmd" from "blueprint".
	p.Prepare()

	// Open a pipe to Stdin.
	in, err := p.StdinPipe()
	if err != nil {
		return 0, stdinPipeErr(err)
	}

	// Start the program; read from the provided io.Reader to program Stdin.
	if err = p.Start(); err != nil {
		return 0, startErr(err)
	}
	if n, err = io.Copy(in, r); err != nil {
		return n, fmt.Errorf("portal: failed to write to Stdin: %v", err)
	}

	// Close program Stdin to signal that we are done with it.
	if err = in.Close(); err != nil {
		return n, closeStdinErr(err)
	}

	// Wait for program to exit.
	if err = p.Wait(); err != nil {
		return n, waitErr(err)
	}
	return n, err
}
