package mod

import "os"

func IsDebug() bool {
	return os.Getenv("LG_DEBUG") == "1"
}
