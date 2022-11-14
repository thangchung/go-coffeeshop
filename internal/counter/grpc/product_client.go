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

type productDomainService struct {
	conn *grpc.ClientConn
}

var _ domain.ProductDomainService = (*productDomainService)(nil)

func NewProductDomainService(conn *grpc.ClientConn) domain.ProductDomainService {
	return &productDomainService{
		conn: conn,
	}
}

func (p *productDomainService) GetItemsByType(
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
