package authHandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/ajayantheeswar/bit.ly/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name,omitempty"`
	Email    string             `json:"email" bson:"email,omitempty"`
	Password string             `json:"password" bson:"password,omitempty"`
}

type POJO_AuthSignIn_Request struct{
	Email    string             `json:"email" bson:"email,omitempty"`
	Password string             `json:"password" bson:"password,omitempty"`
}



func CreateToken(Id string) (string, error) {
	var err error

	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") 
	Claims := jwt.MapClaims{}
	Claims["user_id"] = Id

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
	   return "", err
	}
	return token, nil
  }


func AuthSignup(c *gin.Context) {
	var requestBody User
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		panic(err)
	}

	// Check if the user Already Present
	result, err := database.Users.Find(context.Background(), bson.M{"email": requestBody.Email})

	if err != nil {
		fmt.Print("Error in cursot",err)
		result.Close(context.Background())
	}

	for result.Next(context.Background()) {
		c.JSON(http.StatusConflict, gin.H{"Message": "User Already Registered"})
		return
	}

	if InsertedResult, err := database.Users.InsertOne(context.Background(), &requestBody); err != nil {
		c.JSON(200, gin.H{"msg": "Sorry , the Url not Found"})
	} else {
		token, err:= CreateToken(InsertedResult.InsertedID.(primitive.ObjectID).Hex())
		if err != nil{
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"message": token , "email" : requestBody.Email,"name" : requestBody.Name})
	}
}


func AuthSignIn(c *gin.Context){
	ctx := context.Background()
	var AuthSignInRequest POJO_AuthSignIn_Request
	var user User
	c.ShouldBindJSON(&AuthSignInRequest)


	result, err := database.Users.Find(ctx, bson.M{"email": AuthSignInRequest.Email})
	if err != nil {
		fmt.Print("Error in cursot",err)
		result.Close(ctx)
	} 

	// Email Does not Exists
	if !result.Next(ctx) {
		c.JSON(http.StatusConflict, gin.H{"Message": "Email Does Not Exist"})
		return
	}

	err = result.Decode(&user)
	if err !=  nil {
		panic(err)
	}
	if user.Password == AuthSignInRequest.Password {
		token,err := CreateToken(user.ID.Hex())
		if err != nil{
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"message": token , "email" : user.Email,"name" : user.Name})
	}else{
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid Credientials"})
	}
	
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		AuthHeader, err := GetAuthHeader(c.Request)
		userId,err :=  getAuthClaims(AuthHeader)
		userObjectID, err := primitive.ObjectIDFromHex(userId)
		ctx := context.TODO()
		result, err := database.Links.Find(ctx, bson.M{"_id": userObjectID})
		
		if err != nil || result.Next(ctx) {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Request.Header.Add("UserId",userId)
		c.Next()
	}
}

func GetAuthHeader(req *http.Request) (string,error) {
	bearToken := req.Header.Get("Authorization")

  	if bearToken == "" {
     	return "" , errors.New("Authentication Header Not Found")
	}
	return bearToken , nil  

  	
}

func getAuthClaims (token string) (string,error){   
	claims := jwt.MapClaims{}

	Parsedtoken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("jdnfksdmfksd"), nil
	})
	// ... error handling
	if err != nil{
		return "" ,errors.New("Invalid Token")
	}
	if Parsedtoken.Valid {
		return claims["user_id"].(string) , nil
	}else{
		return "" , errors.New("Invlalid Token")
	}
}


