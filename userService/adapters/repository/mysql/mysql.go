package mysql

import (
	"context"
	"fmt"
	"time"
	"user-svc/internal/entity"
	"user-svc/ports"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type MysqlDB struct {
	metrics ports.DatabaseMetrics
	db      *sqlx.DB
	logger  ports.Logger
}

func New(dbPool *sqlx.DB, logger ports.Logger, metrics ports.DatabaseMetrics) *MysqlDB {
	return &MysqlDB{db: dbPool, logger: logger, metrics: metrics}

}

func (r MysqlDB) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {

	query := "INSERT INTO users (phone_number, email, password) VALUES (?, ?, ?)"

	res, err := r.db.ExecContext(ctx, query, user.PhoneNumber, user.Email, user.Password)
	if err != nil {
		r.logger.Error("Failed to execute command", zap.String("query", query), zap.Error(err))
		return entity.User{}, fmt.Errorf("can not execute command %w", err)
	}

	id, iErr := res.LastInsertId()
	if iErr != nil {
		return entity.User{}, fmt.Errorf("failed to ... %w", iErr)
	}
	user.ID = uint(id)

	r.logger.Info("Data inserted successfully", zap.String("query", query))

	//! metrics
	start := time.Now()
	duration := time.Since(start).Seconds()
	dbDurationHistogram := r.metrics.RegisterDatabaseDurationHistogram().WithLabelValues(query)
	dbDurationHistogram.Observe(duration)

	if err := recover(); err != nil {
		dbErrorCounter := r.metrics.RegisterDatabaseErrorCounter().WithLabelValues(query)
		dbErrorCounter.Inc()
	}

	return user, nil
}

func (r MysqlDB) GetUserByID(ctx context.Context, userID uint) (entity.User, error) {
	var user entity.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = ?", userID)
	return user, err
}

func (r MysqlDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	panic("")
}

// ! helper function
// IsDBConnected checks if the database connection is successful
func (m *MysqlDB) IsDBConnected() bool {
	return m.db != nil
}

// defer func(query string) {
// 	duration := time.Since(start).Seconds()
// 	dbDurationHistogram := r.metrics.RegisterDatabaseDurationHistogram().WithLabelValues(query)
// 	dbDurationHistogram.Observe(duration)

// 	if err := recover(); err != nil {
// 		dbErrorCounter := r.metrics.RegisterDatabaseErrorCounter().WithLabelValues(query)
// 		dbErrorCounter.Inc()
// 	}
// }(query)
