/**
 *
 * @author 半城风雨
 * @since 2021/11/11
 * @File : user
 */
package service

import (
	"easygoadmin/app/dto"
	"easygoadmin/app/model"
	"easygoadmin/app/vo"
	"easygoadmin/utils"
	"easygoadmin/utils/gconv"
	"errors"
	"strings"
	"time"
)

var User = new(userService)

type userService struct{}

func (s *userService) GetList(req *dto.UserPageReq) ([]vo.UserInfoVo, int64, error) {
	// 初始化查询实例
	query := utils.XormDb.Where("mark=1")
	// 查询条件
	if req != nil {
		// 用户名
		if req.Username != "" {
			query = query.Where("username like ?", "%"+req.Username+"%")
		}
		// 性别
		if req.Gender > 0 {
			query = query.Where("gender=?", req.Gender)
		}
	}
	// 排序
	query = query.Asc("id")
	// 分页设置
	offset := (req.Page - 1) * req.Limit
	query = query.Limit(req.Limit, offset)
	// 查询列表
	list := make([]model.User, 0)
	count, err := query.FindAndCount(&list)
	if err != nil {
		return nil, 0, err
	}
	//return list, count, nil
	// 获取职级列表
	levelList := make([]model.Level, 0)
	utils.XormDb.Where("mark=1").OrderBy("sort asc").Cols("id,name").Find(&levelList)
	var levelMap = make(map[int]string)
	for _, v := range levelList {
		levelMap[v.Id] = v.Name
	}
	// 获取岗位列表
	positionList := make([]model.Position, 0)
	utils.XormDb.Where("mark=1").Cols("id,name").Find(&positionList)
	var positionMap = make(map[int]string)
	for _, v := range positionList {
		positionMap[v.Id] = v.Name
	}
	// 部门
	deptList := make([]model.Dept, 0)
	utils.XormDb.Where("mark=1").Cols("id,name").Find(&deptList)
	var deptMap = make(map[int]string)
	for _, v := range deptList {
		deptMap[v.Id] = v.Name
	}

	// 数据处理
	var result []vo.UserInfoVo
	for _, v := range list {
		item := vo.UserInfoVo{}
		item.User = v
		// 头像
		if v.Avatar != "" {
			item.Avatar = utils.GetImageUrl(v.Avatar)
		}
		// 性别
		if v.Gender > 0 {
			item.GenderName = utils.GENDER_LIST[v.Gender]
		}
		// 职级
		if v.LevelId > 0 {
			item.LevelName = levelMap[v.LevelId]
		}
		// 岗位
		if v.PositionId > 0 {
			item.PositionName = positionMap[v.PositionId]
		}
		// 部门
		if v.DeptId > 0 {
			item.DeptName = deptMap[v.DeptId]
		}
		// 角色列表
		roleList := UserRole.GetUserRoleList(v.Id)
		if len(roleList) > 0 {
			item.RoleList = roleList
		} else {
			item.RoleList = make([]model.Role, 0)
		}
		// 省市区
		cityList := make([]string, 0)
		// 省份编号
		cityList = append(cityList, item.ProvinceCode)
		// 城市编号
		cityList = append(cityList, item.CityCode)
		// 县区编号
		cityList = append(cityList, item.DistrictCode)
		item.City = cityList
		// 加入数组
		result = append(result, item)
	}
	return result, count, nil
}

func (s *userService) Add(req *dto.UserAddReq, userId int) (int64, error) {
	var entity model.User
	entity.Realname = req.Realname
	entity.Nickname = req.Nickname
	entity.Gender = req.Gender
	entity.Avatar = req.Avatar
	entity.Mobile = req.Mobile
	entity.Email = req.Email
	entity.Birthday = req.Birthday
	entity.DeptId = req.DeptId
	entity.LevelId = req.LevelId
	entity.PositionId = req.PositionId
	entity.Address = req.Address
	entity.Username = req.Username
	entity.Intro = req.Intro
	entity.Status = req.Status
	entity.Note = req.Note
	entity.Sort = req.Sort

	// 省市区处理
	if len(req.City) == 3 {
		entity.ProvinceCode = req.City[0]
		entity.CityCode = req.City[1]
		entity.DistrictCode = req.City[2]
	}

	// 密码
	if req.Password != "" {
		password, _ := utils.Md5(req.Password + req.Username)
		entity.Password = password
	}

	// 头像处理
	if req.Avatar != "" {
		avatar, err := utils.SaveImage(req.Avatar, "user")
		if err != nil {
			return 0, err
		}
		entity.Avatar = avatar
	}
	entity.CreateUser = userId
	entity.CreateTime = time.Now()
	entity.Mark = 1
	// 插入记录
	return entity.Insert()
}

