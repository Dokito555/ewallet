package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dokito555/ewallet/ewallet-transaction/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ExtWallet struct {
	Log *logrus.Logger
	Config *viper.Viper
}

func NewExtWallet(log *logrus.Logger, conf *viper.Viper) *ExtWallet {
	return &ExtWallet{
		Log: log,
		Config: conf,
	}
}

func (e *ExtWallet) CreditBalance(ctx context.Context, req models.UpdateBalance, token string) (*models.UpdateBalanceResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		e.Log.Warnf("failed to marshal json: ", err)
		return nil, fmt.Errorf("failed to marshal json")
	} 

	// url := helpers.GetEnv("WALLET_HOST", "") + helpers.GetEnv("WALLET_ENDPOINT_CREDIT", "")
	url := e.Config.GetString("WALLET_HOST") + e.Config.GetString("WALLET_ENDPOINT_CREDIT")
	httpReq, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		e.Log.Warnf("failed to create wallet http request: ", err)
		return nil, fmt.Errorf("failed to create wallet http request")
	}
	httpReq.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		e.Log.Warnf("failed to connect wallet service: ", err)
		return nil, fmt.Errorf("failed to connect wallet service")
	}

	if resp.StatusCode != http.StatusOK {
		e.Log.Warnf("got error response from wallet service: %d", resp.StatusCode)
		return nil, fmt.Errorf("got error response from wallet service: %d", resp.StatusCode)
	}
	
	result := &models.UpdateBalanceResponse{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		e.Log.Warnf("failed to read response body: ", err)
		return nil, fmt.Errorf("failed to read response body")
	}
	defer resp.Body.Close()

	return result, nil
}

func (e *ExtWallet) DebitBalance(ctx context.Context, req models.UpdateBalance, token string) (*models.UpdateBalanceResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		e.Log.Warnf("failed to marshal json: ", err)
		return nil, fmt.Errorf("failed to marshal json")
	} 

	url := e.Config.GetString("WALLET_HOST") + e.Config.GetString("WALLET_ENDPOINT_DEBIT")
	httpReq, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		e.Log.Warnf("failed to create wallet http request")
		return nil,fmt.Errorf("failed to create wallet http request")
	}

	httpReq.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		e.Log.Warnf("failed to connect wallet service")
		return nil, fmt.Errorf("failed to connect wallet service")
	}

	if resp.StatusCode != http.StatusOK {
		e.Log.Warnf("got error response from wallet service: %d", resp.StatusCode)
		return nil, fmt.Errorf("got error response from wallet service: %d", resp.StatusCode)
	}
	
	result := &models.UpdateBalanceResponse{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		e.Log.Warnf("failed to read response body")
		return nil, fmt.Errorf("failed to read response body")
	}
	defer resp.Body.Close()

	return result, nil
}