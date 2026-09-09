package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ticket-system/internal/models"
	"ticket-system/internal/repository"
)

type TicketHandler struct {
	Repo *repository.TicketRepository
}

type CreateTicketRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func NewTicketHandler(repo *repository.TicketRepository) *TicketHandler {
	return &TicketHandler{
		Repo: repo,
	}
}

func getUserID(c *gin.Context) (uint, bool) {
	value, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint)
	return userID, ok
}

func (h *TicketHandler) Create(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var req CreateTicketRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "title and description are required",
		})
		return
	}

	ticket := models.Ticket{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      "open",
	}

	if err := h.Repo.Create(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create ticket",
		})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

func (h *TicketHandler) List(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	tickets, err := h.Repo.FindAllByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch tickets",
		})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func (h *TicketHandler) GetByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ticket id",
		})
		return
	}

	ticket, err := h.Repo.FindByIDAndUserID(uint(ticketID), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "ticket not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch ticket",
		})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

func (h *TicketHandler) UpdateStatus(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ticket id",
		})
		return
	}

	var req UpdateStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "status is required",
		})
		return
	}

	if req.Status != "open" &&
		req.Status != "in_progress" &&
		req.Status != "closed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid status",
		})
		return
	}

	ticket, err := h.Repo.FindByIDAndUserID(uint(ticketID), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "ticket not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch ticket",
		})
		return
	}

	if ticket.Status == "closed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "closed ticket cannot be reopened or modified",
		})
		return
	}

	if ticket.Status == "open" && req.Status == "closed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ticket must move to in_progress before closing",
		})
		return
	}

	if ticket.Status == "in_progress" && req.Status == "open" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ticket cannot move back to open",
		})
		return
	}

	ticket.Status = req.Status

	if err := h.Repo.Update(ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update ticket",
		})
		return
	}

	c.JSON(http.StatusOK, ticket)
}