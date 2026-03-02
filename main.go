package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
)

// Uso etiquetas para saber como mapear el json a la estructura de Go
type Tarta struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Masa        string `json:"masa"`
	Imagen      string `json:"imagen"`
	MensajeWP   string `json:"mensajeWP"`
}

func main() {
	// 1. Leo el archivo JSON con la información de las tartas
	file, err := os.ReadFile("tartas.json")
	if err != nil {
		log.Fatal("Error al leer tartas.json: ", err)
	}

	var tartas []Tarta
	err = json.Unmarshal(file, &tartas)
	if err != nil {
		log.Fatal("Error al parsear el JSON: ", err)
	}

	// 2. Codificamos los mensajes de WhatsApp para que las URLs sean válidas
	for i := range tartas {
		tartas[i].MensajeWP = url.QueryEscape(tartas[i].MensajeWP)
	}

	// 3. Servir archivos estáticos (CSS, imágenes)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 4. Ruta principal
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseGlob("templates/*.html")
		if err != nil {
			http.Error(w, "Error al cargar las plantillas", http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "index.html", tartas)
		if err != nil {
			http.Error(w, "Error al renderizar", http.StatusInternalServerError)
		}
	})

	log.Println("Servidor de Tartas del Sol corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
