package commands

import (
	"fmt"
	"notgit/internal/commit"
	"strconv"
	"time"
)

func Log() error {
	head := commit.ParseHead()

	current := head
	for current != nil {
		timestamp := time.Unix(current.Time, 0)

		offset := current.Offset
		offsetHours, _ := strconv.Atoi(offset[:3])              // "-07" -> -7
		offsetMins, _ := strconv.Atoi(offset[0:1] + offset[3:]) // "-00" -> 0
		totalOffset := offsetHours*3600 + offsetMins*60

		loc := time.FixedZone("commit zone", totalOffset)
		adjusted := timestamp.In(loc).Format("2006-01-02")

		fmt.Printf("\033[33m%s\033[0m \033[34m%s\033[0m %s\n", current.Hash()[:7], adjusted, current.Message)
		if len(current.Parents) > 0 {
			current = current.Parents[0]
		} else {
			current = nil
		}
	}

	return nil
}
