package main

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api"
	"github.com/YarKhan02/MahirLearningEngine/internal/config"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/announcement"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/assignment"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/attendance"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/batch"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/course"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/dashboard"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/student"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/timetable"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/user"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/crypto"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/logging"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/metrics"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/postgres/migrations"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/postgres/repository"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"

	"go.uber.org/zap"
)

// main only reports the error — all setup lives in run so deferred
// cleanups actually execute on failure (log.Fatalf skips defers).
func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger, err := logging.New(logging.Config{
		Env: cfg.Env,
		ServiceName: "mahirlearning",
		Version: "1",
		LogFilePath: "/Users/yarkhan/Tech/MahirLearning/MahirLearningEngine/logs/app.log",
	})
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffered log entries on shutdown

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = db.PingContext(ctx)
	cancel()
	if err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// DB pool collector + build info on the default Prometheus registry.
	metrics.Register(db, "mahirlearning", "1", cfg.Env)

	if err := migrations.RunMigration(cfg.MigrationsPath, cfg.DatabaseURL); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	redisClient, err := redis.NewRedisClient(cfg.RedisURL)
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	defer redisClient.Close()

	var key *rsa.PrivateKey
	if cfg.RSAPrivateKeyPEM != "" {
		key, err = crypto.LoadRSAPrivateKeyFromPEM(cfg.RSAPrivateKeyPEM)
		if err != nil {
			return fmt.Errorf("failed to load RSA key from RSA_PRIVATE_KEY_PEM: %w", err)
		}
	}

	// Secure cookies (SameSite=None) only work over HTTPS; in local HTTP dev we
	// must fall back to a plain SameSite=Lax cookie or the browser drops it.
	secureCookies := !strings.EqualFold(cfg.Env, "development")

	roleRepo := repository.NewRoleRepository(db)

	announcementRepo := repository.NewAnnouncementRepository(db)
	announcementCache := announcement.NewCachedRepository(announcementRepo, redisClient)
	announcementSvc := announcement.NewService(announcementCache)

	courseRepo := repository.NewCourseRepository(db)
	courseCache := course.NewCachedRepository(courseRepo, redisClient)
	courseSvc := course.NewService(courseCache)

	batchRepo := repository.NewBatchRepository(db)
	batchCache := batch.NewCachedRepository(batchRepo, redisClient)
	batchSvc := batch.NewService(batchCache)
	
	studentRepo := repository.NewStudentRepository(db)
	studentCache := student.NewCachedRepository(studentRepo, redisClient)
	studentSvc := student.NewService(studentCache)
	
	dashboardRepo := repository.NewDashboardRepository(db)
	dashboardCache := dashboard.NewCachedRepository(dashboardRepo, redisClient)
	dashboardSvc := dashboard.NewService(dashboardCache)

	timetableRepo := repository.NewTimetableRepository(db)
	timetableCache := timetable.NewCachedRepository(timetableRepo, redisClient)
	timetableSvc := timetable.NewService(timetableCache)

	userRepo := repository.NewUserRepository(db)
	userSvc := user.NewService(userRepo, roleRepo)

	assignmentRepo := repository.NewAssignmentRepository(db)
	assignmentSvc := assignment.NewService(assignmentRepo)
	
	attendanceRepo := repository.NewAttendanceRepository(db)
	attendanceSvc := attendance.NewService(attendanceRepo)
	
	tokenRepo := repository.NewTokenRepository(db)
	tokenSvc := token.NewService(key, tokenRepo, cfg.JWTIssuer, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)

	module := []api.Module{
		user.NewModule(userSvc, studentSvc, tokenSvc, redisClient, secureCookies),
		course.NewModule(courseSvc, tokenSvc, redisClient),
		batch.NewModule(batchSvc, tokenSvc, redisClient),
		student.NewModule(studentSvc, userSvc, tokenSvc, redisClient, cfg.TempPassword),
		assignment.NewModule(assignmentSvc, tokenSvc, redisClient),
		attendance.NewModule(attendanceSvc, tokenSvc, redisClient),
		dashboard.NewModule(dashboardSvc, tokenSvc, redisClient),
		timetable.NewModule(timetableSvc, tokenSvc, redisClient),
		announcement.NewModule(announcementSvc, tokenSvc, redisClient),
	}

	srv := api.NewServer(cfg.AllowedOrigin, cfg.Addr, module, logger, cfg.RateLimitRequests, cfg.RateLimitWindow, cfg.PrometheusUsername, cfg.PrometheusPassword)

	logger.Info("server starting", zap.String("event", "server_start"), zap.String("addr", cfg.Addr))
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server stopped unexpectedly", zap.String("event", "server_error"), zap.Error(err))
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
