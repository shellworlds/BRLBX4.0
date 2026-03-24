package repo

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/pgxutil"
)

type Subscription struct {
	ID                   uuid.UUID  `json:"id"`
	ClientID             string     `json:"client_id"`
	Plan                 string     `json:"plan"`
	Status               string     `json:"status"`
	StartDate            time.Time  `json:"start_date"`
	NextBilling          *time.Time `json:"next_billing,omitempty"`
	StripeCustomerID     *string    `json:"stripe_customer_id,omitempty"`
	StripeSubscriptionID *string    `json:"stripe_subscription_id,omitempty"`
}

type MealPaymentRecord struct {
	ID                    uuid.UUID `json:"id"`
	VendorID              uuid.UUID `json:"vendor_id"`
	KitchenID             uuid.UUID `json:"kitchen_id"`
	MealCount             int       `json:"meal_count"`
	Amount                float64   `json:"amount"`
	PaymentMethod         *string   `json:"payment_method,omitempty"`
	StripePaymentIntentID *string   `json:"stripe_payment_intent_id,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
}

type CarbonPurchase struct {
	ID                    uuid.UUID `json:"id"`
	ClientID              string    `json:"client_id"`
	Tonnes                float64   `json:"tonnes"`
	Amount                float64   `json:"amount"`
	StripePaymentIntentID *string  `json:"stripe_payment_intent_id,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
}

type Store struct {
	DB pgxutil.Querier
}

func (s *Store) InsertSubscription(ctx context.Context, sub *Subscription) error {
	const q = `
INSERT INTO subscriptions (client_id, plan, status, start_date, next_billing, stripe_customer_id, stripe_subscription_id)
VALUES ($1,$2,$3,now(),$4,$5,$6)
RETURNING id, start_date`
	return s.DB.QueryRow(ctx, q, sub.ClientID, sub.Plan, sub.Status, sub.NextBilling, sub.StripeCustomerID, sub.StripeSubscriptionID).
		Scan(&sub.ID, &sub.StartDate)
}

func (s *Store) GetActiveSubscription(ctx context.Context, clientID string) (*Subscription, error) {
	const q = `
SELECT id, client_id, plan, status, start_date, next_billing, stripe_customer_id, stripe_subscription_id
FROM subscriptions WHERE client_id=$1 AND status IN ('active','stub_active','trialing')
ORDER BY start_date DESC LIMIT 1`
	var sub Subscription
	err := s.DB.QueryRow(ctx, q, clientID).Scan(
		&sub.ID, &sub.ClientID, &sub.Plan, &sub.Status, &sub.StartDate, &sub.NextBilling, &sub.StripeCustomerID, &sub.StripeSubscriptionID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

func (s *Store) CancelSubscription(ctx context.Context, clientID string) (int64, error) {
	const q = `UPDATE subscriptions SET status='canceled', next_billing=NULL WHERE client_id=$1 AND status IN ('active','stub_active','trialing')`
	ct, err := s.DB.Exec(ctx, q, clientID)
	if err != nil {
		return 0, err
	}
	return ct.RowsAffected(), nil
}

func (s *Store) UpdateSubscriptionByStripeID(ctx context.Context, stripeSubID, status string, next *time.Time) error {
	const q = `UPDATE subscriptions SET status=$2, next_billing=$3 WHERE stripe_subscription_id=$1`
	_, err := s.DB.Exec(ctx, q, stripeSubID, status, next)
	return err
}

func (s *Store) InsertMealRecord(ctx context.Context, r *MealPaymentRecord) error {
	const q = `
INSERT INTO payment_meal_records (vendor_id, kitchen_id, meal_count, amount, payment_method, stripe_payment_intent_id)
VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at`
	return s.DB.QueryRow(ctx, q, r.VendorID, r.KitchenID, r.MealCount, r.Amount, r.PaymentMethod, r.StripePaymentIntentID).
		Scan(&r.ID, &r.CreatedAt)
}

func (s *Store) InsertCarbon(ctx context.Context, c *CarbonPurchase) error {
	const q = `
INSERT INTO carbon_credit_purchases (client_id, tonnes, amount, stripe_payment_intent_id)
VALUES ($1,$2,$3,$4) RETURNING id, created_at`
	return s.DB.QueryRow(ctx, q, c.ClientID, c.Tonnes, c.Amount, c.StripePaymentIntentID).Scan(&c.ID, &c.CreatedAt)
}

// TryMarkWebhookEvent inserts the event id; returns (true, nil) if duplicate.
func (s *Store) TryMarkWebhookEvent(ctx context.Context, eventID string) (duplicate bool, err error) {
	const q = `INSERT INTO stripe_webhook_events (id) VALUES ($1) ON CONFLICT (id) DO NOTHING`
	ct, err := s.DB.Exec(ctx, q, eventID)
	if err != nil {
		return false, err
	}
	return ct.RowsAffected() == 0, nil
}
