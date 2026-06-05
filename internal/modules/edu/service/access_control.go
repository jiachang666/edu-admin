package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	edumodel "edu-admin/internal/modules/edu/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	userStatusEnabled  = "启用"
	userStatusDisabled = "停用"
	roleStatusEnabled  = "启用"
	roleStatusDisabled = "停用"

	superAdminRoleCode = "super_admin"
)

var (
	ErrInvalidCredentials      = errors.New("账号或密码不正确")
	ErrUserDisabled            = errors.New("账号已停用")
	ErrUserNotFound            = errors.New("user not found")
	ErrRoleNotFound            = errors.New("role not found")
	ErrUsernameAlreadyExists   = errors.New("账号已存在")
	ErrRoleCodeAlreadyExists   = errors.New("角色编码已存在")
	ErrInvalidPermissionCodes  = errors.New("存在未识别的权限项")
	ErrCannotDisableCurrentUser = errors.New("不能停用当前登录账号")
	ErrProtectedRole          = errors.New("内置超级管理员角色不能停用或降权")
)

type Operator struct {
	UserID      uint64
	DisplayName string
}

type AuthProfile struct {
	ID           uint64   `json:"id"`
	Username     string   `json:"username"`
	DisplayName  string   `json:"displayName"`
	Mobile       string   `json:"mobile"`
	Status       string   `json:"status"`
	Roles        []string `json:"roles"`
	RoleNames    []string `json:"roleNames"`
	Permissions  []string `json:"permissions"`
	LastLoginAt  string   `json:"lastLoginAt"`
}

type UserRoleItem struct {
	ID   uint64 `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type UserItem struct {
	ID              uint64         `json:"id"`
	Username        string         `json:"username"`
	DisplayName     string         `json:"displayName"`
	Mobile          string         `json:"mobile"`
	Status          string         `json:"status"`
	LastLoginAt     string         `json:"lastLoginAt"`
	Roles           []UserRoleItem `json:"roles"`
	PrimaryRoleCode string         `json:"primaryRoleCode"`
	PrimaryRoleName string         `json:"primaryRoleName"`
}

type UserPayload struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
	Mobile      string `json:"mobile"`
	RoleCode    string `json:"roleCode"`
	Status      string `json:"status"`
}

type RoleItem struct {
	ID              uint64   `json:"id"`
	Name            string   `json:"name"`
	Code            string   `json:"code"`
	Description     string   `json:"description"`
	Status          string   `json:"status"`
	UserCount       int      `json:"userCount"`
	PermissionCount int      `json:"permissionCount"`
	Permissions     []string `json:"permissions"`
}

type RolePayload struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type RolePermissionPayload struct {
	Permissions []string `json:"permissions"`
}

type PermissionDefinition struct {
	Code        string `json:"code"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type PermissionGroup struct {
	Key         string                 `json:"key"`
	Label       string                 `json:"label"`
	Description string                 `json:"description"`
	Permissions []PermissionDefinition `json:"permissions"`
}

type OperationLogItem struct {
	ID          uint64 `json:"id"`
	UserID      uint64 `json:"userId"`
	UserName    string `json:"userName"`
	Module      string `json:"module"`
	Action      string `json:"action"`
	TargetType  string `json:"targetType"`
	TargetID    uint64 `json:"targetId"`
	Content     string `json:"content"`
	CreatedAt   string `json:"createdAt"`
}

type userRoleRecord struct {
	UserID      uint64 `gorm:"column:user_id"`
	RoleID      uint64 `gorm:"column:role_id"`
	RoleCode    string `gorm:"column:role_code"`
	RoleName    string `gorm:"column:role_name"`
}

type roleUserCountRecord struct {
	RoleID    uint64 `gorm:"column:role_id"`
	UserCount int    `gorm:"column:user_count"`
}

var permissionGroups = []PermissionGroup{
	{
		Key:         "system",
		Label:       "系统与权限",
		Description: "控制后台账号、角色和操作记录。",
		Permissions: []PermissionDefinition{
			{Code: "dashboard:view", Label: "查看首页总览", Description: "进入首页查看当天总体情况。"},
			{Code: "users:view", Label: "查看账号管理", Description: "查看后台账号列表和账号信息。"},
			{Code: "users:manage", Label: "维护账号", Description: "新增、编辑、启用和停用账号。"},
			{Code: "roles:view", Label: "查看角色权限", Description: "查看角色列表和权限配置。"},
			{Code: "roles:manage", Label: "维护角色权限", Description: "新增角色、编辑角色和保存权限。"},
			{Code: "operation_logs:view", Label: "查看操作记录", Description: "回看关键操作的执行人和时间。"},
		},
	},
	{
		Key:         "faculty",
		Label:       "老师与学员",
		Description: "管理老师、学员和班级基础台账。",
		Permissions: []PermissionDefinition{
			{Code: "teachers:view", Label: "查看老师", Description: "查看老师列表和师资分布。"},
			{Code: "teachers:manage", Label: "维护老师", Description: "新增和编辑老师资料。"},
			{Code: "students:view", Label: "查看学员", Description: "查看学员列表与学员详情。"},
			{Code: "students:manage", Label: "维护学员", Description: "新增和编辑学员资料，以及维护家长信息。"},
			{Code: "courses:view", Label: "查看课程", Description: "查看课程模板与课程信息。"},
			{Code: "courses:manage", Label: "维护课程", Description: "新增和编辑课程模板资料。"},
			{Code: "classes:view", Label: "查看班级", Description: "查看班级列表和班级详情。"},
			{Code: "classes:manage", Label: "维护班级", Description: "新增和编辑班级资料，以及调整班级学员。"},
		},
	},
	{
		Key:         "operations",
		Label:       "排课与课后流程",
		Description: "覆盖排课、签到、作业反馈和通知。",
		Permissions: []PermissionDefinition{
			{Code: "schedules:view", Label: "查看排课", Description: "查看课程安排和课程安排详情。"},
			{Code: "schedules:manage", Label: "维护排课", Description: "新建、编辑、调课、停课和补课。"},
			{Code: "attendance:view", Label: "查看签到", Description: "查看签到工作台与签到结果。"},
			{Code: "attendance:manage", Label: "维护签到", Description: "对课程安排完成签到保存。"},
			{Code: "homeworks:view", Label: "查看作业反馈", Description: "查看作业与课后反馈列表。"},
			{Code: "homeworks:manage", Label: "维护作业反馈", Description: "保存作业和班级反馈内容。"},
			{Code: "notices:view", Label: "查看通知", Description: "查看通知列表和通知详情。"},
			{Code: "notices:manage", Label: "维护通知", Description: "新建、编辑和发送通知。"},
		},
	},
}

func (s *Service) PermissionGroups() []PermissionGroup {
	groups := make([]PermissionGroup, 0, len(permissionGroups))
	for _, group := range permissionGroups {
		clonedPermissions := make([]PermissionDefinition, 0, len(group.Permissions))
		clonedPermissions = append(clonedPermissions, group.Permissions...)
		groups = append(groups, PermissionGroup{
			Key:         group.Key,
			Label:       group.Label,
			Description: group.Description,
			Permissions: clonedPermissions,
		})
	}

	return groups
}

func (s *Service) AuthenticateUser(username string, password string) (AuthProfile, error) {
	if s.db == nil {
		if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
			return AuthProfile{}, ErrInvalidCredentials
		}

		return AuthProfile{
			ID:          1,
			Username:    strings.TrimSpace(username),
			DisplayName: "System Admin",
			Status:      userStatusEnabled,
			Roles:       []string{superAdminRoleCode},
			RoleNames:   []string{"超级管理员"},
			Permissions: defaultRolePermissions(superAdminRoleCode),
		}, nil
	}

	username = strings.TrimSpace(username)
	if username == "" || strings.TrimSpace(password) == "" {
		return AuthProfile{}, ErrInvalidCredentials
	}

	var profile AuthProfile
	authenticateErr := s.db.Transaction(func(tx *gorm.DB) error {
		var user edumodel.User
		findErr := tx.Where("username = ?", username).Take(&user).Error
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return ErrInvalidCredentials
		}
		if findErr != nil {
			return findErr
		}

		if user.Status != userStatusEnabled {
			return ErrUserDisabled
		}

		compareErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if compareErr != nil {
			return ErrInvalidCredentials
		}

		now := time.Now()
		updateErr := tx.Model(&edumodel.User{}).
			Where("id = ?", user.ID).
			Update("last_login_at", now).Error
		if updateErr != nil {
			return updateErr
		}

		profile = buildAuthProfile(user, now, nil, nil)
		loadedProfile, profileErr := s.authProfileByUserIDTx(tx, user.ID)
		if profileErr != nil {
			return profileErr
		}
		profile = loadedProfile
		profile.LastLoginAt = now.Format(dateTimeLayout)

		return s.recordOperationTx(
			tx,
			Operator{UserID: user.ID, DisplayName: user.DisplayName},
			"auth",
			"login",
			"user",
			user.ID,
			fmt.Sprintf("账号 %s 登录了后台。", user.DisplayName),
		)
	})
	if authenticateErr != nil {
		return AuthProfile{}, authenticateErr
	}

	return profile, nil
}

