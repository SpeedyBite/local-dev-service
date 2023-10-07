package engine

// Ignore some environment variables that only valid on k8s or in k8s clusters
var IgnoreUnusedEnvVars = []string{
	"UWSGI_PORT",
	"WORKER_LOG_LEVEL",
	"APP_DEBUG",
	"DOMAIN",
	"FLASK_LOG_LEVEL",
	"LOG_STDOUT",
	"METRIC_PORT",
	"NUM_OF_PROCESSES",
	"STATSD_HOST",
	"STATSD_PORT",
	"WORKER_LOG_LEVEL",
	"AUTORELOAD_VALUE",
	"APP_ENV",
}
