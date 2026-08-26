package app

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestSingleInstanceMutex(t *testing.T) {
	name := fmt.Sprintf(`Local\Stendoclip.Test.%d.%d`, os.Getpid(), time.Now().UnixNano())
	first, already, err := acquireInstance(name)
	if err != nil || already {
		t.Fatalf("first acquire: already=%v err=%v", already, err)
	}
	defer first.Close()

	second, already, err := acquireInstance(name)
	if err != nil || !already || second != nil {
		t.Fatalf("second acquire: instance=%v already=%v err=%v", second, already, err)
	}
}
