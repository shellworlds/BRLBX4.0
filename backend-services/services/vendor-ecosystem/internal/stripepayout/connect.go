package stripepayout

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/account"
	"github.com/stripe/stripe-go/v81/accountlink"

	"github.com/shellworlds/BRLBX4.0/backend-services/services/vendor-ecosystem/internal/repo"
)

// EnsureExpressAccount creates a Stripe Connect Express account when missing and stores its id on the vendor.
func EnsureExpressAccount(ctx context.Context, st *repo.Store, vendorID uuid.UUID, country, emailFallback string) (accountID string, created bool, err error) {
	v, err := st.GetVendor(ctx, vendorID)
	if err != nil {
		return "", false, err
	}
	if v.StripeConnectAccountID != nil && strings.TrimSpace(*v.StripeConnectAccountID) != "" {
		return strings.TrimSpace(*v.StripeConnectAccountID), false, nil
	}
	email := vendorEmail(v, emailFallback)
	params := &stripe.AccountParams{
		Type:    stripe.String(string(stripe.AccountTypeExpress)),
		Country: stripe.String(country),
		Email:   stripe.String(email),
		Capabilities: &stripe.AccountCapabilitiesParams{
			CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{Requested: stripe.Bool(true)},
			Transfers:    &stripe.AccountCapabilitiesTransfersParams{Requested: stripe.Bool(true)},
		},
	}
	acct, err := account.New(params)
	if err != nil {
		return "", false, err
	}
	if err := st.SetStripeConnectAccountID(ctx, vendorID, acct.ID); err != nil {
		return "", false, err
	}
	return acct.ID, true, nil
}

func vendorEmail(v *repo.Vendor, fallback string) string {
	c := strings.TrimSpace(v.Contact)
	if strings.Contains(c, "@") {
		return c
	}
	if strings.TrimSpace(fallback) != "" {
		return fallback
	}
	return fmt.Sprintf("vendor+%s@onboarding.borelsigma.invalid", v.ID.String())
}

// CreateAccountOnboardingURL returns a hosted Stripe URL for Connect onboarding.
func CreateAccountOnboardingURL(acctID, refreshURL, returnURL string) (string, error) {
	link, err := accountlink.New(&stripe.AccountLinkParams{
		Account:    stripe.String(acctID),
		RefreshURL: stripe.String(refreshURL),
		ReturnURL:  stripe.String(returnURL),
		Type:       stripe.String(string(stripe.AccountLinkTypeAccountOnboarding)),
	})
	if err != nil {
		return "", err
	}
	return link.URL, nil
}
