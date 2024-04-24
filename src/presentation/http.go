package presentation

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
	"github.com/speedianet/control/src/presentation/api"
	"github.com/speedianet/control/src/presentation/ui"
)

type CustomLogger struct{}

func (*CustomLogger) Write(rawMessage []byte) (int, error) {
	messageStr := strings.TrimSpace(string(rawMessage))

	shouldLog := true
	if strings.HasSuffix(messageStr, "tls: unknown certificate") {
		shouldLog = false
	}

	messageLen := len(rawMessage)
	if !shouldLog {
		return messageLen, nil
	}

	return messageLen, log.Output(2, messageStr)
}

func NewCustomLogger() *log.Logger {
	return log.New(&CustomLogger{}, "", 0)
}

func HttpServerInit(
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
) {
	e := echo.New()

	api.ApiInit(e, persistentDbSvc, transientDbSvc)
	ui.UiInit(e, persistentDbSvc, transientDbSvc)

	httpServer := http.Server{
		Addr:     ":3141",
		Handler:  e,
		ErrorLog: NewCustomLogger(),
	}

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
