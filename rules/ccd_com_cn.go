package rules

import (
	"fmt"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"strings"
)

func init() {
	ccd_com_cn.Register()
}

var ccd_com_cn = &spider.Spider{
	Name:        "中国建筑装饰网",
	Description: "http://jiancai.ccd.com.cn/company/index.aspx?Page=1",
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.Aid(map[string]interface{}{
				"loop": [2]int{1, 2612}, //totalPage:2612
				"rule": "企业分页",
			}, "企业分页")
		},
		Trunk: map[string]*spider.Rule{
			"企业分页": {
				AidFunc: func(ctx *spider.Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] <= loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  fmt.Sprintf("http://jiancai.ccd.com.cn/company/index.aspx?Page=%d", loop[0]),
							Rule: aid["rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					dom.Find(".sub2_con2 table").Each(func(i int, selection *goquery.Selection) {
						tr := selection.Find("tbody tr")
						a := tr.Eq(0).Find("td").First().Find("a")
						url, _ := a.Attr("href")
						name := strings.TrimSpace(a.Text())

						major := strings.TrimSpace(tr.Eq(1).Find("td").Eq(1).Find("div").Text())
						major = strings.ReplaceAll(major, "经营范围:", "")

						address := strings.TrimSpace(tr.Eq(2).Find("td").First().Text())
						address = strings.ReplaceAll(address, "地址：", "")

						ctx.AddQueue(&request.Request{
							Url:  url,
							Rule: "企业主页",
							Temp: map[string]interface{}{
								"name":    name,
								"major":   major,
								"address": address,
							},
						})
					})
				},
			},
			"企业主页": {
				ParseFunc: func(ctx *spider.Context) {
					var contact, tel, mobile, email string
					ctx.GetDom().Find(".dabox p").Each(func(i int, selection *goquery.Selection) {
						//e.g 联 系 人：邹先生
						text := strings.TrimSpace(selection.Text())
						arr := strings.Split(text, "：")
						if len(arr) < 2 {
							return
						}
						tmp := arr[1]
						switch arr[0] {
						case "联 系 人", "联系人":
							contact = tmp
						case "座机":
							tel = tmp
						case "手机号":
							mobile = tmp
						case "电子邮箱":
							email = tmp
						}
					})
					ctx.Output(map[int]interface{}{
						0: ctx.GetTemp("name", ""),
						1: ctx.GetTemp("major", ""),
						2: contact,
						3: tel,
						4: mobile,
						5: email,
						6: ctx.GetTemp("address", ""),
					})
				},
				ItemFields: []string{
					"名称",
					"主营",
					"联系人",
					"电话",
					"手机",
					"邮箱",
					"地址",
				},
			},
		},
	},
}
