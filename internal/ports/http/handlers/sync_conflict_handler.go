package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/services"
)

type SyncConflictHandler struct {
	conflictService *services.SyncConflictService
}

func NewSyncConflictHandler(conflictService *services.SyncConflictService) *SyncConflictHandler {
	return &SyncConflictHandler{
		conflictService: conflictService,
	}
}

// GetPendingConflicts returns all pending sync conflicts for the authenticated user
func (h *SyncConflictHandler) GetPendingConflicts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	conflicts, err := h.conflictService.GetPendingConflicts(c.Request.Context(), userID.(entities.UserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pending conflicts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conflicts": conflicts,
		"count":     len(conflicts),
	})
}

// ResolveConflictRequest represents the request body for resolving a conflict
type ResolveConflictRequest struct {
	Action     string                 `json:"action" binding:"required"` // "use_local", "use_google", "merge", "ignore"
	EventData  map[string]interface{} `json:"event_data,omitempty"`      // For merge action
	Resolution string                 `json:"resolution"`                // Human-readable description
}

// ResolveConflict resolves a specific sync conflict
func (h *SyncConflictHandler) ResolveConflict(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	conflictID := c.Param("id")
	if conflictID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Conflict ID is required"})
		return
	}

	var req ResolveConflictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate action
	validActions := map[string]bool{
		"use_local":  true,
		"use_google": true,
		"merge":      true,
		"ignore":     true,
	}
	if !validActions[req.Action] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}

	action := &entities.ConflictResolutionAction{
		Action:     req.Action,
		EventData:  req.EventData,
		Resolution: req.Resolution,
	}

	if action.Resolution == "" {
		action.Resolution = "Manually resolved by user"
	}

	err := h.conflictService.ResolveConflict(
		c.Request.Context(),
		conflictID,
		action,
		string(userID.(entities.UserID)),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve conflict"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Conflict resolved successfully",
		"action":  req.Action,
	})
}

// GetConflictStats returns conflict statistics for the authenticated user
func (h *SyncConflictHandler) GetConflictStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get days parameter (default: 30)
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 365 {
		days = 30
	}

	stats, err := h.conflictService.GetConflictStats(
		c.Request.Context(),
		userID.(entities.UserID),
		days,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conflict statistics"})
		return
	}

	// Calculate total conflicts
	total := 0
	for _, count := range stats {
		total += count
	}

	c.JSON(http.StatusOK, gin.H{
		"stats":  stats,
		"total":  total,
		"period": days,
	})
}

// BulkResolveRequest represents a request to resolve multiple conflicts with the same action
type BulkResolveRequest struct {
	ConflictIDs []string `json:"conflict_ids" binding:"required"`
	Action      string   `json:"action" binding:"required"`
	Resolution  string   `json:"resolution"`
}

// BulkResolveConflicts resolves multiple conflicts with the same action
func (h *SyncConflictHandler) BulkResolveConflicts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req BulkResolveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if len(req.ConflictIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No conflict IDs provided"})
		return
	}

	// Validate action
	validActions := map[string]bool{
		"use_local":  true,
		"use_google": true,
		"ignore":     true,
	}
	if !validActions[req.Action] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action for bulk resolution"})
		return
	}

	action := &entities.ConflictResolutionAction{
		Action:     req.Action,
		Resolution: req.Resolution,
	}

	if action.Resolution == "" {
		action.Resolution = "Bulk resolved by user"
	}

	var results []map[string]interface{}
	successCount := 0

	for _, conflictID := range req.ConflictIDs {
		result := map[string]interface{}{
			"conflict_id": conflictID,
		}

		err := h.conflictService.ResolveConflict(
			c.Request.Context(),
			conflictID,
			action,
			string(userID.(entities.UserID)),
		)

		if err != nil {
			result["success"] = false
			result["error"] = err.Error()
		} else {
			result["success"] = true
			successCount++
		}

		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Bulk resolution completed",
		"total":        len(req.ConflictIDs),
		"successful":   successCount,
		"failed":       len(req.ConflictIDs) - successCount,
		"results":      results,
	})
}

// GetConflictDetails returns detailed information about a specific conflict
func (h *SyncConflictHandler) GetConflictDetails(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	conflictID := c.Param("id")
	if conflictID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Conflict ID is required"})
		return
	}

	// Note: In a real implementation, you'd want to verify that the conflict belongs to the user
	// This would require adding a method to the conflict service or repository

	c.JSON(http.StatusOK, gin.H{
		"message": "Conflict details endpoint - implementation pending",
		"conflict_id": conflictID,
		"user_id": userID,
	})
}