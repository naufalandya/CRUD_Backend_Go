package handlers

import (
	"database/sql"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type newZooSummary struct {
	ZooSummary_nickname  string `json:"nickname" binding:"required"`
	ZooSummary_name_d    string `json:"name_d" binding:"required"`
	ZooSummary_birthdate string `json:"birthdate" binding:"required"`
	ZooSummary_age       string `json:"age" binding:"required"`
	ZooSummary_text      string `json:"text" binding:"required"`
}

func ZooSummaryrowToStruct(rows *sql.Rows, dest interface{}) error {
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

func ZooSummarygetAllHandler(c *gin.Context, db *sql.DB) {
	var newZooSummarys []newZooSummary

	rows, err := db.Query("SELECT * FROM zoo_animal_summary")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	if err := ZooSummaryrowToStruct(rows, &newZooSummarys); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(newZooSummarys) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newZooSummarys})
}

func ZooSummarygetHandler(c *gin.Context, db *sql.DB) {
	var newZooSummary []newZooSummary

	ZooSummaryNickname := c.Param("nickname")

	rows, err := db.Query("SELECT * FROM zoo_animal_summary WHERE nickname = $1", ZooSummaryNickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	if err := ZooSummaryrowToStruct(rows, &newZooSummary); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(newZooSummary) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newZooSummary})
}

func InitZooSummaryRoutes(group *gin.RouterGroup, db *sql.DB) {

	group.GET("", func(ctx *gin.Context) {
		ZooSummarygetAllHandler(ctx, db)
	})

	group.GET("/:nickname", func(ctx *gin.Context) {
		ZooSummarygetHandler(ctx, db)
	})
}
