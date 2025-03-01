package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/evg96/buffer/internal/usecase"
	"github.com/go-chi/chi"
)

type Handler struct {
	usecases *usecase.UseCases
}

func NewHandler(usecases *usecase.UseCases) *Handler {
	return &Handler{usecases: usecases}
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	// роут для записи данных в буфер
	router.Get("/api/v1/facts/get", h.GetFacts)

	// роут для сохранения данных
	router.Get("/api/v1/fact/save", h.SaveFact)

	return router
}

// обработчик сохранения пула записей в буфер
func (h Handler) GetFacts(w http.ResponseWriter, r *http.Request) {
	urlValues := getUrlValues(r.URL.Query())

	offset := getOffset(r)

	newOffset, err := h.usecases.SaveFactsToBuffer(r.Context(), urlValues, offset)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Sprintf("you save %d facts\n", newOffset)
		w.Write([]byte(msg))
		return
	}

	w.Write([]byte("you save all facts\n"))
}

// обработчик извлечения данных из буфера и сохранения их
func (h Handler) SaveFact(w http.ResponseWriter, r *http.Request) {
	err := h.usecases.GetFactFromBuffer(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getUrlValues(values url.Values) map[string]string {
	vals := make(map[string]string, len(values))

	for key, val := range values {
		vals[key] = val[0]
	}
	return vals
}

func getOffset(r *http.Request) int {
	offset := r.URL.Query().Get("offset")

	n, err := strconv.Atoi(offset)

	if err != nil {
		return 0
	}
	return n
}
