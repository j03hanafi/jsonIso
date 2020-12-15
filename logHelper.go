package main

import (
	"fmt"
	"os"
	"time"
)

func writeLog(content string) {
	log, _ := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer log.Close()

	dt := time.Now()

	_, err := log.Write([]byte(dt.Format("2006-01-02 15:04:05 ") + content + "\n"))
	if err != nil {
		fmt.Println(err)
	}
}
