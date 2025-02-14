package urlshortner

import (
	"context"
	"errors"
	"time"

	url_shortnerv1 "github.com/notblinkyet/proto_url_shortner/gen/go"
	my_errors "github.com/notblinkyet/url_shortner/internal/errors"
	"github.com/notblinkyet/url_shortner/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverApi struct {
	url_shortnerv1.UnimplementedUrlShortnerServer
	services services.IServices
	timeout  time.Duration
}

func Register(gRRPCserver *grpc.Server, services services.IServices, timeout time.Duration) {
	url_shortnerv1.RegisterUrlShortnerServer(gRRPCserver, serverApi{services: services, timeout: timeout})
}

func (s serverApi) Create(ctx context.Context, in *url_shortnerv1.CreateRequest) (*url_shortnerv1.CreateResponse, error) {
	if in.Url == "" {
		return nil, status.Error(codes.InvalidArgument, "short url is required")
	}
	shortUrl, err := s.services.Create(in.Url)
	if err != nil {
		if errors.Is(err, my_errors.ErrAlreadyExist) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, "fail to create short url")
	}
	return &url_shortnerv1.CreateResponse{
		ShortUrl: shortUrl,
	}, nil
}

func (s serverApi) Get(ctx context.Context, in *url_shortnerv1.GetRequest) (*url_shortnerv1.GetResponse, error) {

	if in.ShortUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "url is required")
	}
	url, err := s.services.Get(in.ShortUrl)
	if err != nil {
		if errors.Is(err, my_errors.ErrAliaceDontUse) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "fail to get short url")
	}
	return &url_shortnerv1.GetResponse{
		Url: url,
	}, nil
}
