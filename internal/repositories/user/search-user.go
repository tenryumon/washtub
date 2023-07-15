package user

import (
	"context"
	"strings"

	"github.com/nyelonong/boilerplate-go/core/helper"
	"github.com/nyelonong/boilerplate-go/internal/entities"
)

var queryCount = "SELECT COUNT(*) as count "
var querySearchUser = querySelectUser + " FROM users"
var queryCountSearchUser = queryCount + " FROM users"
var queryUserWhereNameLike = "lower(name) LIKE :query"
var queryUserWhereUserIDIn = "id IN (:user_list)"
var queryUserWhereStatus = "status = :status"

func (us *UserObject) getSearchUserWhere(param map[string]interface{}) string {
	where := []string{}
	if value, ok := param["query"]; ok && value != "" {
		where = append(where, queryUserWhereNameLike)
	}
	if _, ok := param["user_list"]; ok {
		where = append(where, queryUserWhereUserIDIn)
	}
	if _, ok := param["status"]; ok {
		where = append(where, queryUserWhereStatus)
	}

	if len(where) > 0 {
		return " WHERE " + strings.Join(where, " AND ")
	}

	return ""
}

func (us *UserObject) SearchUser(ctx context.Context, param map[string]interface{}) ([]entities.User, error) {
	query := querySearchUser
	query += us.getSearchUserWhere(param)
	query += helper.GetSqlCommonClauses(param)

	result := []entities.User{}
	err := us.db.Select(ctx, &result, query, param)
	return result, err
}

func (us *UserObject) countSearchUser(ctx context.Context, param map[string]interface{}) (int64, error) {
	query := queryCountSearchUser
	query += us.getSearchUserWhere(param)

	result := entities.Aggregate{}
	err := us.db.Get(ctx, &result, query, param)
	return result.Count, err
}

func (us *UserObject) PaginationUser(ctx context.Context, param map[string]interface{}) (result []entities.User, pagination entities.Pagination, err error) {
	result = []entities.User{}
	pagination = entities.Pagination{}

	pagination.Total, err = us.countSearchUser(ctx, param)
	if err != nil {
		return
	}

	page, limit := helper.NormalizePageLimit(pagination.Total, &param)
	pagination.Page = page
	pagination.Limit = limit

	result, err = us.SearchUser(ctx, param)
	return
}
