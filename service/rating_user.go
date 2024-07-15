package service

import (
	"context"
	"database/sql"
	pb "ecoswap/genproto/rating"
	"ecoswap/pkg/logger"
	"ecoswap/storage/postgres"
	"fmt"
	"log/slog"
)

type RatingService struct {
	pb.UnimplementedRatingServer
	RatingDB *postgres.ItemRepo
	Logger   *slog.Logger
}

func NewRatingService(db *sql.DB) *RatingService {
	return &RatingService{
		RatingDB: postgres.NewItemRepo(db),
		Logger:   logger.NewLogger(),
	}
}

func (R *RatingService) CreateUserRating(ctx context.Context, req *pb.RatingReq) (*pb.RatingResp, error) {
	resp, err := R.RatingDB.CreateUserRating(req)
	if err != nil {
		R.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func (R *RatingService) GetUserActivity(ctx context.Context, req *pb.FilterActivity) (*pb.Activity, error) {
	resp, err := R.RatingDB.GetUserActivity(req)
	if err != nil {
		R.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func (R *RatingService) GetUserRating(ctx context.Context, req *pb.FilterField) (*pb.UserRating, error) {
	resp, err := R.RatingDB.GetUserRating(req)
	if err != nil {
		R.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}
