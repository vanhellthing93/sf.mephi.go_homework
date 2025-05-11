package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	// "os"
	"strconv"
	"time"

	"github.com/vanhellthing93/sf.mephi.go_homework/internal/models"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

// var hmacSecret = []byte(os.Getenv("HMAC_SECRET"))

var hmacSecret = []byte("HMAC_SECRET")

type CardService struct {
	repo *repositories.CardRepository
}

func NewCardService(repo *repositories.CardRepository) *CardService {
	return &CardService{repo: repo}
}

func (s *CardService) CreateCard(accountID uint) (*models.Card, error) {
	// Генерация данных карты
	cardNumber := generateCardNumber()
	cvv := generateCVV()
	expiry := generateExpiryDate()

	card := &models.Card{
		AccountID: accountID,
		Number:    cardNumber,
		CVV:       cvv,
		Expiry:    expiry,
		CreatedAt: time.Now(),
	}

	// Шифрование данных карты
	if err := s.EncryptCardData(card); err != nil {
		return nil, err
	}

	// Сохранение карты в базе данных
	if err := s.repo.CreateCard(card); err != nil {
		return nil, err
	}

	return card, nil
}

func (s *CardService) GetCardsByAccountID(accountID uint) ([]models.Card, error) {
	cards, err := s.repo.GetCardsByAccountID(accountID)
	if err != nil {
		return nil, err
	}

	// Расшифровка данных карт
	for i := range cards {
		if err := s.DecryptCardData(&cards[i]); err != nil {
			return nil, err
		}
	}

	return cards, nil
}

func (s *CardService) GetCardByID(cardID uint) (*models.Card, error) {
	card, err := s.repo.GetCardByID(cardID)
	if err != nil {
		return nil, err
	}

	// Расшифровка данных карты
	if err := s.DecryptCardData(card); err != nil {
		return nil, err
	}

	return card, nil
}

func (s *CardService) DeleteCard(cardID uint) error {
	return s.repo.DeleteCard(cardID)
}

func (s *CardService) AccountBelongsToUser(accountID, userID uint) bool {
	// Проверяем, что аккаунт принадлежит пользователю
	var count int
	query := `SELECT COUNT(*) FROM accounts WHERE id=$1 AND user_id=$2`
	err := s.repo.DB.QueryRow(query, accountID, userID).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (s *CardService) CardBelongsToUser(cardID, userID uint) bool {
	// Проверяем, что карта принадлежит пользователю
	var count int
	query := `SELECT COUNT(*) FROM cards c JOIN accounts a ON c.account_id = a.id WHERE c.id=$1 AND a.user_id=$2`
	err := s.repo.DB.QueryRow(query, cardID, userID).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (s *CardService) EncryptCardData(card *models.Card) error {
    // Хеширование CVV
    hashedCVV, err := bcrypt.GenerateFromPassword([]byte(card.CVV), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    card.CVV = string(hashedCVV)

    // Генерация HMAC для проверки целостности
    card.HMAC = computeHMAC(card.Number, hmacSecret)

    return nil
}

func (s *CardService) DecryptCardData(card *models.Card) error {
    // Проверка HMAC
    if !verifyHMAC([]byte(card.Number), []byte(card.HMAC), hmacSecret) {
        return fmt.Errorf("HMAC verification failed")
    }

    return nil
}

func generateCardNumber() string {
    // Генерация номера карты по алгоритму Луна
    source := rand.NewSource(time.Now().UnixNano())
    rng := rand.New(source)
    cardNumber := ""
    for i := 0; i < 15; i++ { 
        cardNumber += strconv.Itoa(rng.Intn(10))
    }

    // Вычисляем контрольную цифру по алгоритму Луна
    cardNumber = cardNumber + calculateLuhnChecksum(cardNumber)

    return cardNumber
}

func calculateLuhnChecksum(cardNumber string) string {
	sum := 0
	// Проходим по цифрам с конца (исключая последнюю, которая будет контрольной)
	for i := len(cardNumber) - 2; i >= 0; i-- {
		digit := int(cardNumber[i] - '0')

		// Удваиваем каждую вторую цифру
		if (len(cardNumber)-i)%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit = digit/10 + digit%10
			}
		}

		sum += digit
	}

	// Контрольная цифра - это цифра, которая делает сумму кратной 10
	checksum := (10 - (sum % 10)) % 10
	return strconv.Itoa(checksum)
}

func generateCVV() string {
	// Генерация CVV
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	return strconv.Itoa(rng.Intn(900) + 100)
}

func generateExpiryDate() string {
	// Генерация срока действия карты
	// Добавляем 3 года к текущей дате
	expiryDate := time.Now().AddDate(3, 0, 0)
	return expiryDate.Format("01/06")
}

func computeHMAC(data string, secret []byte) string {
    h := hmac.New(sha256.New, secret)
    h.Write([]byte(data))
    return hex.EncodeToString(h.Sum(nil))
}

func verifyHMAC(data, mac, secret []byte) bool {
    expectedMAC := computeHMAC(string(data), secret)
    return hmac.Equal([]byte(expectedMAC), mac)
}