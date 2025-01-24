package repositories

import (
	"database/sql"

	"github.com/jeagerism/medium-clone/backend/internal/payments/entities"
)

type subscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

// บันทึกข้อมูล Subscription ลงใน `user_subscriptions`
func (r *subscriptionRepository) SaveSubscription(subscription entities.UserSubscription) error {
	query := `
		INSERT INTO user_subscriptions (user_id, plan_id, subscription_id, status, start_date, end_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`
	_, err := r.db.Exec(query, subscription.UserID, subscription.PlanID, subscription.SubscriptionID, subscription.Status, subscription.StartDate, subscription.EndDate)
	return err
}

// บันทึกข้อมูล Payment ลงใน `payment_records`
func (r *subscriptionRepository) SavePayment(payment entities.PaymentRecord) error {
	query := `
		INSERT INTO payment_records (user_id, subscription_id, payment_gateway, payment_id, amount, currency, payment_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, payment.UserID, payment.SubscriptionID, payment.PaymentGateway, payment.PaymentID, payment.Amount, payment.Currency, payment.PaymentDate)
	return err
}
