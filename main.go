package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"

	"github.com/jaswdr/faker"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	UserType string `json:"userType"`
	Tckn     string `json:"tckn"`
	Mobile   string `json:"mobile"`
}

type Vendor struct {
	Name          string `json:"name"`
	TaxNo         string `json:"taxNo"`
	ExternalVCode string `json:"externalVCode"`
}

type DealerAndBuyer struct {
	Name  string `json:"name"`
	TaxNo string `json:"taxNo"`
}

type DealerSite struct {
	Name           string `json:"name"`
	DealerId       int    `json:"dealerId"`
	ExternalVCode  string `json:"externalVCode"`
	ExternalDSCode string `json:"externalDSCode"`
}

type BuyerSite struct {
	Name           string `json:"name"`
	BuyerId        int    `json:"buyerId"`
	ExternalVCode  string `json:"externalVCode"`
	ExternalDSCode string `json:"externalDSCode"`
	ExternalBSCode string `json:"externalBSCode"`
}

func CreateVendor(faker faker.Faker, code string) *Vendor {
	vendor := &Vendor{
		Name:          faker.Person().Name(),
		TaxNo:         strconv.Itoa(faker.RandomNumber(11)),
		ExternalVCode: code,
	}

	return vendor
}

func CreateUser(faker faker.Faker, userType string) *User {

	user := &User{
		Username: faker.Person().Name(),
		Password: "google",
		Email:    faker.Internet().Email(),
		UserType: userType,
		Tckn:     strconv.Itoa(faker.RandomNumber(11)),
		Mobile:   strconv.Itoa(faker.RandomNumber(10)),
	}

	return user

}

func CreateDealerAndBuyer(faker faker.Faker) *DealerAndBuyer {
	dealer := &DealerAndBuyer{
		Name:  faker.Person().FirstName(),
		TaxNo: strconv.Itoa(faker.RandomNumber(11)),
	}
	return dealer
}

func CreateDealerSite(faker faker.Faker, dealerId int, vcode string, dscode string) *DealerSite {
	dealerSite := &DealerSite{
		Name:           faker.Person().Name(),
		DealerId:       dealerId,
		ExternalVCode:  vcode,
		ExternalDSCode: dscode,
	}

	return dealerSite
}

func CreateBuyerSite(faker faker.Faker, bi int, vcode string, dscode string) *BuyerSite {
	buyerSite := &BuyerSite{
		Name:           faker.Company().Name(),
		BuyerId:        bi,
		ExternalVCode:  vcode,
		ExternalDSCode: dscode,
		ExternalBSCode: faker.Lorem().Word(),
	}

	return buyerSite
}

type UserEntity struct {
	BuyerSiteTableRef  any `json:"buyerSiteTableRef"`
	DealerSiteTableRef any `json:"dealerSiteTableRef"`
	VendorTableRef     any `json:"vendorTableRef"`
	UserId             any `json:"userId"`
}

func CreateUserEntity(faker faker.Faker, randomid int) *UserEntity {
	randId := randomId(3)

	cols := []string{
		"BuyerSiteTableRef", "DealerSiteTableRef", "VendorTableRef",
	}
	randomCol := cols[randId-1]
	userEntity := UserEntity{nil, nil, nil, "0"}

	randomColumn := getAttr(&userEntity, randomCol)
	randomColumn.Set(reflect.ValueOf(strconv.Itoa(randomid)))
	userEntity.UserId = strconv.Itoa(randomid)

	return &userEntity

}

type Vds struct {
	VendorId     int `json:"vendorId"`
	DealerSiteId int `json:"dealerSiteId"`
}

func CreateVds(vRandId int, dRandId int) *Vds {
	vds := Vds{
		VendorId:     vRandId,
		DealerSiteId: dRandId,
	}
	return &vds
}

type Vdsbs struct {
	VdsRltnId   int `json:"vdsRltnId"`
	BuyerSiteId int `json:"buyerSiteId"`
}

func CreateVdsbs(vdsbsRandId int, bRandId int) *Vdsbs {
	vdsbs := Vdsbs{
		VdsRltnId:   vdsbsRandId,
		BuyerSiteId: bRandId,
	}
	return &vdsbs
}

