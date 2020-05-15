package commands

import (
	"fmt"
	"hwc-tool/commands/utils"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/auth/token"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/sirupsen/logrus"
)

const (
	CLOUD_PARAM_REGION            = "ap-southeast-3"
	CLOUD_PARAM_CLOUD_DOAMIN_NAME = "myhuaweicloud.com"
	IDENTITY_USER_NAME            = "hw*******"
	IDENTITY_PASSWORD             = "{cipher_a}f462bc9a0f4fd7****************"
	IDENTITY_DOMAIN_ID            = "07b045206c00252***************"
)

var cloudParamProjectIdMap = map[string]string{
	"cn-north-4":     "07e2cc042e00263b2f23c00f7cc111bb",
	"cn-east-3":      "07e2d5c7d40025232f3ac00f6ea2cd34",
	"cn-east-2":      "07b04520760025232fa1c00fbbc7e01b",
	"cn-south-1":     "07b04b0a66000f092f6ec00f79a087c6",
	"cn-southwest-2": "084908187580f2602f7bc00f25129456",
	"ap-southeast-1": "08486ce9f28025942fe6c00fddf84459",
	"ap-southeast-2": "0849082a7d80f39f2fc5c00f058b7721",
	"ap-southeast-3": "08485cf5ad80f2622f6dc00f79fb4ae0",
	"af-south-1":     "084908461180258d2fbac00fa4f81567",
	"la-south-2":     "08490861ac80f39f2fc6c00fb0dc2cc6",
}

func createGopherCloudIMSV2ServiceClient(region string) (*gophercloud.ServiceClient, error) {
	if region == "" {
		region = CLOUD_PARAM_REGION
	}
	identityURL := "https://iam." + region + "." + CLOUD_PARAM_CLOUD_DOAMIN_NAME + "/v3"
	pass, err := utils.AesDePassword("guid", "seed", IDENTITY_PASSWORD)
	if err != nil {
		fmt.Println("Failed to decrypt: ", err)
		logrus.Infof("Failed to decrypt: %v", err)
		return nil, err
	}
	tokenOpts := token.TokenOptions{
		IdentityEndpoint: identityURL,
		Username:         IDENTITY_USER_NAME,
		Password:         pass,
		DomainID:         IDENTITY_DOMAIN_ID,
		ProjectID:        cloudParamProjectIdMap[region],
	}

	//初始化provider client
	provider, authErr := openstack.AuthenticatedClient(tokenOpts)
	if authErr != nil {
		fmt.Println("Failed to get the AuthenticatedClient: ", authErr)
		logrus.Infof("Failed to get the AuthenticatedClient: %v", authErr)
		return nil, authErr
	}
	//初始化service client
	sc, clientErr := openstack.NewIMSV2(provider, gophercloud.EndpointOpts{})

	if clientErr != nil {
		fmt.Println("Failed to get the NewIMSV2 client: ", clientErr)
		logrus.Infof("Failed to get the NewIMSV2 client: %v", clientErr)
		return nil, clientErr
	}
	return sc, nil
}
