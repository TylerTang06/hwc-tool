package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/ims/v2/cloudimages"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type Filter struct {
	Key   string
	Value string
}

var (
	GetImagesListFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "Get image of the id.",
			Value: "",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "Get image of the name.",
			Value: "",
		},
		cli.StringFlag{
			Name:  "ostype",
			Usage: "Get images list of the ostype.",
			Value: "Linux",
		},
		cli.StringFlag{
			Name:  "platform",
			Usage: "Get images list of the platform.",
			Value: "CentOS",
		},
		cli.StringFlag{
			Name:  "region",
			Usage: "Get images list of the region.",
			Value: "ap-southeast-3",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "Get images list, only limit to display.",
			Value: 500,
		},
	}
)

func cmdGetImageListAction(c *cli.Context) error {
	logrus.Infof("c.Args()=%++v", c.Args())
	filters := []Filter{
		Filter{
			Key:   "ostype",
			Value: c.String("ostype"),
		},
		Filter{
			Key:   "platform",
			Value: c.String("platform"),
		},
		Filter{
			Key:   "limit",
			Value: strconv.Itoa(c.Int("limit")),
		},
	}

	logrus.Infof("id=%v", c.String("id"))
	if c.String("id") != "" {
		filters = append(filters, Filter{
			Key:   "id",
			Value: c.String("id"),
		})
	}
	if c.String("name") != "" {
		filters = append(filters, Filter{
			Key:   "id",
			Value: c.String("name"),
		})
	}

	sc, err := createGopherCloudIMSV2ServiceClient(c.String("region"))
	if err != nil {
		logrus.Errorf("create the IMSV2 service client meet error=%v", err)
		return err
	}

	logrus.Infof("fileters=%++v", filters)
	if err = GetImageList(sc, filters); err != nil {
		logrus.Errorf("get images list meet error=%v", err)
		return err
	}

	return nil
}

func buildImagesListOptStruct(filters []Filter) cloudimages.ListOpts {
	listOpts := cloudimages.ListOpts{
		MinDisk:    40,
		Status:     "active",
		Imagetype:  "gold",
		Visibility: "public",
	}
	for _, filter := range filters {
		key := strings.ToLower(filter.Key)
		if key == "id" {
			listOpts.ID = filter.Value
		}
		if key == "ostype" {
			listOpts.OsType = filter.Value
		}
		if key == "platform" {
			listOpts.Platform = filter.Value
		}
		if key == "name" {
			listOpts.Name = filter.Value
		}
		if key == "limit" {
			value, err := strconv.Atoi(filter.Value)
			if err != nil {
				continue
			}
			listOpts.Limit = value
		}
	}
	return listOpts
}

func GetImageList(sc *gophercloud.ServiceClient, filters []Filter) error {
	listOpts := buildImagesListOptStruct(filters)
	logrus.Infof("listOpts=%++v", listOpts)
	pager := cloudimages.List(sc, listOpts)
	count, pages := 0, 0
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		pages++
		images, err := cloudimages.ExtractImages(page)
		if err != nil {
			logrus.Errorf("extract images meet error=%v", err)
			return false, err
		}

		fmt.Println("ID    NAME    PLATFORM    OSTYPE    OSVERSION    STATUS")
		for _, image := range images {
			fmt.Printf("%v   %v   %v   %v   %v   %v\n", image.ID, image.Name, image.Platform, image.OsType, image.OsVersion, image.Status)
			logrus.Infof("%v   %v   %v   %v   %v   %v", image.ID, image.Name, image.Platform, image.OsType, image.OsVersion, image.Status)
			count++
		}

		return true, nil
	})
	if err != nil {
		logrus.Errorf("iterate each page meet error=%v", err)
		return err
	}

	return nil
}
