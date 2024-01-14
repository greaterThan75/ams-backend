package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dwocOrg/student-app-nitt/models"
	"github.com/dwocOrg/student-app-nitt/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Repository struct{
	DB *gorm.DB
}

type Book struct{
	Author string `json:"author"`
	Title string `json:"title"`
	Publisher string `json:"publisher"` 
}

func (r *Repository) CreateBook(context *fiber.Ctx) error{
	book := new(Book)
	err:= context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message":"request failed"})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message":"Could not create book"})
		return err	
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message":"Book created successfully"})
	return nil

}

func (r *Repository) GetBooks(context *fiber.Ctx) error{
	bookModels := &[]models.Book{}

	err := r.DB.Find(bookModels).Error

	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message":"Could not get books"})
		return err	
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message":"Book created successfully",
	"data":bookModels})
	return nil
}

func (r *Repository) DeleteBook(context*fiber.Ctx) error {
	bookModel:= models.Book{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message":"Invalid id"},
		)
		return nil
	}

	err := r.DB.Delete(bookModel,id)

	if err.Error != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message":"Could not delete book"})
		return err.Error	
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message":"Book deleted successfully"})
	return nil
}

func (r *Repository) GetBookById(context *fiber.Ctx) error {

	id := context.Params("id")
	bookModel := &models.Book{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "Invalid id"},
		)
		return nil
	}
	fmt.Println("the ID is ", id)
	err := r.DB.Where("id = ?", id).Find(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message":"Book retrieved successfully","data":bookModel})
	return nil

}

func(r *Repository) SetupRoutes(app *fiber.App){
	api:= app.Group("/api")
	api.Post("/create_books",r.CreateBook)
	api.Delete("/delete_book/:id",r.DeleteBook)
	api.Get("/get_books/:id",r.GetBookById)
	api.Get("/books",r.GetBooks)


}

func main(){
	err:= godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file",err)
	}
	config:= &storage.Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("DB_SSLMODE"),
		
	}
	db,err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Error connecting to database",err)
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Error migrating books",err)
	}
	r:= Repository{
		DB: db,
	}	
	app:= fiber.New()
	r.SetupRoutes(app)
	app.Listen(":4000")	
}