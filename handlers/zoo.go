package handlers

import (
	"database/sql"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type newZoo struct {
	Zoo_id           uint64 `json:"id" binding:"required"`
	Zoo_name         string `json:"name" binding:"required"`
	Zoo_city         string `json:"city" binding:"required"`
	Zoo_country      string `json:"country" binding:"required"`
	Zoo_address      string `json:"address" binding:"required"`
	Zoo_ticket_price uint64 `json:"ticket_price" binding:"required"`
}

func ZoorowToStruct(rows *sql.Rows, dest interface{}) error {
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

func ZoopostHandler(c *gin.Context, db *sql.DB) {
	var newZoo newZoo

	if c.Bind(&newZoo) == nil {
		_, err := db.Exec("INSERT INTO zoo VALUES ($1,$2,$3,$4,$5,$6)", newZoo.Zoo_id, newZoo.Zoo_name, newZoo.Zoo_city, newZoo.Zoo_country, newZoo.Zoo_address, newZoo.Zoo_ticket_price)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success create"})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
}

func ZoogetAllHandler(c *gin.Context, db *sql.DB) {
	var newZoos []newZoo

	rows, err := db.Query("SELECT * FROM zoo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	if err := ZoorowToStruct(rows, &newZoos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(newZoos) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newZoos})
}

func ZoogetHandler(c *gin.Context, db *sql.DB) {
	var newZoo []newZoo

	zooId := c.Param("id")

	rows, err := db.Query("SELECT * FROM Zoo WHERE id = $1", zooId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	if err := ZoorowToStruct(rows, &newZoo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(newZoo) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newZoo})
}

func ZooputHandler(c *gin.Context, db *sql.DB) {
	var newZoo newZoo

	zooId := c.Param("id")

	if c.Bind(&newZoo) == nil {
		_, err := db.Exec("UPDATE zoo SET name=$1, city=$2, country=$3, address=$4, ticket_price=$5 WHERE id=$6", newZoo.Zoo_name, newZoo.Zoo_city, newZoo.Zoo_country, newZoo.Zoo_address, newZoo.Zoo_ticket_price, zooId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success update"})
		return
	}
}

func ZoodelHandler(c *gin.Context, db *sql.DB) {
	zooId := c.Param("id")

	_, err := db.Exec("DELETE FROM zoo WHERE id=$1", zooId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success delete"})
}

func InitZooRoutes(group *gin.RouterGroup, db *sql.DB) {
	group.POST("", func(ctx *gin.Context) {
		ZoopostHandler(ctx, db)
	})

	group.GET("", func(ctx *gin.Context) {
		ZoogetAllHandler(ctx, db)
	})

	group.GET("/:id", func(ctx *gin.Context) {
		ZoogetHandler(ctx, db)
	})

	group.PUT("/:id", func(ctx *gin.Context) {
		ZooputHandler(ctx, db)
	})

	group.DELETE("/:id", func(ctx *gin.Context) {
		ZoodelHandler(ctx, db)
	})
}