func genCode(faker faker.Faker) string {
	return faker.RandomStringWithLength(10)
}

func RandomUserType() string {
	userType := []string{"SA", "VA", "V", "B", "BA", "D", "DA"}
	randomIndex := rand.Intn(len(userType))
	return userType[randomIndex]
}

type ResponseType struct {
	Message string `json:"message"`
}

func helper[T interface{}](obj T, route string, token string) (string, error) {
	baseRoute := "http://localhost:9000"
	postBody, _ := json.Marshal(obj)
	body := bytes.NewBuffer(postBody)

	var baerer = "Bearer " + token

	req, err := http.NewRequest(http.MethodPost, baseRoute+route, body)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", baerer)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	b, e := io.ReadAll(resp.Body)
	if err != nil {
		return "", e
	}

	defer resp.Body.Close()
	return string(b), nil

}

func main() {

	faker := faker.New()
	msg := "already exists"
	for i := 0; i < 1000; i++ {
		userType := RandomUserType()
		randomid := randomId(200)

		token, err := genToken(randomid)
		if err != nil {
			panic(err.Error())
		}
		vCode := genCode(faker)
		dsCode := genCode(faker)

		user := CreateUser(faker, userType)

		userOut, err := helper(user, "/register", token)
		if err != nil {
			panic(err.Error())
		}

		isCUser := IsContains(userOut, msg)
		if isCUser {
			continue
		}
		vendor := CreateVendor(faker, vCode)
		venOut, err := helper(vendor, "/vendor/create-vendor", token)
		if err != nil {
			panic(err.Error())
		}

		isCVendor := IsContains(venOut, msg)
		if isCVendor {
			continue
		}

		buyer := CreateDealerAndBuyer(faker)
		buyerOut, err := helper(buyer, "/buyer/create-buyer", token)

		if err != nil {
			panic(err.Error())
		}
		isCBuyer := IsContains(buyerOut, msg)
		if isCBuyer {
			continue
		}

		dealer := CreateDealerAndBuyer(faker)
		dealerOut, err := helper(dealer, "/dealer/create-dealer", token)

		if err != nil {
			panic(err.Error())
		}

		isCDealer := IsContains(dealerOut, msg)
		if !isCDealer {
			dealerSite := CreateDealerSite(faker, randomid, vCode, dsCode)
			dealerSiteOut, err := helper(dealerSite, "/dealer-site/create-dealersite", token)
			if err != nil {
				panic(err.Error())
			}
			isCDealerSite := IsContains(dealerSiteOut, msg)

			userEntity := CreateUserEntity(faker, randomid)
			outUserEntity, err := helper(userEntity, "/relations/create-user-entity", token)
			if err != nil {
				panic(err.Error())
			}
			fmt.Println(dealerSiteOut)
			fmt.Println(outUserEntity)
			if !isCDealerSite {

				buyerSite := CreateBuyerSite(faker, randomid, vCode, dsCode)
				buyerSiteOut, err := helper(buyerSite, "/buyer-site/create-buyersite", token)

				if err != nil {
					panic(err.Error())
				}
				fmt.Println(buyerSiteOut)
				isCBuyerSite := IsContains(buyerSiteOut, msg)
				if !isCBuyerSite {
					vds := CreateVds(randomid, randomid)
					vdsOut, err := helper(vds, "/relations/vds-relations", token)
					if err != nil {
						panic(err.Error())
					}

					fmt.Println(vdsOut)

					isCVds := IsContains(vdsOut, msg)
					if !isCVds {
						vdsbs := CreateVdsbs(randomid, randomid)
						vdsbsOut, err := helper(vdsbs, "/relations/vdsbs-relations", token)
						if err != nil {
							panic(err.Error())
						}

						fmt.Println(vdsbsOut)
					}

				}
			}
		}

		fmt.Println(userOut)
		fmt.Println(venOut)
		fmt.Println(buyerOut)
		fmt.Println(dealerOut)
	}

}
