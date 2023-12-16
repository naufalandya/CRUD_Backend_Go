package handlers

import (
	"database/sql"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type newTransaction struct {
	Transaction_id           uint64 `json:"id" binding:"required"`
	Transaction_date         string `json:"date" binding:"required"`
	Transaction_ticket_count uint64 `json:"ticket_count" binding:"required"`
	Transaction_visitor_id   uint64 `json:"visitor_id" binding:"required"`
	Transaction_zoo_id       uint64 `json:"zoo_id" binding:"required"`
}

func TransactionrowToStruct(rows *sql.Rows, dest interface{}) error {
	destv := reflect.ValueOf(dest).Elem()

	args := make([]interface{}, destv.Type().Elem().NumField())

	for rows.Next() {
		rowp := reflect.New(destv.Type().Elem())
		rowv := rowp.Elem()

		for i := 0; i < rowv.NumField(); i++ {
			args[i] = rowv.Field(i).Addr().Interface()
		}

		if err := rows.Scan(args...); err != nil {
			return err
		}

		destv.Set(reflect.Append(destv, rowv))
	}

	return nil
}

func TransactionpostHandler(c *gin.Context, db *sql.DB) {
	var newTransaction newTransaction

	if c.Bind(&newTransaction) == nil {
		_, err := db.Exec("INSERT INTO transaction VALUES ($1,$2,$3,$4,$5)", newTransaction.Transaction_id, newTransaction.Transaction_date, newTransaction.Transaction_ticket_count, newTransaction.Transaction_visitor_id, newTransaction.Transaction_zoo_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success create"})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
}

func TransactiongetAllHandler(c *gin.Context, db *sql.DB) {
	var newTransactions []newTransaction

	rows, err := db.Query("SELECT * FROM transaction")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	if err := TransactionrowToStruct(rows, &newTransactions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(newTransactions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newTransactions})
}

func TransactiongetHandler(c *gin.Context, db *sql.DB) {
	var newTransaction []newTransaction

	TransactionId := c.Param("id")

	rows, err := db.Query("SELECT * FROM transaction WHERE id = $1", TransactionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	if err := TransactionrowToStruct(rows, &newTransaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(newTransaction) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newTransaction})
}

func TransactionputHandler(c *gin.Context, db *sql.DB) {
	var newTransaction newTransaction

	TransactionId := c.Param("id")

	if c.Bind(&newTransaction) == nil {
		_, err := db.Exec("UPDATE transaction SET date=$1, ticket_count=$2, visitor_id=$3, zoo_id=$4 WHERE id=$5", newTransaction.Transaction_date, newTransaction.Transaction_ticket_count, newTransaction.Transaction_visitor_id, newTransaction.Transaction_zoo_id, TransactionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success update"})
		return
	}
}

func TransactiondelHandler(c *gin.Context, db *sql.DB) {
	TransactionId := c.Param("id")

	_, err := db.Exec("DELETE FROM transaction WHERE id=$1", TransactionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success delete"})
}

func InitTransactionRoutes(group *gin.RouterGroup, db *sql.DB) {
	group.POST("", func(ctx *gin.Context) {
		TransactionpostHandler(ctx, db)
	})

	group.GET("", func(ctx *gin.Context) {
		TransactiongetAllHandler(ctx, db)
	})

	group.GET("/:id", func(ctx *gin.Context) {
		TransactiongetHandler(ctx, db)
	})

	group.PUT("/:id", func(ctx *gin.Context) {
		TransactionputHandler(ctx, db)
	})

	group.DELETE("/:id", func(ctx *gin.Context) {
		TransactiondelHandler(ctx, db)
	})
}
