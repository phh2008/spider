package rules

import (
	"com.phh/spider/utils"
	"fmt"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"strings"
)

func init() {
	Mzlmcw.Register()
}

var Mzlmcw = &spider.Spider{
	Name:         "满洲里木材网",
	Description:  "http://www.mzlmcw.cn/company",
	EnableCookie: false,
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.Aid(map[string]interface{}{
				"loop": [2]int{1, 86},
				"rule": "公司条目",
			}, "公司条目")
		},
		Trunk: map[string]*spider.Rule{
			"公司条目": {
				AidFunc: func(ctx *spider.Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] <= loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  fmt.Sprintf("http://www.mzlmcw.cn/company/?page=%d", loop[0]),
							Rule: aid["rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					dom.Find(".left_box .list").Each(func(i int, div *goquery.Selection) {
						ctx.SetTemp("div", div)
						ctx.Parse("企业信息")
					})
				},
			},
			"企业信息": {
				ItemFields: []string{
					"企业",
					"主营",
					"省份",
					"城市",
				},
				ParseFunc: func(ctx *spider.Context) {
					div := ctx.GetTemp("div", &goquery.Selection{}).(*goquery.Selection)
					td := div.Find("tr").First().Find("td")
					li := td.Eq(2).Find("li")
					//公司
					corpName := li.Eq(0).Find("a").Text()
					//主营：运输、配送、仓储、包装、搬运装卸、流通制作
					major := li.Eq(1).Text()
					majorArr := strings.Split(major, "：")
					major = utils.If(len(majorArr) > 1, majorArr[1], "").(string)
					//[内蒙古/满洲里市]
					address := td.Eq(4).Text()
					address = strings.ReplaceAll(address, "[", "")
					address = strings.ReplaceAll(address, "]", "")
					address = strings.ReplaceAll(address, " ", "")
					addressArr := strings.Split(address, "/")
					province := utils.If(len(addressArr) > 0, addressArr[0], "").(string)
					city := utils.If(len(addressArr) > 1, addressArr[1], "").(string)

					ctx.Output(map[int]interface{}{
						0: corpName,
						1: major,
						2: province,
						3: city,
					})
				},
			},
		},
	},
}
