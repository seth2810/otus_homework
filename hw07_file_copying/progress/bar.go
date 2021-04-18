package progress

import (
	"fmt"
	"io"
	"math"
	"strings"
	"sync/atomic"
)

type Bar struct {
	r io.Reader
	p *printer
	n int64
	t int64
}

func NewBar(r io.Reader, limit int64) *Bar {
	return &Bar{r: r, t: limit}
}

func (b *Bar) Read(p []byte) (n int, err error) {
	n, err = b.r.Read(p)
	atomic.AddInt64(&b.n, int64(n))
	return
}

func (b *Bar) Close() error {
	defer b.p.StopWatch()
	atomic.StoreInt64(&b.n, b.t)
	return nil
}

func (b *Bar) Render(width int) string {
	value := float64(b.n) / float64(b.t)
	indicator := fmt.Sprintf("%7.2f%%", value*100)
	maxWidth := math.Max(float64(width-len(indicator)), 0)
	barWidth := math.Min(maxWidth*value, maxWidth)
	paddingWidth := maxWidth - barWidth

	return fmt.Sprintf("%s%s%s", strings.Repeat("|", int(barWidth)), strings.Repeat(" ", int(paddingWidth)), indicator)
}

func (b *Bar) Watch() {
	b.p = newPrinter(b)

	go b.p.Watch()
}
