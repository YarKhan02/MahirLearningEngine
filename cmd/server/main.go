package main

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	apihttp "github.com/YarKhan02/MahirLearningEngine/internal/api/http"
	"github.com/YarKhan02/MahirLearningEngine/internal/config"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/assignment"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/attendance"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/announcement"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/batch"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/course"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/dashboard"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/role"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/student"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/timetable"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/user"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/crypto"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/postgres/migrations"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/postgres/repository"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
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

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = db.PingContext(ctx)
	cancel()
	if err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

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

	userRepo := repository.NewUserRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	batchRepo := repository.NewBatchRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	assignmentRepo := repository.NewAssignmentRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)
	timetableRepo := repository.NewTimetableRepository(db)
	announcementRepo := repository.NewAnnouncementRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	userSvc := user.NewService(userRepo, roleRepo)
	courseSvc := course.NewService(courseRepo)
	batchSvc := batch.NewService(batchRepo)
	roleSvc := role.NewService(roleRepo)
	studentSvc := student.NewService(studentRepo)
	assignmentSvc := assignment.NewService(assignmentRepo)
	attendanceSvc := attendance.NewService(attendanceRepo)
	dashboardSvc := dashboard.NewService(dashboardRepo)
	timetableSvc := timetable.NewService(timetableRepo)
	announcementSvc := announcement.NewService(announcementRepo)
	tokenSvc := token.NewService(key, tokenRepo, cfg.JWTIssuer, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)

	srv := apihttp.NewServer(cfg, userSvc, roleSvc, courseSvc, batchSvc, studentSvc, assignmentSvc, attendanceSvc, dashboardSvc, timetableSvc, announcementSvc, tokenSvc, redisClient)

	log.Printf("listening on: %s", cfg.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
