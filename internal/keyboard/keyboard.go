package keyboard

import (
	"runtime"
	"time"

	"github.com/micmonay/keybd_event"
)

// SimulatePaste 模拟粘贴操作 (Ctrl+V / Cmd+V)
func SimulatePaste() error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return err
	}

	// 设置 V 键
	kb.SetKeys(keybd_event.VK_V)

	if runtime.GOOS == "darwin" {
		// Mac: Cmd+V
		kb.HasSuper(true)
	} else {
		// Windows/Linux: Ctrl+V
		kb.HasCTRL(true)
	}

	// Linux 需要等待 uinput 准备
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	return kb.Launching()
}

// SimulateEnter 模拟回车键
func SimulateEnter() error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return err
	}

	kb.SetKeys(keybd_event.VK_ENTER)

	// Linux 需要等待 uinput 准备
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	return kb.Launching()
}
