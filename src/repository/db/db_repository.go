package db

import (
	"github.com/avinashb98/bookstore-oauth-api/src/clients/cassandra"
	"github.com/avinashb98/bookstore-oauth-api/src/domain/access_token"
	"github.com/avinashb98/bookstore-oauth-api/src/utils/errors"
	"github.com/gocql/gocql"
)

const (
	queryGetAccessToken    = "SELECT access_token, user_id, client_id, expires from access_tokens where access_token=?;"
	queryCreateAccessToken = "INSERT into access_tokens(access_token, user_id, client_id, expires) VALUES (?, ?, ?, ?);"
	queryUpdateExpires     = "UPDATE access_tokens SET expires=? where user_id=?;"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, errors.RestErr)
	Create(access_token.AccessToken) errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) errors.RestErr
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

type dbRepository struct{}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, errors.RestErr) {
	session := cassandra.GetSession()
	var result access_token.AccessToken
	if err := session.Query(queryGetAccessToken, id).Scan(
		&result.AccessToken,
		&result.UserId,
		&result.ClientId,
		&result.Expires,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.NewNotFoundError("no access token found for the id")
		}
		return nil, errors.NewInternalServerError(err.Error(), nil)
	}
	return &result, nil
}

func (r *dbRepository) Create(token access_token.AccessToken) errors.RestErr {
	session := cassandra.GetSession()

	if err := session.Query(queryCreateAccessToken,
		token.AccessToken,
		token.UserId,
		token.ClientId,
		token.Expires,
	).Exec(); err != nil {
		return errors.NewInternalServerError(err.Error(), err)
	}
	return nil
}

func (r *dbRepository) UpdateExpirationTime(token access_token.AccessToken) errors.RestErr {
	session := cassandra.GetSession()

	if err := session.Query(queryUpdateExpires,
		token.Expires,
		token.UserId,
	).Exec(); err != nil {
		return errors.NewInternalServerError(err.Error(), err)
	}
	return nil
}
