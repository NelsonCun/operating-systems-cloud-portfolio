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
		return c.SendString("Hola, responde la API3: API3 en la VM2, desarrollada por el estudiante Nelson Cún con carnet: 201222010")
	})

	// Endpoint para llamar a api1
	app.Get("/api3/201222010/llamar-a-api1", func(c *fiber.Ctx) error {
		// Hacer petición a API1
		resp, err := http.Get("http://192.168.122.11:8081/")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error al llamar a API1: %v", err))
		}
		defer resp.Body.Close()

		// Leer respuesta de API1
		buf := make([]byte, 1024)
		n, _ := resp.Body.Read(buf)
		api1Response := string(buf[:n])

		return c.SendString(fmt.Sprintf("Respuesta de API1: %s", api1Response))
	})

	// Endpoint para llamar a api2
	app.Get("/api3/201222010/llamar-a-api2", func(c *fiber.Ctx) error {
		// Hacer petición a API2
		resp, err := http.Get("http://192.168.122.11:8082/")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error al llamar a API2: %v", err))
		}
		defer resp.Body.Close()

		// Leer respuesta de API3
		buf := make([]byte, 1024)
		n, _ := resp.Body.Read(buf)
		api2Response := string(buf[:n])

		return c.SendString(fmt.Sprintf("Respuesta de API3: %s", api2Response))
	})

	// Iniciar servidor en el puerto 8081
	log.Fatal(app.Listen("0.0.0.0:8083"))

}
