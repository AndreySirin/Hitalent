package server

import (
	"encoding/json"
	"hitalent/internal/model"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) HandleAnswer(w http.ResponseWriter, r *http.Request) {
	s.lg.Info("HandleAnswer called", "path", r.URL.Path, "method", r.Method)

	idStr := strings.TrimPrefix(r.URL.Path, "/answers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.lg.Error("Invalid answer ID", "id_str", idStr, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.lg.Info("Parsed answer ID", "id", id)

	switch r.Method {
	case http.MethodGet:
		s.lg.Info("Handling GET for answer", "id", id)
		s.HandleGetAnswer(w, id)
		return
	case http.MethodDelete:
		s.lg.Info("Handling DELETE for answer", "id", id)
		s.HandleDeleteAnswer(w, id)
		return
	default:
		s.lg.Warn("Method not allowed for HandleAnswer", "method", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) HandleDeleteAnswer(w http.ResponseWriter, id int) {
	s.lg.Info("HandleDeleteAnswer called", "id", id)

	var answer model.Answer
	tx := s.storage.Delete(&answer, id)
	if tx.Error != nil {
		s.lg.Error("Failed to delete answer", "id", id, "error", tx.Error)
		http.Error(w, tx.Error.Error(), http.StatusInternalServerError)
		return
	}

	if tx.RowsAffected == 0 {
		s.lg.Warn("No rows affected when deleting answer", "id", id)
		http.Error(w, "Answer not found", http.StatusNotFound)
		return
	}

	s.lg.Info("Answer deleted successfully", "id", id)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "ответ успешно удален",
		"id":      id,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.lg.Error("Failed to encode response", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) HandleGetAnswer(w http.ResponseWriter, id int) {
	s.lg.Info("HandleGetAnswer called", "id", id)

	var answer model.Answer
	tx := s.storage.First(&answer, id)
	if tx.Error != nil {
		s.lg.Error("Failed to fetch answer", "id", id, "error", tx.Error)
		http.Error(w, tx.Error.Error(), http.StatusInternalServerError)
		return
	}

	if tx.RowsAffected == 0 {
		s.lg.Warn("Answer not found", "id", id)
		http.Error(w, "Answer not found", http.StatusNotFound)
		return
	}

	s.lg.Info("Answer fetched successfully", "id", id)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(answer); err != nil {
		s.lg.Error("Failed to encode response", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
