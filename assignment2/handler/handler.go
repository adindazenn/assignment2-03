package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/adindazenn/assignment2-03/assignment2/model"
    "fmt"
)

func CreateOrder(c *gin.Context) {
    var order model.Order
    if err := c.ShouldBindJSON(&order); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db, _ := c.Get("db")
    database := db.(*gorm.DB)

    // Cek apakah item_code sudah ada di database
    var existingItem model.Item
    if err := database.Where("item_code = ?", order.Items[0].ItemCode).First(&existingItem).Error; err == nil {
        // Jika item_code sudah ada, kirim respons error
        errorMessage := "Item code tidak boleh sama"
        fmt.Println(errorMessage)
        c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
        return
    }

    // Buat order baru dalam database jika item_code belum ada
    database.Create(&order)

    fmt.Println("Order dibuat dengan id:", order.OrderID)
    c.JSON(http.StatusCreated, order)
}

func UpdateOrder(c *gin.Context) {
    // Mendapatkan order_id dari parameter URL
    orderID := c.Param("order_id")

    // Mencari order dalam database berdasarkan order_id
    var order model.Order
    db, _ := c.Get("db")
    database := db.(*gorm.DB)
    if err := database.Where("order_id = ?", orderID).First(&order).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Order tidak ditemukan"})
        return
    }

    // Parse data JSON dari permintaan
    var updateData model.UpdateOrderRequest
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Menggunakan GORM untuk pembaruan sebagian atribut pada order
    if err := database.Model(&order).Updates(map[string]interface{}{
        "ordered_at":    updateData.OrderedAt,
        "customer_name": updateData.CustomerName,
    }).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui order"})
        return
    }

    // Untuk memperbarui item, melalui updateData.Items
    for _, itemData := range updateData.Items {
        // Mencari item dalam database berdasarkan item_code
        var item model.Item
        if err := database.Where("item_code = ?", itemData.ItemCode).First(&item).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Item tidak ditemukan"})
            return
        }

        // Menggunakan GORM untuk memperbarui item
        if err := database.Model(&item).Updates(map[string]interface{}{
            "description": itemData.Description,
            "quantity":    itemData.Quantity,
        }).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui item"})
            return
        }
    }

    fmt.Println("Order dan Item berhasil diperbarui")
    c.JSON(http.StatusOK, gin.H{"message": "Order dan Item berhasil diperbarui"})
}

func GetData(c *gin.Context){
    db, _ := c.Get("db")
    database := db.(*gorm.DB)

    var orders []model.Order
    if err := database.Preload("Items").Find(&orders).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data order"})
        return
    }

    c.JSON(http.StatusOK, orders)
}

func DeleteOrder(c *gin.Context) {
    // Dapatkan ID order dari parameter.
    orderID := c.Param("id")

    // Cari order berdasarkan ID.
    db, _ := c.Get("db")
    database := db.(*gorm.DB)

    var order model.Order
    if err := database.Preload("Items").First(&order, orderID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Order tidak ditemukan"})
        return
    }

    // Hapus item terkait order.
    if err := database.Delete(&order.Items).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Hapus order dari database.
    if err := database.Delete(&order).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Kirim respons sukses.
    fmt.Println("Order dan Item berhasil dihapus")
    c.JSON(http.StatusOK, gin.H{"message": "Order dan Item berhasil dihapus"})
}
