package presentation

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/goinfinite/ez/src/infra/db"
	o11yInfra "github.com/goinfinite/ez/src/infra/o11y"
	"github.com/goinfinite/ez/src/presentation/api"
	"github.com/goinfinite/ez/src/presentation/ui"
	"github.com/labstack/echo/v4"
)

func HttpServerInit(
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) {
	e := echo.New()

	api.ApiInit(e, persistentDbSvc, transientDbSvc, trailDbSvc)
	ui.UiInit(e, persistentDbSvc, transientDbSvc, trailDbSvc)

	httpServer := http.Server{Addr: ":3141", Handler: e}

	ezBanner := `Infinite Ez server started on [::]:3141! 🎉`

	o11yQueryRepo := o11yInfra.NewO11yQueryRepo(transientDbSvc)
	o11yOverview, err := o11yQueryRepo.ReadOverview()
	if err == nil {
		ezBanner = `
      INFINITE      |  🔒 HTTPS server started on [::]:3141!
   ▄▄█▀██ █▀▀▀███   |  
  ▄█▀   ██▀  ███    |  🏠 Hostname: ` + o11yOverview.Hostname.String() + `
  ██▀▀▀▀▀▀  ███     |  ⏰ Uptime: ` + o11yOverview.UptimeRelative.String() + `
  ██▄    ▄ ███  ▄   |  🌐 IPs: ` + o11yOverview.PrivateIpAddress.String() + ` ‖ ` + o11yOverview.PublicIpAddress.String() + `
   ▀█████▀███████   |  ⚙️  ` + o11yOverview.HardwareSpecs.String() + `
`
	}
	fmt.Println(ezBanner)

	pkiDir := "/var/infinite/pki"
	certFile := pkiDir + "/ez.crt"
	keyFile := pkiDir + "/ez.key"

	err = httpServer.ListenAndServeTLS(certFile, keyFile)
	if err != http.ErrServerClosed {
		slog.Error("HttpServerError", slog.Any("error", err))
		os.Exit(1)
	}
}
