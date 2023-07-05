package handler

import (
	"context"
	log "github.com/asim/go-micro/v3/logger"
	"github.com/sirupsen/logrus"
	"github.com/yejiabin9/router/domain/model"
	"github.com/yejiabin9/router/domain/service"
	router "github.com/yejiabin9/router/proto/router"
	"github.com/yejiabin9/router/utils"
	"strconv"
)

type RouterHandler struct {
	//注意这里的类型是 IRouterDataService 接口类型
	RouterDataService service.IRouterDataService
}

// Call is a single request handler called via client.Call or the generated client code
func (e *RouterHandler) AddRouter(ctx context.Context, info *router.RouterInfo, rsp *router.Response) error {
	log.Info("Received *router.AddRouter request")
	route := &model.Router{}
	if err := utils.SwapTo(info, route); err != nil {
		logrus.Error(err)
		rsp.Msg = err.Error()
		return err
	}
	//创建route到k8s
	if err := e.RouterDataService.CreateRouterToK8s(info); err != nil {
		logrus.Error(err)
		rsp.Msg = err.Error()
		return err
	} else {
		//写入数据库
		routeID, err := e.RouterDataService.AddRouter(route)
		if err != nil {
			logrus.Error(err)
			rsp.Msg = err.Error()
			return err
		}
		logrus.Info("Route 添加成功 ID 号为：" + strconv.FormatInt(routeID, 10))
		rsp.Msg = "Route 添加成功 ID 号为：" + strconv.FormatInt(routeID, 10)
	}
	return nil
}

func (e *RouterHandler) DeleteRouter(ctx context.Context, req *router.RouterId, rsp *router.Response) error {
	log.Info("Received *router.DeleteRouter request")
	routeModel, err := e.RouterDataService.FindRouterByID(req.Id)
	if err != nil {
		logrus.Error(err)
		return err
	}
	//从k8s中删除，并且删除数据库中数据
	if err := e.RouterDataService.DeleteRouterFromK8s(routeModel); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (e *RouterHandler) UpdateRouter(ctx context.Context, req *router.RouterInfo, rsp *router.Response) error {
	log.Info("Received *router.UpdateRouter request")
	if err := e.RouterDataService.UpdateRouterToK8s(req); err != nil {
		logrus.Error(err)
		return err
	}
	//查询数据库的信息
	routeModel, err := e.RouterDataService.FindRouterByID(req.Id)
	if err != nil {
		logrus.Error(err)
		return err
	}
	//数据更新
	if err := utils.SwapTo(req, routeModel); err != nil {
		logrus.Error(err)
		return err
	}
	return e.RouterDataService.UpdateRouter(routeModel)
}

func (e *RouterHandler) FindRouterByID(ctx context.Context, req *router.RouterId, rsp *router.RouterInfo) error {
	log.Info("Received *router.FindRouterByID request")
	routeModel, err := e.RouterDataService.FindRouterByID(req.Id)
	if err != nil {
		logrus.Error(err)
		return err
	}
	//数据转化
	if err := utils.SwapTo(routeModel, rsp); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (e *RouterHandler) FindAllRouter(ctx context.Context, req *router.FindAll, rsp *router.AllRouter) error {
	log.Info("Received *router.FindAllRouter request")
	allRoute, err := e.RouterDataService.FindAllRouter()
	if err != nil {
		logrus.Error(err)
		return err
	}
	//整理下格式
	for _, v := range allRoute {
		//创建实例
		routeInfo := &router.RouterInfo{}
		//把查询出来的数据进行转化
		if err := utils.SwapTo(v, routeInfo); err != nil {
			logrus.Error(err)
			return err
		}
		//数据合并
		rsp.RouterInfo = append(rsp.RouterInfo, routeInfo)
	}
	return nil
}