func (s *Service) AuthProfileByUserID(userID uint64) (AuthProfile, bool, error) {
	if s.db == nil {
		if userID != 1 {
			return AuthProfile{}, false, nil
		}

		return AuthProfile{
			ID:          1,
			Username:    "admin",
			DisplayName: "System Admin",
			Status:      userStatusEnabled,
			Roles:       []string{superAdminRoleCode},
			RoleNames:   []string{"超级管理员"},
			Permissions: defaultRolePermissions(superAdminRoleCode),
		}, true, nil
	}

	profile, profileErr := s.authProfileByUserIDTx(s.db, userID)
	if errors.Is(profileErr, gorm.ErrRecordNotFound) {
		return AuthProfile{}, false, nil
	}
	if profileErr != nil {
		return AuthProfile{}, false, profileErr
	}

	return profile, true, nil
}

func (s *Service) Users() ([]UserItem, error) {
	if s.db == nil {
		return []UserItem{}, nil
	}

	var users []edumodel.User
	listErr := s.db.Order("id ASC").Find(&users).Error
	if listErr != nil {
		return nil, listErr
	}

	return s.buildUserItems(s.db, users)
}

func (s *Service) User(rawID string) (UserItem, bool, error) {
	if s.db == nil {
		return UserItem{}, false, nil
	}

	userID, parseErr := strconv.ParseUint(rawID, 10, 64)
	if parseErr != nil {
		return UserItem{}, false, nil
	}

	var user edumodel.User
	findErr := s.db.Where("id = ?", userID).Take(&user).Error
	if errors.Is(findErr, gorm.ErrRecordNotFound) {
		return UserItem{}, false, nil
	}
	if findErr != nil {
		return UserItem{}, false, findErr
	}

	items, buildErr := s.buildUserItems(s.db, []edumodel.User{user})
	if buildErr != nil {
		return UserItem{}, false, buildErr
	}
	if len(items) == 0 {
		return UserItem{}, false, nil
	}

	return items[0], true, nil
}

