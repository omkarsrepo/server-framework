package sfk

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/maypok86/otter"
	"github.com/omkarsrepo/server-framework/sfk/boom"
	"github.com/omkarsrepo/server-framework/sfk/json"
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

func (s *secretService) fetchSecretToken() (string, *boom.Exception) {
	expectedBody := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     s.config.GetString("clientIds.hashicorp"),
		"client_secret": os.Getenv("hashicorpSecret"),
		"audience":      "https://api.hashicorp.cloud",
	}
	var responseResult map[string]interface{}

	resp, err := s.restyClient.R().
		SetBody(&expectedBody).
		SetResult(&responseResult).
		Post("https://auth.idp.hashicorp.com/oauth2/token")

	if err != nil || resp.StatusCode() >= 300 {
		s.logger.Error().Err(err).Msg("Failed to fetch secret token")
		return "", boom.InternalServerError()
	}

	val, err := json.GetObjectValue(&responseResult, "access_token")
	if err != nil {
		s.logger.Error().Err(err).Msgf("Failed to destructure value access_token for result: %+v", responseResult)
		return "", boom.InternalServerError()
	}

	return val.(string), nil
}

var secretTokenCacheKey = "FetchSecretToken"

func (s *secretService) getSecretToken() (string, *boom.Exception) {
	secretToken, ok := s.secretTokenCache.Get(secretTokenCacheKey)
	if !ok {
		secretToken, exp := s.fetchSecretToken()
		if exp != nil {
			s.secretTokenCache.Clear()
			return "", exp
		}

		s.secretTokenCache.Set(secretTokenCacheKey, secretToken)
		return secretToken, nil
	}

	return secretToken.(string), nil
}

func (s *secretService) fetchSecret(secretName string) (string, *boom.Exception) {
	organizationId := s.config.GetString("hashicorp.organizationId")
	projectId := s.config.GetString("hashicorp.projectId")
	env := s.config.GetString("env")

	secretToken, exp := s.getSecretToken()
	if exp != nil {
		return "", exp
	}

	var responseResult map[string]interface{}

	baseUrl := "https://api.cloud.hashicorp.com/secrets/2023-06-13/organizations"
	resp, err := s.restyClient.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(secretToken).
		SetResult(&responseResult).
		Get(fmt.Sprintf("%s/%s/projects/%s/apps/%s/open/%s", baseUrl, organizationId, projectId, env, secretName))

	if err != nil || resp.StatusCode() != 200 {
		s.logger.Error().Err(err).Msgf("Failed to fetch secret value for key %s", secretName)
		return "", boom.InternalServerError()
	}

	val, err := json.GetObjectValue(&responseResult, "secret.version.value")
	if err != nil {
		s.logger.Error().Err(err).
			Msgf("Failed to destructure value 'secret.version.value' for result: %+v", responseResult)
		return "", boom.InternalServerError()
	}

	return val.(string), nil
}

func (s *secretService) ValueOf(secretKey string) (string, *boom.Exception) {
	secretName := s.config.GetString(secretKey)

	secret, ok := s.secretCache.Get(secretName)
	if !ok {
		secret, exp := s.fetchSecret(secretName)
		if exp != nil {
			return "", exp
		}

		s.secretCache.Set(secretName, secret, time.Minute*90)

		return secret, nil
	}

	return secret.(string), nil
}
