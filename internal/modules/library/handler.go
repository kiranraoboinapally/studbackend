package librarymod

import (
	"encoding/json"
	"net/http"
	"strconv"

	"university-erp-backend/internal/domain"
	"university-erp-backend/internal/platform/middleware"
	"university-erp-backend/internal/platform/response"

	"github.com/gorilla/mux"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(r *mux.Router, authMW mux.MiddlewareFunc) {
	api := r.PathPrefix("/api/v1/library").Subrouter()
	api.Use(authMW)

	// Stats
	api.HandleFunc("/stats", h.Stats).Methods("GET")

	// Authors
	api.HandleFunc("/authors", h.ListAuthors).Methods("GET")
	api.HandleFunc("/authors", h.CreateAuthor).Methods("POST")
	api.HandleFunc("/authors/{id:[0-9]+}", h.GetAuthor).Methods("GET")
	api.HandleFunc("/authors/{id:[0-9]+}", h.UpdateAuthor).Methods("PUT")

	// Books
	api.HandleFunc("/books", h.ListBooks).Methods("GET")
	api.HandleFunc("/books", h.CreateBook).Methods("POST")
	api.HandleFunc("/books/{id:[0-9]+}", h.GetBook).Methods("GET")
	api.HandleFunc("/books/{id:[0-9]+}", h.UpdateBook).Methods("PUT")
	api.HandleFunc("/books/{id:[0-9]+}/authors", h.GetBookAuthors).Methods("GET")
	api.HandleFunc("/books/{id:[0-9]+}/authors", h.AddBookAuthor).Methods("POST")
	api.HandleFunc("/books/{id:[0-9]+}/authors/{author_id:[0-9]+}", h.RemoveBookAuthor).Methods("DELETE")

	// Book Copies
	api.HandleFunc("/book-copies", h.ListBookCopies).Methods("GET")
	api.HandleFunc("/book-copies", h.CreateBookCopy).Methods("POST")

	// Digital Resources
	api.HandleFunc("/digital-resources", h.ListDigitalResources).Methods("GET")
	api.HandleFunc("/digital-resources", h.CreateDigitalResource).Methods("POST")

	// Circulations
	api.HandleFunc("/circulations/issue", h.IssueBook).Methods("POST")
	api.HandleFunc("/circulations/{id:[0-9]+}/return", h.ReturnBook).Methods("POST")
	api.HandleFunc("/circulations/active", h.ListActiveCirculations).Methods("GET")
	api.HandleFunc("/circulations/mark-overdue", h.MarkOverdue).Methods("POST")

	// Reservations
	api.HandleFunc("/reservations", h.ListReservations).Methods("GET")
	api.HandleFunc("/reservations", h.CreateReservation).Methods("POST")

	// Fines
	api.HandleFunc("/fines", h.ListFines).Methods("GET")
	api.HandleFunc("/fines/{id:[0-9]+}/pay", h.PayFine).Methods("POST")

	// Purchase Requests
	api.HandleFunc("/purchase-requests", h.ListPurchaseRequests).Methods("GET")
	api.HandleFunc("/purchase-requests", h.CreatePurchaseRequest).Methods("POST")
	api.HandleFunc("/purchase-requests/{id:[0-9]+}/approve", h.ApprovePurchaseRequest).Methods("POST")
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetStats(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}

// Authors
func (h *Handler) ListAuthors(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ListAuthors(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetAuthor(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetAuthor(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var a domain.Author
	json.NewDecoder(r.Body).Decode(&a)
	if err := h.service.CreateAuthor(r.Context(), &a); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, a)
}
func (h *Handler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var a domain.Author
	json.NewDecoder(r.Body).Decode(&a)
	if err := h.service.UpdateAuthor(r.Context(), uint(id), &a); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, a)
}

// Books
func (h *Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	data, err := h.service.ListBooks(r.Context(), search)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetBook(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var b domain.Book
	json.NewDecoder(r.Body).Decode(&b)
	if err := h.service.CreateBook(r.Context(), &b); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, b)
}
func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var b domain.Book
	json.NewDecoder(r.Body).Decode(&b)
	if err := h.service.UpdateBook(r.Context(), uint(id), &b); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, b)
}
func (h *Handler) GetBookAuthors(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetBookAuthors(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) AddBookAuthor(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var ba domain.BookAuthor
	json.NewDecoder(r.Body).Decode(&ba)
	ba.BookID = uint(id)
	if err := h.service.AddBookAuthor(r.Context(), &ba); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, ba)
}
func (h *Handler) RemoveBookAuthor(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	aid, _ := strconv.ParseUint(mux.Vars(r)["author_id"], 10, 64)
	if err := h.service.RemoveBookAuthor(r.Context(), uint(id), uint(aid)); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}

// Book Copies
func (h *Handler) ListBookCopies(w http.ResponseWriter, r *http.Request) {
	bid, _ := strconv.ParseUint(r.URL.Query().Get("book_id"), 10, 64)
	data, err := h.service.ListBookCopies(r.Context(), uint(bid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateBookCopy(w http.ResponseWriter, r *http.Request) {
	var bc domain.BookCopy
	json.NewDecoder(r.Body).Decode(&bc)
	if err := h.service.CreateBookCopy(r.Context(), &bc); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, bc)
}

// Digital Resources
func (h *Handler) ListDigitalResources(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	data, err := h.service.ListDigitalResources(r.Context(), search)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateDigitalResource(w http.ResponseWriter, r *http.Request) {
	var dr domain.DigitalResource
	json.NewDecoder(r.Body).Decode(&dr)
	if err := h.service.CreateDigitalResource(r.Context(), &dr); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, dr)
}

// Circulations
func (h *Handler) IssueBook(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BookCopyID uint `json:"book_copy_id"`
		StudentID  uint `json:"student_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	issuedBy := middleware.GetUserID(r.Context())
	data, err := h.service.IssueBook(r.Context(), req.BookCopyID, req.StudentID, &issuedBy)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, data)
}
func (h *Handler) ReturnBook(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.ReturnBook(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) ListActiveCirculations(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(r.URL.Query().Get("student_id"), 10, 64)
	data, err := h.service.ListActiveCirculations(r.Context(), uint(sid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) MarkOverdue(w http.ResponseWriter, r *http.Request) {
	if err := h.service.MarkOverdue(r.Context()); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "overdue check completed"})
}

// Reservations
func (h *Handler) ListReservations(w http.ResponseWriter, r *http.Request) {
	bid, _ := strconv.ParseUint(r.URL.Query().Get("book_id"), 10, 64)
	sid, _ := strconv.ParseUint(r.URL.Query().Get("student_id"), 10, 64)
	data, err := h.service.ListReservations(r.Context(), uint(bid), uint(sid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var res domain.Reservation
	json.NewDecoder(r.Body).Decode(&res)
	if err := h.service.CreateReservation(r.Context(), &res); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, res)
}

// Fines
func (h *Handler) ListFines(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(r.URL.Query().Get("student_id"), 10, 64)
	data, err := h.service.ListFines(r.Context(), uint(sid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) PayFine(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err := h.service.PayFine(r.Context(), uint(id)); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "fine paid"})
}

// Purchase Requests
func (h *Handler) ListPurchaseRequests(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ListPurchaseRequests(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreatePurchaseRequest(w http.ResponseWriter, r *http.Request) {
	var pr domain.PurchaseRequest
	json.NewDecoder(r.Body).Decode(&pr)
	if err := h.service.CreatePurchaseRequest(r.Context(), &pr); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, pr)
}
func (h *Handler) ApprovePurchaseRequest(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	userID := middleware.GetUserID(r.Context())
	if err := h.service.ApprovePurchaseRequest(r.Context(), uint(id), userID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "purchase request approved"})
}
