package db

import (
	"fmt"
	"sync"

	_ "github.com/lib/pq"
	"github.com/mithileshgupta12/velaris/internal/config"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

var (
	repositories *repository.Repository
	engine       *xorm.Engine
	once         sync.Once
	instanceErr  error
)

func NewDB(dbFlags *config.DBFlags) (*repository.Repository, error) {
	once.Do(func() {
		connStr := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			dbFlags.Host, dbFlags.PORT, dbFlags.User, dbFlags.Password, dbFlags.Name, dbFlags.SSLMode,
		)

		engine, instanceErr = xorm.NewEngine("postgres", connStr)
		if instanceErr != nil {
			return
		}

		engine.SetColumnMapper(names.GonicMapper{})

		instanceErr = engine.Ping()
		if instanceErr != nil {
			return
		}

		repositories = repository.NewRepository(engine)
	})

	return repositories, instanceErr
}
