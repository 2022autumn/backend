package initialize

import (
	"os"
)

func InitMedia() {
	_, err := os.Stat("./media")
	if os.IsNotExist(err) {
		_ = os.MkdirAll("./media", 0755)
	}
	_, err = os.Stat("./media/headshot")
	if os.IsNotExist(err) {
		_ = os.MkdirAll("./media/headshot", 0755)
	}
}
