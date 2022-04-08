package dao

import (
	"context"

	"github.com/RafayLabs/rcloud-base/internal/models"
	userv3 "github.com/RafayLabs/rcloud-base/proto/types/userpb/v3"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func GetProjectOrganization(ctx context.Context, db bun.IDB, id uuid.UUID) (string, string, error) {
	type projectOrg struct {
		Project      string
		Organization string
	}
	var r projectOrg
	err := db.NewSelect().Table("authsrv_project").
		ColumnExpr("authsrv_project.name as project").
		ColumnExpr("authsrv_organization.name as organization").
		Join(`JOIN authsrv_organization ON authsrv_project.organization_id=authsrv_organization.id`).
		Where("authsrv_project.id = ?", id).
		Where("authsrv_project.trash = ?", false).
		Where("authsrv_organization.trash = ?", false).
		Scan(ctx, &r)
	if err != nil {
		return "", "", err
	}
	return r.Project, r.Organization, nil
}

func GetFileteredProjects(ctx context.Context, db bun.IDB, account, partner, org uuid.UUID) ([]models.Project, error) {
	ids := []uuid.UUID{}
	sp := []models.AccountPermission{}
	err := db.NewSelect().Model(&sp).
		ColumnExpr("distinct account_id, project_id").
		Where("sap.partner_id = ?", partner).
		Where("sap.organization_id = ?", org).
		Where("sap.account_id = ?", account).
		Where("sap.permission_name IN (?)", bun.In([]string{"project.read", "ops_star.all"})).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	all := false
	for _, p := range sp {
		if p.ProjectId == uuid.Nil {
			all = true
			break
		}
		ids = append(ids, p.ProjectId)
	}

	prjs := []models.Project{}
	if !all && len(ids) == 0 {
		return prjs, nil
	}
	q := db.NewSelect().Model(&prjs).
		Where("project.partner_id = ?", partner).
		Where("project.organization_id = ?", org).
		Where("project.trash = ?", false)
	if !all {
		q = q.Where("project.id IN (?)", bun.In(ids))
	}
	err = q.Scan(ctx)
	return prjs, err
}

func GetProjectGroupRoles(ctx context.Context, db bun.IDB, id uuid.UUID) ([]*userv3.ProjectNamespaceRole, error) {
	var pr = []*userv3.ProjectNamespaceRole{}
	err := db.NewSelect().Table("authsrv_projectgrouprole").
		ColumnExpr("distinct authsrv_resourcerole.name as role, authsrv_project.name as project, authsrv_group.name as group").
		Join(`JOIN authsrv_resourcerole ON authsrv_resourcerole.id=authsrv_projectgrouprole.role_id`).
		Join(`JOIN authsrv_group ON authsrv_group.id=authsrv_projectgrouprole.group_id`).
		Join(`JOIN authsrv_project ON authsrv_project.id=authsrv_projectgrouprole.project_id`).
		Where("authsrv_projectgrouprole.project_id = ?", id).
		Scan(ctx, &pr)
	if err != nil {
		return nil, err
	}

	var pnr = []*userv3.ProjectNamespaceRole{}
	err = db.NewSelect().Table("authsrv_projectgroupnamespacerole").
		ColumnExpr("distinct authsrv_resourcerole.name as role, authsrv_project.name as project, authsrv_group.name as group, namespace_id as namespace").
		Join(`JOIN authsrv_resourcerole ON authsrv_resourcerole.id=authsrv_projectgroupnamespacerole.role_id`).
		Join(`JOIN authsrv_project ON authsrv_project.id=authsrv_projectgroupnamespacerole.project_id`).
		Join(`JOIN authsrv_group ON authsrv_group.id=authsrv_projectgroupnamespacerole.group_id`). // also need a namespace join
		Where("authsrv_projectgroupnamespacerole.project_id = ?", id).
		Scan(ctx, &pnr)
	if err != nil {
		return nil, err
	}

	return append(pr, pnr...), err
}

func GetProjectUserRoles(ctx context.Context, db bun.IDB, id uuid.UUID) ([]*userv3.UserRole, error) {

	var pr = []*userv3.UserRole{}
	err := db.NewSelect().Table("authsrv_projectaccountresourcerole").
		ColumnExpr("distinct authsrv_resourcerole.name as role, identities.traits ->> 'email' as user").
		Join(`JOIN authsrv_resourcerole ON authsrv_resourcerole.id=authsrv_projectaccountresourcerole.role_id`).
		Join(`JOIN identities ON identities.id=authsrv_projectaccountresourcerole.account_id`).
		Where("authsrv_projectaccountresourcerole.project_id = ?", id).
		Scan(ctx, &pr)
	if err != nil {
		return nil, err
	}

	return pr, err
}