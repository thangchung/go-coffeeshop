package grpc

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/thangchung/go-coffeeshop/internal/product/domain"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
	"github.com/thangchung/go-coffeeshop/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ProductServiceServerImpl struct {
	gen.UnimplementedProductServiceServer
	repo   domain.ProductRepo
	logger *mylogger.Logger
}

func NewProductServiceServerGrpc(
	grpcServer *grpc.Server,
	log *mylogger.Logger,
	repo domain.ProductRepo,
) {
	svc := ProductServiceServerImpl{
		logger: log,
		repo:   repo,
	}

	gen.RegisterProductServiceServer(grpcServer, &svc)

	reflection.Register(grpcServer)
}

func (g *ProductServiceServerImpl) GetItemTypes(ctx context.Context, request *gen.GetItemTypesRequest) (*gen.GetItemTypesResponse, error) {
	g.logger.Info("GET: GetItemTypes")

	res := gen.GetItemTypesResponse{}

	results, err := g.repo.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "ProductServiceServerImpl-g.repo.GetAll")
	}

	res.ItemTypes = append(res.ItemTypes, results...)

	return &res, nil
}

func (g *ProductServiceServerImpl) GetItemsByType(ctx context.Context, request *gen.GetItemsByTypeRequest) (*gen.GetItemsByTypeResponse, error) {
	g.logger.Info("GET: GetItemsByType with %s", request.ItemTypes)

	res := gen.GetItemsByTypeResponse{}

	itemTypes := strings.Split(request.ItemTypes, ",")

	results, err := g.repo.GetByTypes(ctx, itemTypes)
	if err != nil {
		return nil, errors.Wrap(err, "ProductServiceServerImpl-g.repo.GetItemsByType")
	}

	res.Items = append(res.Items, results...)

	return &res, nil
}
