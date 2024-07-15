package service

import (
	"context"
	"database/sql"
	pb "ecoswap/genproto/swaps"
	"ecoswap/pkg/logger"
	"ecoswap/storage/postgres"
	"fmt"
	"log/slog"
)

type SwapService struct {
	pb.UnimplementedSwapsServer
	SwapsDB *postgres.ItemRepo
	Logger  *slog.Logger
}

func NewSwapService(db *sql.DB) *SwapService {
	return &SwapService{
		SwapsDB: postgres.NewItemRepo(db),
		Logger:  logger.NewLogger(),
	}
}

func (S *SwapService) SendSwapRequest(ctx context.Context, req *pb.SwapRequest) (*pb.SwapResponce, error) {
	resp, err := S.SwapsDB.SendSwapRequest(req)
	if err != nil {
		S.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func (S *SwapService) AdoptionSwapRequest(ctx context.Context, req *pb.Reason) (*pb.Responce, error) {
	resp, err := S.SwapsDB.AdoptionSwapRequest(req)
	if err != nil {
		S.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func (S *SwapService) RejectionSwapRequest(ctx context.Context, req *pb.Reason) (*pb.Responce, error) {
	resp, err := S.SwapsDB.RejectionSwapRequest(req)
	if err != nil {
		S.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func (S *SwapService) GetAllSwapRequests(ctx context.Context, req *pb.FilterField) (*pb.AllSwaps, error) {
	resp, err := S.SwapsDB.GetAllSwapRequests(req)
	if err != nil {
		S.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}
