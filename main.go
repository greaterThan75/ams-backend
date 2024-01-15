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

type User struct{
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"` 
}

func (r *Repository) CreateUser(context *fiber.Ctx) error{
	user := new(User)
	err:= context.BodyParser(&user)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message":"request failed"})
		return err
	}

	err = r.DB.Create(&user).Error
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message":"Could not an user account"})
		return err	
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message":"user account created successfully"})
	return nil

}

func (r *Repository) GetUsers(context *fiber.Ctx) error{
	Users := &[]models.User{}

	err := r.DB.Find(Users).Error

	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message":"Could not get users"})
		return err	
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message":"User created successfully",
	"data":Users})
	return nil
}

func (r *Repository) DeleteUser(context*fiber.Ctx) error {
	Users := &[]models.User{}

	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message":"Invalid id"},
		)
		return nil
	}

	err := r.DB.Delete(Users,id)

	if err.Error != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message":"Could not delete user"})
		return err.Error	
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message":"user deleted successfully"})
	return nil
}

func (r *Repository) GetUserById(context *fiber.Ctx) error {
	Users := &[]models.User{}

	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "Invalid id"},
		)
		return nil
	}
	fmt.Println("the ID is ", id)
	err := r.DB.Where("id = ?", id).Find(Users).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get User"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message":"User retrieved successfully","data":Users})
	return nil

}

func(r *Repository) SetupRoutes(app *fiber.App){
	api:= app.Group("/api/users")
	api.Post("/create",r.CreateUser)
	api.Delete("/delete/:id",r.DeleteUser)
	api.Get("/:id",r.GetUserById)
	api.Get("/",r.GetUsers)


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

	err = models.MigrateUsers(db)
	if err != nil {
		log.Fatal("Error migrating Users",err)
	}
	r:= Repository{
		DB: db,
	}	
	app:= fiber.New()
	r.SetupRoutes(app)
	app.Listen(":4000")	
}