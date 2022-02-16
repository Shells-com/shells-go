package spicefyne

import (
	"sync"

	"github.com/go-gl/glfw/v3.3/glfw"
)

import "C"

var (
	fynePhysicalScancodeMap     map[int][]byte
	fynePhysicalScancodeMapInit sync.Once
)

func scancodeLookup(sc int) ([]byte, bool) {
	fynePhysicalScancodeMapInit.Do(func() {
		// fill fyneUsbScancodeMap via glfw because fyne exposes scancode rather than glfw abstraction
		fynePhysicalScancodeMap = make(map[int][]byte)

		for k, v := range glfwScancodeMap {
			scancode := glfw.GetKeyScancode(k)
			if scancode != -1 {
				fynePhysicalScancodeMap[scancode] = v
			}
		}
	})

	v, ok := fynePhysicalScancodeMap[sc]
	return v, ok
}

// glfw → XT scancode table
// see: https://www.win.tue.nl/~aeb/linux/kbd/scancodes-1.html
var glfwScancodeMap = map[glfw.Key][]byte{
	glfw.KeySpace:      []byte{0x39},
	glfw.KeyApostrophe: []byte{0x28}, // Quote
	glfw.KeyComma:      []byte{0x33},
	glfw.KeyMinus:      []byte{0x0c},
	glfw.KeyPeriod:     []byte{0x34},
	glfw.KeySlash:      []byte{0x35},

	glfw.Key0: []byte{0x0b},
	glfw.Key1: []byte{0x02},
	glfw.Key2: []byte{0x03},
	glfw.Key3: []byte{0x04},
	glfw.Key4: []byte{0x05},
	glfw.Key5: []byte{0x06},
	glfw.Key6: []byte{0x07},
	glfw.Key7: []byte{0x08},
	glfw.Key8: []byte{0x09},
	glfw.Key9: []byte{0x0a},

	glfw.KeySemicolon: []byte{0x27},
	glfw.KeyEqual:     []byte{0x0d},

	glfw.KeyA: []byte{0x1e},
	glfw.KeyB: []byte{0x30},
	glfw.KeyC: []byte{0x2e},
	glfw.KeyD: []byte{0x20},
	glfw.KeyE: []byte{0x12},
	glfw.KeyF: []byte{0x21},
	glfw.KeyG: []byte{0x22},
	glfw.KeyH: []byte{0x23},
	glfw.KeyI: []byte{0x17},
	glfw.KeyJ: []byte{0x24},
	glfw.KeyK: []byte{0x25},
	glfw.KeyL: []byte{0x26},
	glfw.KeyM: []byte{0x32},
	glfw.KeyN: []byte{0x31},
	glfw.KeyO: []byte{0x18},
	glfw.KeyP: []byte{0x19},
	glfw.KeyQ: []byte{0x10},
	glfw.KeyR: []byte{0x13},
	glfw.KeyS: []byte{0x1f},
	glfw.KeyT: []byte{0x14},
	glfw.KeyU: []byte{0x16},
	glfw.KeyV: []byte{0x2f},
	glfw.KeyW: []byte{0x11},
	glfw.KeyX: []byte{0x2d},
	glfw.KeyY: []byte{0x15},
	glfw.KeyZ: []byte{0x2c},

	glfw.KeyLeftBracket:  []byte{0x1a},
	glfw.KeyBackslash:    []byte{0x2b},
	glfw.KeyRightBracket: []byte{0x1b},
	glfw.KeyGraveAccent:  []byte{0x29}, // Backquote

	glfw.KeyWorld1: nil,
	glfw.KeyWorld2: nil,

	glfw.KeyEscape:    []byte{0x01},
	glfw.KeyEnter:     []byte{0x1c},
	glfw.KeyTab:       []byte{0x0f},
	glfw.KeyBackspace: []byte{0x0e},
	glfw.KeyInsert:    []byte{0xe0, 0x52},
	glfw.KeyDelete:    []byte{0xe0, 0x53},
	glfw.KeyRight:     []byte{0xe0, 0x4d},
	glfw.KeyLeft:      []byte{0xe0, 0x4b},
	glfw.KeyDown:      []byte{0xe0, 0x50},
	glfw.KeyUp:        []byte{0xe0, 0x48},
	glfw.KeyPageUp:    []byte{0xe0, 0x49},
	glfw.KeyPageDown:  []byte{0xe0, 0x51},
	glfw.KeyHome:      []byte{0xe0, 0x47},
	glfw.KeyEnd:       []byte{0xe0, 0x4f},

	glfw.KeyCapsLock:    []byte{0x39},
	glfw.KeyScrollLock:  []byte{0x46},
	glfw.KeyNumLock:     []byte{0x45},
	glfw.KeyPrintScreen: []byte{0xe0, 0x2a, 0xe0, 0x37},
	glfw.KeyPause:       nil, // []byte{0xe1, 0x1d, 0x45, 0xe1, 0x9d, 0xc5},

	glfw.KeyF1:  []byte{0x3b},
	glfw.KeyF2:  []byte{0x3c},
	glfw.KeyF3:  []byte{0x3d},
	glfw.KeyF4:  []byte{0x3e},
	glfw.KeyF5:  []byte{0x3f},
	glfw.KeyF6:  []byte{0x40},
	glfw.KeyF7:  []byte{0x41},
	glfw.KeyF8:  []byte{0x42},
	glfw.KeyF9:  []byte{0x43},
	glfw.KeyF10: []byte{0x44},
	glfw.KeyF11: []byte{0x57},
	glfw.KeyF12: []byte{0x58},
	// up to F25

	glfw.KeyKP0: []byte{0x52},
	glfw.KeyKP1: []byte{0x4f},
	glfw.KeyKP2: []byte{0x50},
	glfw.KeyKP3: []byte{0x51},
	glfw.KeyKP4: []byte{0x4b},
	glfw.KeyKP5: []byte{0x4c},
	glfw.KeyKP6: []byte{0x4d},
	glfw.KeyKP7: []byte{0x47},
	glfw.KeyKP8: []byte{0x48},
	glfw.KeyKP9: []byte{0x49},

	glfw.KeyKPDecimal:  []byte{0x53},
	glfw.KeyKPDivide:   []byte{0xe0, 0x35},
	glfw.KeyKPMultiply: []byte{0x37},
	glfw.KeyKPSubtract: []byte{0x4a},
	glfw.KeyKPAdd:      []byte{0x4e},
	glfw.KeyKPEnter:    []byte{0xe0, 0x1c},
	glfw.KeyKPEqual:    nil, // ?? []byte{},

	glfw.KeyLeftShift:    []byte{0x2a},
	glfw.KeyLeftControl:  []byte{0x1d},
	glfw.KeyLeftAlt:      []byte{0x38},
	glfw.KeyLeftSuper:    []byte{0xe0, 0x5b},
	glfw.KeyRightShift:   []byte{0x36},
	glfw.KeyRightControl: []byte{0xe0, 0x1d},
	glfw.KeyRightAlt:     []byte{0xe0, 0x38},
	glfw.KeyRightSuper:   []byte{0xe0, 0x5c},
	glfw.KeyMenu:         []byte{0xe0, 0x5d},
}
