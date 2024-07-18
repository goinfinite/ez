package presentation

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
	"github.com/speedianet/control/src/presentation/api"
	"github.com/speedianet/control/src/presentation/ui"
)

func HttpServerInit(
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) {
	e := echo.New()

	api.ApiInit(e, persistentDbSvc, transientDbSvc, trailDbSvc)
	ui.UiInit(e, persistentDbSvc, transientDbSvc)

	httpServer := http.Server{Addr: ":3141", Handler: e}

	pkiDir := "/var/speedia/pki"
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
		log.Fatal(err)
	}
}
