package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

const (
	Migration = `
		CREATE TABLE IF NOT EXISTS VAULT_NOTES (
			id SERIAL PRIMARY KEY,
			author TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at timestamp with time zone DEFAULT current_timestamp
		)
	`
)

var db *pgxpool.Pool

type Notes struct {
	Author    string    `json:"author" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
}

func GetAllNotes() ([]Notes, error) {
	const query = `SELECT author, content, created_at FROM VAULT_NOTES ORDER BY created_at DESC LIMIT 100`

	var (
		author    string
		content   string
		createdAt time.Time
	)

	rows, err := db.Query(context.Background(), query)
	notes := make([]Notes, 0)

	for rows.Next() {
		rows.Scan(&author, &content, &createdAt)
		notes = append(notes, Notes{author, content, createdAt})
	}

	if err != nil {
		return nil, err
	}

	return notes, nil
}

func GetNoteById(id int) ([]Notes, error) {

	const query = `SELECT author, content, created_at FROM VAULT_NOTES WHERE id = $1`

	var (
		author    string
		content   string
		createdAt time.Time
	)

	rows, err := db.Query(context.Background(), query, id)
	notes := make([]Notes, 0)

	for rows.Next() {
		rows.Scan(&author, &content, &createdAt)
		notes = append(notes, Notes{author, content, createdAt})
	}

	if err != nil {
		return nil, err
	}

	return notes, nil
}

func DeleteNoteById(id int) error {
	const query = `DELETE FROM VAULT_NOTES WHERE id = $1`

	_, err := db.Exec(context.Background(), query, id)

	if err != nil {
		return err
	}

	return nil

}

func AddNewNote(author string, content string) error {
	const query = `INSERT INTO VAULT_NOTES (author, content) VALUES ($1, $2)`

	_, err := db.Exec(context.Background(), query, author, content)

	if err != nil {
		return err
	}

	return nil
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbStr := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err = pgxpool.Connect(context.Background(), dbStr)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	_, err = db.Query(context.Background(), Migration)

	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
		return
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		c.SendString(`
		GET request to /api/note to get all the latest notes
		GET request to /api/note/id to get a specific note
		POST request to /api/note to put a new note
		DELETE request to/api/note to delete a specific note
		`)

		return nil
	})

	app.Get("/api/note", func(c *fiber.Ctx) error {
		data, err := GetAllNotes()

		if err != nil {
			c.JSON(map[string]string{
				"message": "Failed to get all notes",
				"error":   err.Error(),
			})
		}

		c.JSON(data)

		return nil
	})

	app.Get("/api/note/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))

		if err != nil {
			log.Fatalf("Failed to convert parameter %v", err)
		}

		data, err := GetNoteById(id)

		if err != nil {
			c.JSON(map[string]string{
				"message": "Failed to get note with id of " + c.Params("id"),
				"error":   err.Error(),
			})
		}

		c.JSON(data)

		return nil
	})

	app.Post("/api/note/", func(c *fiber.Ctx) error {
		author := c.Query("author")

		content := c.Query("content")

		if err := AddNewNote(author, content); err != nil {
			c.JSON(map[string]string{
				"message": "failed to add a new note",
				"error":   err.Error(),
			})
		}

		c.JSON(map[string]string{
			"message": "Note Created",
		})

		return nil
	})

	app.Delete("/api/note", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Query("id"))

		if err != nil {
			log.Fatalf("Failed to convert parameter %v", err)
		}

		if err := DeleteNoteById(id); err != nil {
			c.JSON(map[string]string{
				"message": "Failed to get note with id of " + c.Query("id"),
			})
		}

		c.JSON(map[string]string{
			"message": "Note of ID " + c.Query("id") + " has been deleted",
		})

		return nil
	})

	log.Fatal(app.Listen(":8080"))
}
