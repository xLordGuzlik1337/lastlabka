package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/supabase-community/supabase-go"
)

var supabaseClient *supabase.Client

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func main() {
	supabaseURL := "https://wbkqeakwfrhqrsuafakq.supabase.co"
	supabaseKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6India3FlYWt3ZnJocXJzdWFmYWtxIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MzgyNjY1NTIsImV4cCI6MjA1Mzg0MjU1Mn0.kTeIU_bUFAYepxxkbju4FbZxGcYiJcOJsk5Fgd8yMU4"

	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		log.Fatalf("Ошибка подключения к Supabase: %v", err)
	}
	supabaseClient = client

	app := fiber.New()

	
	app.Static("/", "./")

	
	app.Get("/users", func(c *fiber.Ctx) error {
		var users []User
		data, _, err := supabaseClient.From("users").Select("*", "", false).Execute()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		err = json.Unmarshal(data, &users)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Ошибка обработки данных"})
		}

		return c.JSON(users)
	})


	app.Post("/users", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Некорректные данные"})
		}

		_, _, err := supabaseClient.From("users").Insert(user, false, "", "", "").Execute()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Пользователь добавлен"})
	})


	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		_, _, err := supabaseClient.From("users").Delete("", "").Eq("id", id).Execute()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Пользователь удален"})
	})

	log.Fatal(app.Listen(":3001"))
}
