package application_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/identity/internal/adapters/postgres"
	"github.com/isapr/mini-erp/services/identity/internal/application"
)

func TestSignupTenantAdminWithPostgres(t *testing.T) {
	databaseURL := os.Getenv("IDENTITY_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("IDENTITY_DATABASE_URL not set")
	}

	ctx := context.Background()
	pool, err := postgres.OpenPool(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	repo := postgres.NewUserRepository(pool)
	service := application.NewSignupService(repo)
	email := "owner-" + uuid.NewString() + "@example.test"
	businessID := uuid.New()

	output, err := service.SignupTenantAdmin(ctx, application.SignupTenantAdminInput{
		BusinessID: businessID,
		Email:      email,
		Password:   "secret123",
		FullName:   "Owner User",
	})
	if err != nil {
		t.Fatal(err)
	}

	if output.BusinessID != businessID {
		t.Fatalf("expected business id %s, got %s", businessID, output.BusinessID)
	}

	user, err := repo.FindByEmail(ctx, email)
	if err != nil {
		t.Fatal(err)
	}
	if user.Email != email {
		t.Fatalf("expected email %s, got %s", email, user.Email)
	}
}
