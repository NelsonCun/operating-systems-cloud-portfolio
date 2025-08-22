package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Endpoint principal de api1

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hola, responde la API1: API1 en la VM1, desarrollada por el estudiante Nelson Cún con carnet: 201222010")
	})

	// Endpoint para llamar a api2
	app.Get("/api1/201222010/llamar-a-api2", func(c *fiber.Ctx) error {
		// Hacer petición a API2
		resp, err := http.Get("http://192.168.122.11:8082/")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error al llamar a API2: %v", err))
		}
		defer resp.Body.Close()

		// Leer respuesta de API2
		buf := make([]byte, 1024)
		n, _ := resp.Body.Read(buf)
		api2Response := string(buf[:n])

		return c.SendString(fmt.Sprintf("Respuesta de API2: %s", api2Response))
	})

	// Endpoint para llamar a api3
	app.Get("/api1/201222010/llamar-a-api3", func(c *fiber.Ctx) error {
		// Hacer petición a API3
		resp, err := http.Get("http://192.168.122.12:8083/")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error al llamar a API3: %v", err))
		}
		defer resp.Body.Close()

		// Leer respuesta de API3
		buf := make([]byte, 1024)
		n, _ := resp.Body.Read(buf)
		api3Response := string(buf[:n])

		return c.SendString(fmt.Sprintf("Respuesta de API3: %s", api3Response))
	})

	// Iniciar servidor en el puerto 8081
	log.Fatal(app.Listen("0.0.0.0:8081"))

}
