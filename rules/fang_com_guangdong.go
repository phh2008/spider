package rules

import (
	"fmt"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"strings"
)

func init() {
	areaMap["dg"] = [2]int{1, 13}
	areaMap["fs"] = [2]int{1, 16}
	areaMap["gz"] = [2]int{1, 30}
	areaMap["huizhou"] = [2]int{1, 16}
	areaMap["sz"] = [2]int{1, 20}
	areaMap["zh"] = [2]int{1, 9}
	areaMap["zs"] = [2]int{1, 12}
	areaMap["jm"] = [2]int{1, 8}
	areaMap["st"] = [2]int{1, 12}
	areaMap["qingyuan"] = [2]int{1, 9}
	areaMap["zhaoqing"] = [2]int{1, 11}
	areaMap["yangjiang"] = [2]int{1, 7}
	areaMap["maoming"] = [2]int{1, 4}
	areaMap["zj"] = [2]int{1, 9}
	areaMap["meizhou"] = [2]int{1, 6}
	areaMap["jieyang"] = [2]int{1, 5}
	areaMap["heyuan"] = [2]int{1, 6}
	areaMap["shanwei"] = [2]int{1, 3}
	areaMap["yangchun"] = [2]int{1, 2}
	areaMap["kaiping"] = [2]int{1, 3}
	areaMap["yunfu"] = [2]int{1, 6}
	areaMap["chaozhou"] = [2]int{1, 3}
	areaMap["shaoguan"] = [2]int{1, 6}
	areaMap["shunde"] = [2]int{1, 6}
	areaMap["taishan"] = [2]int{1, 3}
	areaMap["enping"] = [2]int{1, 2}
	areaMap["huidong"] = [2]int{1, 3}
	areaMap["gdlm"] = [2]int{1, 1}
	areaMap["boluo"] = [2]int{1, 3}
	areaMap["heshan"] = [2]int{1, 3}
	areaMap["puning"] = [2]int{1, 2}
	areaMap["xf"] = [2]int{1, 1}
	areaMap["leizhou"] = [2]int{1, 1}
	fangComGuangdong.Register()
}

var areaMap = map[string][2]int{}

var fangComGuangdong = &spider.Spider{
	Name:        "房天下-广东",
	Description: "https://dg.newhouse.fang.com/house/s/b91/?ctm=1.dg.xf_search.page.2",
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.Aid(map[string]interface{}{
				"loop": areaMap,
				"rule": "广东省",
			}, "广东省")
		},
		Trunk: map[string]*spider.Rule{
			"广东省": {
				AidFunc: func(ctx *spider.Context, aid map[string]interface{}) interface{} {
					guangdong := aid["loop"].(map[string][2]int)
					for key, val := range guangdong {
						for ; val[0] <= val[1]; val[0]++ {
							ctx.AddQueue(&request.Request{
								Url:  fmt.Sprintf("https://%s.newhouse.fang.com/house/s/b%d/?ctm=1.%s.xf_search.page.%d", key, 90+val[0], key, val[0]),
								Rule: aid["rule"].(string),
								Temp: map[string]interface{}{
									"city": key,
								},
							})
						}
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
									"city": ctx.GetTemp("city", ""),
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
					"城市",
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
						8: ctx.GetTemp("city", ""),
					})
				},
			},
		},
	},
}
