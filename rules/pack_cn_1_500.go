package rules

import (
	"fmt"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"strings"
)

func init() {
	packCn_1_500.Register()
}

var packCn_1_500 = &spider.Spider{
	Name:        "中国包装网1-500",
	Description: "http://www.pack.cn/company/index.php?page=1",
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.Aid(map[string]interface{}{
				"loop": [2]int{1, 500},
				"rule": "企业分页",
			}, "企业分页")
		},
		Trunk: map[string]*spider.Rule{
			"企业分页": {
				AidFunc: func(ctx *spider.Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] <= loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  fmt.Sprintf("http://www.pack.cn/company/index.php?page=%d", loop[0]),
							Rule: aid["rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					dom.Find(".box_body .col").Each(func(i int, selection *goquery.Selection) {
						dd := selection.Find(".contbox dl").First().Find("dd")
						a := dd.Eq(0).Find("a")
						url, _ := a.Attr("href")
						name := strings.TrimSpace(a.Text())
						major := strings.TrimSpace(dd.Eq(1).Text())
						major = strings.ReplaceAll(major, "主营产品：", "")
						address := strings.TrimSpace(dd.Eq(2).Text())
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
					dom := ctx.GetDom()
					var contact, tel, mobile, fax, email string
					dom.Find(".m .side_body").Each(func(i int, div *goquery.Selection) {
						div.Find("ul li").Each(func(i int, selection *goquery.Selection) {
							//联系人：周洋  邮件：zhouyang1000@21cn.com   等等
							text := strings.TrimSpace(selection.Text())
							arr := strings.Split(text, "：")
							if arr == nil || len(arr) < 2 {
								return
							}
							tmp := strings.TrimSpace(arr[1])
							switch arr[0] {
							case "联系人":
								contact = tmp
							case "电话":
								tel = tmp
							case "手机":
								mobile = tmp
							case "传真":
								fax = tmp
							case "邮件", "Email":
								email = tmp
							default:
								//
							}
						})
					})
					ctx.Output(map[int]interface{}{
						0: ctx.GetTemp("name", ""),
						1: ctx.GetTemp("major", ""),
						2: ctx.GetTemp("address", ""),
						3: contact,
						4: tel,
						5: mobile,
						6: email,
						7: fax,
					})
				},
				ItemFields: []string{
					"名称",
					"主营",
					"地址",
					"联系人",
					"电话",
					"手机",
					"邮件",
					"传真",
				},
			},
		},
	},
}
