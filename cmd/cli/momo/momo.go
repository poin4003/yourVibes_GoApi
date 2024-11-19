package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

type MomoConfig struct {
	PartnerCode  string
	AccessKey    string
	SecretKey    string
	RedirectURL  string
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

func NewMomoConfig() MomoConfig {
	return MomoConfig{
		PartnerCode:  "MOMO",
		AccessKey:    "F8BBA842ECF85",
		SecretKey:    "K951B6PE1waDMi640xX08PD3vg6EkVlz",
		RedirectURL:  "http://localhost:8080/bills/paidBill",
		IpnURL:       "http://localhost:8080/bills/paidBill",
		EndpointHost: "test-payment.momo.vn",
		EndpointPath: "/v2/gateway/api/create",
	}
}

func (config *MomoConfig) createSignature(
	billId string,
	orderInfo string,
	amount int,
	subject string,
	text string,
) (string, string) {
	requestId := fmt.Sprintf("%s%d.%s", config.PartnerCode, time.Now().UnixNano(), billId)
	orderId := requestId
	extraData := fmt.Sprintf("%s<splitText>%s", subject, text)

	rawSignature := fmt.Sprintf("accessKey=%s&amount=%d&extraData=%s&ipnUrl=%s&orderId=%s&orderInfo=%s&partnerCode=%s&redirectUrl=%s&requestId=%s&requestType=captureWallet",
		config.AccessKey, amount, extraData, config.IpnURL, orderId, orderInfo, config.PartnerCode, config.RedirectURL, requestId)

	fmt.Println("Raw signature:", rawSignature)

	h := hmac.New(sha256.New, []byte(config.SecretKey))
	h.Write([]byte(rawSignature))
	signature := hex.EncodeToString(h.Sum(nil))

	fmt.Println("Signature:", signature)

	return signature, requestId
}

func (config *MomoConfig) CreateRequestBody(
	billId string,
	orderInfo string,
	amount int,
	subject string,
	text string,
) (*RequestBody, error) {
	signature, requestId := config.createSignature(billId, orderInfo, amount, subject, text)
	return &RequestBody{
		PartnerCode: config.PartnerCode,
		AccessKey:   config.AccessKey,
		RequestId:   requestId,
		Amount:      amount,
		OrderId:     requestId,
		OrderInfo:   orderInfo,
		RedirectURL: config.RedirectURL,
		IpnURL:      config.IpnURL,
		ExtraData:   fmt.Sprintf("%s<splitText>%s", subject, text),
		RequestType: "captureWallet",
		Signature:   signature,
		Lang:        "en",
	}, nil
}

func (config *MomoConfig) SendRequest(body *RequestBody) (string, error) {
	requestBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	fmt.Printf("Request body: %s\n", requestBody)

	req, err := http.NewRequest("POST", "https://"+config.EndpointHost+config.EndpointPath, bytes.NewBuffer(requestBody))
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

	fmt.Printf("Response data: %s\n", responseData)

	var result map[string]interface{}
	err = json.Unmarshal(responseData, &result)
	if err != nil {
		return "", err
	}

	payUrl, ok := result["payUrl"].(string)
	if !ok || payUrl == "" {
		errorCode, _ := result["errorCode"].(string)
		errorMessage, _ := result["message"].(string)
		return "", fmt.Errorf("no payUrl found in response. ErrorCode: %s, Message: %s", errorCode, errorMessage)
	}

	return payUrl, nil
}

func momoCallbackHandler(c *gin.Context) {
	partnerCode := c.DefaultQuery("partnerCode", "")
	orderId := c.DefaultQuery("orderId", "")
	requestId := c.DefaultQuery("requestId", "")
	amount := c.DefaultQuery("amount", "")
	orderInfo := c.DefaultQuery("orderInfo", "")
	orderType := c.DefaultQuery("orderType", "")
	transId := c.DefaultQuery("transId", "")
	resultCode := c.DefaultQuery("resultCode", "")
	message := c.DefaultQuery("message", "")
	payType := c.DefaultQuery("payType", "")
	responseTime := c.DefaultQuery("responseTime", "")
	extraData := c.DefaultQuery("extraData", "")
	signature := c.DefaultQuery("signature", "")

	fmt.Printf("Callback data from Momo:\n")
	fmt.Printf("partnerCode: %s\n", partnerCode)
	fmt.Printf("orderId: %s\n", orderId)
	fmt.Printf("requestId: %s\n", requestId)
	fmt.Printf("amount: %s\n", amount)
	fmt.Printf("orderInfo: %s\n", orderInfo)
	fmt.Printf("orderType: %s\n", orderType)
	fmt.Printf("transId: %s\n", transId)
	fmt.Printf("resultCode: %s\n", resultCode)
	fmt.Printf("message: %s\n", message)
	fmt.Printf("payType: %s\n", payType)
	fmt.Printf("responseTime: %s\n", responseTime)
	fmt.Printf("extraData: %s\n", extraData)
	fmt.Printf("signature: %s\n", signature)

	if resultCode != "0" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid result code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Callback received successfully"})
}

func payByMomoHandler(c *gin.Context) {
	var req struct {
		BillId    string `json:"billId"`
		OrderInfo string `json:"orderInfo"`
		Amount    int    `json:"amount"`
		Subject   string `json:"subject"`
		Text      string `json:"text"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	momoConfig := NewMomoConfig()
	requestBody, err := momoConfig.CreateRequestBody(req.BillId, req.OrderInfo, req.Amount, req.Subject, req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	payUrl, err := momoConfig.SendRequest(requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payUrl": payUrl})
}

func main() {
	r := gin.Default()
	r.POST("/momo/pay", payByMomoHandler)
	r.GET("/bills/paidBill", momoCallbackHandler)
	r.Run(":8080")
}
