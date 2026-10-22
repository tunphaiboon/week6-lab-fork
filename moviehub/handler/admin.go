package handler

import (
	"net/http"

	"moviehub/mw"
)

func (h *Handler) AdminMembers(w http.ResponseWriter, r *http.Request) {
	ms, err := h.Store.ListMembers(r.Context())
	if err != nil {
		internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": ms})
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/me", mw.RequireLogin(h.Tickets, h.Me))
	mux.HandleFunc("POST /api/movies/{id}/reviews", mw.RequireLogin(h.Tickets, h.AddReview))
	mux.HandleFunc("DELETE /api/reviews/{reviewId}", mw.RequireLogin(h.Tickets, h.DeleteReview))
	mux.HandleFunc("GET /api/admin/members", mw.RequireLogin(h.Tickets, mw.RequireRole("admin", h.AdminMembers)))
}
