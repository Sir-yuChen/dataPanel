package code

import (
	"context"
	"dataPanel/serviceend/global"
	"dataPanel/serviceend/model"
	"dataPanel/serviceend/model/stock"
	"errors"
	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	_ "path/filepath"
	"runtime"
	"strings"
	"time"
)

var sqliteDb *gorm.DB

func InitDB(sqlitePath string) {
	dbLogger := NewGormLogger(logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		Colorful:                  true,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      false,
		LogLevel:                  logger.Info,
	})

	if sqlitePath == "" {
		sqlitePath = global.GvaConfig.System.DbPath
		if !strings.Contains(sqlitePath, "?") {
			sqlitePath += "?cache=shared&mode=rwc&_journal_mode=WAL&_busy_timeout=5000"
		}
	}

	var err error
	sqliteDb, err = gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{
		Logger:                                   dbLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: global.GvaConfig.System.DbTablePrefix,
		},
	})

	if err != nil {
		global.GvaLog.Fatal("数据库连接失败", zap.Any("path", sqlitePath), zap.Error(err))
	}

	dbCon, err := sqliteDb.DB()
	if err != nil {
		global.GvaLog.Fatal("数据库连接异常,sqliteDb.DB error", zap.Error(err))
	}

	// 适配 SQLite 并发特性
	dbCon.SetMaxIdleConns(1)
	dbCon.SetMaxOpenConns(1)
	dbCon.SetConnMaxLifetime(time.Hour * 2)

	if err := IntCreateTable(sqliteDb); err != nil {
		global.GvaLog.Fatal("自动迁移失败", zap.Error(err))
	}
	global.GvaSqliteDb = sqliteDb
}

func GetDB() *gorm.DB {
	return sqliteDb
}

func IntCreateTable(db *gorm.DB) error {
	//自动创建表结构
	models := []interface{}{
		&model.AppSetting{},
		&stock.StockBasic{},
	}

	for _, m := range models {
		if err := db.AutoMigrate(m); err != nil {
			return err
		}
	}
	return nil
}

type GormLogger struct {
	SlowThreshold time.Duration
	Log           *zap.Logger
	config        logger.Config
}

func NewGormLogger(config logger.Config) *GormLogger {
	return &GormLogger{
		SlowThreshold: 200 * time.Millisecond,
		Log:           global.GvaLog,
		config:        config,
	}
}

var _ logger.Interface = (*GormLogger)(nil)

func (l *GormLogger) LogMode(lev logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.config.LogLevel = lev
	return &newLogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger().Sugar().Infof(msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger().Sugar().Warnf(msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger().Sugar().Errorf(msg, data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	logFields := make([]zap.Field, 0, 4)
	logFields = append(logFields,
		zap.String("SQL", sql),
		zap.Duration("TIME", elapsed),
		zap.Int64("ROWS", rows),
	)

	// 慢查询判断
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		logFields = append(logFields, zap.String("threshold", l.SlowThreshold.String()))
		l.Log.Warn("[SLOW SQL]", logFields...)
	}

	// 错误处理
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.Log.Warn("Database ErrRecordNotFound", logFields...)
		} else {
			logFields = append(logFields, zap.Error(err))
			l.Log.Error("[SQL]", logFields...)
		}
		return
	}

	// 常规日志
	if l.config.LogLevel <= logger.Info {
		l.Log.Info("[SQL]", logFields...)
	}
}

func (l *GormLogger) logger() *zap.Logger {
	const maxCallerDepth = 8
	var pcs [maxCallerDepth]uintptr

	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "gorm.io/gorm") {
			return l.Log.WithOptions(zap.AddCallerSkip(frame.Line))
		}
		if !more {
			break
		}
	}
	return l.Log
}
