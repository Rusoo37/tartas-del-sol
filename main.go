package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Renombramos "Tarta" a "Producto" porque ahora también incluye Wraps
type Producto struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Masa        string `json:"masa"`
	Imagen      string `json:"imagen"`
	MensajeWP   string `json:"mensajeWP"`
	URLSegura   template.URL
}

// Nueva estructura para agrupar el menú completo
type Menu struct {
	Grandes      []Producto `json:"grandes"`
	Individuales []Producto `json:"individuales"`
	Wraps        []Producto `json:"wraps"`
}

// Función auxiliar para no repetir código al generar los links
func generarLinksWP(productos []Producto) {
	for i := range productos {
		textoCodificado := url.QueryEscape(productos[i].MensajeWP)
		textoParaWP := strings.ReplaceAll(textoCodificado, "+", "%20")
		linkCompleto := "https://wa.me/5492262309986?text=" + textoParaWP
		productos[i].URLSegura = template.URL(linkCompleto)
	}
}

func main() {
	// 1. Leo el archivo JSON (ahora debería tener la estructura dividida)
	file, err := os.ReadFile("tartas.json")
	if err != nil {
		log.Fatal("Error al leer tartas.json: ", err)
	}

	var menu Menu
	err = json.Unmarshal(file, &menu)
	if err != nil {
		log.Fatal("Error al parsear el JSON: ", err)
	}

	// 2. Codificamos los mensajes de WhatsApp para cada categoría
	generarLinksWP(menu.Grandes)
	generarLinksWP(menu.Individuales)
	generarLinksWP(menu.Wraps)

	// 3. Servir archivos estáticos
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Cargamos TODAS las plantillas una sola vez
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Error al cargar las plantillas: ", err)
	}

	// Ruta principal (carga toda la página entera)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", menu)
	})

	// --- RUTAS HTMX (solo devuelven el fragmento de código) ---
	http.HandleFunc("/grandes", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "grandes", menu.Grandes)
	})

	http.HandleFunc("/individuales", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "individuales", menu.Individuales)
	})

	http.HandleFunc("/wraps", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "wraps", menu.Wraps)
	})

	log.Println("Servidor de Tartas del Sol corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
