package main

// This part of the code that detect the keystrokes was copied from the playground website. https://play.golang.org/p/osbinz7DmN
import (
	"bytes"
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/Triballian/tanktool/clipboard"
)

const (
	ModAlt = 1 << iota
	ModCtrl
	ModShift
	ModWin
)

type Hotkey struct {
	Id        int // Unique id
	Modifiers int // Mask of modifiers
	KeyCode   int // Key code, e.g. 'A'
}

// String returns a human-friendly display name of the hotkey
// such as "Hotkey[Id: 1, Alt+Ctrl+O]"
func (h *Hotkey) String() string {
	mod := &bytes.Buffer{}
	if h.Modifiers&ModAlt != 0 {
		mod.WriteString("Alt+")
	}
	if h.Modifiers&ModCtrl != 0 {
		mod.WriteString("Ctrl+")
	}
	if h.Modifiers&ModShift != 0 {
		mod.WriteString("Shift+")
	}
	if h.Modifiers&ModWin != 0 {
		mod.WriteString("Win+")
	}
	return fmt.Sprintf("Hotkey[Id: %d, %s%c]", h.Id, mod, h.KeyCode)
}

func main() {
	clipboard.WriteAll("Juggernaut build; game_stats_build 565656565677777444445675674422222")

	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

	reghotkey := user32.MustFindProc("RegisterHotKey")

	// Hotkeys to listen to:
	keys := map[int16]*Hotkey{
		1: &Hotkey{1, ModCtrl, 'O'},          // ALT+CTRL+O
		2: &Hotkey{2, ModCtrl, 'J'},          // ALT+SHIFT+M
		3: &Hotkey{3, ModCtrl, 'B'},          // ALT+SHIFT+b
		4: &Hotkey{4, ModAlt + ModCtrl, 'X'}, // ALT+CTRL+X

		// add ctgr + b
	}

	// Register hotkeys:
	for _, v := range keys {
		r1, _, err := reghotkey.Call(
			0, uintptr(v.Id), uintptr(v.Modifiers), uintptr(v.KeyCode))
		if r1 == 1 {
			fmt.Println("Registered", v)
		} else {
			fmt.Println("Failed to register", v, ", error:", err)
		}
	}

	peekmsg := user32.MustFindProc("PeekMessageW")

	for {
		var msg = &MSG{}
		peekmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

		// Registered id is in the WPARAM field:
		if id := msg.WPARAM; id != 0 {
			fmt.Println("Hotkey pressed:", keys[id])
			if id == 2 { // CTRL+O = Exit
				fmt.Println("Overseer Configuration Ready")
				clipboard.WriteAll("Overseer 3; game_stats_build 568568568568568564747474474477756")

			}
			if id == 2 { // CTRL+J = Exit
				fmt.Println("Juggernaut Configuration Ready")
				clipboard.WriteAll("Juggernaut build; game_stats_build 565656565677777444445675674422222")

			}
			if id == 3 { // CTRL+B = Exit
				fmt.Println("Boster Configuration Ready")
				clipboard.WriteAll("Flanker 3;  game_stats_build 568568568568568567647477774754444")

			}
			if id == 4 { // CTRL+ALT+X = Exit
				fmt.Println("CTRL+ALT+X pressed, goodbye...")
				return
			}
		}

		time.Sleep(time.Millisecond * 5)
	}
}

type MSG struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}
