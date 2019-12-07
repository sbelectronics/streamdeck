package windows

import (
	"github.com/JamesHovious/w32"
)

// Get the title of the topmost window
func GetTopmostWindowTitle() (string, error) {
	fgw, err := w32.GetForegroundWindow()
	if err != nil {
		return "", err
	}

	windowTitle := make([]uint16, 1024)
	_, err = w32.GetWindowTextW(fgw, &windowTitle[0], 1024)
	if err != nil {
		return "", err
	}

	return w32.UTF16PtrToString(&windowTitle[0]), err
}
