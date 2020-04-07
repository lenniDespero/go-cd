package main

import (
	"fmt"
	"reflect"

	"github.com/lenniDespero/go-cd/internal/models"
	"github.com/mitchellh/mapstructure"
)

func main() {
	var C models.Config
	models.ReadConfig(&C)
	fmt.Printf("%v\n", C)
	fmt.Printf("%v\n", C.Targets["devel"])
	var target models.Target
	target = C.Targets["devel"]
	fmt.Printf("%v\n", target.Pipe)
	for _, pipe := range target.Pipe {
		inter := models.PipeNames[pipe.Type]
		fmt.Printf("%v  - %s \n", inter, pipe.Type)

		for _, args := range pipe.Args {
			//fmt.Printf("%v\n", args)
			err := mapstructure.Decode(args, inter)
			if err != nil {
				panic(err)
			}
			reflect.TypeOf(inter)
			inter.checkArgsConfig()
			fmt.Printf("%v\n", reflect.ValueOf(inter).Interface().(models.ArgsInterface))

		}
	}
	//pipes := target.Pipe
	//
	//for pipe := range C.Targets["devel"].Pipe {
}
