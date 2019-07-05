package rules

import (
	"fmt"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"strings"
)

func init() {
	huangye88_com.Register()
}

var huangye88_com = &spider.Spider{
	Name:        "黄页88网",
	Description: "http://b2b.huangye88.com/guangdong/jiajuwang/pn50/",
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.Aid(map[string]interface{}{
				"loop": [2]int{1, 50},
				"rule": "企业分页",
			}, "企业分页")
		},
		Trunk: map[string]*spider.Rule{
			"企业分页": {
				AidFunc: func(ctx *spider.Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] <= loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  fmt.Sprintf("http://b2b.huangye88.com/guangdong/jiajuwang/pn%d/", loop[0]),
							Rule: aid["rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *spider.Context) {
					ctx.GetDom().Find("#jubao dl").Each(func(i int, selection *goquery.Selection) {
						dt := selection.Find("dt")
						name := strings.TrimSpace(dt.Find("h4 a").Text())
						url, ok := dt.Find("span a").Attr("href")
						if !ok {
							return
						}
						ctx.AddQueue(&request.Request{
							Url:  url,
							Rule: "企业联系",
							Temp: map[string]interface{}{
								"name": name,
							},
						})
					})
				},
			},
			"企业联系": {
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					major := dom.Find(".headercont .title .small").Text()
					var contact, mobile, tel, corName, fax, address string
					dom.Find(".con-txt li").Each(func(i int, selection *goquery.Selection) {
						text := strings.TrimSpace(selection.Text())
						//e.g 联系人：xxx
						arr := strings.Split(text, "：")
						if len(arr) < 2 {
							return
						}
						switch arr[0] {
						case "联系人":
							contact = arr[1]
						case "手机":
							mobile = arr[1]
						case "公司名称":
							corName = arr[1]
						case "电话":
							tel = arr[1]
						case "传真":
							fax = arr[1]
						case "地址":
							address = arr[1]
						}
					})
					ctx.Output(map[int]interface{}{
						0: ctx.GetTemp("name", corName),
						1: major,
						2: contact,
						3: tel,
						4: mobile,
						5: fax,
						6: address,
					})
				},
				ItemFields: []string{
					"名称",
					"主营",
					"联系人",
					"电话",
					"手机",
					"传真",
					"地址",
				},
			},
		},
	},
}
