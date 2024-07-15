package main

import (
	"ecoswap/config"
	"ecoswap/genproto/challenges"
	"ecoswap/genproto/items"
	"ecoswap/genproto/rating"
	"ecoswap/genproto/recycling_center"
	"ecoswap/genproto/swaps"
	"ecoswap/service"
	"ecoswap/storage/postgres"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main(){
	db, err := postgres.ConnectDB()
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()

	listener, err := net.Listen("tcp", config.Load().ITEM_SERVICE)
	if err != nil{
		log.Fatal(err)
	}
	defer listener.Close()

	challengeService := service.NewChallengeService(db)
	itemService := service.NewItemService(db)
	ratingService := service.NewRatingService(db)
	centerService := service.NewCenterService(db)
	swapService := service.NewSwapService(db)

	service := grpc.NewServer()

	challenges.RegisterChallengesServer(service, challengeService)	
	items.RegisterItemsServer(service, itemService)
	rating.RegisterRatingServer(service, ratingService)
	recycling_center.RegisterRecyclingCenterServer(service, centerService)
	swaps.RegisterSwapsServer(service, swapService)

	log.Printf("Server run: %s", config.Load().ITEM_SERVICE)
	if err = service.Serve(listener); err != nil{
		log.Fatal(err)
	}
}