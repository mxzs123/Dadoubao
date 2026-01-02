package clipboard

import (
	"golang.design/x/clipboard"
)

// Init 初始化剪贴板
func Init() error {
	return clipboard.Init()
}

// WriteText 写入文字到剪贴板
func WriteText(text string) {
	clipboard.Write(clipboard.FmtText, []byte(text))
}

