package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings" // <-- ¡No te olvides de agregar este import!
)

type Tarta struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Masa        string `json:"masa"`
	Imagen      string `json:"imagen"`
	MensajeWP   string `json:"mensajeWP"`
	URLSegura   template.URL
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
		// A. Codificamos el texto (esto pone los '+')
		textoCodificado := url.QueryEscape(tartas[i].MensajeWP)

		// B. Reemplazamos todos los '+' por '%20'
		textoParaWP := strings.ReplaceAll(textoCodificado, "+", "%20")

		// C. Armamos el link completo de WhatsApp y lo guardamos como URL segura
		linkCompleto := "https://wa.me/5492262309986?text=" + textoParaWP
		tartas[i].URLSegura = template.URL(linkCompleto)
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
