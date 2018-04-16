package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"strconv"
	"checkout-service/client/rest"
	log "github.com/sirupsen/logrus"
	"errors"
)

type Option int
const (
	CREATE_BASKET  Option = iota +1
	DISPLAY_BASKET
	DELETE_BASKET
	SCAN_VOUCHER
	SCAN_MUG
	SCAN_TSHIRT
	DISPLAY_TOTAL
	EXIT
)


var prompt=fmt.Sprintf(
	`
What action would you like to take?

Create checkout basket:  | %d
Display checkout basket: | %d <basket_id>
Delete checkout basket:  | %d <basket_id>
Scan Voucher for basket: | %d [<basket_id>] (def. last basket used)
Scan Mug for basket:     | %d [<basket_id>] (def. last basket used)
Scan T-Shirt for basket: | %d [<basket_id>] (def. last basket used)
Display total for basket:| %d [<basket_id>] (def. last basket used)
Exit:                    | %d

`, CREATE_BASKET, DISPLAY_BASKET, DELETE_BASKET, SCAN_VOUCHER, SCAN_MUG,
	SCAN_TSHIRT,DISPLAY_TOTAL, EXIT)

func main() {


	driver := Driver{Client:*rest.NewClient()}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(prompt)


	for scanner.Scan() {
		input := scanner.Text()

		option, param, err := parse(input)

		if err != nil {
			fmt.Println(err)
			fmt.Print(prompt)
			continue
		}

		switch  Option(option){
		case CREATE_BASKET:
			basketId, err := driver.CreateBasket()
			driver.resolvedParam = basketId
			driver.handleResponse(err, fmt.Sprintf("Created basket %s\n", basketId))
		case EXIT:
			fmt.Println("Exiting....")
			os.Exit(0)
		default:

			driver.resolveBasketId(param)
			if len(driver.resolvedParam)==0 { break }

			switch Option(option){

			case DISPLAY_BASKET:
				basket,err := driver.GetBasket(driver.resolvedParam)
				driver.handleResponse(err, fmt.Sprintf("Basket content: %s\n", basket))
			case DELETE_BASKET:
				err := driver.DeleteBasket(driver.resolvedParam)
				driver.handleResponse(err, fmt.Sprintf("Deleted basket %s\n",driver.resolvedParam))
			case SCAN_VOUCHER:
				err := driver.ScanProduct("VOUCHER", driver.resolvedParam)
				driver.handleResponse(err, fmt.Sprintf("Added VOUCHER to basket %s\n", driver.resolvedParam) )
			case SCAN_MUG:
				err := driver.ScanProduct("MUG", driver.resolvedParam)
				driver.handleResponse(err, fmt.Sprintf("Added MUG to basket %s\n", driver.resolvedParam))
			case SCAN_TSHIRT:
				err := driver.ScanProduct("TSHIRT", driver.resolvedParam)
				driver.handleResponse(err, fmt.Sprintf("Added TSHIRT to basket %s\n", driver.resolvedParam))
			case DISPLAY_TOTAL:
				total, err := driver.GetTotal(driver.resolvedParam)
				driver.handleResponse(err, fmt.Sprintf("Total for basket %s: %f\n", driver.resolvedParam, total))

			}

		}

		fmt.Print(prompt)

	}


}


func parse(input string) (option Option, param string, err error){
	split := strings.Split(input, " ")

	if len(split) == 0 {
		return -1 ,"", errors.New("Malformed input")
	}

	optionInt, err := strconv.Atoi(split[0])
	if err !=nil {
		return -1 ,"", errors.New("Malformed input")
	}

	option = Option(optionInt)

	if len(split) == 2 {
		param = strings.TrimSpace(split[1])
	}

	return option, param, err

}

type Driver struct {
	rest.Client
	recentBasketId string
	resolvedParam  string
}

func (driver *Driver) resolveBasketId(optional string) {
	if len(optional)==0 {
		if len(driver.recentBasketId) == 0 {
			driver.resolvedParam = ""
		}else{
			driver.resolvedParam = driver.recentBasketId
		}
	}else{
		driver.resolvedParam =optional
	}
}


func (driver *Driver) handleResponse(err error, message string){
	if err == nil {
		driver.recentBasketId = driver.resolvedParam
		fmt.Println(message)
	}else{
		fmt.Println(err)
	}
}

func init(){
	log.SetLevel(log.InfoLevel)
}



