package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

const scopePosition int = 2

// Esta aplicación es un proxy para llamar al binario específico basado en el valor de la variable SCOPE.
// El formato de SCOPE es: [nombre de aplicación] o [nombre de entorno]-[nombre de aplicación]
// Una vez que se obtiene el SCOPE, se cambia al directorio de la aplicación y se ejecuta un binario llamado "app".
// La salida estándar y el error de la nueva aplicación se redirigen a la salida estándar y error del OS.
func run() error {
	log.Println("STARTING PAYVUE APP")
	scope := os.Getenv("SCOPE")
	if scope == "" {
		return errors.New("empty scope: set SCOPE environment variable (reader or writer)")
	}

	parts := strings.Split(scope, "-")
	app := parts[0]
	if len(parts) >= scopePosition {
		app = parts[1]
	}

	log.Println("changes path to: " + app)
	if err := os.Chdir(app); err != nil {
		log.Println("error", err)
		return fmt.Errorf("change dir error on: %s - %s", app, err)
	}

	cmd := exec.Command("./app", os.Args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return err
	}

	if err := cmd.Wait(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
