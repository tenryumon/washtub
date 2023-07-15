package log

import (
	"bytes"
	"fmt"
)

// The type "Buffer" contains a bytes buffer and an error.
// @property buff - `buff` is a field of type `bytes.Buffer` which is a buffer that can be used to read
// and write bytes. It provides methods for appending, reading, and writing bytes. In this case, it is
// embedded within the `Buffer` struct, which means that the methods of `bytes
// @property {error} err - The `err` property is a variable of type `error` that is used to store any
// error that occurs during the execution of a function or method that uses the `Buffer` struct. It
// allows the caller of the function or method to check if an error occurred and handle it
// appropriately.
type Buffer struct {
	buff bytes.Buffer
	err  error
}

// The function creates a new buffer with a given message.
func NewBuffer(msg string) *Buffer {
	var buffer Buffer
	buffer.buff.WriteString(msg)
	return &buffer
}

// The `PrintBuffer` method is a function that takes a string `msg` as input and appends it to the
// `Buffer` struct's `buff` field. It then prints the contents of the `buff` field using the `Println`
// function from the `log` package.
func (b *Buffer) PrintBuffer(msg string) {
	b.WriteString(msg)
	Println(b.buff.String())
}

// The `WriteString` method is a function that takes a string `msg` as input and appends it to the
// `Buffer` struct's `buff` field with a space character added before the message. It uses the
// `WriteString` method of the embedded `bytes.Buffer` field to append the message to the buffer.
func (b *Buffer) WriteString(msg string) {
	b.buff.WriteString(" " + msg)
}

// The `Err` method is a function that sets the `err` property of the `Buffer` struct to the error
// passed as an argument. It takes an `error` value as input and assigns it to the `err` property of
// the `Buffer` struct. This method is used to set the error property of the buffer when an error
// occurs during the execution of a function or method that uses the `Buffer` struct.
func (b *Buffer) Err(err error) {
	b.err = err
}

// The `IsErr` method is a function that checks if an error has occurred during the execution of a
// function or method that uses the `Buffer` struct. It returns a boolean value of `true` if an error
// has occurred and the `err` property of the `Buffer` struct is not `nil`, and `false` otherwise. This
// method can be used by the caller of the function or method to determine if an error occurred and
// handle it appropriately.
func (b *Buffer) IsErr() bool {
	return b.err != nil
}

// The `Println` method is a function that prints the contents of the `Buffer` struct's `buff` field
// using the `Println` function from the `log` package. If an error has occurred during the execution
// of a function or method that uses the `Buffer` struct, it appends the error message to the buffer
// with the string "[ERROR]" and then prints the contents of the buffer using the `Println` function.
// This method allows the caller to print the contents of the buffer and handle any errors that
// occurred during its creation.
func (b *Buffer) Println() {
	if b.IsErr() {
		b.PrintBuffer(fmt.Sprintf("%s [ERROR]", b.err))
		return
	}

	Println(b.buff.String())
}

// The `WriteStringf` method is a function that takes a format string and a variable number of
// arguments as input. It formats the string using the `fmt.Sprintf` function and then calls the
// `WriteString` method of the `Buffer` struct to append the formatted string to the buffer with a
// space character added before the message. This method allows the caller to format a string and
// append it to the buffer in a single step.
func (b *Buffer) WriteStringf(format string, args ...interface{}) {
	b.WriteString(fmt.Sprintf(format, args...))
}
