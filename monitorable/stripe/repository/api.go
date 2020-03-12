package repository

import (
	"strconv"
	"time"

	"github.com/monitoror/monitoror/config"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

type (
	stripeRepository struct {
		stripeClient *client.API
		config       *config.Stripe
	}
)

// NewStripeRepository makes a new Stripe connection from an API key
func NewStripeRepository(config *config.Stripe) *stripeRepository {
	sc := &client.API{}
	sc.Init(config.Token, nil)
	return &stripeRepository{
		stripeClient: sc,
		config:       config,
	}
}

func (r *stripeRepository) GetCount(afterTimestamp string) int {
	if afterTimestamp == "" || afterTimestamp == "today" {
		afterTimestamp = strconv.FormatInt(bod(time.Now().Local()).Unix(), 10)
	}
	params := &stripe.BalanceTransactionListParams{}
	params.Filters.AddFilter("type", "", "charge")
	params.Filters.AddFilter("created", "gte", afterTimestamp)
	params.Filters.AddFilter("limit", "", "100")
	result := r.stripeClient.BalanceTransaction.List(params)
	count := 0
	for result.Next() {
		count = count + 1
	}
	return count
}

func bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}