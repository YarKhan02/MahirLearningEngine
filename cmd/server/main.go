package main

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	apihttp "github.com/YarKhan02/MahirLearningEngine/internal/api/http"
	"github.com/YarKhan02/MahirLearningEngine/internal/config"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/assignment"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/attendance"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/batch"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/course"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/role"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/student"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/user"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/crypto"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/postgres/migrations"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/postgres/repository"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
)

func main () {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		cancel()
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		cancel()
		log.Fatalf("database ping failed: %v", err)
	}
	cancel()

	if err := migrations.RunMigration(cfg.MigrationsPath, cfg.DatabaseURL); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	redis, err := redis.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer redis.Close()

	var key *rsa.PrivateKey
	var keyErr error
	var keySource string

	if cfg.RSAPrivateKeyPEM != "" {
		keySource = "RSA_PRIVATE_KEY_PEM"
		key, keyErr = crypto.LoadRSAPrivateKeyFromPEM(cfg.RSAPrivateKeyPEM)
	}

	if keyErr != nil {
		log.Fatalf("failed to load RSA key from %s: %v", keySource, keyErr)
	}

	userRepo	:= repository.NewUserRepository(db)
	courseRepo 	:= repository.NewCourseRepository(db)
	batchRepo 	:= repository.NewBatchRepository(db)
	roleRepo 	:= repository.NewRoleRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	assignmentRepo := repository.NewAssignmentRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)
	tokenRepo 	:= repository.NewTokenRepository(db)

	userSvc 	:= user.NewService(userRepo, roleRepo)
	courseSvc 	:= course.NewService(courseRepo)
	batchSvc 	:= batch.NewService(batchRepo)
	roleSvc 	:= role.NewService(roleRepo)
	studentSvc 	:= student.NewService(studentRepo)
	assignmentSvc := assignment.NewService(assignmentRepo)
	attendanceSvc := attendance.NewService(attendanceRepo)
	tokenSvc 	:= token.NewService(key, tokenRepo, cfg.JWTIssuer, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)

	srv := apihttp.NewServer(cfg, userSvc, roleSvc, courseSvc, batchSvc, studentSvc, assignmentSvc, attendanceSvc, tokenSvc, redis)

	log.Printf("listening on: %s", cfg.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}