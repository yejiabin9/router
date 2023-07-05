package service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/yejiabin9/router/domain/model"
	"github.com/yejiabin9/router/domain/repository"
	"github.com/yejiabin9/router/proto/router"
	"k8s.io/api/apps/v1"
	v12 "k8s.io/api/networking/v1"
	v14 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
)

// 这里是接口类型
type IRouterDataService interface {
	AddRouter(*model.Router) (int64, error)
	DeleteRouter(int64) error
	UpdateRouter(*model.Router) error
	FindRouterByID(int64) (*model.Router, error)
	FindAllRouter() ([]model.Router, error)

	CreateRouterToK8s(*router.RouterInfo) error
	DeleteRouterFromK8s(*model.Router) error
	UpdateRouterToK8s(*router.RouterInfo) error
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

func (u *RouterDataService) CreateRouterToK8s(info *router.RouterInfo) error {
	ingress := u.setIngress(info)
	if _, err := u.K8sClientSet.NetworkingV1().Ingresses(info.RouterNamespace).Get(context.TODO(), info.RouterName, v14.GetOptions{}); err != nil {
		if _, err := u.K8sClientSet.NetworkingV1().Ingresses(info.RouterNamespace).Create(context.TODO(), ingress, v14.CreateOptions{}); err != nil {
			logrus.Error(err)
			return err
		}
		return nil
	} else {
		logrus.Info("router " + info.RouterName + " has exits")
		return errors.New("router " + info.RouterName + " has exits")
	}
}

func (u *RouterDataService) setIngress(info *router.RouterInfo) *v12.Ingress {
	route := &v12.Ingress{}
	route.TypeMeta = v14.TypeMeta{
		Kind:       "Ingress",
		APIVersion: "v1",
	}
	route.ObjectMeta = v14.ObjectMeta{
		Name:      info.RouterName,
		Namespace: info.RouterNamespace,
		Labels: map[string]string{
			"app-name": info.RouterName,
			"author":   "jiabin",
		},
		Annotations: map[string]string{
			"k8s.generated-by-ye": "create by jiabin Ye",
		},
	}

	className := "nginx"
	route.Spec = v12.IngressSpec{
		IngressClassName: &className,
		DefaultBackend:   nil,
		TLS:              nil,
		Rules:            u.getIngressPath(info),
	}
	return route

}

func (u *RouterDataService) getIngressPath(info *router.RouterInfo) (path []v12.IngressRule) {
	pathRule := v12.IngressRule{Host: info.RouterHost}
	ingressPath := []v12.HTTPIngressPath{}
	for _, v := range info.RouterPath {
		pathType := v12.PathTypePrefix
		ingressPath = append(ingressPath, v12.HTTPIngressPath{
			Path:     v.RouterPathName,
			PathType: &pathType,
			Backend: v12.IngressBackend{Service: &v12.IngressServiceBackend{
				Name: v.RouterBackendService,
				Port: v12.ServiceBackendPort{Number: v.RouterBackendServicePort},
			}},
		})
	}
	pathRule.IngressRuleValue = v12.IngressRuleValue{HTTP: &v12.HTTPIngressRuleValue{
		Paths: ingressPath,
	}}
	path = append(path, pathRule)

	return path
}

func (u *RouterDataService) DeleteRouterFromK8s(route2 *model.Router) error {
	//删除Ingress
	if err := u.K8sClientSet.NetworkingV1().Ingresses(route2.RouterNamespace).Delete(context.TODO(), route2.RouterName, v14.DeleteOptions{}); err != nil {
		//如果删除失败记录下
		logrus.Error(err)
		return err
	} else {
		if err := u.DeleteRouter(route2.ID); err != nil {
			logrus.Error(err)
			return err
		}
		logrus.Info("删除 ingress ID：" + strconv.FormatInt(route2.ID, 10) + " 成功！")
	}
	return nil
}

func (u *RouterDataService) UpdateRouterToK8s(info *router.RouterInfo) error {
	ingress := u.setIngress(info)
	if _, err := u.K8sClientSet.NetworkingV1().Ingresses(info.RouterNamespace).Update(context.TODO(), ingress, v14.UpdateOptions{}); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
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
