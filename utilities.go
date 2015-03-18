package main

import (
	"log"
)

func logerr(err error) {
	if err != nil {
		log.Print(err)
	}
}
