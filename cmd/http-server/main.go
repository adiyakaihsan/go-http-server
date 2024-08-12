package main

import (
	"github.com/adiyakaihsan/go-http-server/pkg/app"
	_ "github.com/lib/pq"
)

// 1. Structure: https://github.com/golang-standards/project-layout
// - cmd
// - pkg
// 2.

func main() {
	app.Run()
}
