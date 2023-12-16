package handlers

import (
	"database/sql"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type newVisitor struct {
	Visitor_id   uint64 `json:"id" binding:"required"`
	Visitor_name string `json:"name" binding:"required"`
}

func VisitorrowToStruct(rows *sql.Rows, dest interface{}) error {
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

func VisitorpostHandler(c *gin.Context, db *sql.DB) {
	var newVisitor newVisitor

	if c.Bind(&newVisitor) == nil {
		_, err := db.Exec("INSERT INTO visitor VALUES ($1,$2)", newVisitor.Visitor_id, newVisitor.Visitor_name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success create"})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
}

func VisitorgetAllHandler(c *gin.Context, db *sql.DB) {
	var newVisitors []newVisitor

	rows, err := db.Query("SELECT * FROM visitor")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	if err := VisitorrowToStruct(rows, &newVisitors); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(newVisitors) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newVisitors})
}

func VisitorgetHandler(c *gin.Context, db *sql.DB) {
	var newVisitor []newVisitor

	visitorId := c.Param("id")

	rows, err := db.Query("SELECT * FROM visitor WHERE id = $1", visitorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	if err := VisitorrowToStruct(rows, &newVisitor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(newVisitor) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newVisitor})
}

func VisitorputHandler(c *gin.Context, db *sql.DB) {
	var newVisitor newVisitor

	visitorId := c.Param("id")

	if c.Bind(&newVisitor) == nil {
		_, err := db.Exec("UPDATE visitor SET Name=$1 WHERE id=$2", newVisitor.Visitor_name, visitorId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success update"})
		return
	}
}

func VisitordelHandler(c *gin.Context, db *sql.DB) {
	visitorId := c.Param("id")

	_, err := db.Exec("DELETE FROM visitor WHERE id=$1", visitorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success delete"})
}

func InitVisitorRoutes(group *gin.RouterGroup, db *sql.DB) {
	group.POST("", func(ctx *gin.Context) {
		VisitorpostHandler(ctx, db)
	})

	group.GET("", func(ctx *gin.Context) {
		VisitorgetAllHandler(ctx, db)
	})

	group.GET("/:id", func(ctx *gin.Context) {
		VisitorgetHandler(ctx, db)
	})

	group.PUT("/:id", func(ctx *gin.Context) {
		VisitorputHandler(ctx, db)
	})

	group.DELETE("/:id", func(ctx *gin.Context) {
		VisitordelHandler(ctx, db)
	})
}
