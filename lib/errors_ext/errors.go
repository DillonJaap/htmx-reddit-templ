package errors_ext

import "fmt"

// Join is like errors.Join, except it doesn't concat with a newline
func Join(err1, err2 error) error {
	return fmt.Errorf("%w%w", err1, err2)
}