func (s *userService) Update(req *dto.UserUpdateReq, userId int) (int64, error) {
	entity := &model.User{Id: req.Id}
	has, err := entity.Get()
	if err != nil || !has {
		return 0, errors.New("记录不存在")
	}
	entity.Realname = req.Realname
	entity.Nickname = req.Nickname
	entity.Gender = req.Gender
	entity.Avatar = req.Avatar
	entity.Mobile = req.Mobile
	entity.Email = req.Email
	entity.Birthday = req.Birthday
	entity.DeptId = req.DeptId
	entity.LevelId = req.LevelId
	entity.PositionId = req.PositionId
	entity.Address = req.Address
	entity.Username = req.Username
	entity.Intro = req.Intro
	entity.Status = req.Status
	entity.Note = req.Note
	entity.Sort = req.Sort

	// 省市区处理
	if len(req.City) == 3 {
		entity.ProvinceCode = req.City[0]
		entity.CityCode = req.City[1]
		entity.DistrictCode = req.City[2]
	}

	// 密码
	if req.Password != "" {
		password, _ := utils.Md5(req.Password + req.Username)
		entity.Password = password
	}

	// 头像处理
	if req.Avatar != "" {
		avatar, err := utils.SaveImage(req.Avatar, "user")
		if err != nil {
			return 0, err
		}
		entity.Avatar = avatar
	}
	entity.CreateUser = userId
	entity.CreateTime = time.Now()
	entity.Mark = 1
	// 更新记录
	return entity.Update()
}

func (s *userService) Delete(ids string) (int64, error) {
	// 记录ID
	idsArr := strings.Split(ids, ",")
	if len(idsArr) == 1 {
		// 单个删除
		entity := &model.User{Id: gconv.Int(ids)}
		rows, err := entity.Delete()
		if err != nil || rows == 0 {
			return 0, errors.New("删除失败")
		}
		return rows, nil
	} else {
		// 批量删除
		return 0, nil
	}
}

func (s *userService) Status(req *dto.UserStatusReq, userId int) (int64, error) {
	// 查询记录
	info := &model.User{Id: req.Id}
	has, err := info.Get()
	if err != nil || !has {
		return 0, errors.New("记录不存在")
	}

	entity := &model.User{}
	entity.Id = info.Id
	entity.Status = req.Status
	entity.UpdateUser = userId
	entity.UpdateTime = time.Now()
	return entity.Update()
}

func (s *userService) ResetPwd(id int, userId int) (int64, error) {
	if utils.AppDebug() {
		return 0, errors.New("演示环境，暂无权限操作")
	}
	// 查询记录
	info := &model.User{Id: id}
	has, err := info.Get()
	if err != nil || !has {
		return 0, err
	}
	if info == nil {
		return 0, errors.New("记录不存在")
	}
	// 设置初始密码
	password, err := utils.Md5("123456" + info.Username)
	if err != nil {
		return 0, err
	}

	// 初始化密码
	rows, err := utils.XormDb.Id(id).Update(&model.User{
		Password:   password,
		UpdateUser: userId,
		UpdateTime: time.Now(),
	})
	if err != nil {
		return 0, err
	}

	// 获取受影响行数
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func (s *userService) CheckUser(req *dto.CheckUserReq) (*model.User, error) {
	user := &model.User{Username: req.Username, Mark: 1}
	has, err := user.Get()
	if err != nil || !has {
		return nil, err
	}
	return user, nil
}