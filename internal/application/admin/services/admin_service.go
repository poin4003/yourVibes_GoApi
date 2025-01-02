package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/query"
)

type (
	IAdminAuth interface {
		Login(ctx context.Context, command *command.LoginCommand) (result *command.LoginCommandResult, err error)
		ChangeAdminPassword(ctx context.Context, command *command.ChangeAdminPasswordCommand) (result *command.ChangeAdminPasswordCommandResult, err error)
		ForgotAdminPassword(ctx context.Context, command *command.ForgotAdminPasswordCommand) (result *command.ForgotAdminPasswordCommandResult, err error)
	}
	IAdminInfo interface {
		UpdateAdmin(ctx context.Context, command *command.UpdateAdminInfoCommand) (result *command.UpdateAdminInfoCommandResult, err error)
		GetAdminStatusById(ctx context.Context, id uuid.UUID) (status bool, err error)
	}
	ISuperAdmin interface {
		CreateAdmin(ctx context.Context, command *command.CreateAdminCommand) (result *command.CreateAdminCommandResult, err error)
		UpdateAdmin(ctx context.Context, command *command.UpdateAdminForSuperAdminCommand) (result *command.UpdateAdminForSuperAdminCommandResult, err error)
		GetOneAdmin(ctx context.Context, query *query.GetOneAdminQuery) (result *query.AdminQueryResult, err error)
		GetManyAdmin(ctx context.Context, query *query.GetManyAdminQuery) (result *query.AdminQueryListResult, err error)
	}
)

var (
	localAdminAuth  IAdminAuth
	localAdminInfo  IAdminInfo
	localSuperAdmin ISuperAdmin
)

func AdminAuth() IAdminAuth {
	if localAdminAuth == nil {
		panic("repository_implement localAdminAuth not found for interface IAdminAuth")
	}

	return localAdminAuth
}

func AdminInfo() IAdminInfo {
	if localAdminInfo == nil {
		panic("repository_implement localAdminInfo not found for interface IAdminInfo")
	}

	return localAdminInfo
}

func SuperAdmin() ISuperAdmin {
	if localSuperAdmin == nil {
		panic("repository_implement localSuperAdmin not found for interface IAdminInfo")
	}

	return localSuperAdmin
}

func InitAdminAuth(i IAdminAuth) {
	localAdminAuth = i
}

func InitAdminInfo(i IAdminInfo) {
	localAdminInfo = i
}

func InitSuperAdmin(i ISuperAdmin) {
	localSuperAdmin = i
}
