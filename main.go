package main

import (
	"belajar_REST/db"
	"belajar_REST/handlers"
	"log"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func setupRouter() *gin.Engine {
	db, err := db.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3000/dashboard"} 
	r.Use(cors.New(config))

	studentGroup := r.Group("/data/API/student")
	handlers.InitStudentRoutes(studentGroup, db)

	zooGroup := r.Group("/data/API/zoo")
	handlers.InitZooRoutes(zooGroup, db)

	visitorGroup := r.Group("data/API/visitor")
	handlers.InitVisitorRoutes(visitorGroup, db)

	areaGroup := r.Group("data/API/area")
	handlers.InitAreaRoutes(areaGroup, db)

	transactionGroup := r.Group("data/API/transaction")
	handlers.InitTransactionRoutes(transactionGroup, db)

	L_animalGroup := r.Group("data/API/l_animal")
	handlers.InitL_AnimalRoutes(L_animalGroup, db)

	D_animalGroup := r.Group("data/API/d_animal")
	handlers.InitD_animalRoutes(D_animalGroup, db)

	//views
	ZooSummaryGroup := r.Group("data/API/nickname_l_animal")
	handlers.InitZooSummaryRoutes(ZooSummaryGroup, db)

	return r
}

func main() {
	r := setupRouter()

	r.Run()

}

//parameter harus sama ama kolom di tabel database
