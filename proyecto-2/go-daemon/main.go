package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	procContFile = "/proc/continfo_so1_201222010"
	logFile      = "./logs/daemon.log"
	dbFile       = "./db/contenedores.sqlite"
	cronScript   = "../bash/generar_contenedores.sh"
	kernelScript = "../bash/cargar_modulo_kernel.sh"

	minAlto = 2
	minBajo = 3
	loopSec = 20
)

type Contenedor struct {
	PID     int
	Nombre  string
	Tipo    string // "alto" o "bajo"
	RAM     float64
	CPU     float64
	Cmdline string
}

// Ejecutar comando bash
func ejecutarComando(comando string, args ...string) error {
	cmd := exec.Command(comando, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error ejecutando %s %v: %v, output: %s", comando, args, err, string(output))
	}
	return err
}

// Crear cronjob
func crearCronJob() {
	cronCmd := fmt.Sprintf("* * * * * %s", cronScript)
	cmd := exec.Command("bash", "-c", fmt.Sprintf("(crontab -l 2>/dev/null; echo \"%s\") | crontab -", cronCmd))
	if err := cmd.Run(); err != nil {
		log.Printf("Error creando cronjob: %v", err)
	} else {
		log.Println("Cronjob creado:", cronCmd)
	}
}

// Eliminar cronjob
func eliminarCronJob() {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("crontab -l | grep -v '%s' | crontab -", cronScript))
	if err := cmd.Run(); err != nil {
		log.Printf("Error eliminando cronjob: %v", err)
	} else {
		log.Println("Cronjob eliminado")
	}
}

// Leer contenedores desde /proc
func leerContInfo() ([]Contenedor, error) {
	data, err := os.ReadFile(procContFile)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var conts []Contenedor
	for _, line := range lines[2:] {
		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}
		pid, _ := strconv.Atoi(fields[0])
		tipo := "bajo"
		if strings.Contains(fields[1], "alto") {
			tipo = "alto"
		}
		conts = append(conts, Contenedor{
			PID:     pid,
			Nombre:  fields[1],
			Tipo:    tipo,
			Cmdline: strings.Join(fields[6:], " "),
		})
	}
	return conts, nil
}

// Guardar log en SQLite
func guardarLog(db *sql.DB, cont Contenedor, accion string) {
	stmt, _ := db.Prepare("INSERT INTO logs(timestamp, contenedor, accion, cpu, ram) VALUES (?, ?, ?, ?, ?)")
	_, err := stmt.Exec(time.Now(), cont.Nombre, accion, cont.CPU, cont.RAM)
	if err != nil {
		log.Printf("Error guardando log: %v", err)
	}
}

// Eliminar contenedores sobrantes
func eliminarSobrantes(conts []Contenedor) {
	// Contar contenedores por tipo
	var alto, bajo []Contenedor
	for _, c := range conts {
		if c.Tipo == "alto" {
			alto = append(alto, c)
		} else {
			bajo = append(bajo, c)
		}
	}

	// Eliminar alto sobrante
	if len(alto) > minAlto {
		for _, c := range alto[minAlto:] {
			ejecutarComando("docker", "rm", "-f", c.Nombre)
			log.Println("Contenedor alto eliminado:", c.Nombre)
		}
	}

	// Eliminar bajo sobrante
	if len(bajo) > minBajo {
		for _, c := range bajo[minBajo:] {
			ejecutarComando("docker", "rm", "-f", c.Nombre)
			log.Println("Contenedor bajo eliminado:", c.Nombre)
		}
	}
}

func main() {
	// Inicializar logging
	os.MkdirAll("./logs", 0755)
	logFileObj, _ := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer logFileObj.Close()
	log.SetOutput(logFileObj)

	// Inicializar DB
	os.MkdirAll("./db", 0755)
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Error abriendo DB: %v", err)
	}
	defer db.Close()
	db.Exec(`CREATE TABLE IF NOT EXISTS logs(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME,
		contenedor TEXT,
		accion TEXT,
		cpu REAL,
		ram REAL
	)`)

	log.Println("Daemon iniciado")

	// Señales para terminar
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Daemon finalizando...")
		eliminarCronJob()

		// Eliminar todos los contenedores creados por el daemon (excepto Grafana)
		conts, err := leerContInfo()
		if err == nil {
			for _, c := range conts {
				ejecutarComando("docker", "rm", "-f", c.Nombre)
				log.Println("Contenedor eliminado al cerrar daemon:", c.Nombre)
			}
		}

		os.Exit(0)
	}()

	// 1. Cargar módulos kernel
	ejecutarComando("bash", kernelScript)

	// 2. Crear cronjob
	crearCronJob()

	// 3. Loop principal
	for {
		conts, err := leerContInfo()
		if err != nil {
			log.Printf("Error leyendo /proc: %v", err)
			time.Sleep(loopSec * time.Second)
			continue
		}

		// Contar y crear contenedores faltantes
		countAlto, countBajo := 0, 0
		for _, c := range conts {
			if c.Tipo == "alto" {
				countAlto++
			} else {
				countBajo++
			}
		}

		// Crear alto consumo faltante
		for i := countAlto; i < minAlto; i++ {
			nombre := fmt.Sprintf("alto_%d", rand.Intn(10000))
			ejecutarComando("docker", "run", "-d", "--name", nombre, "alto_cpu_ram_1")
			log.Println("Contenedor alto creado:", nombre)
		}

		// Crear bajo consumo faltante
		for i := countBajo; i < minBajo; i++ {
			nombre := fmt.Sprintf("bajo_%d", rand.Intn(10000))
			ejecutarComando("docker", "run", "-d", "--name", nombre, "bajo_consumo")
			log.Println("Contenedor bajo creado:", nombre)
		}

		// Eliminar sobrantes
		eliminarSobrantes(conts)

		// Guardar logs
		for _, cont := range conts {
			guardarLog(db, cont, "existente")
		}

		time.Sleep(loopSec * time.Second)
	}
}
