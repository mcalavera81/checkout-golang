package rest

import (
	"fmt"
	"errors"
	"github.com/go-resty/resty"
	"encoding/json"
	"strings"
	"net/http"
	"checkout-service/client/config"
	"net"
	"strconv"
	log "github.com/sirupsen/logrus"

)

type Client struct {
	BaseUri string
}

type Response struct {
	Id string
	Total float64
}

func NewClient() *Client{
	c := config.GetConfig()
	serverAddr := strings.Join([]string{c.Host, strconv.Itoa(c.Port)},":")
	_, err:=net.Dial("tcp", serverAddr)
	if err !=nil {
		log.Fatal(err)
	}else{
		log.Infof("Server found on %s", serverAddr)
	}
	return &Client{c.BaseURI}
}

func (c *Client) CreateBasket() (string, error) {
	resp, err := resty.SetRESTMode().R().
		Post(c.BaseUri)

	switch {
	case err != nil:
		return "", err
	case resp.StatusCode()!=http.StatusCreated:
		return  "", errors.New(fmt.Sprintf("Error: %s", resp))
	}

	var body Response
	json.NewDecoder(strings.NewReader(resp.String())).Decode(&body)
	return body.Id, nil
}

type BasketClient struct {
	Id   string
	Items map[string]int
}


func (c *Client) GetBasket(basketId string) (string,error) {
	resp, err := resty.SetRESTMode().R().
		SetPathParams(map[string]string{
		"basketId": basketId,
	}).Get(strings.Join([]string{c.BaseUri,"{basketId}"},""))

	switch {
		case err != nil:
			return "", err
		case resp.StatusCode()!=http.StatusOK:
			return  "", errors.New(fmt.Sprintf("Error: %s", resp))
	}


	b := BasketClient{Items:make(map[string]int)}
	if err:=json.Unmarshal(resp.Body(), &b);err != nil{
		return "", err
	}
	if bytes, err := json.MarshalIndent(b,"","\t");err!=nil{
		return "", err
	}else{
		return string(bytes), nil
	}


}

func (c *Client) DeleteBasket(basketId string) error {

	resp, err := resty.SetRESTMode().R().
		SetPathParams(map[string]string{
		"basketId": basketId,
	}).Delete(strings.Join([]string{c.BaseUri,"{basketId}"},""))

	switch {
		case err != nil:
			return  err
		case resp.StatusCode()!=http.StatusOK:
			return  errors.New(fmt.Sprintf("Error: %s", resp))
	}

	return err

}

func (c *Client) ScanProduct(item string, basketId string) error {
	resp, err := resty.SetRESTMode().R().
		SetPathParams(map[string]string{
		"basketId": basketId,
	}).
		SetBody(map[string]interface{}{"code": item}).
		Put(strings.Join([]string{c.BaseUri,"{basketId}"},""))

	switch {
		case err != nil:
			return  err
		case resp.StatusCode()!=http.StatusOK:
			return  errors.New(fmt.Sprintf("Error: %s", resp))
	}


	return err

}

func (c *Client) GetTotal(basketId string) (float64,error){
	resp, err := resty.SetRESTMode().R().
		SetPathParams(map[string]string{
		"basketId": basketId,
	}).Get(strings.Join([]string{c.BaseUri,"{basketId}/total"},""))

	switch {
		case err != nil:
			return -1.0, err
		case resp.StatusCode()!=http.StatusOK:
			return  -1.0, errors.New(fmt.Sprintf("Error: %s", resp))
	}


	var body Response
	json.NewDecoder(strings.NewReader(resp.String())).Decode(&body)
	return body.Total, nil

}
