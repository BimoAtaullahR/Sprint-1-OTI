package controllers

import (
	"net/http"

	"school-vault/config"
	"school-vault/models"

	"github.com/gin-gonic/gin"
)

// Mencatat barang masuk/keluar dan update stok otomatis
func CreateTransaction(c *gin.Context) {
	var trx models.Transaction
	if err := c.ShouldBindJSON(&trx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	// Mulai Database Transaction agar aman (jika satu gagal, semua dibatalkan)
	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi"})
		return
	}

	// 1. Simpan riwayat transaksi
	err = tx.QueryRow(
		`INSERT INTO transactions (item_id, type, quantity, notes) VALUES ($1, $2, $3, $4) RETURNING id`,
		trx.ItemID, trx.Type, trx.Quantity, trx.Notes,
	).Scan(&trx.ID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan riwayat transaksi"})
		return
	}

	// 2. Update stok di tabel items berdasarkan tipe (IN/OUT)
	var updateQuery string
	if trx.Type == "IN" {
		updateQuery = `UPDATE items SET stock = stock + $1 WHERE id = $2`
	} else if trx.Type == "OUT" {
		updateQuery = `UPDATE items SET stock = stock - $1 WHERE id = $2 AND stock >= $1`
	} else {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipe transaksi harus IN atau OUT"})
		return
	}

	res, err := tx.Exec(updateQuery, trx.Quantity, trx.ItemID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update stok"})
		return
	}

	// Cek apakah stok cukup untuk transaksi OUT
	affected, _ := res.RowsAffected()
	if affected == 0 && trx.Type == "OUT" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stok tidak mencukupi atau barang tidak ditemukan"})
		return
	}

	tx.Commit() // Simpan semua perubahan permanen
	c.JSON(http.StatusCreated, trx)
}
