package router

import (
	"github.com/abdou-1614/go-rest-api/common"
	"github.com/abdou-1614/go-rest-api/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddTodoGroupe(app *fiber.App) {
	booksGroup := app.Group("/book")

	booksGroup.Get("/", GetBooks)
	booksGroup.Get("/:id", GetBookById)
	booksGroup.Post("/", createBook)
	booksGroup.Put("update/:id", updateBook)
	booksGroup.Delete("delete/:id", deleteBook)
}

func GetBooks(c *fiber.Ctx) error {
	cell := common.GetDbCollection("book")

	books := make([]model.Book, 0)

	cursor, err := cell.Find(c.Context(), bson.M{})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"ERROR": err.Error(),
		})
	}

	for cursor.Next(c.Context()) {
		book := model.Book{}
		err := cursor.Decode(&book)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"ERROR": err.Error(),
			})
		}

		books = append(books, book)
	}

	return c.Status(200).JSON(fiber.Map{"data": books})
}

func GetBookById(c *fiber.Ctx) error {
	coll := common.GetDbCollection("book")

	id := c.Params("id")

	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"ERROR": "ID Is Required",
		})
	}
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"ERROR": "INVALID ID",
		})
	}

	book := model.Book{}

	err = coll.FindOne(c.Context(), bson.M{"_id": objectId}).Decode(&book)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"ERROR": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{"DATA": book})
}

type createDTO struct {
	Title  string `json:"title" bson:"title" validate:"required"`
	Author string `json:"author" bson:"author" validate:"required"`
	Year   string `json:"year" bson:"year" validate:"required"`
}

func createBook(c *fiber.Ctx) error {
	b := new(createDTO)

	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"ERROR": "Invalid Body",
		})
	}

	validate := validator.New()

	if err := validate.Struct(b); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Message": "Validation Error",
			"ERROR":   err.Error(),
		})
	}

	coll := common.GetDbCollection("book")

	result, err := coll.InsertOne(c.Context(), b)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Message": "Failed To Create",
			"ERROR":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"Message": "Book Created Successfuly",
		"DATA":    result,
	})
}

type updateDto struct {
	Title  string `json:"title,omitempty" bson:"title,omitempty"`
	Author string `json:"author,omitempty" bson:"author,omitempty"`
	Year   string `json:"year,omitempty" bson:"year,omitempty"`
}

func updateBook(c *fiber.Ctx) error {
	b := new(updateDto)

	if err := c.BodyParser(b); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"ERROR": "Invalid Body",
		})
	}

	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"ERROR": "ID Is Required",
		})
	}

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"ERROR": "Invalid ID",
		})
	}

	coll := common.GetDbCollection("book")

	result, err := coll.UpdateOne(c.Context(), bson.M{"_id": objectId}, bson.M{"$set": b})

	if err != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{
			"Message": err.Error(),
			"ERROR":   "Failed To Update BOOK",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Result": result,
	})

}

func deleteBook(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ERROR": "ID is Required",
		})
	}

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ERROR": "INVALID ID",
		})
	}

	coll := common.GetDbCollection("book")

	result, err := coll.DeleteOne(c.Context(), bson.M{"_id": objectId})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Message": "Failed To Delete Book",
			"ERROR":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Result": result,
	})
}
