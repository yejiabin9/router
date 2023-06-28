package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/yejiabin9/router/domain/model"
)

// 创建需要实现的接口
type IRouterRepository interface {
	//初始化表
	InitTable() error
	//根据ID查处找数据
	FindRouterByID(int64) (*model.Router, error)
	//创建一条 router 数据
	CreateRouter(*model.Router) (int64, error)
	//根据ID删除一条 router 数据
	DeleteRouterByID(int64) error
	//修改更新数据
	UpdateRouter(*model.Router) error
	//查找router所有数据
	FindAll() ([]model.Router, error)
}

// 创建routerRepository
func NewRouterRepository(db *gorm.DB) IRouterRepository {
	return &RouterRepository{mysqlDb: db}
}

type RouterRepository struct {
	mysqlDb *gorm.DB
}

// 初始化表
func (u *RouterRepository) InitTable() error {
	return u.mysqlDb.CreateTable(&model.Router{}, &model.RouterPath{}).Error
}

// 根据ID查找Router信息
func (u *RouterRepository) FindRouterByID(routerID int64) (router *model.Router, err error) {
	router = &model.Router{}
	return router, u.mysqlDb.Preload("RouterPath").First(router, routerID).Error
}

// 创建Router信息
func (u *RouterRepository) CreateRouter(router *model.Router) (int64, error) {
	return router.ID, u.mysqlDb.Create(router).Error
}

// 根据ID删除Router信息
func (u *RouterRepository) DeleteRouterByID(routerID int64) error {
	tx := u.mysqlDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		logrus.Error(tx.Error)
		return tx.Error
	}
	if err := u.mysqlDb.Where("id = ?", routerID).Delete(&model.Router{}).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return err
	}

	if err := u.mysqlDb.Where("router_id = ?", routerID).Delete(&model.RouterPath{}).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return err
	}

	return tx.Commit().Error
}

// 更新Router信息
func (u *RouterRepository) UpdateRouter(router *model.Router) error {
	return u.mysqlDb.Model(router).Update(router).Error
}

// 获取结果集
func (u *RouterRepository) FindAll() (routerAll []model.Router, err error) {
	return routerAll, u.mysqlDb.Preload("RouterPath").Find(&routerAll).Error
}
