package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/barun-bash/human/human-studio/server/middleware"
	"github.com/barun-bash/human/human-studio/server/models"
)

func ListConnections(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		rows, err := db.Query(
			`SELECT id, service, connected_at FROM mcp_connections WHERE user_id = $1 ORDER BY service`,
			userID,
		)
		if err != nil {
			jsonError(w, "failed to list connections", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var connections []models.MCPConnection
		for rows.Next() {
			var c models.MCPConnection
			if err := rows.Scan(&c.ID, &c.Service, &c.ConnectedAt); err != nil {
				continue
			}
			connections = append(connections, c)
		}

		if connections == nil {
			connections = []models.MCPConnection{}
		}
		jsonResponse(w, connections, http.StatusOK)
	}
}

func CreateConnection(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		var body struct {
			Service      string `json:"service"`
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if body.Service == "" || body.AccessToken == "" {
			jsonError(w, "service and access_token are required", http.StatusBadRequest)
			return
		}

		_, err := db.Exec(
			`INSERT INTO mcp_connections (user_id, service, access_token, refresh_token)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (user_id, service)
			 DO UPDATE SET access_token = $3, refresh_token = $4, connected_at = NOW()`,
			userID, body.Service, body.AccessToken, body.RefreshToken,
		)
		if err != nil {
			jsonError(w, "failed to create connection", http.StatusInternalServerError)
			return
		}

		jsonResponse(w, map[string]string{"message": "connected"}, http.StatusCreated)
	}
}

func DeleteConnection(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		service := r.PathValue("service")

		_, err := db.Exec(
			`DELETE FROM mcp_connections WHERE user_id = $1 AND service = $2`,
			userID, service,
		)
		if err != nil {
			jsonError(w, "failed to delete connection", http.StatusInternalServerError)
			return
		}

		jsonResponse(w, map[string]string{"message": "disconnected"}, http.StatusOK)
	}
}
