package example

import (
	"fmt"
	"github.com/ameise84/time"
	"log"
	"testing"
)

func TestTime(t *testing.T) {
	fmt.Println("---TestTime---")
	now1 := time.Now()
	log.Printf("now1:%s\n ", now1.Format(time.Layout))
	now2 := now1.Add(20 * time.Second)
	log.Printf("now2:%s\n", now2.Format(time.Layout))
	if err := time.FastForwardToLocal(now2.Format(time.Layout)); err != nil {
		log.Println(err)
	}
	now3 := time.Now()
	log.Printf("now3:%s\n", now3.Format(time.Layout))
	now4 := time.Unix(now3.Unix()+10, 0)
	log.Printf("now4:%s\n", now4.Format(time.Layout))
	log.Printf("now4:%s\n", now4.UTC().Format(time.Layout))
	if err := time.FastForwardToUTC(now4.UTC().Format(time.Layout)); err != nil {
		log.Println(err)
	}
	now5 := time.Now()
	log.Printf("now5:%s\n", now5.Format(time.Layout))
}
