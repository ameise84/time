package example

import (
	"fmt"
	"github.com/ameise84/time"
	"log"
	"sync"
	"sync/atomic"
	"testing"
)

type handler struct {
}

func (h *handler) OnTimer(time.Timer, int, time.Time) {
	log.Printf("t3 triggger by handler[%s]\n", time.Now().Format(time.LayoutSec))
}

func TestTimer(t *testing.T) {
	fmt.Println("---TestTimer---")
	var a int64
	atomic.AddInt64(&a, 1)
	atomic.AddInt64(&a, 1)
	wrap := time.NewTimerHandlerWrap(2)
	_ = time.NewTick(wrap.GetTimerHandler(), 1, 10*time.Second, time.InfiniteTimes)
	_ = time.NewTimer(wrap.GetTimerHandler(), 2, 1*time.Second)

	h := &handler{}
	_ = time.NewTick(h, 3, 1*time.Second, time.InfiniteTimes)

	wg := sync.WaitGroup{}
	maxC := 100
	v := 0
	wg.Add(maxC)
	go func() {
		log.Println("start")
	loop:
		for {
			select {
			case trigger := <-wrap.C():
				if trigger.Timer().Context().(int) == 1 {
					log.Printf("t1 triggger on time[%s]\n", trigger.Now().Format(time.LayoutSec))
					_ = time.FastForward(5 * time.Second)
				} else {
					log.Printf("t2 triggger on time[%s]->n[%d]\n", trigger.Now().Format(time.LayoutSec), trigger.DoTimes())
					_ = time.FastForward(3 * time.Second)
				}
				wg.Add(trigger.DoTimes() * -1)
				v += trigger.DoTimes()
				if v >= maxC {
					break loop
				}
			}
		}
	}()
	wg.Wait()
}
