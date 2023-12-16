package handlers

import (
	"database/sql"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type newD_animal struct {
	D_animal_id          uint64 `json:"id" binding:"required"`
	D_animal_name        string `json:"name" binding:"required"`
	D_animal_information string `json:"information" binding:"required"`
	D_animal_species     string `json:"species" binding:"required"`
	D_animal_family      string `json:"family" binding:"required"`
	D_animal_genus       string `json:"genus" binding:"required"`
	D_animal_order       string `json:"order" binding:"required"`
	D_animal_class       string `json:"class" binding:"required"`
	D_animal_phylum      string `json:"phylum" binding:"required"`
}

func D_animalrowToStruct(rows *sql.Rows, dest interface{}) error {
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

func D_animalpostHandler(c *gin.Context, db *sql.DB) {
	var newD_animal newD_animal

	if c.Bind(&newD_animal) == nil {
		_, err := db.Exec("insert into d_animal values ($1,$2,$3,$4,$5,$6,$7,$8,$9)", newD_animal.D_animal_id, newD_animal.D_animal_name, newD_animal.D_animal_information, newD_animal.D_animal_species, newD_animal.D_animal_family, newD_animal.D_animal_genus, newD_animal.D_animal_order, newD_animal.D_animal_class, newD_animal.D_animal_phylum)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": "success create"})
	}

	c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
}

func D_animalgetAllHandler(c *gin.Context, db *sql.DB) {
	var newD_animal []newD_animal

	row, err := db.Query("select * from d_animal")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	D_animalrowToStruct(row, &newD_animal)

	if newD_animal == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newD_animal})

}

func D_animalgetHandler(c *gin.Context, db *sql.DB) {
	var newD_animal []newD_animal

	D_animalId := c.Param("id")

	row, err := db.Query("select * from d_animal where id = $1", D_animalId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	D_animalrowToStruct(row, &newD_animal)

	if newD_animal == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newD_animal})
}

func D_animalputHandler(c *gin.Context, db *sql.DB) {
	var newD_animal newD_animal

	D_animalId := c.Param("id")

	if c.Bind(&newD_animal) == nil {
		_, err := db.Exec("UPDATE d_animal SET name=$1, information=$2, species=$3, family=$4, genus=$5, order=$6, class=$7, phylum=$8 WHERE id=$6", newD_animal.D_animal_name, newD_animal.D_animal_information, newD_animal.D_animal_species, newD_animal.D_animal_family, newD_animal.D_animal_genus, newD_animal.D_animal_order, newD_animal.D_animal_class, newD_animal.D_animal_phylum, D_animalId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": "success update"})
	}

}

func D_animaldelHandler(c *gin.Context, db *sql.DB) {
	D_animalId := c.Param("id")

	_, err := db.Exec("delete from d_animal where id=$1", D_animalId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success delete"})

}

func InitD_animalRoutes(group *gin.RouterGroup, db *sql.DB) {
	group.POST("", func(ctx *gin.Context) {
		D_animalpostHandler(ctx, db)
	})

	group.GET("", func(ctx *gin.Context) {
		D_animalgetAllHandler(ctx, db)
	})

	group.GET("/:id", func(ctx *gin.Context) {
		D_animalgetHandler(ctx, db)
	})

	group.PUT("/:id", func(ctx *gin.Context) {
		D_animalputHandler(ctx, db)
	})

	group.DELETE("/:id", func(ctx *gin.Context) {
		D_animaldelHandler(ctx, db)
	})
}
