package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	sweb "github.com/sanchpet/sweb-go-sdk"
)

func resolveTestClient(t *testing.T) *sweb.Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"result":[
			{"billingId":"login_vps_10","name":"infra-01"},
			{"billingId":"login_vps_12","name":"infra-03"},
			{"billingId":"login_vps_6","name":"dup"},
			{"billingId":"login_vps_7","name":"dup"}
		]}`))
	}))
	t.Cleanup(srv.Close)
	return sweb.New(sweb.WithBaseURL(srv.URL), sweb.WithHTTPClient(srv.Client()), sweb.WithToken("t"))
}

func TestResolveVPS(t *testing.T) {
	c := resolveTestClient(t)
	ctx := context.Background()

	if id, err := resolveVPS(ctx, c, "infra-01"); err != nil || id != "login_vps_10" {
		t.Fatalf("by name: got %q, %v; want login_vps_10", id, err)
	}
	if id, err := resolveVPS(ctx, c, "login_vps_12"); err != nil || id != "login_vps_12" {
		t.Fatalf("by billing id: got %q, %v; want login_vps_12", id, err)
	}
	if _, err := resolveVPS(ctx, c, "does-not-exist"); err == nil {
		t.Fatal("not-found: want error, got nil")
	}
	if _, err := resolveVPS(ctx, c, "dup"); err == nil {
		t.Fatal("ambiguous name: want error, got nil")
	}
}
