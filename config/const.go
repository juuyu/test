package config

const (
	DockerRegistry = "60.188.46.2:20012"

	minioEndpoint = "http://60.188.46.2:20002/"
	minioBucket   = "conf"

	dockerComposePath         = "/docker/docker-compose"
	serverComposeTemplatePath = "/docker/compose-file/server-docker-compose.yml"
	appComposeTemplatePath    = "/docker/compose-file/app-docker-compose.yml"
	redisConfPath             = "/server/redis/redis.conf"
	backConfPath              = "/app/backend/application.yml"
	frontConfPath             = "/app/front/default.conf"
	frontPemPath              = "/app/cert/%s.i-metalab.com.pem"
	frontKeyPath              = "/app/cert/%s.i-metalab.com.key"
	ffmpegPath                = "/ffmpeg/ffmpeg"
	ffmpegConfPath            = "/app/ffmpeg-go/config.toml"
	sqlPath                   = "/server/mysql/metalabs-teaching.sql"
	srsConfPath               = "/server/srs/docker.conf"
)

type url struct {
	DockerCompose         string
	ServerComposeTemplate string
	AppComposeTemplate    string
	RedisConf             string
	BackConfPath          string
	FrontConfPath         string
	FrontPemPath          string
	FrontKeyPath          string
	FfmpegPath            string
	FfmpegConfPath        string
	SQLPath               string
	SrsConfPath           string
}

var Url *url

func init() {
	Url = new(url)
	Url.DockerCompose = getBasePath() + dockerComposePath
	Url.ServerComposeTemplate = getBasePath() + serverComposeTemplatePath
	Url.AppComposeTemplate = getBasePath() + appComposeTemplatePath
	Url.RedisConf = getBasePath() + redisConfPath
	Url.BackConfPath = getBasePath() + backConfPath
	Url.FrontConfPath = getBasePath() + frontConfPath
	Url.FrontPemPath = getBasePath() + frontPemPath
	Url.FrontKeyPath = getBasePath() + frontKeyPath
	Url.FfmpegPath = getBasePath() + ffmpegPath
	Url.FfmpegConfPath = getBasePath() + ffmpegConfPath
	Url.SQLPath = getBasePath() + sqlPath
	Url.SrsConfPath = getBasePath() + srsConfPath
}

func getBasePath() string {
	return minioEndpoint + minioBucket
}
