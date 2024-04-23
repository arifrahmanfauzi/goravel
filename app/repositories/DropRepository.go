package repositories

import (
	"context"
	"goravel/app/helpers"
	"goravel/app/models"
	"goravel/app/transformers"
	"time"

	"github.com/goravel/framework/facades"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DropRepository struct {
	collection *mongo.Collection
}

func NewDropRepository() *DropRepository {
	return &DropRepository{
		collection: Connection(models.Drop{}),
	}
}

func (Drop *DropRepository) GetAll(page int64, pageSize int64, total *int64, totalPage *int64) ([]*models.Drop, error) {
	ctx := context.Background()
	// Calculate the number of documents to skip
	skip := (page - 1) * pageSize
	// Get total count of trips
	totalRecord, err := Drop.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	*total = totalRecord
	// Calculate total pages
	totalPages := totalRecord / pageSize
	if totalRecord%pageSize != 0 {
		totalPages++
	}
	*totalPage = totalPages

	cursor, err := Drop.collection.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &pageSize,
	})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			facades.Log().Error(err)
		}
	}(cursor, ctx)

	var Drops []*models.Drop
	for cursor.Next(ctx) {
		var customerOrder models.Drop
		if err := cursor.Decode(&customerOrder); err != nil {
			return nil, err
		}
		Drops = append(Drops, &customerOrder)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return Drops, nil
}
func (Drop *DropRepository) Update(ID string, update map[string]interface{}) *models.Drop {
	id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		facades.Log().Error(err)
		return nil
	}
	filter := bson.M{"_id": id}
	update = bson.M{"$set": update}
	// Set returnDocument for returning the updated document
	returnDocument := options.FindOneAndUpdate().SetReturnDocument(options.After)
	//Update content
	var Drops *models.Drop
	err = Drop.collection.FindOneAndUpdate(context.Background(), filter, update, returnDocument).Decode(&Drops)
	if err != nil {
		facades.Log().Error(err)
		return nil
	}

	return Drops
}
func (Drop *DropRepository) Delete(ID string) (*mongo.DeleteResult, error) {
	id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		facades.Log().Error(err)
		return nil, err
	}
	filter := bson.M{"_id": id}
	res, err := Drop.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		facades.Log().Error(err)
	}
	return res, err
}
func (Drop *DropRepository) FindByTripNumber(TripNumber string, Page int64, PageSize int64) ([]*models.Drop, int64, int64, error) {
	ctx := context.Background()
	// Calculate the number of documents to skip
	skip := (Page - 1) * PageSize
	rootStage := bson.D{
		{"$project",
			bson.D{
				{"cDT", "$$ROOT"},
				{"_id", 0},
			},
		},
	}
	lookUpStage := bson.D{
		{"$lookup",
			bson.D{
				{"localField", "cDT.tripIdObject"},
				{"from", "customerTripPlanning"},
				{"foreignField", "_id"},
				{"as", "cTP"},
			},
		},
	}
	unwindStage := bson.D{
		{"$unwind",
			bson.D{
				{"path", "$cTP"},
				{"preserveNullAndEmptyArrays", false},
			},
		},
	}
	matchStage := bson.D{{"$match", bson.D{{"cTP.tripNumber", bson.D{{"$eq", TripNumber}}}}}}
	countStage := bson.D{{"$count", "count"}}
	projectStage := bson.D{
		{"$project", bson.D{
			{"_id", "$cDT._id"},
			{"dispatchNumber", "$cDT.dispatchNumber"},
			{"tripId", "$cDT.tripId"},
			{"tripNumber", "$cTP.tripNumber"},
			{"tripStatus", "$cTP.tripStatus"},
		}},
	}
	countPipeline := bson.A{
		rootStage,
		lookUpStage,
		unwindStage,
		matchStage,
		countStage,
	}
	countCursor, err := Drop.collection.Aggregate(ctx, countPipeline)
	defer func(countCursor *mongo.Cursor, ctx context.Context) {
		err := countCursor.Close(ctx)
		if err != nil {
			facades.Log().Error(err)
		}
	}(countCursor, context.Background())
	if err != nil {
		facades.Log().Error(err)
	}
	if err != nil {
		facades.Log().Error(err)
		return nil, 0, 0, err
	}
	// Get total count of trips
	var countResult []bson.M
	if err := countCursor.All(context.Background(), &countResult); err != nil {
		facades.Log().Error(err)
		return nil, 0, 0, err
	}
	var totalRecord int64
	if len(countResult) > 0 {
		totalRecord = int64(countResult[0]["count"].(int32))
	}
	// Calculate total pages
	totalPages := totalRecord / PageSize
	if totalRecord%PageSize != 0 {
		totalPages++
	}
	totalPage := totalPages

	pipeline := mongo.Pipeline{
		rootStage,
		lookUpStage,
		unwindStage,
		matchStage,
		bson.D{{"$skip", skip}},
		bson.D{{"$limit", PageSize}},
		projectStage}
	// Find drops query
	cursor, err := Drop.collection.Aggregate(ctx, pipeline)
	if err != nil {
		facades.Log().Error(err)
		return nil, totalRecord, 0, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			facades.Log().Error(err)
		}
	}(cursor, ctx)
	var Drops []*models.Drop
	if err := cursor.All(context.Background(), &Drops); err != nil {
		return nil, totalRecord, totalPage, err
	}
	return Drops, totalRecord, totalPage, nil
}
func (Drop *DropRepository) FindByTripStatus(TripStatus string, Page int, Limit int) ([]*models.Drop, *transformers.Pagination) {
	opts := options.Aggregate().SetMaxTime(10 * time.Second)
	pipeline := mongo.Pipeline{
		bson.D{{"$project", bson.D{{"cDT", "$$ROOT"}, {"_id", 0}}}},
		bson.D{{"$lookup", bson.D{{"localField", "cDT.tripIdObject"}, {"from", "customerTripPlanning"}, {"foreignField", "_id"}, {"as", "cTP"}}}},
		bson.D{{"$unwind", bson.D{{"path", "$cTP"}, {"preserveNullAndEmptyArrays", false}}}},
		bson.D{{"$match", bson.D{{"cTP.tripStatus", bson.D{{"$regex", primitive.Regex{Pattern: TripStatus}}}}}}},
		bson.D{{"$replaceRoot", bson.D{{"newRoot", bson.D{{"$mergeObjects", bson.A{"$cDT", "$cTP", "$$ROOT"}}}}}}},
		bson.D{{"$project", bson.D{{"cDT", 0}, {"cTP", 0}}}},
	}

	cursor, err := Drop.collection.Aggregate(context.Background(), pipeline, opts)
	if err != nil {
		facades.Log().Error(err)
		return nil, nil
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			facades.Log().Error(err)
		}
	}(cursor, context.Background())

	var drops []*models.Drop
	if err := cursor.All(context.Background(), &drops); err != nil {
		facades.Log().Error(err)
		return nil, nil
	}

	total := len(drops)
	Skip, TotalPages := helpers.Paginate(Page, total, Limit)

	pipeline = append(pipeline, bson.D{{"$limit", int64(Limit)}}, bson.D{{"$skip", int64(Skip)}})

	cursor, err = Drop.collection.Aggregate(context.Background(), pipeline, opts)
	if err != nil {
		facades.Log().Error(err)
		return nil, nil
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			facades.Log().Error(err)
		}
	}(cursor, context.Background())

	drops = []*models.Drop{}
	if err := cursor.All(context.Background(), &drops); err != nil {
		facades.Log().Error(err)
		return nil, nil
	}

	pagination := &transformers.Pagination{
		Total:       int64(total),
		Count:       int64(Limit),
		PerPage:     int64(Limit),
		TotalPages:  int64(TotalPages),
		CurrentPage: int64(Page),
	}

	return drops, pagination
}
func (Drop *DropRepository) FetchFilter(TripStatus string, AssignType string, distance int, Page int, Limit int) ([]*models.Drop, *transformers.Pagination) {
	opts := options.Aggregate().SetMaxTime(10 * time.Second)
	var filter bson.A
	if TripStatus != "" {
		filter = append(filter, bson.D{{"cTP.tripStatus", bson.D{{"$eq", "Mencari Driver"}}}})
	}
	if AssignType != "" {
		filter = append(filter, bson.D{{"cTP.assignType", bson.D{{"$eq", "Manual Assign"}}}})
	}
	if distance != 0 {
		search := bson.A{
			bson.D{
				{
					"cTP.distance", bson.D{
						{
							"$eq", 18,
						},
					},
				},
			},
			bson.D{
				{
					"cTP.distance", bson.D{
						{
							"$eq", 20,
						},
					},
				},
			},
		}
		filter = bson.A{
			bson.D{{"$and", filter}},
			bson.D{{"$or", search}},
		}
	}
	rootStage := bson.D{
		{"$project",
			bson.D{
				{"customerTripPlanningDt", "$$ROOT"},
				{"_id", 0},
			},
		},
	}
	lookUpStage := bson.D{
		{"$lookup",
			bson.D{
				{"localField", "customerTripPlanningDt.tripIdObject"},
				{"from", "customerTripPlanning"},
				{"foreignField", "_id"},
				{"as", "cTP"},
			},
		},
	}
	unwindStage := bson.D{
		{"$unwind",
			bson.D{
				{"path", "$cTP"},
				{"preserveNullAndEmptyArrays", false},
			},
		},
	}
	match := bson.D{
		{"$match",
			bson.D{
				{"$and",
					filter,
				},
			},
		},
	}
	replaceStage := bson.D{
		{"$replaceRoot",
			bson.D{
				{"newRoot",
					bson.D{
						{"$mergeObjects",
							bson.A{
								"$customerTripPlanningDt",
								"$cTP",
								"$$ROOT",
							},
						},
					},
				},
			},
		},
	}
	projectionStage := bson.D{
		{"$project",
			bson.D{
				{"customerTripPlanningDt", 0},
				{"cTP", 0},
			},
		},
	}
	countStage := bson.D{{"$count", "count"}}
	limitStage := bson.D{{"$limit", Limit}}
	sortStage := bson.D{{"$sort", bson.D{{"created_at", -1}}}}
	countCursor, err := Drop.collection.Aggregate(context.Background(), mongo.Pipeline{rootStage, lookUpStage, unwindStage, match, countStage})
	var countRecord []bson.M
	if err = countCursor.All(context.Background(), &countRecord); err != nil {
		facades.Log().Error(err)
	}
	Skip, TotalPages := helpers.Paginate(Page, int(countRecord[0]["count"].(int32)), Limit)
	skipStage := bson.D{{"$skip", Skip}}
	defer func(countCursor *mongo.Cursor, ctx context.Context) {
		err := countCursor.Close(ctx)
		if err != nil {
			facades.Log().Error(err)
		}
	}(countCursor, context.Background())
	cursor, err := Drop.collection.Aggregate(context.Background(), mongo.Pipeline{rootStage, lookUpStage, unwindStage, match, replaceStage, projectionStage, skipStage, limitStage, sortStage}, opts)
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			facades.Log().Error()
		}
	}(cursor, context.Background())
	if err != nil {
		facades.Log().Error(err)
		return nil, nil
	}
	var Drops []*models.Drop
	if err := cursor.All(context.Background(), &Drops); err != nil {
		facades.Log().Error(err)
	}
	Pagination := &transformers.Pagination{
		Total:       int64(countRecord[0]["count"].(int32)),
		Count:       int64(Limit),
		PerPage:     int64(Limit),
		TotalPages:  int64(TotalPages),
		CurrentPage: int64(Page),
	}
	return Drops, Pagination
}
