//go:build windows

package console

// setNonblock is a no-op on Windows
func setNonblock(_ int, _ bool) error {
	return nil
}
