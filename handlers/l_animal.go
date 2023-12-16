package handlers

import (
	"database/sql"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type newL_animal struct {
	L_Animal_id          uint64 `json:"id" binding:"required"`
	L_Animal_nickname    string `json:"nickname" binding:"required"`
	L_Animal_birthdate   string `json:"birthdate" binding:"required"`
	L_Animal_age         string `json:"age" binding:"required"`
	L_Animal_d_animal_id uint64 `json:"d_animal_id" binding:"required"`
	L_Animal_area_id     uint64 `json:"area_id" binding:"required"`
}

func L_AnimalrowToStruct(rows *sql.Rows, dest interface{}) error {
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

func L_AnimalpostHandler(c *gin.Context, db *sql.DB) {
	var newL_Animal newL_animal

	if c.Bind(&newL_Animal) == nil {
		_, err := db.Exec("insert into l_Animal values ($1,$2,$3,$4,$5,$6)", newL_Animal.L_Animal_id, newL_Animal.L_Animal_nickname, newL_Animal.L_Animal_birthdate, newL_Animal.L_Animal_age, newL_Animal.L_Animal_d_animal_id, newL_Animal.L_Animal_area_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": "success create"})
	}

	c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
}

func L_AnimalgetAllHandler(c *gin.Context, db *sql.DB) {
	var newL_animal []newL_animal

	row, err := db.Query("select * from l_Animal")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	L_AnimalrowToStruct(row, &newL_animal)

	if newL_animal == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newL_animal})

}

func L_AnimalgetHandler(c *gin.Context, db *sql.DB) {
	var newL_Animal []newL_animal

	L_AnimalId := c.Param("id")

	row, err := db.Query("select * from l_Animal where id = $1", L_AnimalId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	L_AnimalrowToStruct(row, &newL_Animal)

	if newL_Animal == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newL_Animal})
}

func L_AnimalputHandler(c *gin.Context, db *sql.DB) {
	var newL_Animal newL_animal

	L_AnimalId := c.Param("id")

	if c.Bind(&newL_Animal) == nil {
		_, err := db.Exec("UPDATE l_Animal SET nickname=$1, birthdate=$2, age=$3, d_animal_id=$4, area_id=$5 WHERE id=$6", newL_Animal.L_Animal_nickname, newL_Animal.L_Animal_birthdate, newL_Animal.L_Animal_age, newL_Animal.L_Animal_d_animal_id, newL_Animal.L_Animal_area_id, L_AnimalId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": "success update"})
	}

}

func L_AnimaldelHandler(c *gin.Context, db *sql.DB) {
	L_AnimalId := c.Param("id")

	_, err := db.Exec("delete from l_Animal where id=$1", L_AnimalId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success delete"})

}

func InitL_AnimalRoutes(group *gin.RouterGroup, db *sql.DB) {
	group.POST("", func(ctx *gin.Context) {
		L_AnimalpostHandler(ctx, db)
	})

	group.GET("", func(ctx *gin.Context) {
		L_AnimalgetAllHandler(ctx, db)
	})

	group.GET("/:id", func(ctx *gin.Context) {
		L_AnimalgetHandler(ctx, db)
	})

	group.PUT("/:id", func(ctx *gin.Context) {
		L_AnimalputHandler(ctx, db)
	})

	group.DELETE("/:id", func(ctx *gin.Context) {
		L_AnimaldelHandler(ctx, db)
	})
}
