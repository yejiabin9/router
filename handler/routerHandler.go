package handler

import (
	"context"
	log "github.com/asim/go-micro/v3/logger"
	"github.com/yejiabin9/router/domain/service"
	router "github.com/yejiabin9/router/proto/router"
)

type RouterHandler struct {
	//注意这里的类型是 IRouterDataService 接口类型
	RouterDataService service.IRouterDataService
}

// Call is a single request handler called via client.Call or the generated client code
func (e *RouterHandler) AddRouter(ctx context.Context, info *router.RouterInfo, rsp *router.Response) error {
	log.Info("Received *router.AddRouter request")

	return nil
}

func (e *RouterHandler) DeleteRouter(ctx context.Context, req *router.RouterId, rsp *router.Response) error {
	log.Info("Received *router.DeleteRouter request")

	return nil
}

func (e *RouterHandler) UpdateRouter(ctx context.Context, req *router.RouterInfo, rsp *router.Response) error {
	log.Info("Received *router.UpdateRouter request")

	return nil
}

func (e *RouterHandler) FindRouterByID(ctx context.Context, req *router.RouterId, rsp *router.RouterInfo) error {
	log.Info("Received *router.FindRouterByID request")

	return nil
}

func (e *RouterHandler) FindAllRouter(ctx context.Context, req *router.FindAll, rsp *router.AllRouter) error {
	log.Info("Received *router.FindAllRouter request")

	return nil
}
