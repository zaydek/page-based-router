package main

import (
	"fmt"
	"log"
)

const CONF_PATH = "go/conf.json"

func main() {
	route, err := newPageBasedRoute("test.mdx")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", route)

	// conf, err := config.ReadConfigurationFile(CONF_PATH)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _, err = os.Stat(CONF_PATH)
	// if os.IsNotExist(err) {
	// 	// TODO: Should write the initialized configuration defaults to disk; change
	// 	// message to `initialized from recommended defaults and created ...`.
	// 	log.Print("no such configuration file; initialized from recommended defaults")
	// }
	// fmt.Printf("%+v\n", conf)
}
