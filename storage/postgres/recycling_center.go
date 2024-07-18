package postgres

import (
	"ecoswap/genproto/items"
	pb "ecoswap/genproto/recycling_center"
	"ecoswap/genproto/users"
	"ecoswap/pkg"
	"encoding/json"
	"fmt"
	"log"
	"math"
)

func (R *ItemRepo) CreateRecyclingCenter(req *pb.ResCenter) (*pb.ResponceResCenter, error) {
	resp := pb.ResponceResCenter{}
	query := `
				INSERT INTO item_service_recycling_centers(
					name, address, accepted_materials, working_hours, contact_number)
				VALUES
					($1, $2, $3, $4, $5)
				RETURNING
					id, name, address, accepted_materials, working_hours, contact_number, created_at`
	data, err := json.Marshal(req.AcceptedMaterials)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var scan []byte
	err = R.Db.QueryRow(query, req.Name, req.Address, data, req.WorkingHours, req.ContactNumber).
		Scan(&resp.Id, &resp.Name, &resp.Address, &scan, &resp.WorkingHours, &resp.ContactNumber, &resp.CreatedAt)
	if err != nil{
		log.Println(err)
		return nil, err
	}
	err = json.Unmarshal(scan, &resp.AcceptedMaterials)
	if err != nil{
		log.Println(err)
		return nil, err
	}
	return &resp, nil
}

func (R *ItemRepo) SearchRecyclingCenter(req *pb.FilterField) (*pb.ResAllCenter, error) {
	centers := []*pb.Center{}
	var total int32
	query := `
				SELECT 
					id, name, address, accepted_materials, working_hours, contact_number
				FROM 
					item_service_recycling_centers
				WHERE
					TRUE`
	queryTotal := `
				SELECT
					count(*)
				FROM 
					item_service_recycling_centers
				WHERE 
					TRUE`
		
	arr := []interface{}{}
	if len(req.Material) > 0 {
		query += " AND material = $1"
		queryTotal += " AND material = $1"
		arr = append(arr, req.Material)
	}
	err := R.Db.QueryRow(queryTotal, arr...).Scan(&total)
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

	rows, err := R.Db.Query(query, arr...)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var scan []byte
	for rows.Next() {
		var center pb.Center
		err = rows.Scan(&center.Id, &center.Name, &center.Address, &scan,
			&center.WorkingHours, &center.ContactNumber)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		err := json.Unmarshal(scan, &center.AcceptedMaterials)
		if err != nil{
			log.Println(err)
			return nil, err
		}
		centers = append(centers, &center)
	}

	return &pb.ResAllCenter{
		Centers: centers,
		Total:   total,
		Page:    int32(math.Ceil(float64(total) / float64(req.Limit))),
		Limit:   req.Limit,
	}, nil
}

func (R *ItemRepo) ProductDelivery(req *pb.Submission) (*pb.SubmissionResp, error) {
	resp := pb.SubmissionResp{}
	query := `
				INSERT INTO item_service_recycling_submissions(
					center_id, user_id, items, eco_points_earned)
				VALUES
					($1, $2, $3, $4)
				RETURNING
					id, center_id, user_id, items, eco_points_earned, created_at`
	data, err := json.Marshal(req.JsonDataItems)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = R.Db.QueryRow(query, req.CenterId, req.UserId, data, len(data)).
		Scan(&resp.Id, &resp.CenterId, &resp.UserId, &resp.JsonDataItems, &resp.EcoPointsEarned, &resp.CreatedAt)
	return &resp, err
}

//----------------------------------------------?
func (R *ItemRepo) GetStatistics(req *pb.StatisticField) (*pb.Statistics, error) {
	resp := pb.Statistics{}
	query := `
		SELECT 
			count(*)
		FROM 
			%s
		WHERE
			created_at BETWEEN to_timestamp(%s) AND to_timestamp(%s)`

	err := R.Db.QueryRow(fmt.Sprintf(query, "item_service_swaps", req.StartDate, req.EndDate)).Scan(&resp.TotalSwaps)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = R.Db.QueryRow(fmt.Sprintf(query, "item_service_recycling_submissions", req.StartDate, req.EndDate)).Scan(&resp.TotalRecycledItems)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	users, err := pkg.GetAllUsers(&users.FilterField{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for _, i := range users.Users {
		resp.TotalEarnedPoints += i.EcoPoints
	}

	query = `
		SELECT 
			id, name
		FROM
			(SELECT 
				category_id
			FROM
				item_service_items
			GROUP BY 
				category_id
			ORDER BY 
				COUNT(id) DESC
			LIMIT 1) as T1
		JOIN
			item_service_item_categories as T2
		ON 
			T1.category_id = T2.id`
	var category items.CategoryResponse
	err = R.Db.QueryRow(query).Scan(&category.Id, &category.Name)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resp.JsonDataTopCateg = category.Name

	query = `
		SELECT 
			id, name
		FROM
			(SELECT 
				center_id
			FROM
				item_service_recycling_submissions
			GROUP BY 
				center_id
			ORDER BY 
				COUNT(id) DESC
			LIMIT 1) as T1
		JOIN
			item_service_recycling_centers as T2
		ON 
			T1.center_id = T2.id`
	var center pb.Center
	err = R.Db.QueryRow(query).Scan(&center.Id, &center.Name)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resp.JsonDateTopCenter = center.Name

	return &resp, nil
}
