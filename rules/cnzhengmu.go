package rules

import (
	"fmt"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"strings"
)

func init() {
	cnZhengMu.Register()
}

var cnZhengMu = &spider.Spider{
	Name:        "整木网",
	Description: "https://www.cnzhengmu.com/company/list/index?page=1",
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.Aid(map[string]interface{}{
				"loop": [2]int{1, 6},
				"rule": "企业分页",
			}, "企业分页")
		},
		Trunk: map[string]*spider.Rule{
			"企业分页": {
				AidFunc: func(ctx *spider.Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] <= loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  fmt.Sprintf("https://www.cnzhengmu.com/company/list/index?page=%d", loop[0]),
							Rule: aid["rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					dom.Find(".media-body").Each(func(i int, selection *goquery.Selection) {
						a := selection.Find("h5 a")
						url, _ := a.Attr("href")
						name := a.Text()
						ctx.AddQueue(&request.Request{
							Url:  url,
							Rule: "企业主页",
							Temp: map[string]interface{}{
								"name": name,
							},
						})
					})
				},
			},
			"企业主页": {
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					li := dom.Find(".navbar-nav").First().Find("li").Last()
					url, _ := li.Find("a").Attr("href")
					introduction := strings.TrimSpace(dom.Find("main .content").Text())
					temp := ctx.GetTemps()
					temp["introduction"] = introduction
					ctx.AddQueue(&request.Request{
						Url:  url,
						Rule: "企业信息",
						Temp: temp,
					})
				},
			},
			"企业信息": {
				ItemFields: []string{
					"名称",
					"简介",
					"联系人",
					"地址",
					"电话",
					"邮箱",
					"网址",
				},
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					rows := dom.Find(".border- .row")
					contact := strings.TrimSpace(rows.Eq(0).Find(".col-10").Text())
					address := strings.TrimSpace(rows.Eq(1).Find(".col-10").Text())
					mobile := strings.TrimSpace(rows.Eq(2).Find(".col-10").Text())
					email := strings.TrimSpace(rows.Eq(3).Find(".col-10").Text())
					net := strings.TrimSpace(rows.Eq(4).Find(".col-10").Text())

					ctx.Output(map[int]interface{}{
						0: ctx.GetTemp("name", ""),
						1: ctx.GetTemp("introduction", ""),
						2: contact,
						3: address,
						4: mobile,
						5: email,
						6: net,
					})
				},
			},
		},
	},
}
