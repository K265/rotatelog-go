package daily

import (
	"log"
	"testing"
)

func Test1(t *testing.T) {
	w := New("r:/test.", ".log", 30, 1024, nil)
	log.SetOutput(w)
	msg := "application running ..."
	for i := 0; i < 10000; i++ {
		log.Println(msg)
	}
}
