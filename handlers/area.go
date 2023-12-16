package handlers

import (
	"database/sql"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type newArea struct {
	Area_id   uint64 `json:"id" binding:"required"`
	Area_name string `json:"name" binding:"required"`
	Zoo_id    uint64 `json:"zoo_id" binding:"required"`
}

func ArearowToStruct(rows *sql.Rows, dest interface{}) error {
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

func AreapostHandler(c *gin.Context, db *sql.DB) {
	var newArea newArea

	if c.Bind(&newArea) == nil {

		_, err := db.Exec("INSERT INTO area VALUES ($1, $2, $3)", newArea.Area_id, newArea.Area_name, newArea.Zoo_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success create"})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request data"})
}

func AreagetAllHandler(c *gin.Context, db *sql.DB) {
	var newArea []newArea

	rows, err := db.Query("SELECT * FROM area")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	if err := ArearowToStruct(rows, &newArea); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(newArea) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newArea})
}

func AreagetHandler(c *gin.Context, db *sql.DB) {
	var newArea []newArea

	areaId := c.Param("id")

	rows, err := db.Query("SELECT * FROM area WHERE id = $1", areaId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	if err := ArearowToStruct(rows, &newArea); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(newArea) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newArea})
}

func AreaputHandler(c *gin.Context, db *sql.DB) {
	var newArea newArea

	areaId := c.Param("id")

	if c.Bind(&newArea) == nil {
		_, err := db.Exec("UPDATE area SET name=$1, id=$2, WHERE id=$3", newArea.Area_name, newArea.Zoo_id, areaId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success update"})
		return
	}
}

func AreadelHandler(c *gin.Context, db *sql.DB) {
	areaId := c.Param("id")

	_, err := db.Exec("DELETE FROM area WHERE id=$1", areaId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success delete"})
}

func InitAreaRoutes(group *gin.RouterGroup, db *sql.DB) {
	group.POST("", func(ctx *gin.Context) {
		AreapostHandler(ctx, db)
	})

	group.GET("", func(ctx *gin.Context) {
		AreagetAllHandler(ctx, db)
	})

	group.GET("/:id", func(ctx *gin.Context) {
		AreagetHandler(ctx, db)
	})

	group.PUT("/:id", func(ctx *gin.Context) {
		AreaputHandler(ctx, db)
	})

	group.DELETE("/:id", func(ctx *gin.Context) {
		AreadelHandler(ctx, db)
	})
}
