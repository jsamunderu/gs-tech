package main

import (
	"fmt"
	"math/rand"
	"question3/questions"
	"sync"
)

func main() {
	reqMan := questions.NewRequestManager()
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			switch rand.Intn(3) {
			case 0:
				reqMan.Set(&questions.Request{Value: 0})
				fmt.Println("Set")
			case 1:
				reqMan.Update(&questions.Request{Value: 1})
				fmt.Println("Update")
			case 2:
				reqMan.Delete()
				fmt.Println("Delete")
			}

		}()
	}
	wg.Wait()
}
