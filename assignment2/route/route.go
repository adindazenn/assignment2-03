package route

import (
    "github.com/gin-gonic/gin"
    "github.com/adindazenn/assignment2-03/assignment2/db"
    "github.com/adindazenn/assignment2-03/assignment2/handler"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // Inisialisasi koneksi database PostgreSQL
    db, err := db.InitDB()
    if err != nil {
        panic("Failed to connect to database")
    }

    // Mengirimkan koneksi database sebagai middleware ke handler
    r.Use(func(c *gin.Context) {
        c.Set("db", db)
        c.Next()
    })

    // Definisi route
    r.POST("/order", handler.CreateOrder)
    r.PUT("/order/:order_id", handler.UpdateOrder)
    r.PATCH("/order/:order_id", handler.UpdateOrder)
    r.GET("/all", handler.GetData)
    r.DELETE("/delete/:order_id", handler.DeleteOrder)

    return r
}

