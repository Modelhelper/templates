package mssql

import (
	"context"
	"database/sql"
	"{{ .Name | kebab }}/service"

	"github.com/pkg/errors"

	_ "embed"

	_ "github.com/denisenkom/go-mssqldb"
)

//go:embed queries/get-hesp-data.sql
var sqlHepsConfigQuery string

func NewHespDataRepository(cm SqlConnectionManager) service.HespDataRepository {
	return &dataEventLoader{cm}
}

type dataEventLoader struct {
	cm SqlConnectionManager
}

func (d *dataEventLoader) GetFromWorkout(ctx context.Context, id int) (*service.HespConfigData, error) {
	conn, err := d.cm.Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "openConnection")
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, sqlHepsConfigQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(sqlHepsConfigQuery, sql.Named("id", id))

	var hespConfig service.HespConfigData

	if err := row.Scan(
		&hespConfig.Id,
		&hespConfig.StreamKey,
		&hespConfig.Player,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &hespConfig, nil
}
