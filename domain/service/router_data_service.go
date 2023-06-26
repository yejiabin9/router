package service

import (
	"github.com/yejiabin9/router/domain/model"
	"github.com/yejiabin9/router/domain/repository"
	"k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
)

// 这里是接口类型
type IRouterDataService interface {
	AddRouter(*model.Router) (int64, error)
	DeleteRouter(int64) error
	UpdateRouter(*model.Router) error
	FindRouterByID(int64) (*model.Router, error)
	FindAllRouter() ([]model.Router, error)
}

// 创建
// 注意：返回值 IRouterDataService 接口类型
func NewRouterDataService(routerRepository repository.IRouterRepository, clientSet *kubernetes.Clientset) IRouterDataService {
	return &RouterDataService{RouterRepository: routerRepository, K8sClientSet: clientSet, deployment: &v1.Deployment{}}
}

type RouterDataService struct {
	//注意：这里是 IRouterRepository 类型
	RouterRepository repository.IRouterRepository
	K8sClientSet     *kubernetes.Clientset
	deployment       *v1.Deployment
}

// 插入
func (u *RouterDataService) AddRouter(router *model.Router) (int64, error) {
	return u.RouterRepository.CreateRouter(router)
}

// 删除
func (u *RouterDataService) DeleteRouter(routerID int64) error {
	return u.RouterRepository.DeleteRouterByID(routerID)
}

// 更新
func (u *RouterDataService) UpdateRouter(router *model.Router) error {
	return u.RouterRepository.UpdateRouter(router)
}

// 查找
func (u *RouterDataService) FindRouterByID(routerID int64) (*model.Router, error) {
	return u.RouterRepository.FindRouterByID(routerID)
}

// 查找
func (u *RouterDataService) FindAllRouter() ([]model.Router, error) {
	return u.RouterRepository.FindAll()
}
