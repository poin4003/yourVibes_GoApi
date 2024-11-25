package payment

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/global"
	"io"
	"net/http"
	"time"
)

type MomoConfig struct {
	PartnerCode  string
	AccessKey    string
	SecretKey    string
	RedirectUrl  string
	IpnURL       string
	EndpointHost string
	EndpointPath string
}

type RequestBody struct {
	PartnerCode string `json:"partnerCode"`
	AccessKey   string `json:"accessKey"`
	RequestId   string `json:"requestId"`
	Amount      int    `json:"amount"`
	OrderId     string `json:"orderId"`
	OrderInfo   string `json:"orderInfo"`
	RedirectURL string `json:"redirectUrl"`
	IpnURL      string `json:"ipnUrl"`
	ExtraData   string `json:"extraData"`
	RequestType string `json:"requestType"`
	Signature   string `json:"signature"`
	Lang        string `json:"lang"`
}

func SendRequestToMomo(
	billId string,
	orderInfo string,
	amount int,
	subject string,
	text string,
	redirectUrl string,
) (string, error) {
	// 1. Init momo config
	config := global.Config.MomoSetting
	momoConfig := &MomoConfig{
		PartnerCode:  config.PartnerCode,
		AccessKey:    config.AccessKey,
		SecretKey:    config.SecretKey,
		RedirectUrl:  config.RedirectUrl,
		IpnURL:       config.IpnURL,
		EndpointHost: config.EndpointHost,
		EndpointPath: config.EndpointPath,
	}

	// 2. Create requestId and orderId
	requestId := fmt.Sprintf("%s%d.%s", momoConfig.PartnerCode, time.Now().UnixNano(), billId)
	orderId := requestId
	extraData := fmt.Sprintf("%s<splitText>%s<splitText>%s", subject, text, redirectUrl)

	// 3. Create raw signature
	rawSignature := fmt.Sprintf("accessKey=%s&amount=%d&extraData=%s&ipnUrl=%s&orderId=%s&orderInfo=%s&partnerCode=%s&redirectUrl=%s&requestId=%s&requestType=captureWallet",
		momoConfig.AccessKey, amount, extraData, momoConfig.IpnURL, orderId, orderInfo, momoConfig.PartnerCode, momoConfig.RedirectUrl, requestId)

	// 4. Create HMAC-SHA256 signature
	h := hmac.New(sha256.New, []byte(momoConfig.SecretKey))
	h.Write([]byte(rawSignature))
	signature := hex.EncodeToString(h.Sum(nil))

	// 5. Create Request body
	body := &RequestBody{
		PartnerCode: momoConfig.PartnerCode,
		AccessKey:   momoConfig.AccessKey,
		RequestId:   requestId,
		Amount:      amount,
		OrderId:     requestId,
		OrderInfo:   orderInfo,
		RedirectURL: momoConfig.RedirectUrl,
		IpnURL:      momoConfig.IpnURL,
		ExtraData:   extraData,
		RequestType: "captureWallet",
		Signature:   signature,
		Lang:        "en",
	}

	fmt.Println(body)

	// 6. Send http request to momo api
	requestBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://"+momoConfig.EndpointHost+momoConfig.EndpointPath, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(responseData, &result)
	if err != nil {
		return "", err
	}

	// 6. Check payment url
	payUrl, ok := result["payUrl"].(string)
	if !ok || payUrl == "" {
		errorCode, _ := result["errorCode"].(string)
		errorMessage, _ := result["message"].(string)
		return "", fmt.Errorf("no payUrl found in response. ErrorCode: %s, Message: %s", errorCode, errorMessage)
	}

	return payUrl, nil
}
