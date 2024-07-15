package postgres

import (
	pb "ecoswap/genproto/rating"
	"ecoswap/genproto/users"
	"ecoswap/pkg"
	"fmt"
	"log"
	"math"
)

func (R *ItemRepo) CreateUserRating(req *pb.RatingReq) (*pb.RatingResp, error) {
	resp := pb.RatingResp{}
	query := `
				INSERT INTO item_service_ratings(
					user_id, rater_id, rating, comment, swap_id)
				VALUES
					($1, $2, $3, $4, $5)
				RETURNING
					id, user_id, rater_id, rating, comment, swap_id, created_at`
	err := R.Db.QueryRow(query, req.UserId, req.RaterId, req.Rating, req.Comment, req.SwapId).
		Scan(&resp.Id, &resp.UserId, &resp.RaterId, &resp.Rating, &resp.Comment, &resp.SwapId, &resp.CreatedAt)
	return &resp, err
}

func (R *ItemRepo) GetUserActivity(req *pb.FilterActivity) (*pb.Activity, error) {
	resp := pb.Activity{UserId: req.UserId}
	query := `
				SELECT 
					count(id)
				FROM 
					%s
				WHERE 
					%s = $1 AND 
					created_at BETWEEN %s AND %s`

	err := R.Db.QueryRow(fmt.Sprintf(query, "item_service_swaps", "requester_id", req.StartDate, req.EndDate), req.UserId).
		Scan(&resp.SwapsInitiated)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = R.Db.QueryRow(fmt.Sprintf(query, "item_service_swaps", "owner_id", req.StartDate, req.EndDate), req.UserId).
		Scan(&resp.SwapsCompleted)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = R.Db.QueryRow(fmt.Sprintf(query, "item_service_items", "owner_id", req.StartDate, req.EndDate), req.UserId).
		Scan(&resp.ItemsListed)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = R.Db.QueryRow(fmt.Sprintf(query, "item_service_recycling_submissions", "user_id", req.StartDate, req.EndDate), req.UserId).
		Scan(&resp.RecyclingSubmissions)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	userPoint, err := pkg.GetUserPoints(&users.UserId{
		Id: req.UserId,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resp.EcoPointsEarned = userPoint.EcoPoints

	return &resp, nil
}

func (R *ItemRepo) GetUserRating(req *pb.FilterField) (*pb.UserRating, error) {
	ratings := []*pb.RatingUser{}
	var total int32
	query := `
				SELECT 
					id, count(id), rater_id, rating, comment, swap_id, created_at
				FROM 
					item_service_ratings
				WHERE 
					user_id = $1`
	err := R.Db.QueryRow(query, req.UserId).Scan(nil, &total, nil, nil, nil, nil, nil)
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
	rows, err := R.Db.Query(query, req.UserId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var sumRatings float32
	for rows.Next() {
		var rating pb.RatingUser
		err = rows.Scan(&rating.Id, nil, &rating.RaterId, &rating.Rating, &rating.Comment, &rating.SwapId, &rating.CreatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		sumRatings += rating.Rating
		ratings = append(ratings, &rating)
	}

	return &pb.UserRating{
		UserId:        req.UserId,
		Ratings:       ratings,
		AverageRating: sumRatings / float32(total),
		TotalRatings:  total,
		Page:          int32(math.Ceil(float64(total) / float64(req.Limit))),
		Limit:         req.Limit,
	}, nil
}
