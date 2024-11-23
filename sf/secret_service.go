package sf

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/maypok86/otter"
	"github.com/omkarsrepo/server-framework/sf/boom"
	"github.com/omkarsrepo/server-framework/sf/json"
	"github.com/rs/zerolog"
	"os"
	"sync"
	"time"
)

var singletonSecretService *secretService
var once sync.Once

type SecretService interface {
	ValueOf(secretKey string) (string, *boom.Exception)
}

type secretService struct {
	secretTokenCache *otter.Cache[string, any]
	secretCache      *otter.CacheWithVariableTTL[string, any]
	restyClient      *resty.Client
	config           ConfigService
	logger           *zerolog.Logger
}

func SecretServiceInstance() SecretService {
	once.Do(func() {
		cache := Cache()
		restyClient := resty.New().
			SetJSONMarshaler(jsoniter.ConfigCompatibleWithStandardLibrary.Marshal).
			SetJSONUnmarshaler(jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal)

		secretTokenCache := cache.New(2, time.Minute*59)
		secretCache := cache.NewVariable(50)

		loggerInstance := LoggerServiceInstance()

		singletonSecretService = &secretService{
			secretTokenCache: &secretTokenCache,
			secretCache:      &secretCache,
			restyClient:      restyClient,
			config:           ConfigServiceInstance(),
			logger:           loggerInstance.GetZeroLogger(),
		}
	})

	return singletonSecretService
}

func (props *secretService) fetchSecretToken() (string, *boom.Exception) {
	expectedBody := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     props.config.GetString("clientIds.hashicorp"),
		"client_secret": os.Getenv("hashicorpSecret"),
		"audience":      "https://api.hashicorp.cloud",
	}
	var responseResult map[string]interface{}

	resp, err := props.restyClient.R().
		SetBody(&expectedBody).
		SetResult(&responseResult).
		Post("https://auth.idp.hashicorp.com/oauth2/token")

	if err != nil || resp.StatusCode() >= 300 {
		props.logger.Error().Err(err).Msg("Failed to fetch secret token")
		return "", boom.InternalServerError()
	}

	val, err := json.GetObjectValue(&responseResult, "access_token")
	if err != nil {
		props.logger.Error().Err(err).Msgf("Failed to destructure value access_token for result: %+v", responseResult)
		return "", boom.InternalServerError()
	}

	return val.(string), nil
}

var secretTokenCacheKey = "FetchSecretToken"

func (props *secretService) getSecretToken() (string, *boom.Exception) {
	secretToken, ok := props.secretTokenCache.Get(secretTokenCacheKey)
	if !ok {
		secretToken, exp := props.fetchSecretToken()
		if exp != nil {
			props.secretTokenCache.Clear()
			return "", exp
		}

		props.secretTokenCache.Set(secretTokenCacheKey, secretToken)
		return secretToken, nil
	}

	return secretToken.(string), nil
}

func (props *secretService) fetchSecret(secretName string) (string, *boom.Exception) {
	organizationId := props.config.GetString("hashicorp.organizationId")
	projectId := props.config.GetString("hashicorp.projectId")
	env := props.config.GetString("env")

	secretToken, exp := props.getSecretToken()
	if exp != nil {
		return "", exp
	}

	var responseResult map[string]interface{}

	baseUrl := "https://api.cloud.hashicorp.com/secrets/2023-06-13/organizations"
	resp, err := props.restyClient.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(secretToken).
		SetResult(&responseResult).
		Get(fmt.Sprintf("%s/%s/projects/%s/apps/%s/open/%s", baseUrl, organizationId, projectId, env, secretName))

	if err != nil || resp.StatusCode() != 200 {
		props.logger.Error().Err(err).Msgf("Failed to fetch secret value for key %s", secretName)
		return "", boom.InternalServerError()
	}

	val, err := json.GetObjectValue(&responseResult, "secret.version.value")
	if err != nil {
		props.logger.Error().Err(err).
			Msgf("Failed to destructure value 'secret.version.value' for result: %+v", responseResult)
		return "", boom.InternalServerError()
	}

	return val.(string), nil
}

func (props *secretService) ValueOf(secretKey string) (string, *boom.Exception) {
	secretName := props.config.GetString(secretKey)

	secret, ok := props.secretCache.Get(secretName)
	if !ok {
		secret, exp := props.fetchSecret(secretName)
		if exp != nil {
			return "", exp
		}

		props.secretCache.Set(secretName, secret, time.Minute*90)

		return secret, nil
	}

	return secret.(string), nil
}
