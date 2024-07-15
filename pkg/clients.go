package pkg

import (
	"context"
	"ecoswap/config"
	"ecoswap/genproto/users"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetUserPoints(id *users.UserId) (*users.UserEcoPoints, error) {
	conn, err := grpc.NewClient(config.Load().USER_SERVICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	client := users.NewUsersServiceClient(conn)
	resp, err := client.GetEcoPointsByUser(context.Background(), id)
	return resp, err
}

func GetAllUsers(req *users.FilterField) (*users.Users, error) {
	conn, err := grpc.NewClient(config.Load().USER_SERVICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		log.Println(err)
		return nil, err
	}

	client := users.NewUsersServiceClient(conn)
	resp, err := client.GetAllUsers(context.Background(), req)
	return resp, err
}
