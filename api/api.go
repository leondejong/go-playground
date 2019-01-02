package api

import (
	"log"
	"net/http"
	"time"

	"github.com/leondejong/go-playground/database"
	"github.com/leondejong/go-playground/network"
	_ "github.com/lib/pq"
)

type Item struct {
	Id      int       `json:"id,omitempty"`
	Name    string    `json:"name,omitempty"`
	Content string    `json:"content,omitempty"`
	Date    time.Time `json:"date,omitempty"`
}

func MapToItem(m map[string]interface{}) *Item {
	id, i := m["id"].(int64)
	name, n := m["name"].(string)
	content, c := m["content"].(string)
	date, d := m["date"].(time.Time)

	if i && n && c && d {
		return &Item{int(id), name, content, date}
	}

	return nil
}

func ListToItems(l []map[string]interface{}) []*Item {
	list := []*Item{}

	for _, item := range l {
		list = append(list, MapToItem(item))
	}

	return list
}

func all(w http.ResponseWriter, r *http.Request) {
	res, err := database.All("item")
	if err != nil {
		network.Status(w, http.StatusInternalServerError, err.Error())
		return
	}

	network.JSON(w, ListToItems(res))
}

func read(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		network.Status(w, http.StatusBadRequest, "not enough parameters")
		return
	}

	res, err := database.Select("item", "id", id)
	if err != nil {
		network.Status(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(res) < 1 {
		network.Status(w, http.StatusNotFound, "zero results")
		return
	}

	network.JSON(w, MapToItem(res[0]))
}

func create(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	content := r.FormValue("content")
	if name == "" || content == "" {
		network.Status(w, http.StatusBadRequest, "not enough parameters")
		return
	}

	cols := []string{"name", "content"}
	vals := []interface{}{name, content}

	ra, err := database.Insert("item", cols, vals)
	if err != nil {
		network.Status(w, http.StatusInternalServerError, err.Error())
		return
	}
	if ra < 1 {
		network.Status(w, http.StatusInternalServerError, "no rows affected")
		return
	}

	network.Status(w, http.StatusOK)
}

func update(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	name := r.FormValue("name")
	content := r.FormValue("content")
	if id == "" || name == "" || content == "" {
		network.Status(w, http.StatusBadRequest, "not enough parameters")
		return
	}

	cols := []string{"name", "content"}
	vals := []interface{}{name, content}

	ra, err := database.Update("item", "id", id, cols, vals)
	if err != nil {
		network.Status(w, http.StatusInternalServerError, err.Error())
		return
	}
	if ra < 1 {
		network.Status(w, http.StatusNotFound, "no rows affected")
		return
	}

	network.Status(w, http.StatusOK)
}

func delete(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		network.Status(w, http.StatusBadRequest, "not enough parameters")
		return
	}

	ra, err := database.Delete("item", "id", id)
	if err != nil {
		network.Status(w, http.StatusInternalServerError, err.Error())
		return
	}
	if ra < 1 {
		network.Status(w, http.StatusNotFound, "no rows affected")
		return
	}

	network.Status(w, http.StatusOK)
}

func Init() {
	database.Connect("user", "password", "localhost", "list")

	network.Get("/", all)
	network.Get("/item", read)
	network.Post("/item", create)
	network.Put("/item", update)
	network.Delete("/item", delete)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
