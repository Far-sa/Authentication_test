package mysql

import (
	"context"
	"fmt"
	"user-svc/internal/entity"
	"user-svc/ports"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type MysqlDB struct {
	//metrics ports.DatabaseMetrics
	config ports.Config
	db     *sqlx.DB
	logger ports.Logger
}

func New(config ports.Config, logger ports.Logger) *MysqlDB {
	dbConfig := config.GetDatabaseConfig()
	//db, err := sqlx.Connect("mysql", "root:password@(localhost:3306)/mysql_app")

	fmt.Println("DB host is:", dbConfig.Host)

	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@%s:%d/%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName))
	if err != nil {
		logger.Error("Failed to open MySQL database", zap.Error(err))
		fmt.Errorf("can not open mysql :%v", err)
	}

	// db.SetConnMaxLifetime(time.Minute * 3)
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(10)

	return &MysqlDB{config: config, db: db, logger: logger}

}

func (r MysqlDB) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {

	query := "INSERT INTO users (phone_number, email, password) VALUES (?, ?, ?)"

	res, err := r.db.ExecContext(ctx, query, user.PhoneNumber, user.Email, user.Password)
	if err != nil {
		r.logger.Error("Failed to execute command", zap.String("query", query), zap.Error(err))
		return entity.User{}, fmt.Errorf("can not execute command %w", err)
	}

	id, _ := res.LastInsertId()
	user.ID = uint(id)

	r.logger.Info("Data inserted successfully", zap.String("query", query))

	//! metrics
	// start := time.Now()
	// duration := time.Since(start).Seconds()
	// dbDurationHistogram := r.metrics.RegisterDatabaseDurationHistogram().WithLabelValues(query)
	// dbDurationHistogram.Observe(duration)

	// if err := recover(); err != nil {
	// 	dbErrorCounter := r.metrics.RegisterDatabaseErrorCounter().WithLabelValues(query)
	// 	dbErrorCounter.Inc()
	// }

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
