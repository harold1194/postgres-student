package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/harold/postgres-student/models"
	"github.com/harold/postgres-student/storage"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Student struct {
	FullName string `json:"fullname"`
	Address  string `json:"address"`
	Degree   string `json:"degree"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateStudent(context *fiber.Ctx) error {
	student := Student{}

	err := context.BodyParser(&student)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&student).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create student"})
		return nil
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "student has been added"})
	return nil
}

func (r *Repository) DeleteStudent(context *fiber.Ctx) error {
	studentModel :=
		models.Student{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(studentModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete student",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "student deleted successfully",
	})
	return nil
}

func (r *Repository) GetStudent(context *fiber.Ctx) error {
	studentModel := &[]models.Student{}

	err := r.DB.Find(studentModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not get student data",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "student successfully fetch", "data": studentModel})
	return nil
}

func (r *Repository) GetStudentByID(context *fiber.Ctx) error {
	id := context.Params("id")
	studentModel := &models.Student{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot found",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(studentModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not get student",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "student id fetched successfully",
		"data":    studentModel,
	})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_students", r.CreateStudent)
	api.Delete("delete_student/:id", r.DeleteStudent)
	api.Get("/get_students/:id", r.GetStudentByID)
	api.Get("/students", r.GetStudent)

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}
	err = models.MigrateStudents(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
