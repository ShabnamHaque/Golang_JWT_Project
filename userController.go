package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ShabnamHaque/go-jwt/database"
	helper "github.com/ShabnamHaque/go-jwt/helpers"
	"github.com/ShabnamHaque/go-jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	//"gopkg.in/mgo.v2/bson"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

// collectionName="user"
var validate = validator.New()

func HashPassword(password string) string {
	//encrpyt the password provided
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func verifyPassword(userPass string, providedPass string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPass), []byte(userPass))
	//compare entered form of password and their hashvalues.
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("Email or password is incorrect")
		check = false
	}
	return check, msg
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil { //only admin has to be have this access.
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10 //by default we need 10 records per page/
		}
		page, err1 := strconv.Atoi(c.Query("page")) // Retrieve the value of the "page" query parameter from the URL

		if err1 != nil || page < 1 {
			page = 1 // by default one page
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		// aggregation in MongoDB!!
		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}},        // grouped by the unique identity - _id
			{"total_count", bson.D{{"$sum", 1}}},    //to see how many users
			{"data", bson.D{{"$push", "$$ROOT"}}}}}} // to not just show count but also show the users' data
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}}}}}
		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}
		var allusers []bson.M // a slice created to return.
		if err = result.All(ctx, &allusers); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allusers[0])
	}
}
func GetUser() gin.HandlerFunc {
	//Only for admin-use
	return func(c *gin.Context) {
		userId := c.Param("user_id") //to fetch user_id from the route.
		if err := helper.MatchUserTypetoUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // returns an error if not ADMIN or The USER him/herself.
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var User models.User //creating a variable user of structure models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&User)
		//decode from JSON because golang doesnot understand JSON.
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, User)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser) //stored in foundUser with email inputed.
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "entry not made"})
			return
		}

		passwordIsValid, msg := verifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid == false {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, *foundUser.User_type, foundUser.User_id)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id) // when the user logins from time to time,token is refreshed. Hence it needs to be updated.

		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//we are validating the data fields from the user struct.
		validationErr := validate.Struct(user) //from validator package - checks if key - values are intact with the struct defined
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		var count, err = userCollection.CountDocuments(ctx, bson.M{"email": user.Email}) //count users with same email
		defer cancel()
		if err != nil {
			log.Panic(err)
			ErrorData := gin.H{"error": "Error occurred while checking for email"}
			c.JSON(http.StatusInternalServerError, ErrorData)
		}
		if count > 0 { //because signup we need to maintain both email and phone as unique for each user..
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "This email or/and phone number already exists"})

		}
		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone}) //count documents with the entered phone we want unique phone and email
		defer cancel()
		if err != nil {
			log.Panic(err)
			ErrorData := gin.H{"error": "Error occurred while checking for Phone Number"}
			c.JSON(http.StatusInternalServerError, ErrorData)

		}
		if count > 0 { //because signup we need to maintain both email and phone as unique for each user..
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "This email or/and phone number already exists"})

		}

		// create Timestamps - while creating new accounts through Signups.
		// w e
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID() // format = BSON objectId
		user.User_id = (user.ID).Hex()    // converts from BSON to Hexadec str rep
		token, refreshToken, err := helper.GenerateAllTokens(*user.Email, *user.First_Name, *user.Last_Name, *user.User_type, user.User_id)
		user.Token = &token
		user.Password = &password
		user.Refresh_Token = &refreshToken
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not able to tokenise"})
			return
		}
		//we have created the new user.

		//now insert into the database
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

//go get github.com/go-playground/validator/v10
