package service

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	eduservice "edu-admin/internal/modules/edu/service"
)

var (
	ErrInvalidCredentials = eduservice.ErrInvalidCredentials
	ErrUserDisabled       = eduservice.ErrUserDisabled
)

type Service struct {
	devToken   string
	eduService *eduservice.Service
}

func New(devToken string, eduService *eduservice.Service) *Service {
	return &Service{
		devToken:   devToken,
		eduService: eduService,
	}
}

func (s *Service) Login(username string, password string) (map[string]any, error) {
	profile, profileErr := s.eduService.AuthenticateUser(username, password)
	if profileErr != nil {
		return nil, profileErr
	}

	return map[string]any{
		"accessToken": s.buildAccessToken(profile.ID),
		"expiresIn":   7200,
		"user": map[string]any{
			"id":          profile.ID,
			"username":    profile.Username,
			"displayName": profile.DisplayName,
		},
		"roles":       profile.Roles,
		"roleNames":   profile.RoleNames,
		"permissions": profile.Permissions,
	}, nil
}

func (s *Service) Me(userID uint64) (map[string]any, bool, error) {
	profile, found, profileErr := s.eduService.AuthProfileByUserID(userID)
	if profileErr != nil {
		return nil, false, profileErr
	}
	if !found {
		return nil, false, nil
	}

	return map[string]any{
		"id":          profile.ID,
		"username":    profile.Username,
		"displayName": profile.DisplayName,
		"mobile":      profile.Mobile,
		"status":      profile.Status,
		"roles":       profile.Roles,
		"roleNames":   profile.RoleNames,
		"lastLoginAt": profile.LastLoginAt,
	}, true, nil
}

func (s *Service) Permissions(userID uint64) (map[string]any, bool, error) {
	profile, found, profileErr := s.eduService.AuthProfileByUserID(userID)
	if profileErr != nil {
		return nil, false, profileErr
	}
	if !found {
		return nil, false, nil
	}

	return map[string]any{
		"permissions":      profile.Permissions,
		"permissionGroups": s.eduService.PermissionGroups(),
	}, true, nil
}

func (s *Service) AuthProfile(userID uint64) (eduservice.AuthProfile, bool, error) {
	return s.eduService.AuthProfileByUserID(userID)
}

func (s *Service) buildAccessToken(userID uint64) string {
	payload := fmt.Sprintf("%s:%d", s.devToken, userID)
	return base64.StdEncoding.EncodeToString([]byte(payload))
}

func (s *Service) ParseAccessToken(token string) (uint64, bool) {
	decoded, decodeErr := base64.StdEncoding.DecodeString(strings.TrimSpace(token))
	if decodeErr != nil {
		return 0, false
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 || parts[0] != s.devToken {
		return 0, false
	}

	userID, parseErr := strconv.ParseUint(parts[1], 10, 64)
	if parseErr != nil || userID == 0 {
		return 0, false
	}

	return userID, true
}
