package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	gen "github.com/thangchung/go-coffeeshop/proto/gen"
	"google.golang.org/grpc"
)

type productGRPCClient struct {
	conn *grpc.ClientConn
}

var _ domain.ProductDomainService = (*productGRPCClient)(nil)

func NewGRPCProductClient(conn *grpc.ClientConn) domain.ProductDomainService {
	return &productGRPCClient{
		conn: conn,
	}
}

func (p *productGRPCClient) GetItemsByType(
	ctx context.Context,
	request *gen.PlaceOrderRequest,
	isBarista bool,
) (*gen.GetItemsByTypeResponse, error) {
	c := gen.NewProductServiceClient(p.conn)

	itemTypes := ""
	if isBarista {
		itemTypes = lo.Reduce(request.BaristaItems, func(agg string, item *gen.CommandItem, _ int) string {
			return fmt.Sprintf("%s,%s", agg, item.ItemType)
		}, "")
	} else {
		itemTypes = lo.Reduce(request.KitchenItems, func(agg string, item *gen.CommandItem, _ int) string {
			return fmt.Sprintf("%s,%s", agg, item.ItemType)
		}, "")
	}

	return c.GetItemsByType(ctx, &gen.GetItemsByTypeRequest{ItemTypes: strings.TrimLeft(itemTypes, ",")})
}
