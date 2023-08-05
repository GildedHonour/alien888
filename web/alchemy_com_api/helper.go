package alchemy_com_api

import (
	"log"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
