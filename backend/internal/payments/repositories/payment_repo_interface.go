package repositories

import "github.com/jeagerism/medium-clone/backend/internal/payments/entities"

type SubscriptionRepository interface {
	SaveSubscription(subscription entities.UserSubscription) error
	SavePayment(payment entities.PaymentRecord) error
}
