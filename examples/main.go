package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func csvreader() (chan []string, <-chan bool) {
	queue := make(chan []string)
	done := make(chan bool)

	f, _ := os.Open("data.csv")
	reader := csv.NewReader(f)

	go func() {
		// defer close(queue)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				log.Println("Reached EOF")
				done <- true
				return
			}
			if err != nil {
				log.Fatal(err)
			}

			queue <- record
		}
	}()

	return queue, done
}

func filterQueue(inQueue chan []string, offset int, searchString string) <-chan []string {
	outQueue := make(chan []string)

	go func() {
		// defer close(outQueue)

		for {
			select {
			case row := <-inQueue:
				if row[offset] == searchString {
					//send it to the filtered queue
					outQueue <- row
				} else {
					//just publish back onto the main queue
					inQueue <- row
				}
			}
		}
	}()

	return outQueue
}

func fanout(inQueue <-chan []string, outQueues ...chan<- []string) {
	go func() {
		for {
			select {
			case row := <-inQueue:
				for _, outQueue := range outQueues {
					outQueue <- row
				}
			}

		}
	}()
}

func notbasic() {
	mainQueue, doneQueue := csvreader()

	for {
		select {
		case row := <-mainQueue:
			log.Println(row)
		}
	}

	// LGfilteredQueue := filterQueue(mainQueue, 6, "LG")

	// for {
	// 	select {
	// 	case row := <-LGfilteredQueue:
	// 		log.Println(row)
	// 	}
	// }

}

func main() {
	notbasic()
}
