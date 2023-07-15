package main

import (
	"github.com/xendit/xendit-go"

	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/nyelonong/boilerplate-go/core/config"
	"github.com/nyelonong/boilerplate-go/core/cron"
	"github.com/nyelonong/boilerplate-go/core/database"
	"github.com/nyelonong/boilerplate-go/core/environment"
	"github.com/nyelonong/boilerplate-go/core/log"
	"github.com/nyelonong/boilerplate-go/core/mailer"
	"github.com/nyelonong/boilerplate-go/core/meta"
	"github.com/nyelonong/boilerplate-go/core/monitor"
	"github.com/nyelonong/boilerplate-go/core/pinger"
	"github.com/nyelonong/boilerplate-go/core/qiscus"
	"github.com/nyelonong/boilerplate-go/core/redis"
	"github.com/nyelonong/boilerplate-go/core/storage"
	"github.com/nyelonong/boilerplate-go/core/validate"
	"github.com/nyelonong/boilerplate-go/internal/handlers/admin_dashboard"
	"github.com/nyelonong/boilerplate-go/internal/handlers/general"
	staff_dashboard "github.com/nyelonong/boilerplate-go/internal/handlers/staff-dashboard"
	"github.com/nyelonong/boilerplate-go/internal/models"
	"github.com/nyelonong/boilerplate-go/internal/repositories/export"
	"github.com/nyelonong/boilerplate-go/internal/repositories/notification"
	"github.com/nyelonong/boilerplate-go/internal/repositories/role"
	"github.com/nyelonong/boilerplate-go/internal/repositories/session"
	"github.com/nyelonong/boilerplate-go/internal/repositories/upload"
	"github.com/nyelonong/boilerplate-go/internal/repositories/user"
	"github.com/nyelonong/boilerplate-go/internal/usecases/admin_dash"
	"github.com/nyelonong/boilerplate-go/internal/usecases/authorization"
	general_uc "github.com/nyelonong/boilerplate-go/internal/usecases/general"
	staff_dash "github.com/nyelonong/boilerplate-go/internal/usecases/staff-dash"
)

type Configuration struct {
	Database  DatabaseConfig
	Redis     redis.Config
	Log       log.Config
	Http      HttpConfig
	Templates TemplateConfig
	Mailer    mailer.Config
	Whatsapp  meta.WhatsappConfig
	Qiscus    qiscus.QiscusConfig
	Storage   storage.Config
	Xendit    XenditConfig
}

type XenditConfig struct {
	SecretKey string
}
type RedisConfig struct {
	Address  string
	Password string
}
type DatabaseConfig struct {
	Driver     string
	Connection string
}
type HttpConfig struct {
	BaseDomain string
	Port       string
}
type TemplateConfig struct {
	Sms      string
	Whatsapp string
	Email    string
}

