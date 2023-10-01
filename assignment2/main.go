package main

import (
    "github.com/adindazenn/assignment2-03/assignment2/route"
    "gorm.io/gorm"
)

var db *gorm.DB

func main() {
    r := route.SetupRouter()
    r.Run(":8080")
}
