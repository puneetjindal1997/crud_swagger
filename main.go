package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	docs "goguru/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Welcome to your channel go guruji

// CRUD => create, read, update, delete

// swagger=>
// what is swagger?

type Data struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" `
	Name  string             `json:"name" bson:"name"`
	Email string             `json:"email" bson:"email"`
}

type Data2 struct {
	ID    string `json:"id,omitempty" bson:"id,omitempty" `
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}

type manager struct {
	connection *mongo.Client
	ctx        context.Context
	cancel     context.CancelFunc
}

var Mgr Manager

type Manager interface {
	Insert(interface{}) error
	GetAll() ([]Data, error)
	DeleteData(primitive.ObjectID) error
	UpdateData(Data) error
}

func connectDb() {
	uri := "localhost:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("%s%s", "mongodb://", uri)))
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connected!!!")
	Mgr = &manager{connection: client, ctx: ctx, cancel: cancel}
}

func close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {
	defer cancel()

	defer func() {

		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func init() {
	connectDb()
}

func main() {
	r := gin.Default()

	//Routes for swagger
	swagger := r.Group("swagger")
	{
		docs.SwaggerInfo.Title = "CRUD"
		docs.SwaggerInfo.Description = "Some description"
		docs.SwaggerInfo.Version = "1"

		swagger.GET("/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	r.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	r.POST("/data", insertData)
	r.GET("/data1", getAll)
	r.DELETE("/data", deleteData)
	r.PUT("/data", updateData)
	r.Run(":9090")

}

func insertData(c *gin.Context) {
	var d Data
	err := c.BindJSON(&d)

	if err != nil {
		fmt.Println(err)
		return
	}
	Mgr.Insert(d)
	c.JSON(http.StatusOK, gin.H{"message": d})
}

func getAll(c *gin.Context) {
	data, err := Mgr.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": data})
}

func deleteData(c *gin.Context) {
	id := c.Query("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = Mgr.DeleteData(objectId)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted success"})
}

func updateData(c *gin.Context) {
	var d Data2
	var final Data
	err := c.BindJSON(&d)
	if err != nil {
		fmt.Println(err)
		return
	}
	objectId, err := primitive.ObjectIDFromHex(d.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	final.ID = objectId
	final.Name = d.Name
	final.Email = d.Email
	err = Mgr.UpdateData(final)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated success"})
}

func (mgr *manager) Insert(data interface{}) error {
	orgCollection := mgr.connection.Database("goguru").Collection("collectiongoguru")
	result, err := orgCollection.InsertOne(context.TODO(), data)
	fmt.Println(result.InsertedID)
	return err
}

func (mgr *manager) GetAll() (data []Data, err error) {

	orgCollection := mgr.connection.Database("goguru").Collection("collectiongoguru")

	// Pass these options to the Find method
	findOptions := options.Find()

	cur, err := orgCollection.Find(context.TODO(), bson.M{}, findOptions)
	for cur.Next(context.TODO()) {
		var d Data
		err := cur.Decode(&d)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, d)
	} // close for

	if err := cur.Err(); err != nil {
		return nil, err
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	return data, nil
}

func (mgr *manager) DeleteData(id primitive.ObjectID) error {
	orgCollection := mgr.connection.Database("goguru").Collection("collectiongoguru")

	filter := bson.D{{"_id", id}}
	_, err := orgCollection.DeleteOne(context.TODO(), filter)
	return err
}

func (mgr *manager) UpdateData(data Data) error {
	orgCollection := mgr.connection.Database("goguru").Collection("collectiongoguru")

	filter := bson.D{{"_id", data.ID}}
	update := bson.D{{"$set", data}}

	_, err := orgCollection.UpdateOne(context.TODO(), filter, update)

	return err

}
