package controllers

import (
	"net/http"

	"school-vault/config"
	"school-vault/models"

	"github.com/gin-gonic/gin"
)

// Menambahkan barang baru
func CreateItem(c *gin.Context) {
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	query := `INSERT INTO items (name, category, stock, minimum_stock, unit) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := config.DB.QueryRow(query, item.Name, item.Category, item.Stock, item.MinimumStock, item.Unit).Scan(&item.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan barang"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// Mengambil semua barang
func GetItems(c *gin.Context) {
	rows, err := config.DB.Query(`SELECT id, name, category, stock, minimum_stock, unit FROM items`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var i models.Item
		if err := rows.Scan(&i.ID, &i.Name, &i.Category, &i.Stock, &i.MinimumStock, &i.Unit); err == nil {
			items = append(items, i)
		}
	}
	c.JSON(http.StatusOK, items)
}

// Mendapatkan daftar barang dengan stok kritis (Fitur Utama)
func GetCriticalStocks(c *gin.Context) {
	rows, err := config.DB.Query(`SELECT id, name, category, stock, minimum_stock, unit FROM items WHERE stock <= minimum_stock`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kritis"})
		return
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var i models.Item
		if err := rows.Scan(&i.ID, &i.Name, &i.Category, &i.Stock, &i.MinimumStock, &i.Unit); err == nil {
			items = append(items, i)
		}
	}
	c.JSON(http.StatusOK, items)
}

// Mengubah data barang (Update)
func UpdateItem(c *gin.Context) {
	id := c.Param("id") // Mengambil ID dari URL
	var item models.Item

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	// Kita biasanya tidak meng-update 'stock' secara manual, melainkan lewat transaksi.
	// Jadi kita hanya update profil barangnya saja.
	query := `UPDATE items SET name = $1, category = $2, minimum_stock = $3, unit = $4 WHERE id = $5`
	res, err := config.DB.Exec(query, item.Name, item.Category, item.MinimumStock, item.Unit, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengubah data barang"})
		return
	}

	// Cek apakah barang dengan ID tersebut ada
	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data barang berhasil diubah"})
}

// Menghapus barang (Delete)
func DeleteItem(c *gin.Context) {
	id := c.Param("id") // Mengambil ID dari URL

	query := `DELETE FROM items WHERE id = $1`
	res, err := config.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus barang"})
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Barang berhasil dihapus"})
}
