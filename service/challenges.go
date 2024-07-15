package service

import (
	"context"
	"database/sql"
	pb "ecoswap/genproto/challenges"
	"ecoswap/pkg/logger"
	"ecoswap/storage/postgres"
	"fmt"
	"log/slog"
)


type ChallengeService struct{
	pb.UnimplementedChallengesServer
	ChallengeDB *postgres.ItemRepo
	Logger *slog.Logger
}

func NewChallengeService(db *sql.DB)*ChallengeService{
	return &ChallengeService{
		ChallengeDB: postgres.NewItemRepo(db),
		Logger: logger.NewLogger(),
	}
}

func(C *ChallengeService) CreateChallenge(ctx context.Context, req *pb.Challenge)(*pb.RespChallenge, error){
	resp, err := C.ChallengeDB.CreateChallenge(req)
	if err != nil{
		C.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func(C *ChallengeService) AttendChallenge(ctx context.Context, req *pb.Attend)(*pb.AttendResp, error){
	resp, err := C.ChallengeDB.AttendChallenge(req)
	if err != nil{
		C.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func(C *ChallengeService) UpdateChallengeResult(ctx context.Context, req *pb.ChallengeUpdate)(*pb.RespUpdate, error){
	resp, err := C.ChallengeDB.UpdateChallengeResult(req)
	if err != nil{
		C.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func(C *ChallengeService) CreateEcoTips(ctx context.Context, req *pb.EcoTip)(*pb.RespEcoTip, error){
	resp, err := C.ChallengeDB.CreateEcoTips(req)
	if err != nil{
		C.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func(C *ChallengeService) GetAllEcoTips(ctx context.Context, req *pb.FilterTip)(*pb.Tips, error){
	resp, err := C.ChallengeDB.GetAllEcoTips(req)
	if err != nil{
		C.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}