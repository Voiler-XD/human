package billing

import (
	"github.com/barun-bash/human/human-studio/server/models"
)

// Service manages Stripe billing integration.
// In production, this would use the Stripe Go SDK.
type Service struct {
	secretKey string
}

func NewService(secretKey string) *Service {
	return &Service{secretKey: secretKey}
}

// GetSubscription retrieves the user's current subscription.
func (s *Service) GetSubscription(userID string) (*models.Subscription, error) {
	// TODO: Query from database, sync with Stripe
	return &models.Subscription{
		Plan:   "free",
		Status: "active",
	}, nil
}

// CreateCheckoutSession creates a Stripe Checkout session for upgrading.
func (s *Service) CreateCheckoutSession(userID, priceID string) (string, error) {
	// TODO: Use stripe-go to create a checkout session
	// Returns the checkout URL
	return "https://checkout.stripe.com/placeholder", nil
}

// GetBillingHistory returns recent invoices for the user.
func (s *Service) GetBillingHistory(userID string) ([]models.BillingRecord, error) {
	// TODO: Query from database, synced via Stripe webhooks
	return []models.BillingRecord{}, nil
}

// HandleWebhook processes incoming Stripe webhook events.
func (s *Service) HandleWebhook(payload []byte, signature string) error {
	// TODO: Verify signature with Stripe webhook secret
	// TODO: Handle events: checkout.session.completed, invoice.paid,
	//       customer.subscription.updated, customer.subscription.deleted
	return nil
}
