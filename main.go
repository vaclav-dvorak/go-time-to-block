package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const inputDate = "22.02.2022 22:22+GMT"

type resp struct {
	Top    int
	Bottom int
	Middle float64
}

func main() {
	date, err := time.Parse("02.01.2006 15:04+MST", inputDate)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("looking for time %d", date.Unix())
	r, err := getBlockData(date)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("👆 After block: %d\n👉 Exact blocktime: %.2f\n👇 Before block: %d\n", r.Top, r.Middle, r.Bottom)
}
