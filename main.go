package main

import (
	"log"
	_ "myproject/routers"
	"os"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory: ", err)
	}
	log.Println("Current working directory:", dir)

	// Load the app.conf file
	err = beego.LoadAppConfig("ini", "./conf/app.conf")
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	// Run the Beego application
	beego.Run()
}
