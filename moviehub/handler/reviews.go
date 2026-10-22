package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"moviehub/mw"
	"moviehub/store"
)

type addReviewReq struct {
	Text string `json:"text"`
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	member, err := h.Store.MemberByID(r.Context(), mw.WhoFrom(r.Context()).ID)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "MEMBER_NOT_FOUND", "ไม่พบบัญชีสมาชิก")
		return
	}
	if err != nil {
		internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, member)
}

func (h *Handler) AddReview(w http.ResponseWriter, r *http.Request) {
	movieID, ok := pathID(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "BAD_ID", "เลขหนังต้องเป็นตัวเลข")
		return
	}
	var req addReviewReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "อ่านกล่อง JSON ไม่ออก")
		return
	}
	req.Text = strings.TrimSpace(req.Text)
	if utf8.RuneCountInString(req.Text) < 10 {
		writeError(w, http.StatusBadRequest, "BAD_REVIEW", "รีวิวต้องมีอย่างน้อย 10 ตัวอักษร")
		return
	}

	review, err := h.Store.AddReview(r.Context(), mw.WhoFrom(r.Context()).ID, movieID, req.Text)
	switch {
	case errors.Is(err, store.ErrDuplicate):
		writeError(w, http.StatusConflict, "REVIEW_EXISTS", "คุณรีวิวหนังเรื่องนี้แล้ว")
		return
	case errors.Is(err, store.ErrNoMovie):
		writeError(w, http.StatusNotFound, "MOVIE_NOT_FOUND", "ไม่พบหนังเรื่องนี้")
		return
	case err != nil:
		internal(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, review)
}

func (h *Handler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	reviewID, ok := pathID(r, "reviewId")
	if !ok {
		writeError(w, http.StatusBadRequest, "BAD_ID", "เลขรีวิวต้องเป็นตัวเลข")
		return
	}
	who := mw.WhoFrom(r.Context())
	err := h.Store.DeleteReview(r.Context(), reviewID, who.ID, who.Role == "admin")
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "REVIEW_NOT_FOUND", "ไม่พบรีวิวนี้")
		return
	}
	if err != nil {
		internal(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
