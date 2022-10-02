package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	db "github.com/MahadevMaahi/gin-jwt-authentication/database"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email			string
	First_name		string
	Last_name		string
	Uid				string
	User_type 		string
	jwt.StandardClaims
}

var userDbClient db.DbKit = &db.DbClient{}
var userCollection = userDbClient.OpenCollection(db.Client, "user")

var ACCESS_KEY string = os.Getenv("ACCESS_KEY")
var REFRESH_KEY string = os.Getenv("REFRESH_KEY")

func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid string) (signedToken string, signedRefereshToken string, err error) {
	fmt.Println("Generating tokens")
	claims := &SignedDetails{
		Email : email,
		First_name : firstName,
		Last_name : lastName,
		User_type : userType,
		Uid : uid,
		StandardClaims : jwt.StandardClaims {
			ExpiresAt : time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims : jwt.StandardClaims {
			ExpiresAt : time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(ACCESS_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(REFRESH_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	fmt.Println("Validating token")
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(ACCESS_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("Token is expired")
		msg = err.Error()
		return
	}
	return claims, msg
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	fmt.Println("Updating All Tokens")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	
	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})
	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", Updated_at})

	upsert := true
	filter := bson.M{"user_id" : userId}

	opt := options.UpdateOptions{
		Upsert : &upsert,
	}

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj},}, &opt)
	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
}