package infraEnvs

const (
	InfiniteEzVersion    string = "0.1.0"
	InfiniteEzMainDir    string = "/var/infinite"
	MarketplaceDir       string = InfiniteEzMainDir + "/marketplace"
	InfiniteEzBinary     string = InfiniteEzMainDir + "/ez"
	AccessTokenCookieKey string = "control-access-token"
	UserDataDirectory    string = "/var/data"
	BackupCronFilePath   string = "/etc/cron.d/ez-backup"
	RestoreTaskTmpDir    string = InfiniteEzMainDir + "/tasks-restores-tmp"
)
