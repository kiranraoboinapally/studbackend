package financemod

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
	api := r.PathPrefix("/api/v1/finance").Subrouter()
	api.Use(authMW)

	// Summary
	api.HandleFunc("/summary", h.Summary).Methods("GET")

	// Fee Heads
	api.HandleFunc("/fee-heads", h.ListFeeHeads).Methods("GET")
	api.HandleFunc("/fee-heads", h.CreateFeeHead).Methods("POST")
	api.HandleFunc("/fee-heads/{id:[0-9]+}", h.UpdateFeeHead).Methods("PUT")

	// Fee Structures
	api.HandleFunc("/fee-structures", h.ListFeeStructures).Methods("GET")
	api.HandleFunc("/fee-structures", h.CreateFeeStructure).Methods("POST")
	api.HandleFunc("/fee-structures/{id:[0-9]+}", h.UpdateFeeStructure).Methods("PUT")

	// Invoices
	api.HandleFunc("/invoices", h.ListAllInvoices).Methods("GET")
	api.HandleFunc("/invoices/me", h.MyInvoices).Methods("GET")
	api.HandleFunc("/invoices/generate", h.GenerateInvoice).Methods("POST")
	api.HandleFunc("/invoices/{id:[0-9]+}", h.GetInvoice).Methods("GET")
	api.HandleFunc("/invoices/student/{student_id:[0-9]+}", h.StudentInvoices).Methods("GET")

	// Payments
	api.HandleFunc("/payments", h.ListPayments).Methods("GET")
	api.HandleFunc("/payments", h.ProcessPayment).Methods("POST")
	api.HandleFunc("/payments/{id:[0-9]+}", h.GetPayment).Methods("GET")

	// Scholarships
	api.HandleFunc("/scholarships", h.ListScholarships).Methods("GET")
	api.HandleFunc("/scholarships", h.CreateScholarship).Methods("POST")
	api.HandleFunc("/scholarships/{id:[0-9]+}", h.UpdateScholarship).Methods("PUT")
	api.HandleFunc("/scholarships/assign", h.AssignScholarship).Methods("POST")
	api.HandleFunc("/scholarships/student/{student_id:[0-9]+}", h.StudentScholarships).Methods("GET")

	// Discounts
	api.HandleFunc("/discounts/student/{student_id:[0-9]+}", h.StudentDiscounts).Methods("GET")
	api.HandleFunc("/discounts", h.ApplyDiscount).Methods("POST")

	// Installments
	api.HandleFunc("/installments/student/{student_id:[0-9]+}", h.StudentInstallments).Methods("GET")
	api.HandleFunc("/installments", h.CreateInstallment).Methods("POST")

	// Refunds
	api.HandleFunc("/refunds", h.RequestRefund).Methods("POST")
	api.HandleFunc("/refunds", h.ListRefunds).Methods("GET")
	api.HandleFunc("/refunds/{id:[0-9]+}/approve", h.ApproveRefund).Methods("POST")
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetSummary(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) ListFeeHeads(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ListFeeHeads(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateFeeHead(w http.ResponseWriter, r *http.Request) {
	var fh domain.FeeHead
	json.NewDecoder(r.Body).Decode(&fh)
	if err := h.service.CreateFeeHead(r.Context(), &fh); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, fh)
}
func (h *Handler) UpdateFeeHead(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var fh domain.FeeHead
	json.NewDecoder(r.Body).Decode(&fh)
	if err := h.service.UpdateFeeHead(r.Context(), uint(id), &fh); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, fh)
}
func (h *Handler) ListFeeStructures(w http.ResponseWriter, r *http.Request) {
	pid, _ := strconv.ParseUint(r.URL.Query().Get("program_id"), 10, 64)
	year := r.URL.Query().Get("academic_year")
	data, err := h.service.ListFeeStructures(r.Context(), uint(pid), year)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateFeeStructure(w http.ResponseWriter, r *http.Request) {
	var fs domain.FeeStructure
	json.NewDecoder(r.Body).Decode(&fs)
	if err := h.service.CreateFeeStructure(r.Context(), &fs); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, fs)
}
func (h *Handler) UpdateFeeStructure(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var fs domain.FeeStructure
	json.NewDecoder(r.Body).Decode(&fs)
	if err := h.service.UpdateFeeStructure(r.Context(), uint(id), &fs); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, fs)
}
func (h *Handler) MyInvoices(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var studentID uint
	h.service.db.Table("student.students").Where("user_id = ?", userID).Select("id").Scan(&studentID)
	data, err := h.service.GetStudentInvoices(r.Context(), studentID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) StudentInvoices(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(mux.Vars(r)["student_id"], 10, 64)
	data, err := h.service.GetStudentInvoices(r.Context(), uint(sid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetInvoiceWithItems(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) ListAllInvoices(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	data, total, err := h.service.ListAllInvoices(r.Context(), page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 20 }
	response.List(w, data, total, page, pageSize)
}
func (h *Handler) GenerateInvoice(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StudentID uint `json:"student_id"`
		TermID    uint `json:"term_id"`
		ProgramID uint `json:"program_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	data, err := h.service.GenerateInvoiceForStudent(r.Context(), req.StudentID, req.TermID, req.ProgramID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, data)
}
func (h *Handler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}
	data, err := h.service.ProcessPayment(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, data)
}
func (h *Handler) ListPayments(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(r.URL.Query().Get("student_id"), 10, 64)
	data, err := h.service.ListPayments(r.Context(), uint(sid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetPayment(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetPayment(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) ListScholarships(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ListScholarships(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateScholarship(w http.ResponseWriter, r *http.Request) {
	var s domain.Scholarship
	json.NewDecoder(r.Body).Decode(&s)
	if err := h.service.CreateScholarship(r.Context(), &s); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, s)
}
func (h *Handler) UpdateScholarship(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var s domain.Scholarship
	json.NewDecoder(r.Body).Decode(&s)
	if err := h.service.UpdateScholarship(r.Context(), uint(id), &s); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, s)
}
func (h *Handler) AssignScholarship(w http.ResponseWriter, r *http.Request) {
	var ss domain.StudentScholarship
	json.NewDecoder(r.Body).Decode(&ss)
	if err := h.service.AssignScholarship(r.Context(), &ss); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, ss)
}
func (h *Handler) StudentScholarships(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(mux.Vars(r)["student_id"], 10, 64)
	data, err := h.service.GetStudentScholarships(r.Context(), uint(sid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) StudentDiscounts(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(mux.Vars(r)["student_id"], 10, 64)
	data, err := h.service.GetStudentDiscounts(r.Context(), uint(sid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) ApplyDiscount(w http.ResponseWriter, r *http.Request) {
	var d domain.StudentDiscount
	json.NewDecoder(r.Body).Decode(&d)
	if err := h.service.ApplyDiscount(r.Context(), &d); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, d)
}
func (h *Handler) StudentInstallments(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(mux.Vars(r)["student_id"], 10, 64)
	tid, _ := strconv.ParseUint(r.URL.Query().Get("term_id"), 10, 64)
	data, err := h.service.GetStudentInstallments(r.Context(), uint(sid), uint(tid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateInstallment(w http.ResponseWriter, r *http.Request) {
	var ip domain.InstallmentPlan
	json.NewDecoder(r.Body).Decode(&ip)
	if err := h.service.CreateInstallmentPlan(r.Context(), &ip); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, ip)
}
func (h *Handler) RequestRefund(w http.ResponseWriter, r *http.Request) {
	var ref domain.Refund
	json.NewDecoder(r.Body).Decode(&ref)
	if err := h.service.RequestRefund(r.Context(), &ref); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, ref)
}
func (h *Handler) ListRefunds(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(r.URL.Query().Get("student_id"), 10, 64)
	data, err := h.service.ListRefunds(r.Context(), uint(sid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) ApproveRefund(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	userID := middleware.GetUserID(r.Context())
	if err := h.service.ApproveRefund(r.Context(), uint(id), userID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "refund approved"})
}
