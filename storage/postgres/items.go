package postgres

import (
	pb "ecoswap/genproto/items"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

func (I *ItemRepo) CreateItem(item *pb.Item) (*pb.ItemResponce, error) {
	resp := pb.ItemResponce{}
	query := `
				INSERT INTO item_service_items(
					name, description, category_id, condition, swap_prefernce, owner_id)
				VALUES
					($1, $2, $3, $4, $5, $6)
				RETURNING
					id, name, description, category_id, condition, swap_prefernce,
					owner_id, status, created_at`
	err := I.Db.QueryRow(query, item.Name, item.Description, item.CategoryId,
		item.Condition, item.SwapPreference, item.OwnerId).
		Scan(
			&resp.Id, &resp.Name, &resp.Description, &resp.CategoryId, &resp.Condition,
			&resp.SwapPreference, &resp.OwnerId, &resp.Status, &resp.CreatedAt)
	return &resp, err
}

func (I *ItemRepo) UpdateItem(req *pb.ItemUpdate) (*pb.UpdateResponse, error) {
	resp := pb.UpdateResponse{}
	query := `
				UPDATE item_service_items SET 
					name = $1, condition = $2, status = $3, updated_at = $4
				WHERE 
					id = $5 AND deleted_at is null
				RETURNING
					id, name, description, category_id, condition, swap_prefernce, 
					owner_id, status, updated_at`
	err := I.Db.QueryRow(query, req.Name, req.Condition, req.Status, req.Id).
		Scan(
			&resp.Id, &resp.Name, &resp.Description, &resp.CategoryId, &resp.Condition,
			&resp.SwapPreference, &resp.OwnerId, &resp.Status, &resp.UpdatedAt)
	return &resp, err
}

func (I *ItemRepo) DeleteItem(req *pb.ItemId) (*pb.Status, error) {
	query := `
				UPDATE item_service_items SET
					deleted_at = $1
				WHERE
					deleted_at is null AND id = $2`
	result, err := I.Db.Exec(query, time.Now(), req.Id)
	if err != nil {
		return &pb.Status{
			Status:  false,
			Message: fmt.Sprintf("Ma'lumotlar o'chirilmadi: %v", err),
		}, err
	}
	num, err := result.RowsAffected()
	if err != nil || num == 0 {
		return &pb.Status{
			Status:  false,
			Message: fmt.Sprintf("Bazada bunda id mavjud emas: %v", err),
		}, err
	}
	return &pb.Status{
		Status:  true,
		Message: "Ma'lumotlar muvaffaqiyatli o'chirildi",
	}, nil
}

func (I *ItemRepo) GetItem(itemId *pb.ItemId) (*pb.GetItemResp, error) {
	resp := pb.GetItemResp{}
	query := `
				SELECT 
					id, name, description, catrgory_id, condetion, swap_prefernce, 
					owner_id, status, created_at, updated_at
				FROM 
					item_service_items
				WHERE 
					id = $1 AND deleted_at is null`
	err := I.Db.QueryRow(query, itemId.Id).Scan(
		&resp.Id, &resp.Name, &resp.Description, &resp.CategoryId, &resp.Condition,
		&resp.SwapPreference, &resp.OwnerId, &resp.Status, &resp.CreatedAt, &resp.UpdatedAt)
	return &resp, err
}

func (I *ItemRepo) SearchItems(req *pb.FilterField) (*pb.AllItems, error) {
	items := []*pb.AllItem{}
	var total int32

	query := `
				SELECT 
					id, count(id), name, category_id, condition, owner_id, status
				FROM 
					item_service_items
				WHERE
					deleted_at is null`
	param := []string{}
	arr := []interface{}{}

	if len(req.Name) > 0 {
		query += " AND name = :name"
		param = append(param, ":name")
		arr = append(arr, req.Name)
	}
	if len(req.CategoryId) > 0 {
		query += " AND category_id = :category_id"
		param = append(param, ":category_id")
		arr = append(arr, req.CategoryId)
	}
	if len(req.Condition) > 0 {
		query += " AND condition = :condition"
		param = append(param, ":condition")
		arr = append(arr, req.Condition)
	}

	for i, j := range param {
		query = strings.Replace(query, j, "$"+strconv.Itoa(i+1), 1)
	}

	err := I.Db.QueryRow(query, arr...).Scan(nil, &total, nil, nil, nil, nil, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if req.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", req.Limit)
	}else{
		req.Limit = total
	}
	if req.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	rows, err := I.Db.Query(query, arr...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		var item pb.AllItem
		err = rows.Scan(&item.Id, nil, &item.Name, &item.CategoryId, &item.Condition,
			&item.OwnerId, &item.Status)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		items = append(items, &item)
	}

	return &pb.AllItems{
		Items: items,
		Total: total,
		Page:  int32(math.Ceil(float64(total) / float64(req.Limit))),
		Limit: req.Limit,
	}, nil
}

func(I *ItemRepo) GetAllItems(req *pb.Limit) (*pb.AllItems, error){
	items := []*pb.AllItem{}
	var total int32
	query := `
				SELECT 
					id, count(id), name, category_id, condition, owner_id, status
				FROM 
					item_service_items
				WHERE 
					deleted_at is null`
	
	err := I.Db.QueryRow(query).Scan(nil, &total, nil, nil, nil, nil, nil)
	if err != nil{
		log.Println(err)
		return nil, err
	}
	
	if req.Limit > 0{
		query += fmt.Sprintf(" LIMIT %d", req.Limit)
	}else{
		req.Limit = total
	}
	if req.Offset > 0{
		query += fmt.Sprintf(" OFFSET %d", req.Offset)
	}


	rows, err := I.Db.Query(query)
	if err != nil{
		log.Println(err)
		return nil, err
	}
	for rows.Next(){
		var item pb.AllItem
		err := rows.Scan(&item.Id, nil, &item.Name, &item.CategoryId, &item.Condition, &item.OwnerId, &item.Status)
		if err != nil{
			log.Println(err)
			return nil, err
		}
		items = append(items, &item)
	}
	return &pb.AllItems{
		Items: items,
		Total: total,
		Page: int32(math.Ceil(float64(total) / float64(req.Limit))),
		Limit: req.Limit,
	}, nil
}

func(I *ItemRepo) CreateCategory(req *pb.Category)(*pb.CategoryResponse, error){
	resp := pb.CategoryResponse{}
	query := `
				INSERT INTO item_service_item_categories(
					name, description)
				VALUES
					($1, $2)
				RETURNING
					id, name, description, created_at`
	err := I.Db.QueryRow(query, req.Name, req.Description).
		Scan(&resp.Id, &resp.Name, &resp.Description, &resp.CreatedAt)
	return &resp, err
}
