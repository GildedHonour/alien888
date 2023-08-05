// TODO rename to "common"
package helper

import (
	"log"
)

const (
	ContentTypeHeader = "Content-Type"
	JsonHeader        = "application/json"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
