package progress

import (
	"fmt"
	"time"

	"github.com/cheggaaa/pb/v3/termutil"
)

const (
	refreshRate   = 200 * time.Millisecond
	prevLineReset = "\033[F"
)

type printer struct {
	b      *Bar
	stopCh chan struct{}
}

func newPrinter(b *Bar) *printer {
	return &printer{b, make(chan struct{})}
}

func (w *printer) Watch() {
	select {
	case <-w.stopCh:
		return
	default:
	}

	t := time.NewTicker(refreshRate)

	defer t.Stop()

	defer w.Print()

	fmt.Println()

	for {
		select {
		case <-w.stopCh:
			return
		case <-t.C:
			w.Print()
		}
	}
}

func (w *printer) StopWatch() {
	close(w.stopCh)
}

func (w *printer) Print() {
	width, _ := termutil.TerminalWidth()

	fmt.Print(prevLineReset)
	fmt.Println(w.b.Render(width))
}
