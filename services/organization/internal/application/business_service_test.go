package application_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/organization/internal/adapters/postgres"
	"github.com/isapr/mini-erp/services/organization/internal/application"
)

func TestCreateBusinessWithPostgres(t *testing.T) {
	databaseURL := os.Getenv("ORGANIZATION_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("ORGANIZATION_DATABASE_URL not set")
	}

	ctx := context.Background()
	pool, err := postgres.OpenPool(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	repo := postgres.NewBusinessRepository(pool)
	service := application.NewBusinessService(repo)
	name := "Acme " + uuid.NewString()

	business, err := service.CreateBusiness(ctx, application.CreateBusinessInput{
		Name:     name,
		Timezone: "Asia/Jakarta",
	})
	if err != nil {
		t.Fatal(err)
	}

	stored, err := repo.FindByCode(ctx, business.Code)
	if err != nil {
		t.Fatal(err)
	}
	if stored.ID != business.ID {
		t.Fatalf("expected business id %s, got %s", business.ID, stored.ID)
	}
}
