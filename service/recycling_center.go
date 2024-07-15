package service

import (
	"context"
	"database/sql"
	pb "ecoswap/genproto/recycling_center"
	"ecoswap/pkg/logger"
	"ecoswap/storage/postgres"
	"fmt"
	"log/slog"
)

type CenterService struct {
	pb.UnimplementedRecyclingCenterServer
	CenterDB *postgres.ItemRepo
	Logger   *slog.Logger
}

func NewCenterService(db *sql.DB) *CenterService {
	return &CenterService{
		CenterDB: postgres.NewItemRepo(db),
		Logger:   logger.NewLogger(),
	}
}

func (C *CenterService) CreateRecyclingCenter(ctx context.Context, req *pb.ResCenter) (*pb.ResponceResCenter, error) {
	resp, err := C.CenterDB.CreateRecyclingCenter(req)
	if err != nil {
		C.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func (C *CenterService) SearchRecyclingCenter(ctx context.Context, req *pb.FilterField) (*pb.ResAllCenter, error) {
	resp, err := C.CenterDB.SearchRecyclingCenter(req)
	if err != nil {
		C.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func (C *CenterService) ProductDelivery(ctx context.Context, req *pb.Submission) (*pb.SubmissionResp, error) {
	resp, err := C.CenterDB.ProductDelivery(req)
	if err != nil {
		C.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func (C *CenterService) GetStatistics(ctx context.Context, req *pb.StatisticField) (*pb.Statistics, error) {
	resp, err := C.CenterDB.GetStatistics(req)
	if err != nil {
		C.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}
