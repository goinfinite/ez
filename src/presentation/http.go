package presentation

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/goinfinite/ez/src/infra/db"
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

	pkiDir := "/var/infinite/pki"
	certFile := pkiDir + "/control.crt"
	keyFile := pkiDir + "/control.key"

	controlBanner := `
                                             ▒▓▓▓▒                        ▓▓▓▓▒
                                             ████▒                        ████▓
   ▓▒   ▒████▓▓▓█▒▓████▓▓████▒▓████▓▓█████ ▓█████▓▓ █████▓▓▓▓▒████▓▓█████ ████ 
 ▒█▓    ████▓    ▒████  ▒████ ▓████  ▓████  ████▒   ████▓    ▓████  ████▓▓████ 
▓▓█▒▒▒ ▒████     ▓███▓  ████▓ ████▒  ████▒ ▒████   ▒████     ████▒ ▒████ ████▓ 
  ▓█▓  ████▓     ████▒ ▒████ ▒████  ▓████  ████▓   ████▓    ▓████  ▓███▓▒████  
 ▒█    ████▒    ▓████  ████▓ ████▓  ████▒ ▒████▒   ████     ████▓ ▒████ ▓████  
 ▒     ▒▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▒  ▓▓▓▓   ▓▓▓▓   ▓▓▓▓▓▓ ▒▓▓▓▓     ▒▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▒ 
______________________________________________________________________________

⇨ HTTPS server started on [::]:3141 and is ready to serve! 🎉
`

	fmt.Println(controlBanner)

	err := httpServer.ListenAndServeTLS(certFile, keyFile)
	if err != http.ErrServerClosed {
		slog.Error("HttpServerError", slog.Any("error", err))
		os.Exit(1)
	}
}
