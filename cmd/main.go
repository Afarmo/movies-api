package main

import (
	"context"
	"fmt"
	"log"
	"movies-api/internal/database"
	"movies-api/internal/models"
	"movies-api/internal/repository"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := database.NewDataBase("movieapi.db")
	if err != nil {
		log.Fatal(err)
	}
	if err := database.InitializeQuery(db); err != nil {
		log.Fatal(err)
	}
	bd, _ := time.Parse("2006-01-02", "1964-09-02")
	actor := &models.Actor{Name: "keanau",
		BirthDate: bd}
	r := repository.NewActorRepo(db)

	ctx := context.Background()
	er := r.InsertActor(ctx, actor)
	if er != nil {
		fmt.Println(er)
	}

	actors, e := r.ListAllActors(ctx)
	if e != nil {
		fmt.Println(e)
	}

	for _, a := range actors {
		fmt.Println(a)
	}
	fmt.Println("********************************")
	fmt.Println(r.ListOneActor(ctx, 1000))
	fmt.Println("************************")
	fmt.Println(r.ActorsByName(ctx, "keanau"))

}
