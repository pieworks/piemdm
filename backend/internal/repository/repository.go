// Package repository provides the data access layer for the PieMDM system.
// It contains repository interfaces and implementations that handle database
// operations, data queries, and persistence logic using GORM.
package repository

import (
	"context"
	"fmt"
	"time"

	"piemdm/pkg/log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Repository struct {
	db     *gorm.DB
	rdb    *redis.Client
	logger *log.Logger
}

func NewRepository(db *gorm.DB, rdb *redis.Client, logger *log.Logger) *Repository {
	return &Repository{
		db:     db,
		rdb:    rdb,
		logger: logger,
	}
}

// DB
func NewDB(conf *viper.Viper) *gorm.DB {

	dsn := conf.GetString("data.mysql.dsn")

	// 1. 基础配置
	gormConfig := &gorm.Config{
		// 【性能】跳过默认事务，除非必须，否则手动开启事务
		SkipDefaultTransaction: true,

		// 【性能】缓存预编译语句，提高重复 SQL 执行速度
		PrepareStmt: true,

		// 【规范】命名策略
		NamingStrategy: schema.NamingStrategy{
			// 在现代微服务和云原生架构下，不建议给表名加 tbl_ 前缀。
			// 视图 (View)：建议加 v_ 或 view_
			// 索引 (Index)：必须有规范的前缀/后缀。
			// 普通索引：idx_字段名 (如 idx_email)
			// 唯一索引：uniq_字段名 (如 uniq_username)
			// GORM 在创建索引时会自动处理命名，或者你可以通过 tag gorm:"index:idx_name" 手动指定。
			TablePrefix:   "",    // conf.GetString("data.mysql.table_prefix"),
			SingularTable: false, // true, // 建议：使用单数表名 (user 而不是 users)，代码映射更直观
		},

		// 【架构】禁用物理外键，依靠代码逻辑维护一致性
		DisableForeignKeyConstraintWhenMigrating: true,

		// 【日志】生产环境建议自定义 Logger，不要直接打印所有 SQL
		// 建议集成日志记录，并只在 Slow SQL 时打印
		Logger: logger.Default.LogMode(logger.Warn),
	}

	// 2. 打开数据库连接
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 256, // 如果数据库较新，可以直接用 256；兼容老库用 191
		// 参考https://cloud.tencent.com/developer/article/1917039
	}), gormConfig)

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// 3. 【关键】配置连接池
	// GORM 获取底层的 sql.DB 对象
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get sql.DB: " + err.Error())
	}

	// db, err := gorm.Open(
	// 	mysql.New(mysql.Config{
	// 		DSN:               conf.GetString("data.mysql.dsn"),
	// 		DefaultStringSize: 191, // utf8mb4的字符串长度应该设置为191字节,
	// 		// 参考https://cloud.tencent.com/developer/article/1917039
	// 	}), &gorm.Config{
	// 		// SkipDefaultTransaction: true,
	// 		// Logger: logger.Default.LogMode(logger.Info), // 打印SQL语句
	// 		NamingStrategy: schema.NamingStrategy{
	// 			TablePrefix: conf.GetString("data.mysql.table_prefix"), // 表名前缀
	// 			// SingularTable: true,                                      // 使用单数表名
	// 			// NoLowerCase:   true,
	// 		},
	// 		DisableForeignKeyConstraintWhenMigrating: true, // 不自动创建外键约束
	// 	})
	// if err != nil {
	// 	panic(err)
	// }
	// // 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	// sqlDB, err := db.DB()
	// if err != nil {
	// 	panic("connect db server failed.")
	// }
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	// 根据并发量设置，一般设为 MaxOpenConns 的 30%-50%
	sqlDB.SetMaxIdleConns(conf.GetInt("data.mysql.max_conn"))
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	// 根据数据库规格和微服务实例数计算，防止打爆数据库
	sqlDB.SetMaxOpenConns(conf.GetInt("data.mysql.max_open"))
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	// 必须小于 MySQL 服务端 wait_timeout (默认 8小时)，通常设置为 1小时或更短
	// 如果不设置，由于网络波动或防火墙切断，连接可能会在服务端断开，导致 "invalid connection" 错误
	sqlDB.SetConnMaxLifetime(time.Hour)
	// data, _ := json.Marshal(sqlDB.Stats()) //获得当前的SQL配置情况
	// r.logger.Info("sqlDB", "sqlDB", string(data))
	return db
}

// Redis
func NewRedis(conf *viper.Viper) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.GetString("data.redis.addr"),
		Password: conf.GetString("data.redis.password"),
		DB:       conf.GetInt("data.redis.db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	}

	return rdb
}
