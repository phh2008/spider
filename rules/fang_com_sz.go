package rules

import (
	"fmt"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"strings"
)

func init() {
	fangComSz.Register()
}

var fangComSz = &spider.Spider{
	Name:        "房天下-深圳",
	Description: "https://sz.newhouse.fang.com/house/s/b91/?ctm=1.sz.xf_search.page.2",
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.Aid(map[string]interface{}{
				"loop": [2]int{1, 20},
				"rule": "分页",
			}, "分页")
		},
		Trunk: map[string]*spider.Rule{
			"分页": {
				AidFunc: func(ctx *spider.Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] <= loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  fmt.Sprintf("https://sz.newhouse.fang.com/house/s/b%d/?ctm=1.sz.xf_search.page.%d", 90+loop[0], loop[0]),
							Rule: aid["rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					dom.Find("#newhouse_loupai_list").Find("li").Each(func(i int, li *goquery.Selection) {
						a := li.Find(".nlcd_name a")
						url, ok := a.Attr("href")
						name := strings.TrimSpace(a.Text())
						area := strings.TrimSpace(li.Find(".address .sngrey").Text())
						if ok {
							if strings.HasPrefix(url, "//") {
								url = "https:" + url
							}
							ctx.AddQueue(&request.Request{
								Url:  url,
								Rule: "楼盘首页",
								Temp: map[string]interface{}{
									"name": name,
									"area": area,
								},
							})
						}
					})
				},
			},
			"楼盘首页": {
				ParseFunc: func(ctx *spider.Context) {
					ctx.GetDom().Find("#orginalNaviBox").Find("a").Each(func(i int, a *goquery.Selection) {
						if strings.TrimSpace(a.Text()) == "楼盘详情" {
							url, _ := a.Attr("href")
							if strings.HasPrefix(url, "//") {
								url = "https:" + url
							}
							ctx.AddQueue(&request.Request{
								Url:  url,
								Rule: "楼盘详情",
								Temp: ctx.GetTemps(),
							})
						}
					})
				},
			},
			"楼盘详情": {
				ItemFields: []string{
					"楼盘",
					"价格",
					"户型",
					"开盘时间",
					"交房时间",
					"装修状态",
					"地区",
					"地址",
				},
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					price := strings.TrimSpace(dom.Find(".main-info-price em").Text())
					decorate := ""
					address := ""
					dom.Find(".main-info").SiblingsFiltered("ul").Find("li").Each(func(i int, li *goquery.Selection) {
						label := li.Find(".list-left").Text()
						if label == "装修状况：" {
							decorate = strings.TrimSpace(li.Find(".list-right").Text())
						} else if label == "楼盘地址：" {
							address = strings.TrimSpace(li.Find(".list-right-text").Text())
						}
					})
					openTime := ""
					deliveryTime := ""
					houseType := ""
					dom.Find(".main-item").Each(func(i int, div *goquery.Selection) {
						if div.Find("h3").Text() == "销售信息" {
							div.Find("ul li").Each(func(i int, li *goquery.Selection) {
								label := li.Find(".list-left").Text()
								if label == "开盘时间：" {
									atext := li.Find(".list-right").Find("a").Text()
									val := li.Find(".list-right").Text()
									if atext != "" {
										val = strings.ReplaceAll(val, atext, "")
									}
									openTime = strings.TrimSpace(val)
								} else if label == "交房时间：" {
									deliveryTime = strings.TrimSpace(li.Find(".list-right").Text())
								} else if label == "主力户型：" {
									houseType = strings.TrimSpace(li.Find(".list-right-text").Text())
								}
							})
						}
					})
					ctx.Output(map[int]interface{}{
						0: ctx.GetTemp("name", ""),
						1: price,
						2: houseType,
						3: openTime,
						4: deliveryTime,
						5: decorate,
						6: ctx.GetTemp("area", ""),
						7: address,
					})
				},
			},
		},
	},
}
