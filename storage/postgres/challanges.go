package postgres

import (
	pb "ecoswap/genproto/challenges"
	"fmt"
	"log"
	"math"
	"time"
)

func (C *ItemRepo) CreateChallenge(req *pb.Challenge) (*pb.RespChallenge, error) {
	resp := pb.RespChallenge{}
	query := `
				INSERT INTO item_service_eco_challenges(
					title, description, start_date, end_date, reward_points)
				VALUES
					($1, $2, $3, $4, $5)
				RETURNING
					id, title, description, start_date, end_date, reward_points, created_at`
	err := C.Db.QueryRow(query, req.Title, req.Description, req.StartDate, req.EndDate, req.RewardPoints).
		Scan(
			&resp.Id, &resp.Title, &resp.Description, &resp.StartDate, &resp.EndDate, &resp.RewardPoints, &resp.CreatedAt)
	return &resp, err
}

func (C *ItemRepo) AttendChallenge(req *pb.Attend) (*pb.AttendResp, error) {
	resp := pb.AttendResp{}
	query := `
				INSERT INTO item_service_challenge_participations(
					challenge_id, user_id)
				VALUES
					($1, $2)
				RETURNING
					challenge_id, user_id, status, joined_at`
	err := C.Db.QueryRow(query, req.ChallengeId, req.UserId).
		Scan(&resp.ChallengeId, &resp.UserId, &resp.Status, &resp.JoinedAt)
	return &resp, err
}

func (C *ItemRepo) UpdateChallengeResult(req *pb.ChallengeUpdate) (*pb.RespUpdate, error) {
	resp := pb.RespUpdate{}
	query := `
				UPDATE item_service_challenge_participations SET
					recycled_items_count = $1,
					updated_at = $2
				WHERE
					challenge_id = $3 AND user_id = $4
				RETURNING
					challenge_id, user_id, status, recycled_items_count, updated_at`
	err := C.Db.QueryRow(query, req.ResItemCount, time.Now(), req.ChallengeId, req.UserId).
		Scan(&resp.ChallengeId, &resp.UserId, &resp.Status, &resp.ResItemCount, &resp.UpdatedAt)
	return &resp, err
}

func (C *ItemRepo) CreateEcoTips(req *pb.EcoTip) (*pb.RespEcoTip, error) {
	resp := pb.RespEcoTip{}
	query := `
				INSERT INTO item_service_eco_tips(
					title, content)
				VALUES
					($1, $2)
				RETURNING
					id, title, content, created_at`
	err := C.Db.QueryRow(query, req.Title, req.Content).Scan(&resp.Id, &resp.Title, &resp.Content, &resp.CreatedAt)
	return &resp, err
}

func (C *ItemRepo) GetAllEcoTips(req *pb.FilterTip) (*pb.Tips, error) {
	tips := []*pb.RespEcoTip{}
	var total int32
	query := `
				SELECT 
					id, title, content, created_at
				FROM 
					item_service_eco_tips
				WHERE 
					TRUE`
	queryTotal := `
					SELECT 
						count(*)
					FROM 
						item_service_eco_tips
					WHERE 
						TRUE`
	arr := []interface{}{}
	if len(req.Title) > 0 {
		query += " AND title = $1"
		queryTotal += " AND title = $1"
		arr = append(arr, req.Title)
	}

	err := C.Db.QueryRow(queryTotal, arr...).Scan(&total)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if req.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", req.Limit)
	} else {
		req.Limit = total
	}
	if req.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	rows, err := C.Db.Query(query, arr...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		var tip pb.RespEcoTip
		err = rows.Scan(
			&tip.Id, &tip.Title, &tip.Content, &tip.CreatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		tips = append(tips, &tip)
	}
	return &pb.Tips{
		Tips:  tips,
		Total: total,
		Page:  int32(math.Ceil(float64(total) / float64(req.Limit))),
		Limit: req.Limit,
	}, nil
}
