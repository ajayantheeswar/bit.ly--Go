package linkHandler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ajayantheeswar/bit.ly/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Link struct {
	ID           primitive.ObjectID `json:"ID" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"UserID" bson:"user_id,omitempty"`
	Name         string             `json:"Name" bson:"name,omitempty"`
	Date         string             `json:"Date" bson:"date,omitempty"`
	ShortenedURL string             `json:"ShortenedURL" bson:"shortenedURL,omitempty"`
	OriginalURL  string             `json:"OriginalURL" bson:"originalURL,omitempty"`
	Count        int64              `json:"Count" bson:"count,omitempty"`
}

func CreateLink(c *gin.Context) {
	var InputRequest Link
	userId := c.Request.Header.Get("UserId")
	err := c.ShouldBindJSON(&InputRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	ctx := context.Background()
	result, err := database.Links.Find(ctx, bson.M{"shortenedURL": InputRequest.ShortenedURL})
	if result.Next(ctx) {
		c.JSON(http.StatusConflict, gin.H{"Message": "The Url Already Exists"})
		return
	}
	InputRequest.UserID, err = primitive.ObjectIDFromHex(userId)
	InsertedResult, err := database.Links.InsertOne(ctx, InputRequest)
	c.JSON(http.StatusAccepted, gin.H{"Message": InsertedResult.InsertedID.(primitive.ObjectID).Hex()})
}


type aggregateResult struct {
	ID string `bson:"_id"`
	Count int64 `bson:"count"`
}

func GetAllLinks(c *gin.Context) {
	var InputRequest Link
	userID := c.Request.Header.Get("UserId")
	err := c.ShouldBindJSON(&InputRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	ctx := context.Background()
	var Links []Link
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	

	result, err := database.Links.Find(ctx, bson.M{"user_id": userObjectID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, "err - 52")
		return
	}
	err = result.All(ctx, &Links)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "err - 53")
		return
	}
	var LinksSlice []string
	for i := 0; i < len(Links); i++ {
		LinksSlice = append(LinksSlice, Links[i].ShortenedURL)
	}

	matchStage := bson.D{{"$match", bson.D{{"shortenedURL", bson.D{{"$in", LinksSlice}}}}}}
	groupStage := bson.D{{"$group" , bson.D{{ "_id" ,"$shortenedURL"} , {"count" , bson.D{{ "$sum" , 1 }} }} }}
	sortStage  := bson.D{{"$sort" , bson.D{{"shortenedURL", 1}} }}

	result, err = database.Visits.Aggregate(ctx, mongo.Pipeline{matchStage , groupStage,sortStage})
	
	var countLinks[] aggregateResult

	if err != nil {
		panic(err)
	} else {
		result.All(ctx,&countLinks)
		for i :=0;i<len(countLinks) ;i++ {
			for j := 0 ; j<len(Links) ;j++ {
				if(countLinks[i].ID == Links[j].ShortenedURL){
					Links[j].Count = countLinks[i].Count
					continue
				}
			}	
		}

	} 
/*
	for i := 0; i < len(Links); i++ {
		resultcount,_ := database.Visits.CountDocuments(ctx,bson.M{"shortenedURL" : Links[i].ShortenedURL})
		Links[i].Count = resultcount
	} 
*/
	payload, err := json.Marshal(Links)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "err - 54")
		return
	}
	c.JSON(http.StatusOK, gin.H{"Length": len(Links), "List": string(payload)})
}

type Visit struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Time         string
	LinkID       primitive.ObjectID `bson:"linkID,omitempty"`
	ShortenedURL string             `bson:"shortenedURL,omitempty"`
	IP           string             `bson:"ip,omitempty"`
	Location     string             `bson:"locatoin,omitempty"`
}

func RedirectLink(c *gin.Context) {
	paramURL, _ := c.Params.Get("url")
	ctx := context.Background()
	result, err := database.Links.Find(ctx, bson.M{"shortenedURL": paramURL})

	if err != nil {
		c.JSON(http.StatusNotFound, "URL NOT FOUND")
		return
	}

	if !result.Next(ctx) {
		c.JSON(http.StatusNotFound, "URL NOT FOUND")
		return
	}

	var resultLink Link
	result.Decode(&resultLink)

	cCon := c.Copy()
	go func() {
		cCtx := context.Background()
		// Send the Data the Link Data to Server and Setup Cache
		ipValue := cCon.Request.RemoteAddr
		ip := strings.Split(ipValue, ":")[0]
		var visit = Visit{
			LinkID:       resultLink.ID,
			ShortenedURL: paramURL,
			Time:         strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10),
			IP:           ip,
			Location:     "Salem"}

		_, err = database.Visits.InsertOne(cCtx, &visit)

		if err != nil {
			log.Print(err)
		}
	}()

	c.Redirect(http.StatusPermanentRedirect, resultLink.OriginalURL)
}
