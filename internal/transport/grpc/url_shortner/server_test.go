package urlshortner

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	url_shortnerv1 "github.com/notblinkyet/proto_url_shortner/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type Test struct {
	name           string
	request        *url_shortnerv1.CreateRequest
	expectedStatus codes.Code
}

func NewTest(name string, request *url_shortnerv1.CreateRequest, expectedStatus codes.Code) *Test {
	return &Test{
		name:           name,
		request:        request,
		expectedStatus: expectedStatus,
	}
}

func TestCreate(t *testing.T) {
	l := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	mockSrv := NewMockServices()
	Register(s, mockSrv, 5*time.Second)

	go func() {
		if err := s.Serve(l); err != nil {
			t.Errorf("failed to serve: %v", err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
		return l.Dial()
	}))
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := url_shortnerv1.NewUrlShortnerClient(conn)
	url := gofakeit.URL()

	tests := []*Test{
		NewTest("valid request", &url_shortnerv1.CreateRequest{Url: url}, codes.OK),
		NewTest("empty url", &url_shortnerv1.CreateRequest{Url: ""}, codes.InvalidArgument),
		NewTest("url already exists", &url_shortnerv1.CreateRequest{Url: url}, codes.AlreadyExists),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Create(context.Background(), tt.request)

			if status.Code(err) != tt.expectedStatus {
				t.Errorf("expected error code %v, got %v", tt.expectedStatus, status.Code(err))
			}
			if status.Code(err) == codes.OK {
				if resp == nil {
					t.Error("expected response to be non-nil")
				} else {
					fmt.Println(resp)
				}
			}
		})
	}

	s.Stop()
}
