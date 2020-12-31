// Code generated by protoc-gen-micro. DO NOT EDIT.

package backend

import (
  "context"
  "github.com/ofavor/micro-lite/server"
  "github.com/ofavor/micro-lite/client"
)

type BackendService interface {
  Connect(ctx context.Context, in *StatusRequest, opts ...client.CallOption) (*StatusResponse, error)
  Disconnect(ctx context.Context, in *StatusRequest, opts ...client.CallOption) (*StatusResponse, error)
  Data(ctx context.Context, in *DataRequest, opts ...client.CallOption) (*DataResponse, error)
}

type backendService struct {
  serviceName string
  c client.Client
}

func NewBackendService(name string, c client.Client) BackendService {
  return &backendService {
    serviceName: name,
    c: c,
  }
}

func (s *backendService)Connect(ctx context.Context, in *StatusRequest, opts ...client.CallOption) (*StatusResponse, error) {
  req := client.NewRequest(s.serviceName, "Backend.Connect", in)
  rsp := new(StatusResponse)
  err := s.c.Call(ctx, req, rsp, opts...)
  if err != nil {
    return nil, err
  }
  return rsp, nil
}

func (s *backendService)Disconnect(ctx context.Context, in *StatusRequest, opts ...client.CallOption) (*StatusResponse, error) {
  req := client.NewRequest(s.serviceName, "Backend.Disconnect", in)
  rsp := new(StatusResponse)
  err := s.c.Call(ctx, req, rsp, opts...)
  if err != nil {
    return nil, err
  }
  return rsp, nil
}

func (s *backendService)Data(ctx context.Context, in *DataRequest, opts ...client.CallOption) (*DataResponse, error) {
  req := client.NewRequest(s.serviceName, "Backend.Data", in)
  rsp := new(DataResponse)
  err := s.c.Call(ctx, req, rsp, opts...)
  if err != nil {
    return nil, err
  }
  return rsp, nil
}

type BackendHandler interface {
  Connect(ctx context.Context, in *StatusRequest, out *StatusResponse) error
  Disconnect(ctx context.Context, in *StatusRequest, out *StatusResponse) error
  Data(ctx context.Context, in *DataRequest, out *DataResponse) error
}

func RegisterBackendHandler(s server.Server, h BackendHandler) {
  hdr := server.NewHandler("Backend", h)
  s.Handle(hdr)
}

