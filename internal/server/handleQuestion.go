package server

import (
	"encoding/json"
	"hitalent/internal/model"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) HandleQuestion(w http.ResponseWriter, r *http.Request) {
	s.lg.Info("HandleQuestion called", "method", r.Method, "path", r.URL.Path)

	path := strings.TrimPrefix(r.URL.Path, "/questions")
	path = strings.Trim(path, "/") // нормализуем

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			s.lg.Info("Routing to HandleGetAllQuestion")
			s.HandleGetAllQuestion(w)
		case http.MethodPost:
			s.lg.Info("Routing to HandleCreateQuestion")
			s.HandleCreateQuestion(w, r)
		default:
			s.lg.Warn("Method not allowed", "method", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	parts := strings.Split(path, "/")
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		s.lg.Warn("Invalid question ID", "value", parts[0], "error", err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if len(parts) == 1 {
		switch r.Method {
		case http.MethodGet:
			s.lg.Info("Routing to HandleGetQuestion", "id", id)
			s.HandleGetQuestion(w, id)
		case http.MethodDelete:
			s.lg.Info("Routing to HandleDeleteQuestion", "id", id)
			s.HandleDeleteQuestion(w, id)
		default:
			s.lg.Warn("Method not allowed", "method", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	if len(parts) == 2 && parts[1] == "answers" {
		if r.Method == http.MethodPost {
			s.lg.Info("Routing to HandleAddAnswerForQuestion", "question_id", id)
			s.HandleAddAnswerForQuestion(w, r, id)
			return
		}
		s.lg.Warn("Method not allowed for answers", "method", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.lg.Warn("Path not found", "path", r.URL.Path)
	http.NotFound(w, r)
}

func (s *Server) HandleGetAllQuestion(w http.ResponseWriter) {
	s.lg.Info("HandleGetAllQuestion called")

	var questions []model.Question
	tx := s.storage.Find(&questions)
	if tx.Error != nil {
		s.lg.Error("Failed to fetch all questions", "error", tx.Error)
		http.Error(w, tx.Error.Error(), http.StatusInternalServerError)
		return
	}

	s.lg.Info("Fetched all questions", "count", len(questions))

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(questions)
	if err != nil {
		s.lg.Error("Failed to encode response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.lg.Info("Response successfully sent for HandleGetAllQuestion")
}

func (s *Server) HandleCreateQuestion(w http.ResponseWriter, r *http.Request) {
	s.lg.Info("HandleCreateQuestion called")

	var question model.Question

	err := json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		s.lg.Error("Failed to decode request body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.lg.Info("Request body decoded", "question_text", question.Text)

	err = s.validator.Struct(question)
	if err != nil {
		s.lg.Error("Validation failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.lg.Info("Validation passed")

	tx := s.storage.Create(&question)
	if tx.Error != nil {
		s.lg.Error("Failed to create question in database", "error", tx.Error)
		http.Error(w, tx.Error.Error(), http.StatusInternalServerError)
		return
	}
	s.lg.Info("Question successfully created", "id", question.Id)

	response := map[string]interface{}{
		"message": "Вопрос успешно создан",
		"id":      question.Id,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		s.lg.Error("Failed to encode response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.lg.Info("Response successfully sent for HandleCreateQuestion", "id", question.Id)
}

func (s *Server) HandleGetQuestion(w http.ResponseWriter, id int) {
	s.lg.Info("HandleGetQuestion called", "question_id", id)

	var question model.Question
	if err := s.storage.First(&question, id).Error; err != nil {
		s.lg.Error("Question not found", "id", id, "error", err)
		http.Error(w, "Question not found: "+err.Error(), http.StatusNotFound)
		return
	}
	s.lg.Info("Question fetched from database", "id", question.Id, "text", question.Text)

	var answers []model.Answer
	if err := s.storage.Where("question_id = ?", id).Find(&answers).Error; err != nil {
		s.lg.Error("Failed to fetch answers", "question_id", id, "error", err)
		http.Error(w, "Failed to fetch answers: "+err.Error(), http.StatusInternalServerError)
		return
	}
	s.lg.Info("Answers fetched from database", "question_id", id, "answers_count", len(answers))

	response := map[string]interface{}{
		"question": question,
		"answers":  answers,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.lg.Error("Failed to encode response", "question_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.lg.Info("Response successfully sent", "question_id", id)
}

func (s *Server) HandleDeleteQuestion(w http.ResponseWriter, id int) {
	s.lg.Info("HandleDeleteQuestion called", "question_id", id)

	tx := s.storage.Exec("DELETE FROM questions WHERE id = $1", id)
	if tx.Error != nil {
		s.lg.Error("Failed to delete question", "id", id, "error", tx.Error)
		http.Error(w, tx.Error.Error(), http.StatusInternalServerError)
		return
	}

	if tx.RowsAffected == 0 {
		s.lg.Warn("No question deleted, possibly non-existent ID", "id", id)
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	s.lg.Info("Question deleted successfully", "id", id)

	response := map[string]interface{}{
		"message": "успешное удаление",
		"id":      id,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.lg.Error("Failed to encode response", "question_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.lg.Info("Response successfully sent", "question_id", id)
}

func (s *Server) HandleAddAnswerForQuestion(w http.ResponseWriter, r *http.Request, id int) {
	s.lg.Info("HandleAddAnswerForQuestion called", "question_id", id)

	var answer model.Answer
	var question model.Question

	err := json.NewDecoder(r.Body).Decode(&answer)
	if err != nil {
		s.lg.Error("Failed to decode request body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.lg.Info("Request body decoded", "answer_text", answer.Text)

	tx := s.storage.Find(&question, id)
	if tx.Error != nil {
		s.lg.Error("Failed to find question", "id", id, "error", tx.Error)
		http.Error(w, tx.Error.Error(), http.StatusInternalServerError)
		return
	}
	if tx.RowsAffected == 0 {
		s.lg.Warn("Question not found", "id", id)
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}
	s.lg.Info("Question found", "id", id)

	answer.QuestionId = id
	tx = s.storage.Create(&answer)
	if tx.Error != nil {
		s.lg.Error("Failed to create answer", "question_id", id, "error", tx.Error)
		http.Error(w, tx.Error.Error(), http.StatusInternalServerError)
		return
	}

	s.lg.Info("Answer created successfully", "answer_id", answer.Id, "question_id", id)

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]interface{}{
		"message": "ответ успешно добавлен",
		"id":      answer.Id,
	}
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		s.lg.Error("Failed to encode response", "answer_id", answer.Id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.lg.Info("Response successfully sent", "answer_id", answer.Id)
}
