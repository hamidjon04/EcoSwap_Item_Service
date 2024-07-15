package service

import (
	"context"
	"database/sql"
	pb "ecoswap/genproto/items"
	"ecoswap/pkg/logger"
	"ecoswap/storage/postgres"
	"fmt"
	"log/slog"
)

type ItemsService struct {
	pb.UnimplementedItemsServer
	ItemDB *postgres.ItemRepo
	Logger *slog.Logger
}

func NewItemService(db *sql.DB) *ItemsService {
	return &ItemsService{
		ItemDB: postgres.NewItemRepo(db),
		Logger: logger.NewLogger(),
	}
}

func (I *ItemsService) CreateItem(ctx context.Context, req *pb.Item) (*pb.ItemResponce, error) {
	resp, err := I.ItemDB.CreateItem(req)
	if err != nil {
		I.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func (I *ItemsService) UpdateItem(ctx context.Context, req *pb.ItemUpdate) (*pb.UpdateResponse, error) {
	resp, err := I.ItemDB.UpdateItem(req)
	if err != nil {
		I.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func (I *ItemsService) DeleteItem(ctx context.Context, req *pb.ItemId) (*pb.Status, error) {
	resp, err := I.ItemDB.DeleteItem(req)
	if err != nil {
		I.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func (I *ItemsService) GetItem(ctx context.Context, req *pb.ItemId) (*pb.GetItemResp, error) {
	resp, err := I.ItemDB.GetItem(req)
	if err != nil {
		I.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func (I *ItemsService) SearchItems(ctx context.Context, req *pb.FilterField) (*pb.AllItems, error) {
	resp, err := I.ItemDB.SearchItems(req)
	if err != nil {
		I.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func (I *ItemsService) GetAllItems(ctx context.Context, req *pb.Limit) (*pb.AllItems, error) {
	resp, err := I.ItemDB.GetAllItems(req)
	if err != nil {
		I.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}

func (I *ItemsService) CreateCategory(ctx context.Context, req *pb.Category) (*pb.CategoryResponse, error) {
	resp, err := I.ItemDB.CreateCategory(req)
	if err != nil {
		I.Logger.Error(fmt.Sprintf("Databazadan ma'lumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return resp, nil
}
