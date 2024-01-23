package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"example.com/configs"
	"example.com/models"
	"example.com/responses"

	"github.com/go-playground/validator/v10" // Importing the validator package for data validation
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive" // Importing the MongoDB primitive package
	"go.mongodb.org/mongo-driver/mongo"          // Importing the MongoDB driver
)

// Retrieving the MongoDB collection for users using the configs package
var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

// Creating a validator instance for data validation
var validate = validator.New()


// CreateUser is an HTTP handler function for creating a new user
func CreateUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()
	
		//validate the request body
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
	
		//use the validator library to validate required fields
		if validationErr := validate.Struct(&user); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
	
		newUser := models.User{
			Id:       primitive.NewObjectID(),
			Name:     user.Name,
			Location: user.Location,
			Title:    user.Title,
		}
		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
	
		rw.WriteHeader(http.StatusCreated)
		response := responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(rw).Encode(response)
	}
	}

// GetAUser is an HTTP handler function for retrieving a user by ID
func GetAUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Setting a timeout for the context to prevent hanging requests
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// Extracting parameters (URL variables) from the request
		params := mux.Vars(r)
		userId := params["userId"]

		// Declaring a variable to hold the user data
		var user models.User

		// Deferring the cancellation of the context to ensure it is canceled when the function returns
		defer cancel()

		// Converting the user ID string to a MongoDB ObjectID
		objId, _ := primitive.ObjectIDFromHex(userId)

		// Querying the userCollection to find a user by the specified ObjectID
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
		if err != nil {
			// Handling internal server errors with a 500 status code and providing an error response
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Handling successful retrieval with a 200 status code and providing a success response
		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}}
		json.NewEncoder(rw).Encode(response)
	}
}

// EditAUser is an HTTP handler function for editing a user by ID
func EditAUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Setting a timeout for the context to prevent hanging requests
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// Extracting parameters (URL variables) from the request
		params := mux.Vars(r)
		userId := params["userId"]

		// Declaring a variable to hold the user data
		var user models.User

		// Deferring the cancellation of the context to ensure it is canceled when the function returns
		defer cancel()

		// Converting the user ID string to a MongoDB ObjectID
		objId, _ := primitive.ObjectIDFromHex(userId)

		// Validating the request body by decoding it into the user variable
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			// Handling bad requests with a 400 status code and providing an error response
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Using the validator library to validate required fields
		if validationErr := validate.Struct(&user); validationErr != nil {
			// Handling bad requests with a 400 status code and providing an error response
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Constructing an update document with fields to be modified
		update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}

		// Updating the user in the MongoDB collection based on the user ID
		result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			// Handling internal server errors with a 500 status code and providing an error response
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Initializing a variable to hold updated user details
		var updatedUser models.User

		// Checking if the update was successful and retrieving the updated user details
		if result.MatchedCount == 1 {
			err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

			if err != nil {
				// Handling internal server errors with a 500 status code and providing an error response
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(rw).Encode(response)
				return
			}
		}

		// Handling successful edit with a 200 status code and providing a success response
		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedUser}}
		json.NewEncoder(rw).Encode(response)
	}
}


// DeleteAUser is an HTTP handler function for deleting a user by ID
func DeleteAUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Setting a timeout for the context to prevent hanging requests
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// Extracting parameters (URL variables) from the request
		params := mux.Vars(r)
		userId := params["userId"]

		// Deferring the cancellation of the context to ensure it is canceled when the function returns
		defer cancel()

		// Converting the user ID string to a MongoDB ObjectID
		objId, _ := primitive.ObjectIDFromHex(userId)

		// Deleting the user from the MongoDB collection based on the user ID
		result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})

		// Handling errors during the deletion operation
		if err != nil {
			// Handling internal server errors with a 500 status code and providing an error response
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Checking if the user was found and deleted successfully
		if result.DeletedCount < 1 {
			// Handling not found errors with a 404 status code and providing an error response
			rw.WriteHeader(http.StatusNotFound)
			response := responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Handling successful deletion with a 200 status code and providing a success response
		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}}
		json.NewEncoder(rw).Encode(response)
	}
}

// GetAllUser is an HTTP handler function for retrieving all users
func GetAllUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Setting a timeout for the context to prevent hanging requests
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// Declaring a slice to hold multiple user data
		var users []models.User

		// Deferring the cancellation of the context to ensure it is canceled when the function returns
		defer cancel()

		// Querying the userCollection to find all users in the MongoDB collection
		results, err := userCollection.Find(ctx, bson.M{})

		// Handling errors during the query operation
		if err != nil {
			// Handling internal server errors with a 500 status code and providing an error response
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Reading from the database in an optimal way and decoding each user
		defer results.Close(ctx)
		for results.Next(ctx) {
			// Declaring a variable to hold a single user's data
			var singleUser models.User

			// Decoding the user data from the MongoDB result
			if err = results.Decode(&singleUser); err != nil {
				// Handling internal server errors with a 500 status code and providing an error response
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(rw).Encode(response)
			}

			// Appending the single user data to the slice of users
			users = append(users, singleUser)
		}

		// Handling successful retrieval with a 200 status code and providing a success response
		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": users}}
		json.NewEncoder(rw).Encode(response)
	}
}

