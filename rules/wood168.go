package rules

import (
	"com.phh/spider/utils"
	"encoding/json"
	"fmt"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"github.com/henrylee2cn/pholcus/logs"
	"strings"
)

func init() {
	Wood168.Register()
}

var Wood168 = &spider.Spider{
	Name:         "中国木业信息网",
	Description:  "http://www.wood168.net",
	EnableCookie: true,
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.Aid(map[string]interface{}{
				"loop": [2]int{1, 5385},
				"rule": "企业大全",
			}, "企业大全")
		},
		Trunk: map[string]*spider.Rule{
			"企业大全": {
				AidFunc: func(ctx *spider.Context, aid map[string]interface{}) interface{} {
					//userAgent := "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.25 Safari/537.36 Core/1.70.3676.400 QQBrowser/10.4.3505.400"
					for loop := aid["loop"].([2]int); loop[0] <= loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  fmt.Sprintf("http://www.wood168.net/com_find.asp?nextpage=%d", loop[0]),
							Rule: aid["rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					dom.Find("#sidebar_right table").Eq(4).Find("tr").First().Find("font[size='3']").Each(func(i int, font *goquery.Selection) {
						href, ok := font.Find("a").Attr("href")
						if ok {
							ctx.AddQueue(&request.Request{
								Url:  "http://www.wood168.net/" + href,
								Rule: "企业主页",
							})
						} else {
							logs.Log.Error("未找到企业主页url")
						}
					})
				},
			},
			"企业主页": {
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					table := dom.Find("body").ChildrenFiltered("table")
					b := table.Eq(2).Find("tr").Eq(1).Find("td").Eq(-2).Find("b")
					corName := b.Eq(0).Find("font").Text()
					//主要产品：生产国产、进口三聚氢氨饰面板
					major := b.Eq(1).Find("font").Text()
					if major != "" {
						majorArr := strings.Split(major, "：")
						major = utils.If(len(majorArr) > 1, majorArr[1], major).(string)
					}
					link, ok := table.Eq(3).Find("tr").First().Find("td").Eq(-2).Find("a").Attr("href")
					if ok {
						link = strings.ReplaceAll(link, "..", "")
						ctx.AddQueue(&request.Request{
							Url:  "http://www.wood168.net" + link,
							Rule: "企业信息",
							Temp: map[string]interface{}{
								"major":   major,
								"corName": corName,
							},
						})
					} else {
						logs.Log.Error("未找到联系url")
					}
				},
			},
			"企业信息": {
				ItemFields: []string{
					"企业名称",
					"主营",
					"联系人",
					"性别",
					"地址",
					"电话",
					"传真",
					"手机",
					"邮箱",
				},
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					table := dom.Find("table[width='400']")
					tr := table.Find("tr")
					//于秀利  先生 (市场部经理 )   苏军  先生 (鲁南木业机械市场 经理)  丁健卫  先生 ( )
					contact := tr.Eq(1).Find("td").First().Find("div").Text()
					logs.Log.Informational(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>[%s]", contact)
					contact = strings.ReplaceAll(contact, "(", "")
					contact = strings.ReplaceAll(contact, " )", "")
					contact = strings.ReplaceAll(contact, ")", "")
					contactArr := strings.Split(contact, " ")
					contact = utils.If(len(contactArr) > 0, contactArr[0], "").(string)
					sex := utils.If(len(contactArr) > 1, contactArr[1], "").(string)

					address := tr.Eq(2).Find("td").First().Text()
					addressArr := strings.Split(address, "：")
					address = utils.If(len(addressArr) > 1, addressArr[1], "").(string)

					tel := tr.Eq(4).Find("td").First().Text()
					telArr := strings.Split(tel, "：")
					tel = utils.If(len(telArr) > 1, telArr[1], "").(string)

					fax := tr.Eq(5).Find("td").First().Text()
					faxArr := strings.Split(fax, "：")
					fax = utils.If(len(faxArr) > 1, faxArr[1], "").(string)

					mobile := tr.Eq(6).Find("td").First().Text()
					mobileArr := strings.Split(mobile, "：")
					mobile = utils.If(len(mobileArr) > 1, mobileArr[1], "").(string)

					email := tr.Eq(7).Find("td").First().Text()
					emailArr := strings.Split(email, "：")
					email = utils.If(len(emailArr) > 1, emailArr[1], "").(string)

					info := map[int]interface{}{
						0: ctx.GetTemp("corName", ""),
						1: ctx.GetTemp("major", ""),
						2: contact,
						3: sex,
						4: address,
						5: tel,
						6: fax,
						7: mobile,
						8: email,
					}
					json, _ := json.Marshal(info)
					logs.Log.Informational("企业信息>>>>>:%s", string(json))
					ctx.Output(info)
				},
			},
		},
	},
}
