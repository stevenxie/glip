package portal

import (
	"fmt"
	"io"
)

// Read allows for the reading of data from a command's standard output.
func (p *Portal) Read(dst []byte) (n int, err error) {
	// Prepare a fresh "Cmd" from "blueprint".
	p.Prepare()

	// Open an pipe to Stdout.
	out, err := p.StdoutPipe()
	if err != nil {
		return 0, stdoutPipeErr(err)
	}

	// Start "Cmd"; read data to destination.
	if err = p.Start(); err != nil {
		return 0, startErr(err)
	}
	if n, err = out.Read(dst); err != nil {
		return n, fmt.Errorf("portal: failed to read from Stdout: %v", err)
	}

	// Wait for "Cmd" to complete.
	if err = p.Wait(); err != nil {
		return n, waitErr(err)
	}
	return n, nil
}

// WriteTo allows for the piping of data from a command's standard output into
// an io.Writer.
func (p *Portal) WriteTo(w io.Writer) (n int64, err error) {
	// Prepare a fresh "Cmd" from "blueprint".
	p.Prepare()

	// Open a pipe to Stdout.
	out, err := p.StdoutPipe()
	if err != nil {
		return 0, stdoutPipeErr(err)
	}

	// Start "Cmd"; copy data from program Stdout to the provided io.Writer.
	if err = p.Start(); err != nil {
		return 0, startErr(err)
	}
	if n, err = io.Copy(w, out); err != nil {
		return n, fmt.Errorf("portal: could not to copy from Stdout: %v", err)
	}

	// Wait for "Cmd" to complete.
	if err = p.Wait(); err != nil {
		return n, waitErr(err)
	}
	return n, nil
}
