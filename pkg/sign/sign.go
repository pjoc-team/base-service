package sign

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/pjoc-team/base-service/pkg/logger"
	"github.com/pjoc-team/base-service/pkg/model"
	"github.com/pjoc-team/base-service/pkg/service"
	"github.com/pjoc-team/base-service/pkg/util"
	"github.com/pjoc-team/pay-proto/go"
)

func init() {
	initCheckSignMap()
}

type CheckSignValidator struct {
	paramsCompacter ParamsCompacter
}

func (validator *CheckSignValidator) Validate(request pay.PayRequest, cfg service.GatewayConfig) (e error) {
	paramsString := validator.paramsCompacter.ParamsToString(request)
	if logger.Log.IsDebugEnabled() {
		logger.Log.Debugf("Build interface: %v to string: %v", request, paramsString)
	}
	if config, exists := (*cfg.AppIdAndMerchantMap)[request.AppId]; !exists || config.MerchantRSAPublicKey == "" {
		e = fmt.Errorf("could'nt found config of appId: %v", request.AppId)
		logger.Log.Errorf("could'nt found config of appId: %v request: %v", request.AppId, request)
		return e
	} else {
		e = CheckSign(request.GetCharset(), paramsString, request.GetSign(), config, request.SignType)
	}

	return
}

func NewCheckSignValidator() *CheckSignValidator {
	validator := &CheckSignValidator{}
	validator.paramsCompacter = NewParamsCompacter(&pay.PayRequest{}, "json", []string{"sign"}, true, "&", "=")
	return validator
}

type CheckSignInterface interface {
	checkSign(source []byte, signMsg string, key string) error
	sign(source []byte, key string) (string, error)
	getCheckSignKey(config model.MerchantConfig) string
	getSignKey(config model.MerchantConfig) string
	signType() string
}

var checkSignMap = make(map[string]CheckSignInterface)

func initCheckSignMap() {
	checkSignMap[SIGN_TYPE_MD5] = &Md5{}
	checkSignMap[SIGN_TYPE_SHA256_WITH_RSA] = &Sha256WithRSA{}
}

func CheckSign(charset string, source string, signMsg string, config model.MerchantConfig, signType string) (err error) {
	if signType == "" {
		signType = SIGN_TYPE_SHA256_WITH_RSA
	}
	signFunc := checkSignMap[signType]
	var sourceBytes []byte
	if key := signFunc.getCheckSignKey(config); key == "" {
		err = errors.New("could'nt found key!")
		logger.Log.Errorf("Could'nt get key from config: %v", config)
		return err
	} else if sourceBytes, err = stringToBytes(source, charset); err != nil {
		logger.Log.Errorf("Failed to get charset: %s, error: %s", charset, err.Error())
		return fmt.Errorf("unknown charset: %s", charset)
	} else if signFunc == nil {
		logger.Log.Errorf("Failed to get signType: %s, error: %s", signType, err.Error())
		e := fmt.Errorf("unknown signtype: %s", charset)
		return e
	} else if err = signFunc.checkSign(sourceBytes, signMsg, key); err != nil {
		logger.Log.Errorf("Failed to check sign! error: %s", err.Error())
		e := fmt.Errorf("failed to check sign!")
		return e
	} else {
		return nil
	}
}

func GenerateSign(charset string, source string, config model.MerchantConfig, signType string) (sign string, err error) {
	signFunc := checkSignMap[signType]
	var sourceBytes []byte
	if key := signFunc.getSignKey(config); key == "" {
		err = errors.New("could'nt found key!")
		logger.Log.Errorf("Could'nt get key from config: %v", config)
		return
	} else if sourceBytes, err = stringToBytes(source, charset); err != nil {
		logger.Log.Errorf("Failed to get charset: %s, error: %s", charset, err.Error())
		err = fmt.Errorf("unknown charset: %s", charset)
		return
	} else if signFunc == nil {
		logger.Log.Errorf("Failed to get signType: %s, error: %s", signType, err.Error())
		err = fmt.Errorf("unknown signtype: %s", charset)
		return
	} else if sign, err = signFunc.sign(sourceBytes, key); err != nil {
		logger.Log.Errorf("Failed to sign! error: %s", err.Error())
		err = fmt.Errorf("failed to sign")
		return
	} else {
		return
	}
}

type Md5 struct {
}

func (m *Md5) getCheckSignKey(config model.MerchantConfig) string {
	return config.Md5Key
}

func (m *Md5) getSignKey(config model.MerchantConfig) string {
	return config.Md5Key
}

func (m *Md5) sign(source []byte, key string) (string, error) {
	buffer := bytes.NewBuffer(source)
	buffer.Write([]byte(key))
	b := buffer.Bytes()
	sum := md5.Sum(b)
	s := hex.EncodeToString(sum[:])
	return s, nil
}

func (m *Md5) checkSign(source []byte, signMsg string, key string) error {
	generated, e := m.sign(source, key)
	if e != nil {
		logger.Log.Errorf("Failed to generate sign! error: %v", e.Error())
		return e
	}
	if !util.EqualsIgnoreCase(generated, signMsg) {
		e := errors.New("check sign error")
		logger.Log.Warnf("Failed to check sign! ours: %v actual: %v", generated, signMsg)
		return e
	}

	return nil
}

func (*Md5) signType() string {
	return SIGN_TYPE_MD5
}

type Sha256WithRSA struct {
}

func (s *Sha256WithRSA) getCheckSignKey(config model.MerchantConfig) string {
	return config.MerchantRSAPublicKey
}

func (s *Sha256WithRSA) getSignKey(config model.MerchantConfig) string {
	return config.GatewayRSAPrivateKey
}

func (s *Sha256WithRSA) sign(source []byte, key string) (sign string, err error) {
	signBytes, err := SignPKCS1v15WithStringKey(source, key, crypto.SHA256)
	if err != nil {
		logger.Log.Errorf("Failed to sign! error: %v key: %v", err.Error(), key)
		return
	}
	sign = base64.StdEncoding.EncodeToString(signBytes)
	logger.Log.Debugf("Encode source: %v to sign: %v", string(source), sign)
	return
}

func (*Sha256WithRSA) checkSign(source []byte, signMsg string, key string) (err error) {
	sign, err := base64.StdEncoding.DecodeString(signMsg)
	if err != nil {
		logger.Log.Errorf("Failed to check sign! decode sign: %v with error: %v", signMsg, err.Error())
		return
	}
	err = VerifyPKCS1v15WithStringKey(source, sign, key, crypto.SHA256)
	if err != nil {
		logger.Log.Errorf("Failed to check sign! check source: %v sign: %v with error: %v", string(source), signMsg, err.Error())
		return
	}
	return err
}

func (*Sha256WithRSA) signType() string {
	return SIGN_TYPE_SHA256_WITH_RSA
}
