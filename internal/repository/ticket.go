package repository

import (
	"ticket-system/internal/models"

	"gorm.io/gorm"
)

type TicketRepository struct {
	DB *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{
		DB: db,
	}
}

func (r *TicketRepository) Create(ticket *models.Ticket) error {
	return r.DB.Create(ticket).Error
}

func (r *TicketRepository) FindAllByUserID(userID uint) ([]models.Ticket, error) {
	var tickets []models.Ticket

	err := r.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tickets).Error

	return tickets, err
}

func (r *TicketRepository) FindByIDAndUserID(
	ticketID uint,
	userID uint,
) (*models.Ticket, error) {
	var ticket models.Ticket

	err := r.DB.
		Where("id = ? AND user_id = ?", ticketID, userID).
		First(&ticket).Error

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *TicketRepository) Update(ticket *models.Ticket) error {
	return r.DB.Save(ticket).Error
}