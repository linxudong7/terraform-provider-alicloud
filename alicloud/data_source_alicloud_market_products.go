package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/market"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudProducts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudProductsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"sort": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"user_count-desc", "created_on-desc", "price-desc", "score-desc"}, false),
			},
			"category_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"product_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"APP", "SERVICE", "MIRROR", "DOWNLOAD", "API_SERVICE"}, false),
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values.
			"products": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudProductsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := market.CreateDescribeProductsRequest()
	request.RegionId = client.RegionId
	var productsFilter []market.DescribeProductsFilter
	var product market.DescribeProductsFilter
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[vv.(string)] = vv.(string)
		}
	}
	if v, ok := d.GetOk("sort"); ok && v.(string) != "" {
		product.Key = "sort"
		product.Value = v.(string)
		productsFilter = append(productsFilter, product)
	}
	if v, ok := d.GetOk("category_id"); ok && v.(string) != "" {
		product.Key = "categoryId"
		product.Value = v.(string)
		productsFilter = append(productsFilter, product)
	}
	if v, ok := d.GetOk("product_type"); ok && v.(string) != "" {
		product.Key = "productType"
		product.Value = v.(string)
		productsFilter = append(productsFilter, product)
	}
	request.Filter = &productsFilter
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var allProduct []market.ProductItem
	for {
		raw, err := client.WithMarketClient(func(marketClient *market.Client) (interface{}, error) {
			return marketClient.DescribeProducts(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_market_products", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*market.DescribeProductsResponse)

		if len(response.ProductItems.ProductItem) < 1 {
			break
		}

		for _, item := range response.ProductItems.ProductItem {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.Name) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.Code]; !ok {
					continue
				}
			}
			allProduct = append(allProduct, item)
		}

		if len(response.ProductItems.ProductItem) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	return productsDescriptionAttributes(d, allProduct)
}

func productsDescriptionAttributes(d *schema.ResourceData, allProduct []market.ProductItem) error {
	var ids []string
	var s []map[string]interface{}
	for _, p := range allProduct {
		mapping := map[string]interface{}{
			"code":       p.Code,
			"name":       p.Name,
			"target_url": p.TargetUrl,
		}

		ids = append(ids, p.Code)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("products", s); err != nil {
		return err
	}
	if err := d.Set("ids", ids); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}