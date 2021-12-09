package main

import (
	"context"
	"encoding/json"
	"fmt"
	"link-shortner/src/database/models"
	"link-shortner/src/structs"
	"link-shortner/src/utils"
	"link-shortner/src/validation"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

type Handler struct {
	http.Handler
	Client *mongo.Client
}

func (handler *Handler) CreateLink(w http.ResponseWriter, r *http.Request) {
	body := utils.ReadBody(r.Body)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var parsedBody structs.CreateLink
	err := json.Unmarshal(body, &parsedBody)
	if err != nil {
		utils.ThrowError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("error while reading the request body: %v", err),
		)
		return
	}

	errors := validation.Validate(parsedBody)
	if len(errors) > 0 {
		utils.ThrowError(w, http.StatusBadRequest, errors...)
		return
	}

	if strings.Trim(parsedBody.Name, " ") == "" {
		parsedBody.Name = randstr.String(10)
	}
	collection := handler.Client.Database("link-shortner").Collection("links")

	var doc models.Link
	possibleDoc := collection.FindOne(ctx, struct{ Name string }{parsedBody.Name})
	if possibleDoc.Decode(&doc) == nil {
		utils.ThrowError(w, http.StatusUnauthorized, "A link with this name already exist")
	}

	insertedDoc := models.Link{
		Id:   primitive.NewObjectID(),
		Name: parsedBody.Name,
		URL:  parsedBody.URL,
	}

	collection.InsertOne(ctx, insertedDoc)
	json.NewEncoder(w).Encode(insertedDoc)

}

func (handler *Handler) GetLink(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	splitedURL := strings.Split(r.URL.Path, "/")
	name := splitedURL[len(splitedURL)-1]

	collection := handler.Client.Database("link-shortner").Collection("links")

	var doc models.Link

	result := collection.FindOne(ctx, struct{ Name string }{name})
	if result.Decode(&doc) != nil {
		utils.ThrowError(w, http.StatusBadRequest, fmt.Sprintf("couldn't find a link named '%s'", name))
		return
	}

	json.NewEncoder(w).Encode(struct{ URL string }{doc.URL})
}

func NewHandler(client *mongo.Client) *Handler {
	router := chi.NewRouter()
	handler := &Handler{
		Handler: router,
		Client:  client,
	}

	router.Post("/", handler.CreateLink)
	router.Get("/{name}", handler.GetLink)

	return handler
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	Client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DATABASE_URL")))
	if err != nil {
		log.Fatal(err)
	}

	defer Client.Disconnect(ctx)

	router := chi.NewRouter()
	router.Mount("/", NewHandler(Client))
	router.Handle("/", NewHandler(Client))

	err = http.ListenAndServe(":3000", router)
	if err != nil {
		panic(err)
	}
}
