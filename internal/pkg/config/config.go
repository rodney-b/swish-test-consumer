package config

import (
	"reflect"
	"strings"
	"sync"

	"github.com/rodney-b/swish-test-consumer/pkg/utilities/env"
)

type ConfigProvider interface {
	IsDevelopment() bool
	GetAppName() string
	GetConsumerCA() []byte
	GetConsumerCert() []byte
	GetConsumerCertKey() []byte
	GetHealthcheckPort() string
	GetHealthcheckServicePrefix() string
	GetMessageQueueClientCA() []byte
	GetMessageQueueClientCert() []byte
	GetMessageQueueClientCertKey() []byte
	GetMessageQueueGroupID() string
	GetMessageQueueTopics() []string
	GetMessageQueueURL() string
	GetOTelHTTPReceiverURL() string
	GetOtelStdoutExporterEnabled() bool
	GetStage() string
}

// appCofnig implements ConfigProvider. It "provides" all its values from environment variables.
// each data type must have a case statement in utilities.env.Get()
// note: all types that can be casted to from int, are already covered
type appConfig struct {
	appName                   string `envname:"APP_NAME"`
	consumerCA                string `envname:"CONSUMER_CA"`
	consumerCert              string `envname:"CONSUMER_CRT"`
	consumerCertKey           string `envname:"CONSUMER_KEY"`
	healthcheckPort           string `envname:"HEALTHCHECK_PORT"`
	healthcheckServicePrefix  string `envname:"HEALTHCHECK_SERVICE_PREFIX"`
	messageQueueClientCA      string `envname:"MESSAGE_QUEUE_CA"`
	messageQueueClientCert    string `envname:"MESSAGE_QUEUE_CRT"`
	messageQueueClientCertKey string `envname:"MESSAGE_QUEUE_KEY"`
	messageQueueGroupID       string `envname:"MESSAGE_QUEUE_GROUP_ID"`
	messageQueueTopics        string `envname:"MESSAGE_QUEUE_TOPICS"`
	messageQueueURL           string `envname:"MESSAGE_QUEUE_URL"`
	otelHTTPReceiverURL       string `envname:"OTEL_HTTP_RECEIVER_URL"`
	otelStdoutExporterEnabled string `envname:"OTEL_STDOUT_EXPORTER_ENABLED"`
	stage                     string `envname:"STAGE"`
}

func (ac *appConfig) GetAppName() string {
	return ac.appName
}

func (ac *appConfig) GetConsumerCA() []byte {
	return []byte(ac.consumerCA)
}

func (ac *appConfig) GetConsumerCert() []byte {
	return []byte(ac.consumerCert)
}

func (ac *appConfig) GetConsumerCertKey() []byte {
	return []byte(ac.consumerCertKey)
}

func (ac *appConfig) GetHealthcheckPort() string {
	return ac.healthcheckPort
}

func (ac *appConfig) GetHealthcheckServicePrefix() string {
	return ac.healthcheckServicePrefix
}

func (ac *appConfig) GetMessageQueueClientCA() []byte {
	return []byte(ac.messageQueueClientCA)
}

func (ac *appConfig) GetMessageQueueClientCert() []byte {
	return []byte(ac.messageQueueClientCert)
}

func (ac *appConfig) GetMessageQueueClientCertKey() []byte {
	return []byte(ac.messageQueueClientCertKey)
}

func (ac *appConfig) GetMessageQueueGroupID() string {
	return ac.messageQueueGroupID
}

func (ac *appConfig) GetMessageQueueTopics() []string {
	funcOnce := sync.OnceValue(func() []string {
		return strings.Split(ac.messageQueueTopics, ",")
	})

	return funcOnce()
}

func (ac *appConfig) GetMessageQueueURL() string {
	return ac.messageQueueURL
}

func (ac *appConfig) GetOTelHTTPReceiverURL() string {
	return ac.otelHTTPReceiverURL
}

func (ac *appConfig) GetOtelStdoutExporterEnabled() bool {
	funcOnce := sync.OnceValue(func() bool {
		if ac.otelStdoutExporterEnabled == "true" {
			return true
		}

		return false
	})

	return funcOnce()
}

func (ac *appConfig) GetStage() string {
	return ac.stage
}

// helper funcs
func (ac *appConfig) IsDevelopment() bool {
	return ac.GetStage() != Production && ac.GetStage() != Staging
}

// These are for config values the app shouldn't start without.
var initAppConfig = sync.OnceValues(func() (*appConfig, error) {
	ac := appConfig{}
	appConfVal := reflect.ValueOf(&ac).Elem()
	appConfType := reflect.TypeOf(ac)

	for i := range appConfVal.NumField() {
		fieldInfo := appConfType.Field(i)
		fieldTag := fieldInfo.Tag
		fieldValue := appConfVal.Field(i)
		fieldPtr := fieldValue.Addr().UnsafePointer()
		unsafeFieldValue := reflect.NewAt(fieldValue.Type(), fieldPtr).Elem()

		err := env.Get(fieldTag.Get("envname"), unsafeFieldValue)
		if err != nil {
			return nil, err
		}
	}

	return &ac, nil
})

func InitAppConfig() (*appConfig, error) {
	return initAppConfig()
}
