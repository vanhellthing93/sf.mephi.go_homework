package repositories

import (
	"database/sql"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
)

type CardRepository struct {
	DB *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{DB: db}
}

func (r *CardRepository) CreateCard(card *models.Card) error {
    query := `INSERT INTO cards (account_id, number, cvv, expiry, hmac, created_at)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

    return r.DB.QueryRow(query, card.AccountID, card.Number, card.CVV, card.Expiry, card.HMAC, card.CreatedAt).Scan(&card.ID)
}

func (r *CardRepository) GetCardsByAccountID(accountID uint) ([]models.Card, error) {
    var cards []models.Card
    query := `SELECT id, account_id, number, expiry, hmac, created_at FROM cards WHERE account_id=$1`
    rows, err := r.DB.Query(query, accountID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var card models.Card
        if err := rows.Scan(&card.ID, &card.AccountID, &card.Number, &card.Expiry, &card.HMAC, &card.CreatedAt); err != nil {
            return nil, err
        }
        cards = append(cards, card)
    }
    return cards, nil
}

func (r *CardRepository) GetCardByID(cardID uint) (*models.Card, error) {
    var card models.Card
    query := `SELECT id, account_id, number, cvv, expiry, hmac, created_at FROM cards WHERE id=$1`
    err := r.DB.QueryRow(query, cardID).Scan(&card.ID, &card.AccountID, &card.Number, &card.CVV, &card.Expiry, &card.HMAC, &card.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &card, nil
}

func (r *CardRepository) DeleteCard(cardID uint) error {
	query := `DELETE FROM cards WHERE id=$1`
	_, err := r.DB.Exec(query, cardID)
	return err
}