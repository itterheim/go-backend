package cli

import (
	"fmt"
	"os"
)

func GetArg(key string) (string, error) {
	for i, value := range os.Args {
		if value == "--"+key {
			if (len(os.Args) - 1) == i {
				return "", nil
			}
			return os.Args[i+1], nil
		}
	}

	return "", fmt.Errorf("argument %s not found", key)
}