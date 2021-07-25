package main

import (
	"fmt"
	"go-blog/conf"
	"go-blog/router"
)

func main() {
	conf.DefaultInit()
	//csrf

	r := router.RoutersInit()
	fmt.Println("开始运行")
	_ = r.Run(":8081")

}