func (s *Service) CreateUser(payload UserPayload, operator Operator) (UserItem, error) {
	if s.db == nil {
		return UserItem{}, nil
	}

	passwordHash, hashErr := hashPassword(payload.Password)
	if hashErr != nil {
		return UserItem{}, hashErr
	}

	var createdID uint64
	createErr := s.db.Transaction(func(tx *gorm.DB) error {
		role, roleErr := findRoleByCode(tx, payload.RoleCode)
		if roleErr != nil {
			return roleErr
		}
		if role.Status != roleStatusEnabled {
			return ErrRoleNotFound
		}

		exists, existsErr := userExistsByUsername(tx, payload.Username, 0)
		if existsErr != nil {
			return existsErr
		}
		if exists {
			return ErrUsernameAlreadyExists
		}

		now := time.Now()
		user := edumodel.User{
			Username:     payload.Username,
			PasswordHash: passwordHash,
			DisplayName:  payload.DisplayName,
			Mobile:       payload.Mobile,
			Status:       payload.Status,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		createUserErr := tx.Create(&user).Error
		if createUserErr != nil {
			return createUserErr
		}

		userRole := edumodel.UserRole{
			UserID:    user.ID,
			RoleID:    role.ID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		createRoleErr := tx.Create(&userRole).Error
		if createRoleErr != nil {
			return createRoleErr
		}

		createdID = user.ID

		return s.recordOperationTx(
			tx,
			operator,
			"users",
			"create",
			"user",
			user.ID,
			fmt.Sprintf("创建账号 %s，并分配角色 %s。", user.DisplayName, role.Name),
		)
	})
	if createErr != nil {
		return UserItem{}, createErr
	}

	item, found, itemErr := s.User(strconv.FormatUint(createdID, 10))
	if itemErr != nil {
		return UserItem{}, itemErr
	}
	if !found {
		return UserItem{}, ErrUserNotFound
	}

	return item, nil
}

func (s *Service) UpdateUser(rawID string, payload UserPayload, operator Operator) (UserItem, bool, error) {
	if s.db == nil {
		return UserItem{}, false, nil
	}

	userID, parseErr := strconv.ParseUint(rawID, 10, 64)
	if parseErr != nil {
		return UserItem{}, false, nil
	}

	updateErr := s.db.Transaction(func(tx *gorm.DB) error {
		var user edumodel.User
		findErr := tx.Where("id = ?", userID).Take(&user).Error
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		if findErr != nil {
			return findErr
		}

		role, roleErr := findRoleByCode(tx, payload.RoleCode)
		if roleErr != nil {
			return roleErr
		}
		if role.Status != roleStatusEnabled {
			return ErrRoleNotFound
		}

		exists, existsErr := userExistsByUsername(tx, payload.Username, userID)
		if existsErr != nil {
			return existsErr
		}
		if exists {
			return ErrUsernameAlreadyExists
		}

		updates := map[string]any{
			"username":     payload.Username,
			"display_name": payload.DisplayName,
			"mobile":       payload.Mobile,
			"status":       payload.Status,
			"updated_at":   time.Now(),
		}

		if payload.Password != "" {
			passwordHash, hashErr := hashPassword(payload.Password)
			if hashErr != nil {
				return hashErr
			}
			updates["password_hash"] = passwordHash
		}

		saveErr := tx.Model(&edumodel.User{}).
			Where("id = ?", userID).
			Updates(updates).Error
		if saveErr != nil {
			return saveErr
		}

		replaceErr := replaceUserRole(tx, userID, role.ID)
		if replaceErr != nil {
			return replaceErr
		}

		return s.recordOperationTx(
			tx,
			operator,
			"users",
			"update",
			"user",
			userID,
			fmt.Sprintf("更新账号 %s，并调整角色为 %s。", payload.DisplayName, role.Name),
		)
	})
	if errors.Is(updateErr, ErrUserNotFound) {
		return UserItem{}, false, nil
	}
	if updateErr != nil {
		return UserItem{}, false, updateErr
	}

	item, found, itemErr := s.User(rawID)
	if itemErr != nil {
		return UserItem{}, false, itemErr
	}

	return item, found, nil
}

func (s *Service) EnableUser(rawID string, operator Operator) (UserItem, bool, error) {
	return s.updateUserStatus(rawID, userStatusEnabled, operator)
}

func (s *Service) DisableUser(rawID string, operator Operator) (UserItem, bool, error) {
	return s.updateUserStatus(rawID, userStatusDisabled, operator)
}

func (s *Service) Roles() ([]RoleItem, error) {
	if s.db == nil {
		return buildDefaultRoleItems(), nil
	}

	var roles []edumodel.Role
	listErr := s.db.Order("id ASC").Find(&roles).Error
	if listErr != nil {
		return nil, listErr
	}

	return s.buildRoleItems(s.db, roles)
}

func (s *Service) Role(rawID string) (RoleItem, bool, error) {
	if s.db == nil {
		for _, item := range buildDefaultRoleItems() {
			if strconv.FormatUint(item.ID, 10) == rawID {
				return item, true, nil
			}
		}
		return RoleItem{}, false, nil
	}

	roleID, parseErr := strconv.ParseUint(rawID, 10, 64)
	if parseErr != nil {
		return RoleItem{}, false, nil
	}

	var role edumodel.Role
	findErr := s.db.Where("id = ?", roleID).Take(&role).Error
	if errors.Is(findErr, gorm.ErrRecordNotFound) {
		return RoleItem{}, false, nil
	}
	if findErr != nil {
		return RoleItem{}, false, findErr
	}

	items, buildErr := s.buildRoleItems(s.db, []edumodel.Role{role})
	if buildErr != nil {
		return RoleItem{}, false, buildErr
	}
	if len(items) == 0 {
		return RoleItem{}, false, nil
	}

	return items[0], true, nil
}

func (s *Service) CreateRole(payload RolePayload, operator Operator) (RoleItem, error) {
	if s.db == nil {
		return RoleItem{}, nil
	}

	var createdID uint64
	createErr := s.db.Transaction(func(tx *gorm.DB) error {
		exists, existsErr := roleExistsByCode(tx, payload.Code, 0)
		if existsErr != nil {
			return existsErr
		}
		if exists {
			return ErrRoleCodeAlreadyExists
		}

		now := time.Now()
		role := edumodel.Role{
			Name:                payload.Name,
			Code:                payload.Code,
			Description:         payload.Description,
			Status:              payload.Status,
			PermissionCodesJSON: serializePermissionCodes(defaultRolePermissions(payload.Code)),
			CreatedAt:           now,
			UpdatedAt:           now,
		}

		createRoleErr := tx.Create(&role).Error
		if createRoleErr != nil {
			return createRoleErr
		}

		createdID = role.ID

		return s.recordOperationTx(
			tx,
			operator,
			"roles",
			"create",
			"role",
			role.ID,
			fmt.Sprintf("创建角色 %s。", role.Name),
		)
	})
	if createErr != nil {
		return RoleItem{}, createErr
	}

	item, found, itemErr := s.Role(strconv.FormatUint(createdID, 10))
	if itemErr != nil {
		return RoleItem{}, itemErr
	}
	if !found {
		return RoleItem{}, ErrRoleNotFound
	}

	return item, nil
}

func (s *Service) UpdateRole(rawID string, payload RolePayload, operator Operator) (RoleItem, bool, error) {
	if s.db == nil {
		return RoleItem{}, false, nil
	}

	roleID, parseErr := strconv.ParseUint(rawID, 10, 64)
	if parseErr != nil {
		return RoleItem{}, false, nil
	}

	updateErr := s.db.Transaction(func(tx *gorm.DB) error {
		var role edumodel.Role
		findErr := tx.Where("id = ?", roleID).Take(&role).Error
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		if findErr != nil {
			return findErr
		}

		if role.Code == superAdminRoleCode && (payload.Status != roleStatusEnabled || payload.Code != superAdminRoleCode) {
			return ErrProtectedRole
		}

		exists, existsErr := roleExistsByCode(tx, payload.Code, roleID)
		if existsErr != nil {
			return existsErr
		}
		if exists {
			return ErrRoleCodeAlreadyExists
		}

		updateMap := map[string]any{
			"name":        payload.Name,
			"code":        payload.Code,
			"description": payload.Description,
			"status":      payload.Status,
			"updated_at":  time.Now(),
		}

		saveErr := tx.Model(&edumodel.Role{}).
			Where("id = ?", roleID).
			Updates(updateMap).Error
		if saveErr != nil {
			return saveErr
		}

		return s.recordOperationTx(
			tx,
			operator,
			"roles",
			"update",
			"role",
			roleID,
			fmt.Sprintf("更新角色 %s。", payload.Name),
		)
	})
	if errors.Is(updateErr, ErrRoleNotFound) {
		return RoleItem{}, false, nil
	}
	if updateErr != nil {
		return RoleItem{}, false, updateErr
	}

	item, found, itemErr := s.Role(rawID)
	if itemErr != nil {
		return RoleItem{}, false, itemErr
	}

	return item, found, nil
}

func (s *Service) SaveRolePermissions(rawID string, payload RolePermissionPayload, operator Operator) (RoleItem, bool, error) {
	if s.db == nil {
		return RoleItem{}, false, nil
	}

	roleID, parseErr := strconv.ParseUint(rawID, 10, 64)
	if parseErr != nil {
		return RoleItem{}, false, nil
	}

	normalizedPermissions, valid := normalizePermissionCodes(payload.Permissions)
	if !valid {
		return RoleItem{}, false, ErrInvalidPermissionCodes
	}

	saveErr := s.db.Transaction(func(tx *gorm.DB) error {
		var role edumodel.Role
		findErr := tx.Where("id = ?", roleID).Take(&role).Error
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		if findErr != nil {
			return findErr
		}

		if role.Code == superAdminRoleCode && !sameStringSlice(normalizedPermissions, normalizePermissionCodesMust(defaultRolePermissions(superAdminRoleCode))) {
			return ErrProtectedRole
		}

		updateErr := tx.Model(&edumodel.Role{}).
			Where("id = ?", roleID).
			Updates(map[string]any{
				"permission_codes_json": serializePermissionCodes(normalizedPermissions),
				"updated_at":            time.Now(),
			}).Error
		if updateErr != nil {
			return updateErr
		}

		return s.recordOperationTx(
			tx,
			operator,
			"roles",
			"save_permissions",
			"role",
			roleID,
			fmt.Sprintf("更新角色 %s 的权限配置。", role.Name),
		)
	})
	if errors.Is(saveErr, ErrRoleNotFound) {
		return RoleItem{}, false, nil
	}
	if saveErr != nil {
		return RoleItem{}, false, saveErr
	}

	item, found, itemErr := s.Role(rawID)
	if itemErr != nil {
		return RoleItem{}, false, itemErr
	}

	return item, found, nil
}

func (s *Service) OperationLogs() ([]OperationLogItem, error) {
	if s.db == nil {
		return []OperationLogItem{}, nil
	}

	var logs []edumodel.OperationLog
	listErr := s.db.Order("id DESC").Find(&logs).Error
	if listErr != nil {
		return nil, listErr
	}

	items := make([]OperationLogItem, 0, len(logs))
	for _, logItem := range logs {
		items = append(items, OperationLogItem{
			ID:         logItem.ID,
			UserID:     logItem.UserID,
			UserName:   logItem.UserDisplayName,
			Module:     logItem.Module,
			Action:     logItem.Action,
			TargetType: logItem.TargetType,
			TargetID:   logItem.TargetID,
			Content:    logItem.Content,
			CreatedAt:  logItem.CreatedAt.Format(dateTimeLayout),
		})
	}

	return items, nil
}

func (s *Service) seedAccessControlIfEmpty(tx *gorm.DB, now time.Time) error {
	var roleCount int64
	roleCountErr := tx.Model(&edumodel.Role{}).Count(&roleCount).Error
	if roleCountErr != nil {
		return roleCountErr
	}

	if roleCount == 0 {
		roles := []edumodel.Role{
			{
				Name:                "超级管理员",
				Code:                superAdminRoleCode,
				Description:         "负责系统初始化、账号管理和全部模块配置。",
				Status:              roleStatusEnabled,
				PermissionCodesJSON: serializePermissionCodes(defaultRolePermissions(superAdminRoleCode)),
				CreatedAt:           now,
				UpdatedAt:           now,
			},
			{
				Name:                "机构负责人",
				Code:                "campus_owner",
				Description:         "查看全局业务情况，重点关注老师、学员、班级和通知。",
				Status:              roleStatusEnabled,
				PermissionCodesJSON: serializePermissionCodes(defaultRolePermissions("campus_owner")),
				CreatedAt:           now,
				UpdatedAt:           now,
			},
			{
				Name:                "教务/前台",
				Code:                "front_desk",
				Description:         "处理日常录入、排课、签到和通知相关工作。",
				Status:              roleStatusEnabled,
				PermissionCodesJSON: serializePermissionCodes(defaultRolePermissions("front_desk")),
				CreatedAt:           now,
				UpdatedAt:           now,
			},
			{
				Name:                "老师",
				Code:                "teacher",
				Description:         "围绕排课、签到和课后反馈完成自己的教学任务。",
				Status:              roleStatusEnabled,
				PermissionCodesJSON: serializePermissionCodes(defaultRolePermissions("teacher")),
				CreatedAt:           now,
				UpdatedAt:           now,
			},
		}

		createRoleErr := tx.Create(&roles).Error
		if createRoleErr != nil {
			return createRoleErr
		}
	}

	syncBuiltinRoleErr := syncBuiltinRoles(tx, now)
	if syncBuiltinRoleErr != nil {
		return syncBuiltinRoleErr
	}

	var userCount int64
	userCountErr := tx.Model(&edumodel.User{}).Count(&userCount).Error
	if userCountErr != nil {
		return userCountErr
	}

	if userCount == 0 {
		adminHash, adminHashErr := hashPassword("123456")
		if adminHashErr != nil {
			return adminHashErr
		}
		bossHash, bossHashErr := hashPassword("123456")
		if bossHashErr != nil {
			return bossHashErr
		}
		staffHash, staffHashErr := hashPassword("123456")
		if staffHashErr != nil {
			return staffHashErr
		}
		teacherHash, teacherHashErr := hashPassword("123456")
		if teacherHashErr != nil {
			return teacherHashErr
		}

		users := []edumodel.User{
			{
				Username:     "admin",
				PasswordHash: adminHash,
				DisplayName:  "系统管理员",
				Mobile:       "13800001000",
				Status:       userStatusEnabled,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			{
				Username:     "boss",
				PasswordHash: bossHash,
				DisplayName:  "机构负责人",
				Mobile:       "13800001001",
				Status:       userStatusEnabled,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			{
				Username:     "staff",
				PasswordHash: staffHash,
				DisplayName:  "教务前台",
				Mobile:       "13800001002",
				Status:       userStatusEnabled,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			{
				Username:     "teacher",
				PasswordHash: teacherHash,
				DisplayName:  "周老师",
				Mobile:       "13800001003",
				Status:       userStatusEnabled,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		}

		createUserErr := tx.Create(&users).Error
		if createUserErr != nil {
			return createUserErr
		}
	}

	roleMap, roleMapErr := loadRoleCodeMap(tx)
	if roleMapErr != nil {
		return roleMapErr
	}

	userMap, userMapErr := loadUsernameMap(tx)
	if userMapErr != nil {
		return userMapErr
	}

	var userRoleCount int64
	userRoleCountErr := tx.Model(&edumodel.UserRole{}).Count(&userRoleCount).Error
	if userRoleCountErr != nil {
		return userRoleCountErr
	}

	if userRoleCount == 0 {
		assignments := []struct {
			username string
			roleCode string
		}{
			{username: "admin", roleCode: superAdminRoleCode},
			{username: "boss", roleCode: "campus_owner"},
			{username: "staff", roleCode: "front_desk"},
			{username: "teacher", roleCode: "teacher"},
		}

		userRoles := make([]edumodel.UserRole, 0, len(assignments))
		for _, assignment := range assignments {
			user, userFound := userMap[assignment.username]
			role, roleFound := roleMap[assignment.roleCode]
			if !userFound || !roleFound {
				continue
			}

			userRoles = append(userRoles, edumodel.UserRole{
				UserID:    user.ID,
				RoleID:    role.ID,
				CreatedAt: now,
				UpdatedAt: now,
			})
		}

		if len(userRoles) > 0 {
			createUserRoleErr := tx.Create(&userRoles).Error
			if createUserRoleErr != nil {
				return createUserRoleErr
			}
		}
	}

	var logCount int64
	logCountErr := tx.Model(&edumodel.OperationLog{}).Count(&logCount).Error
	if logCountErr != nil {
		return logCountErr
	}

	if logCount == 0 {
		adminUser, adminFound := userMap["admin"]
		if adminFound {
			seedLogs := []edumodel.OperationLog{
				{
					UserID:          adminUser.ID,
					UserDisplayName: adminUser.DisplayName,
					Module:          "users",
					Action:          "create",
					TargetType:      "user",
					TargetID:        adminUser.ID,
					Content:         "初始化演示账号 admin。",
					CreatedAt:       now.Add(-8 * time.Hour),
				},
				{
					UserID:          adminUser.ID,
					UserDisplayName: adminUser.DisplayName,
					Module:          "roles",
					Action:          "save_permissions",
					TargetType:      "role",
					TargetID:        roleMap[superAdminRoleCode].ID,
					Content:         "初始化默认角色权限模板。",
					CreatedAt:       now.Add(-7 * time.Hour),
				},
				{
					UserID:          adminUser.ID,
					UserDisplayName: adminUser.DisplayName,
					Module:          "auth",
					Action:          "login",
					TargetType:      "user",
					TargetID:        adminUser.ID,
					Content:         "系统管理员已完成首次登录检查。",
					CreatedAt:       now.Add(-6 * time.Hour),
				},
			}

			createLogErr := tx.Create(&seedLogs).Error
			if createLogErr != nil {
				return createLogErr
			}
		}
	}

	return nil
}

func (s *Service) authProfileByUserIDTx(tx *gorm.DB, userID uint64) (AuthProfile, error) {
	var user edumodel.User
	findErr := tx.Where("id = ?", userID).Take(&user).Error
	if findErr != nil {
		return AuthProfile{}, findErr
	}

	assignments, assignmentErr := loadUserRolesForUsers(tx, []uint64{user.ID})
	if assignmentErr != nil {
		return AuthProfile{}, assignmentErr
	}

	roleIDs := make([]uint64, 0, len(assignments))
	for _, assignment := range assignments {
		roleIDs = append(roleIDs, assignment.RoleID)
	}

	roleMap, roleErr := loadRolesByID(tx, roleIDs)
	if roleErr != nil {
		return AuthProfile{}, roleErr
	}

	roles := make([]string, 0, len(assignments))
	roleNames := make([]string, 0, len(assignments))
	permissions := make([]string, 0, 12)

	for _, assignment := range assignments {
		role, found := roleMap[assignment.RoleID]
		if !found {
			continue
		}

		if role.Status != roleStatusEnabled {
			continue
		}

		roles = append(roles, role.Code)
		roleNames = append(roleNames, role.Name)
		permissions = append(permissions, permissionCodesFromRole(role)...)
	}

	return AuthProfile{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Mobile:      user.Mobile,
		Status:      user.Status,
		Roles:       uniqueSortedStrings(roles),
		RoleNames:   uniqueSortedStrings(roleNames),
		Permissions: uniqueSortedStrings(permissions),
		LastLoginAt: formatTimePointer(user.LastLoginAt),
	}, nil
}

func (s *Service) buildUserItems(tx *gorm.DB, users []edumodel.User) ([]UserItem, error) {
	if len(users) == 0 {
		return []UserItem{}, nil
	}

	userIDs := make([]uint64, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	assignments, assignmentErr := loadUserRolesForUsers(tx, userIDs)
	if assignmentErr != nil {
		return nil, assignmentErr
	}

	roleIDs := make([]uint64, 0, len(assignments))
	for _, assignment := range assignments {
		roleIDs = append(roleIDs, assignment.RoleID)
	}

	roleMap, roleErr := loadRolesByID(tx, roleIDs)
	if roleErr != nil {
		return nil, roleErr
	}

	assignmentsByUserID := make(map[uint64][]edumodel.UserRole)
	for _, assignment := range assignments {
		assignmentsByUserID[assignment.UserID] = append(assignmentsByUserID[assignment.UserID], assignment)
	}

	items := make([]UserItem, 0, len(users))
	for _, user := range users {
		roles := make([]UserRoleItem, 0, len(assignmentsByUserID[user.ID]))
		for _, assignment := range assignmentsByUserID[user.ID] {
			role, found := roleMap[assignment.RoleID]
			if !found {
				continue
			}
			roles = append(roles, UserRoleItem{
				ID:   role.ID,
				Code: role.Code,
				Name: role.Name,
			})
		}
		sort.Slice(roles, func(i int, j int) bool {
			return roles[i].ID < roles[j].ID
		})

		item := UserItem{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Mobile:      user.Mobile,
			Status:      user.Status,
			LastLoginAt: formatTimePointer(user.LastLoginAt),
			Roles:       roles,
		}
		if len(roles) > 0 {
			item.PrimaryRoleCode = roles[0].Code
			item.PrimaryRoleName = roles[0].Name
		}

		items = append(items, item)
	}

	return items, nil
}

func (s *Service) buildRoleItems(tx *gorm.DB, roles []edumodel.Role) ([]RoleItem, error) {
	if len(roles) == 0 {
		return []RoleItem{}, nil
	}

	roleIDs := make([]uint64, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	var userCounts []roleUserCountRecord
	countErr := tx.Model(&edumodel.UserRole{}).
		Select("role_id, COUNT(*) AS user_count").
		Where("role_id IN ?", roleIDs).
		Group("role_id").
		Scan(&userCounts).Error
	if countErr != nil {
		return nil, countErr
	}

	countMap := make(map[uint64]int, len(userCounts))
	for _, countItem := range userCounts {
		countMap[countItem.RoleID] = countItem.UserCount
	}

	items := make([]RoleItem, 0, len(roles))
	for _, role := range roles {
		permissions := permissionCodesFromRole(role)
		items = append(items, RoleItem{
			ID:              role.ID,
			Name:            role.Name,
			Code:            role.Code,
			Description:     role.Description,
			Status:          role.Status,
			UserCount:       countMap[role.ID],
			PermissionCount: len(permissions),
			Permissions:     permissions,
		})
	}

	return items, nil
}

func (s *Service) updateUserStatus(rawID string, targetStatus string, operator Operator) (UserItem, bool, error) {
	if s.db == nil {
		return UserItem{}, false, nil
	}

	userID, parseErr := strconv.ParseUint(rawID, 10, 64)
	if parseErr != nil {
		return UserItem{}, false, nil
	}

	statusErr := s.db.Transaction(func(tx *gorm.DB) error {
		var user edumodel.User
		findErr := tx.Where("id = ?", userID).Take(&user).Error
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		if findErr != nil {
			return findErr
		}

		if targetStatus == userStatusDisabled && operator.UserID == userID {
			return ErrCannotDisableCurrentUser
		}

		updateErr := tx.Model(&edumodel.User{}).
			Where("id = ?", userID).
			Updates(map[string]any{
				"status":     targetStatus,
				"updated_at": time.Now(),
			}).Error
		if updateErr != nil {
			return updateErr
		}

		action := "enable"
		content := fmt.Sprintf("启用账号 %s。", user.DisplayName)
		if targetStatus == userStatusDisabled {
			action = "disable"
			content = fmt.Sprintf("停用账号 %s。", user.DisplayName)
		}

		return s.recordOperationTx(tx, operator, "users", action, "user", userID, content)
	})
	if errors.Is(statusErr, ErrUserNotFound) {
		return UserItem{}, false, nil
	}
	if statusErr != nil {
		return UserItem{}, false, statusErr
	}

	item, found, itemErr := s.User(rawID)
	if itemErr != nil {
		return UserItem{}, false, itemErr
	}

	return item, found, nil
}

func (s *Service) recordOperationTx(
	tx *gorm.DB,
	operator Operator,
	module string,
	action string,
	targetType string,
	targetID uint64,
	content string,
) error {
	if tx == nil {
		return nil
	}

	displayName := strings.TrimSpace(operator.DisplayName)
	if displayName == "" {
		displayName = "系统"
	}

	logItem := edumodel.OperationLog{
		UserID:          operator.UserID,
		UserDisplayName: displayName,
		Module:          module,
		Action:          action,
		TargetType:      targetType,
		TargetID:        targetID,
		Content:         content,
		CreatedAt:       time.Now(),
	}

	return tx.Create(&logItem).Error
}

func permissionCodesFromRole(role edumodel.Role) []string {
	permissions := deserializePermissionCodes(role.PermissionCodesJSON)
	if len(permissions) == 0 && role.Code != "" {
		return defaultRolePermissions(role.Code)
	}

	return permissions
}

func deserializePermissionCodes(raw string) []string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return []string{}
	}

	var permissions []string
	if json.Unmarshal([]byte(trimmed), &permissions) == nil {
		return uniqueSortedStrings(permissions)
	}

	return uniqueSortedStrings(strings.Split(trimmed, ","))
}

func serializePermissionCodes(permissions []string) string {
	normalized := uniqueSortedStrings(permissions)
	encoded, marshalErr := json.Marshal(normalized)
	if marshalErr != nil {
		return "[]"
	}

	return string(encoded)
}

func normalizePermissionCodes(input []string) ([]string, bool) {
	validCodes := permissionCodeSet()
	normalized := make([]string, 0, len(input))
	for _, item := range input {
		code := strings.TrimSpace(item)
		if code == "" {
			continue
		}
		if !validCodes[code] {
			return nil, false
		}
		normalized = append(normalized, code)
	}

	return uniqueSortedStrings(normalized), true
}

func normalizePermissionCodesMust(input []string) []string {
	normalized, _ := normalizePermissionCodes(input)
	return normalized
}

func permissionCodeSet() map[string]bool {
	codeMap := make(map[string]bool)
	for _, group := range permissionGroups {
		for _, permission := range group.Permissions {
			codeMap[permission.Code] = true
		}
	}

	return codeMap
}

func defaultRolePermissions(roleCode string) []string {
	switch roleCode {
	case superAdminRoleCode:
		allPermissions := make([]string, 0, 16)
		for _, group := range permissionGroups {
			for _, permission := range group.Permissions {
				allPermissions = append(allPermissions, permission.Code)
			}
		}
		return uniqueSortedStrings(allPermissions)
	case "campus_owner":
		return uniqueSortedStrings([]string{
			"dashboard:view",
			"teachers:view",
			"teachers:manage",
			"students:view",
			"students:manage",
			"courses:view",
			"courses:manage",
			"classes:view",
			"classes:manage",
			"schedules:view",
			"attendance:view",
			"homeworks:view",
			"notices:view",
			"operation_logs:view",
		})
	case "front_desk":
		return uniqueSortedStrings([]string{
			"dashboard:view",
			"students:view",
			"students:manage",
			"courses:view",
			"courses:manage",
			"classes:view",
			"classes:manage",
			"teachers:view",
			"schedules:view",
			"schedules:manage",
			"attendance:view",
			"attendance:manage",
			"homeworks:view",
			"notices:view",
			"notices:manage",
		})
	case "teacher":
		return uniqueSortedStrings([]string{
			"dashboard:view",
			"schedules:view",
			"attendance:view",
			"attendance:manage",
			"homeworks:view",
			"homeworks:manage",
			"notices:view",
		})
	default:
		return []string{"dashboard:view"}
	}
}

func buildDefaultRoleItems() []RoleItem {
	roleCodes := []struct {
		ID          uint64
		Name        string
		Code        string
		Description string
	}{
		{ID: 1, Name: "超级管理员", Code: superAdminRoleCode, Description: "负责系统初始化、账号管理和全部模块配置。"},
		{ID: 2, Name: "机构负责人", Code: "campus_owner", Description: "查看全局业务情况，重点关注老师、学员、班级和通知。"},
		{ID: 3, Name: "教务/前台", Code: "front_desk", Description: "处理日常录入、排课、签到和通知相关工作。"},
		{ID: 4, Name: "老师", Code: "teacher", Description: "围绕排课、签到和课后反馈完成自己的教学任务。"},
	}

	items := make([]RoleItem, 0, len(roleCodes))
	for _, role := range roleCodes {
		permissions := defaultRolePermissions(role.Code)
		items = append(items, RoleItem{
			ID:              role.ID,
			Name:            role.Name,
			Code:            role.Code,
			Description:     role.Description,
			Status:          roleStatusEnabled,
			UserCount:       0,
			PermissionCount: len(permissions),
			Permissions:     permissions,
		})
	}

	return items
}

func syncBuiltinRoles(tx *gorm.DB, now time.Time) error {
	roleDefinitions := []struct {
		Name        string
		Code        string
		Description string
	}{
		{Name: "超级管理员", Code: superAdminRoleCode, Description: "负责系统初始化、账号管理和全部模块配置。"},
		{Name: "机构负责人", Code: "campus_owner", Description: "查看全局业务情况，重点关注老师、学员、班级和通知。"},
		{Name: "教务/前台", Code: "front_desk", Description: "处理日常录入、排课、签到和通知相关工作。"},
		{Name: "老师", Code: "teacher", Description: "围绕排课、签到和课后反馈完成自己的教学任务。"},
	}

	roleMap, roleMapErr := loadRoleCodeMap(tx)
	if roleMapErr != nil {
		return roleMapErr
	}

	for _, definition := range roleDefinitions {
		role, exists := roleMap[definition.Code]
		if !exists {
			role = edumodel.Role{
				Name:                definition.Name,
				Code:                definition.Code,
				Description:         definition.Description,
				Status:              roleStatusEnabled,
				PermissionCodesJSON: serializePermissionCodes(defaultRolePermissions(definition.Code)),
				CreatedAt:           now,
				UpdatedAt:           now,
			}

			createRoleErr := tx.Create(&role).Error
			if createRoleErr != nil {
				return createRoleErr
			}

			roleMap[definition.Code] = role
			continue
		}

		targetPermissions := defaultRolePermissions(definition.Code)
		if definition.Code != superAdminRoleCode {
			targetPermissions = uniqueSortedStrings(append(permissionCodesFromRole(role), targetPermissions...))
		}

		updateMap := map[string]any{
			"name":                  definition.Name,
			"description":           definition.Description,
			"permission_codes_json": serializePermissionCodes(targetPermissions),
			"updated_at":            now,
		}

		updateErr := tx.Model(&edumodel.Role{}).
			Where("id = ?", role.ID).
			Updates(updateMap).Error
		if updateErr != nil {
			return updateErr
		}
	}

	return nil
}

func loadUserRolesForUsers(tx *gorm.DB, userIDs []uint64) ([]edumodel.UserRole, error) {
	if len(userIDs) == 0 {
		return []edumodel.UserRole{}, nil
	}

	var assignments []edumodel.UserRole
	findErr := tx.Where("user_id IN ?", userIDs).Find(&assignments).Error
	if findErr != nil {
		return nil, findErr
	}

	return assignments, nil
}

func loadRolesByID(tx *gorm.DB, roleIDs []uint64) (map[uint64]edumodel.Role, error) {
	if len(roleIDs) == 0 {
		return map[uint64]edumodel.Role{}, nil
	}

	var roles []edumodel.Role
	findErr := tx.Where("id IN ?", uniqueUint64s(roleIDs)).Find(&roles).Error
	if findErr != nil {
		return nil, findErr
	}

	roleMap := make(map[uint64]edumodel.Role, len(roles))
	for _, role := range roles {
		roleMap[role.ID] = role
	}

	return roleMap, nil
}

func loadRoleCodeMap(tx *gorm.DB) (map[string]edumodel.Role, error) {
	var roles []edumodel.Role
	findErr := tx.Find(&roles).Error
	if findErr != nil {
		return nil, findErr
	}

	roleMap := make(map[string]edumodel.Role, len(roles))
	for _, role := range roles {
		roleMap[role.Code] = role
	}

	return roleMap, nil
}

func loadUsernameMap(tx *gorm.DB) (map[string]edumodel.User, error) {
	var users []edumodel.User
	findErr := tx.Find(&users).Error
	if findErr != nil {
		return nil, findErr
	}

	userMap := make(map[string]edumodel.User, len(users))
	for _, user := range users {
		userMap[user.Username] = user
	}

	return userMap, nil
}

func userExistsByUsername(tx *gorm.DB, username string, excludeID uint64) (bool, error) {
	query := tx.Model(&edumodel.User{}).Where("username = ?", username)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}

	var count int64
	countErr := query.Count(&count).Error
	if countErr != nil {
		return false, countErr
	}

	return count > 0, nil
}

func roleExistsByCode(tx *gorm.DB, code string, excludeID uint64) (bool, error) {
	query := tx.Model(&edumodel.Role{}).Where("code = ?", code)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}

	var count int64
	countErr := query.Count(&count).Error
	if countErr != nil {
		return false, countErr
	}

	return count > 0, nil
}

func replaceUserRole(tx *gorm.DB, userID uint64, roleID uint64) error {
	deleteErr := tx.Where("user_id = ?", userID).Delete(&edumodel.UserRole{}).Error
	if deleteErr != nil {
		return deleteErr
	}

	now := time.Now()
	userRole := edumodel.UserRole{
		UserID:    userID,
		RoleID:    roleID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return tx.Create(&userRole).Error
}

func findRoleByCode(tx *gorm.DB, roleCode string) (edumodel.Role, error) {
	var role edumodel.Role
	findErr := tx.Where("code = ?", roleCode).Take(&role).Error
	if errors.Is(findErr, gorm.ErrRecordNotFound) {
		return edumodel.Role{}, ErrRoleNotFound
	}
	if findErr != nil {
		return edumodel.Role{}, findErr
	}

	return role, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		return "", hashErr
	}

	return string(hashedPassword), nil
}

func formatTimePointer(value *time.Time) string {
	if value == nil {
		return ""
	}

	return value.Format(dateTimeLayout)
}

func uniqueUint64s(input []uint64) []uint64 {
	uniqueMap := make(map[uint64]struct{}, len(input))
	result := make([]uint64, 0, len(input))
	for _, item := range input {
		if _, exists := uniqueMap[item]; exists {
			continue
		}
		uniqueMap[item] = struct{}{}
		result = append(result, item)
	}

	sort.Slice(result, func(i int, j int) bool {
		return result[i] < result[j]
	})

	return result
}

func uniqueSortedStrings(input []string) []string {
	uniqueMap := make(map[string]struct{}, len(input))
	result := make([]string, 0, len(input))
	for _, item := range input {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if _, exists := uniqueMap[trimmed]; exists {
			continue
		}
		uniqueMap[trimmed] = struct{}{}
		result = append(result, trimmed)
	}

	sort.Strings(result)
	return result
}

func sameStringSlice(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}

	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}

	return true
}

func buildAuthProfile(
	user edumodel.User,
	lastLoginAt time.Time,
	roleCodes []string,
	permissions []string,
) AuthProfile {
	return AuthProfile{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Mobile:      user.Mobile,
		Status:      user.Status,
		Roles:       uniqueSortedStrings(roleCodes),
		Permissions: uniqueSortedStrings(permissions),
		LastLoginAt: lastLoginAt.Format(dateTimeLayout),
	}
}
