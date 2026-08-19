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
	actor := &models.Actor{BirthDate: bd}
	r := repository.NewActorRepo(db)

	ctx := context.Background()
	err = r.InsertActor(ctx, actor)
	if err != nil {
		fmt.Println(err)
	}

	actors, e := r.ListAllActors(ctx)
	if e != nil {
		fmt.Println(e)
	}

	for _, a := range actors {
		fmt.Println(a)
	}
	fmt.Println("********************************")
	fmt.Println(r.ListOneActor(ctx, 1))
	fmt.Println("************************")
	actorz, _ := r.ActorsByName(ctx, "  keanau")
	for _, a := range actorz {
		fmt.Println(a.ID, a.Name, a.BirthDate.Format("2006-01-02"))
	}

}
