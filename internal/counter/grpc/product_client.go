package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"
	gen "github.com/thangchung/go-coffeeshop/proto/gen"
	"google.golang.org/grpc"
)

type ProductServiceClient struct {
	ctx  context.Context
	conn *grpc.ClientConn
}

func NewProductServiceClient(ctx context.Context, conn *grpc.ClientConn) *ProductServiceClient {
	return &ProductServiceClient{
		ctx:  ctx,
		conn: conn,
	}
}

func (p *ProductServiceClient) GetItemsByType(request *gen.PlaceOrderRequest, isBarista bool) (*gen.GetItemsByTypeResponse, error) {
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

	return c.GetItemsByType(p.ctx, &gen.GetItemsByTypeRequest{ItemTypes: strings.TrimLeft(itemTypes, ",")})
}
