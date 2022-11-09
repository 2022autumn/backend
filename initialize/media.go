package initialize

import (
	"os"
)

func InitMedia() {
	_, err := os.Stat("./media")
	if os.IsNotExist(err) {
		_ = os.MkdirAll("./media", 0755)
	}
	return
}
