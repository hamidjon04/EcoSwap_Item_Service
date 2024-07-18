package postgres

import (
	pb "ecoswap/genproto/swaps"
	"fmt"
	"log"
	"math"
)

func (S *ItemRepo) SendSwapRequest(req *pb.SwapRequest) (*pb.SwapResponce, error) {
	resp := pb.SwapResponce{}
	query := `
				INSERT INTO item_service_swaps(
					requester_id, owner_id, offered_item_id,
					requested_item_id, message, status)
				VALUES
					($1, $2, $3, $4, $5, $6)
				RETURNING
					id, requester_id, owner_id, offered_item_id,
					requested_item_id, message, status, created_at`
	err := S.Db.QueryRow(query, req.RequesterId, req.OwnerId, req.OfferedItemId,
		req.RequestedItemId, req.Message, "activ").
		Scan(
			&resp.Id, &resp.RequesterId, &resp.OwnerId, &resp.OfferedItemId,
			&resp.RequestedItemId, &resp.Message, &resp.Status, &resp.CreatedAt)
	return &resp, err
}

func (S *ItemRepo) AdoptionSwapRequest(req *pb.Reason) (*pb.Responce, error) {
	resp := pb.Responce{}
	query := `
				UPDATE item_service_swaps SET
					status = 'accepted',
					message = $2
				WHERE 
					deleted_at is null AND id = $1
				RETURNING
					id, requester_id, owner_id, offered_item_id,
					requested_item_id, status, message, updated_at`
	err := S.Db.QueryRow(query, req.SwapId, req.Reason).
		Scan(
			&resp.Id, &resp.RequesterId, &resp.OwnerId, &resp.OfferedItemId,
			&resp.RequestedItemId, &resp.Status, &resp.Reason, &resp.UpdatedAt)
	return &resp, err
}

func (S *ItemRepo) RejectionSwapRequest(req *pb.Reason) (*pb.Responce, error) {
	resp := pb.Responce{}
	query := `
				UPDATE item_service_swaps SET
					status = 'rejected',
					message = $2
				WHERE 
					deleted_at is null AND id = $1
				RETURNING
					id, requester_id, owner_id, offered_item_id,
					requested_item_id, status, message, updated_at`
	err := S.Db.QueryRow(query, req.SwapId, req.Reason).
		Scan(
			&resp.Id, &resp.RequesterId, &resp.OwnerId, &resp.OfferedItemId,
			&resp.RequestedItemId, &resp.Status, &resp.Reason, &resp.UpdatedAt)
	return &resp, err
}

func (S *ItemRepo) GetAllSwapRequests(req *pb.FilterField) (*pb.AllSwaps, error) {
	swaps := []*pb.Swap{}
	var total int32
	query := `
				SELECT 
					id, requester_id, owner_id, offered_item_id,
					requested_item_id, message, created_at
				FROM 
					item_service_swaps
				WHERE 
					deleted_at is null`
	queryTotal := `
					SELECT 
						count(*)
					FROM
						item_service_swaps
					WHERE 
						deleted_at is null`
	arr := []interface{}{}
	if len(req.Status) > 0 {
		query += " AND status = $1"
		queryTotal += " AND status = $1"
		arr = append(arr, req.Status)
	}
	
	err := S.Db.QueryRow(queryTotal, arr...).Scan(&total)
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

	rows, err := S.Db.Query(query, arr...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		var swap pb.Swap
		err = rows.Scan(
			&swap.Id, &swap.RequesterId, &swap.OwnerId, &swap.OfferedItemId,
			&swap.RequestedItemId, &swap.Status, &swap.CreatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		swaps = append(swaps, &swap)
	}
	return &pb.AllSwaps{
		Swaps: swaps,
		Total: total,
		Page:  int32(math.Ceil(float64(total) / float64(req.Limit))),
		Limit: req.Limit,
	}, nil
}