func main() {
	// Get Flag for service
	var configFile string
	configPath := "devops/configuration/backend-development.ini"
	if environment.IsStaging() {
		configPath = "devops/configuration/backend-staging.ini"
	}
	if environment.IsProduction() {
		configPath = "devops/configuration/backend-PRODUCTION.ini"
	}
	flag.StringVar(&configFile, "config", configPath, "Configuration File Location")
	flag.Parse()

	// Read and parse config
	conf := Configuration{}
	err := config.ReadFile(&conf, configFile)
	if err != nil {
		log.Fatalf("Failed to read configuration because %s", err)
	}

	// Initialize Base Domain
	environment.SetBaseDomain(conf.Http.BaseDomain)

	// Initialize Log
	log.Init(conf.Log)

	// Add Monitoring Metrics
	err = addMonitoring()
	if err != nil {
		log.Fatalf("Failed to add monitoring metrics because %s", err)
	}

	// Add common validation
	addCommonValidation()

	// Connect to Database Master
	masterDB, err := database.Connect(conf.Database.Driver, conf.Database.Connection)
	if err != nil {
		log.Fatalf("Failed to connect to database because %s", err)
	}
	pinger.AddService("Database Master", masterDB)

	// Connect to Redis
	rdb, err := redis.New(conf.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to redis because %s", err)
	}
	pinger.AddService("Redis", rdb)

	// Setup smtp mail
	mail := mailer.New(conf.Mailer)

	// Setup whatsapp
	whatsapp := meta.NewWhatsapp(conf.Whatsapp)

	// Setup qiscus
	qiscus := qiscus.NewQiscus(conf.Qiscus)

	// Setup S3 Uploader
	storage, err := storage.New(conf.Storage)
	if err != nil {
		log.Fatalf("Failed to connect to aws s3 because %s", err)
	}

	// Setup cron
	cron := cron.New()

	// Start Initialize Repository
	repoNotifConfig := notification.Configuration{SmsTemplateLoc: conf.Templates.Sms, WhatsappTemplateLoc: conf.Templates.Whatsapp, EmailTemplateLoc: conf.Templates.Email, Mailer: mail, Whatsapp: whatsapp, Qiscus: qiscus}
	// Setup smtp mail
	if err := addRepoNotifConfig(&repoNotifConfig); err != nil {
		log.Fatalf("Failed to add additional configuration for notification repository because %s", err)
	}

	// Xendit
	xendit.Opt.SecretKey = conf.Xendit.SecretKey

	repoUser := user.New(user.Configuration{Database: masterDB, Redis: rdb})
	repoRole := role.New(role.Configuration{Database: masterDB, Redis: rdb})
	repoUpload := upload.New(upload.Configuration{Database: masterDB, Redis: rdb, Storage: storage})
	repoSession := session.New(session.Configuration{Database: masterDB, Redis: rdb, ShortExpDuration: 24 * time.Hour, LongExpDuration: 7 * 24 * time.Hour})
	repoNotification := notification.New(repoNotifConfig)
	repoExport := export.New(export.Configuration{Database: masterDB, Redis: rdb, Storage: storage})
	ucGeneral := general_uc.New(general_uc.Configuration{Upload: repoUpload, Notification: repoNotification})
	ucAuthorization := authorization.New(authorization.Configuration{Session: repoSession, User: repoUser, Role: repoRole, Notification: repoNotification})
	ucAdminDashboard := admin_dash.New(admin_dash.Configuration{User: repoUser, Notification: repoNotification})

	ucStaffDashboard := staff_dash.New(staff_dash.Configuration{
		User:         repoUser,
		Role:         repoRole,
		Upload:       repoUpload,
		Notification: repoNotification,
		Export:       repoExport,
	})

	adminDash := admin_dashboard.New(admin_dashboard.Configuration{DashboardUC: ucAdminDashboard, AuthorizationUC: ucAuthorization})
	staffDash := staff_dashboard.New(staff_dashboard.Configuration{DashboardUC: ucStaffDashboard, AuthorizationUC: ucAuthorization})
	genView := general.New(general.Configuration{AuthorizationUC: ucAuthorization, GeneralUC: ucGeneral})

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Get("/ping", ping)
	r.Get("/health", health)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Handle("/metrics", promhttp.Handler())

	// webhook
	webhookRouter := r.Group(nil)
	webhookRouter.Route("/webhook/", func(r chi.Router) {
		r.Post("/webhook/xenplatform/invoices", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	authRouter := r.Group(nil)
	authRouter.Route("/auth/", func(r chi.Router) {
		r.Post("/auth/login", staffDash.DoLoginFromApp)
		r.Post("/auth/refresh-token", staffDash.MustNotLogin(staffDash.DoRefreshToken))

		r.Post("/email-verification/view", genView.ShowEmailVerification)
		r.Post("/forgot-password/request", genView.SendForgotPassword)
		r.Post("/new-password/submit", genView.NewPassword)
	})

	adminDashboard := r.Group(nil)
	adminDashboard.Post("/admin/login", adminDash.MustNotLogin(adminDash.DoLogin))
	adminDashboard.Post("/admin/logout", adminDash.MustLogin(adminDash.DoLogout))
	adminDashboard.Route("/admin/", func(r chi.Router) {
		r.Use(adminDash.UseMustLogin)
		r.Use(adminDash.UseAddContextValue)

		r.Post("/check", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	// Jobs, all jobs need to be listed before start
	cron.AddJob("@daily", func(ctx context.Context) error {
		log.Println("Hello")
		return nil
	})
	cron.Start(context.Background())

	if environment.IsDevelopment() {
		// Cron Test
	}

	srv := &http.Server{
		Addr:    conf.Http.Port,
		Handler: r,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Infof("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Infof("Start HTTP Server in %s", conf.Http.Port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write(pinger.CheckHealth())
}

func addCommonValidation() {
	validate.AddStringCommon(models.ValidRequiredName, validate.WithRequired(), validate.WithMinLength(3),
		validate.WithMaxLength(100), validate.WithAlphaNumeric())
	validate.AddStringCommon(models.ValidRequiredEmail, validate.WithRequired(), validate.WithMinLength(5),
		validate.WithMaxLength(100), validate.WithEmailFormat())
	validate.AddStringCommon(models.ValidRequiredPassword, validate.WithRequired(), validate.WithMinLength(6),
		validate.WithMaxLength(50), validate.WithAlphaNumeric())
	validate.AddStringCommon(models.ValidRequiredPhone, validate.WithRequired(), validate.WithMinLength(8),
		validate.WithMaxLength(20), validate.WithPhoneFormat())
	validate.AddStringCommon(models.ValidRequiredOTP, validate.WithRequired(), validate.WithMinLength(6),
		validate.WithMaxLength(6), validate.WithOnlyNumeric())
	validate.AddStringCommon(models.ValidRequiredAddress, validate.WithRequired(), validate.WithMinLength(5),
		validate.WithMaxLength(255), validate.WithAlphaNumeric())
	validate.AddStringCommon(models.ValidRequiredZipCode, validate.WithRequired(), validate.WithMinLength(5),
		validate.WithMaxLength(5), validate.WithOnlyNumeric())
	validate.AddStringCommon(models.ValidRequiredAnyDate, validate.WithRequired(), validate.WithTimeFormat("YYYY-MM-DD"))

	// Non Spesific Validation
	validate.AddStringCommon(models.ValidRequiredShortStr, validate.WithRequired(), validate.WithMinLength(3),
		validate.WithMaxLength(50), validate.WithAlphaNumeric())
	validate.AddStringCommon(models.ValidRequiredMedStr, validate.WithRequired(), validate.WithMinLength(3),
		validate.WithMaxLength(100), validate.WithAlphaNumeric())
	validate.AddStringCommon(models.ValidRequiredLongStr, validate.WithRequired(), validate.WithMinLength(3),
		validate.WithMaxLength(255), validate.WithAlphaNumeric())

	validate.AddStringCommon(models.ValidRequiredText, validate.WithRequired(), validate.WithMinLength(3),
		validate.WithMaxLength(1024), validate.WithAlphaNumeric())
}

func addMonitoring() error {
	err := monitor.Init(monitor.Config{
		Prefix:       "backend",
		Engine:       monitor.EnginePrometheus,
		DefaultLabel: map[string]string{"env": "development"},
	})
	if err != nil {
		return err
	}

	// Measure API Call Latency and Count
	monitor.NewHistogram("api_call", []string{"api", "status", "usecase"})

	return nil
}

func addRepoNotifConfig(config *notification.Configuration) (err error) {
	if err = config.AddEmailSender(models.EmailSenderNoReply, "No Reply", "no-reply@domain.com"); err != nil {
		return err
	}

	if err = config.AddEmailSender(models.EmailSenderSupport, "Support", "support@domain.com"); err != nil {
		return err
	}
	return nil
}
