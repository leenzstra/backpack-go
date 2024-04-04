package auth

import (
	"crypto"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

const (
	DefaultWindow = 10000
)

var _ Authenticator = (*AuthenticatorImpl)(nil)

type AuthHeaders struct {
	XTimestamp int64  `json:"X-Timestamp" mapstructure:"X-Timestamp"`
	XWindow    int    `json:"X-Window" mapstructure:"X-Window"`
	XAPIKey    string `json:"X-API-Key" mapstructure:"X-API-Key"`
	XSignature string `json:"X-Signature" mapstructure:"X-Signature"`
}

func (h *AuthHeaders) decode() (map[string]interface{}, error) {
	m := make(map[string]interface{})

	err := mapstructure.Decode(h, &m)

	return m, err
}

func (h *AuthHeaders) Map() map[string]string {
	m := make(map[string]string)
	md, _ := h.decode()

	for k, v := range md {
		m[k] = fmt.Sprintf("%v", v)
	}

	return m
}

type Authenticator interface {
	Authenticate(instruction Instruction, body interface{}) (*AuthHeaders, error)
	SetWindow(window int)
}

func NewAuthenticator(window int, secretKey, apiKey string) (Authenticator, error) {
	impl := &AuthenticatorImpl{
		window: window,
		apiKey: apiKey,
	}

	pk, err := decodePrivateKey(secretKey)
	if err != nil {
		return nil, err
	}

	impl.privateKey = pk

	return impl, nil
}

type AuthenticatorImpl struct {
	window     int
	privateKey ed25519.PrivateKey
	apiKey     string
}

func (impl *AuthenticatorImpl) SetWindow(window int) {
	if window < 5000 || window > 60000 {
		impl.window = DefaultWindow
	}

	impl.window = window
}

// Body is flat struct with json tags or map[string]interface{}
func (impl *AuthenticatorImpl) Authenticate(instruction Instruction, body interface{}) (*AuthHeaders, error) {
	bodyQuery, err := createQuery(body)
	if err != nil {
		return nil, fmt.Errorf("auth query err: %v", err)
	}

	ts := time.Now().UTC().UnixMilli()

	query := fmt.Sprintf("instruction=%s&%s&timestamp=%d&window=%d", instruction, bodyQuery, ts, impl.window)
	query = strings.ReplaceAll(query, "&&", "&")

	// fmt.Println(query, body, bodyQuery)

	signature, err := impl.sign([]byte(query))
	if err != nil {
		return nil, fmt.Errorf("auth signature err: %v", err)
	}

	signatureb64 := base64.StdEncoding.EncodeToString(signature)

	return &AuthHeaders{XTimestamp: ts, XWindow: impl.window, XAPIKey: impl.apiKey, XSignature: signatureb64}, nil
}

func (impl *AuthenticatorImpl) sign(data []byte) ([]byte, error) {
	return impl.privateKey.Sign(nil, data, crypto.Hash(0))
}

// decode private key from base64 string (API Secret)
func decodePrivateKey(privateKey string) (ed25519.PrivateKey, error) {
	pk, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}

	return ed25519.NewKeyFromSeed(pk), nil
}

func createQuery(body interface{}) (string, error) {
	if body == nil {
		return "", nil
	}

	bodyType := reflect.TypeOf(body)
	bodyValue := reflect.ValueOf(body)

	switch bodyType.Kind() {
	case reflect.Map:
		nbody := make(map[string]interface{})
		for _, key := range bodyValue.MapKeys() {
			nbody[key.String()] = bodyValue.MapIndex(key).Interface()
		}
		return mapToQuery(nbody)

	case reflect.Struct:
		return structToQuery(body)

	default:
		return "", fmt.Errorf("unsupported body type: %v", bodyType)
	}
}

func mapToQuery(params map[string]interface{}) (string, error) {
	values := url.Values{}

	for key, value := range params {
		valStr := fmt.Sprintf("%v", value)
		values.Add(key, valStr)
	}

	// automaticaly sorted by key
	return values.Encode(), nil
}

func structToQuery(v interface{}) (string, error) {
	values := url.Values{}

	val := reflect.ValueOf(v)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		jsonTag := fieldType.Tag.Get("json")

		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		values.Add(jsonTag, fmt.Sprintf("%v", field.Interface()))
	}

	// automaticaly sorted by key
	return values.Encode(), nil
}
