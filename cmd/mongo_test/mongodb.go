package main

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	/*
		"go.mongodb.org/mongo-driver/mongo/readconcern"
		"go.mongodb.org/mongo-driver/mongo/writeconcern"
	*/
	"go_db/config"
)

type prop struct {
	Name string
	Val  string
}

type brand struct {
	Name    string
	Country string
}

type product struct {
	Id    string `bson:"_id"`
	Name  string
	Price int
	Props []prop
	Brand brand
	Tags  []string
}

func main() {
	fmt.Println("Go MONGO")
	ctx := context.Background()
	opt := options.Client().ApplyURI(config.GetMongoString())

	cl, err := mongo.Connect(ctx, opt)
	if err != nil {
		fmt.Println(err)
		return
	}

	filter := bson.D{}

	listDb, err := cl.ListDatabases(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("list db", listDb.Databases)

	cl.Database("Empl").Collection("Product")

	listCol, err := cl.Database("Empl").ListCollectionNames(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return
	}
	coll := cl.Database("Empl").Collection("Products")

	fmt.Println("list col", strings.Join(listCol, ","))
	coll.Drop(ctx)

	brand1 := brand{Name: "Br1", Country: "USA"}
	brand2 := brand{Name: "Br2", Country: "China"}

	p1 := product{Id: "1", Name: "kartoshka", Price: 91, Brand: brand1, Tags: []string{"A", "B"}}
	p1.Props = []prop{{Name: "color", Val: "red"}, {Name: "size", Val: "big"}}
	p2 := product{Id: "2", Name: "kartoshka", Price: 57, Brand: brand1, Tags: []string{"A"}}
	p2.Props = []prop{{Name: "color", Val: "white"}, {Name: "size", Val: "small"}}

	p3 := product{Id: "3", Name: "apple", Price: 31, Brand: brand2, Tags: []string{"B", "A"}}
	p3.Props = []prop{{Name: "color", Val: "red"}, {Name: "weight", Val: "1"}}
	listProducts := []interface{}{p1, p2}

	insertResult, err := coll.InsertOne(ctx, p3)
	if err != nil {
		fmt.Println("insertOneError:", err)
	} else {
		fmt.Println("insertedID", insertResult.InsertedID)
	}

	_, err = coll.InsertMany(ctx, listProducts)
	if err != nil {
		fmt.Println("insertManyError:", err)
	}

	cursor, err := coll.Find(ctx, filter)
	//err = cursor.Decode(&rx)
	printList(ctx, cursor, err)

	fmt.Println("NAME: kartoshka")
	filter = bson.D{{Key: "name", Value: "kartoshka"}}
	cursor, err = coll.Find(ctx, filter)
	printList(ctx, cursor, err)

	fmt.Println("Tags : A")
	filter = bson.D{{Key: "tags", Value: "B"}}
	cursor, err = coll.Find(ctx, filter)
	printList(ctx, cursor, err)

	fmt.Println("COLOR : RED")
	filter = bson.D{{"props",
		bson.D{{"$elemMatch", bson.D{{"val", "red"}, {"name", "color"}}}}}}
	cursor, err = coll.Find(ctx, filter)
	printList(ctx, cursor, err)

	fmt.Println("COLOR : RED2")
	filter = bson.D{{"props", bson.D{{"name", "color"}, {"val", "red"}}}}
	cursor, err = coll.Find(ctx, filter)
	printList(ctx, cursor, err)

	fmt.Println("COLOR : RED2 FIND ONE & PROJECTION")
	filter = bson.D{{Key: "name", Value: "kartoshka"}}
	projection := bson.D{{"name", 1}}
	opts := options.FindOne().SetProjection(projection)
	singleResult := coll.FindOne(ctx, filter, opts)
	err = singleResult.Err()
	if err != nil {
		fmt.Println(singleResult.Err().Error())
	} else {
		p4 := product{}
		singleResult.Decode(&p4)
		fmt.Println("FindOne:", p4)
	}

	fmt.Println("COLOR : RED3 wrong order")
	filter = bson.D{{"props", bson.D{{"val", "red"}, {"name", "color"}}}}
	cursor, err = coll.Find(ctx, filter)
	printList(ctx, cursor, err)

	fmt.Println("Brand : br2")
	filter2 := bson.M{"brand.name": "Br2"}
	//filter2 := bson.D{{"brand.name", "Br2"}}
	cursor, err = coll.Find(ctx, filter2)
	printList(ctx, cursor, err)

	/* Transaction numbers are only allowed on a replica set member or mongos */
	/* нужен набор реплик а не одиночный сервер что бы использовать транзакции */
	err = transaction(ctx, cl)
	if err != nil {
		fmt.Println("transactoin err:", err.Error())
	} else {
		fmt.Println("transaction: OK")
	}

}

func transaction(ctx context.Context, cl *mongo.Client) error {
	p0 := product{Name: "ogurec", Price: 99}
	/*
		wc := writeconcern.New(writeconcern.WMajority())
		rc := readconcern.Snapshot()
		txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	*/
	session, err := cl.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	coll := cl.Database("Empl").Collection("Products")

	trFunc := func(sessionContext mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			return err
		}
		_, err := coll.InsertOne(sessionContext, p0)
		if err != nil {
			return err
		}
		_, err = coll.InsertOne(sessionContext, p0)
		if err != nil {
			return err
		}

		if err = session.CommitTransaction(sessionContext); err != nil {
			return err
		}
		return nil
	}

	err = mongo.WithSession(ctx, session, trFunc)

	if err != nil {
		if abortErr := session.AbortTransaction(ctx); abortErr != nil {
			return abortErr
		}
		return err
	}
	return nil
}

func printList(ctx context.Context, cursor *mongo.Cursor, err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	for cursor.Next(ctx) {
		var elem product
		err := cursor.Decode(&elem)
		if err != nil {
			fmt.Println("ListErr:", err)
		}
		fmt.Println(elem)
	}
}
