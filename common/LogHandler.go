package common

import (
	"log"
	"reflect"
)


func LogDebugData(data string, debug bool)  {
	if debug{
		log.Println("Debug:" + data)
	}

}
//Logger
func LogData(data string)  {
	log.Println("Debug:" + data)
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}