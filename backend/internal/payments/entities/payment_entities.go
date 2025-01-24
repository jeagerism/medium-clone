package entities

import "time"

// ข้อมูลการ Subscription Request
type SubscriptionRequest struct {
	UserID         int     `json:"user_id" binding:"required"`
	PlanID         string  `json:"plan_id" binding:"required"`
	SubscriptionID string  `json:"subscription_id" binding:"required"`
	PaymentGateway string  `json:"payment_gateway" binding:"required"`
	PaymentID      string  `json:"payment_id" binding:"required"`
	Amount         float64 `json:"amount" binding:"required"`
	Currency       string  `json:"currency" binding:"required"`
}

// ข้อมูลตาราง user_subscriptions
type UserSubscription struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	PlanID         string    `json:"plan_id"`
	SubscriptionID string    `json:"subscription_id"`
	Status         string    `json:"status"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
}

// ข้อมูลตาราง payment_records
type PaymentRecord struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	SubscriptionID string    `json:"subscription_id"`
	PaymentGateway string    `json:"payment_gateway"`
	PaymentID      string    `json:"payment_id"`
	Amount         float64   `json:"amount"`
	Currency       string    `json:"currency"`
	PaymentDate    time.Time `json:"payment_date"`
}
