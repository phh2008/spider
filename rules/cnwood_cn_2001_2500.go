package rules

import (
	"fmt"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"strings"
)

func init() {
	cnwood_2001_2500.Register()
}

var cnwood_2001_2500 = &spider.Spider{
	Name:        "中国木业网2001-2500",
	Description: "http://www.cnwood.cn/company/search.php?page=1",
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.Aid(map[string]interface{}{
				"loop": [2]int{2001, 2500},
				"rule": "企业分页",
			}, "企业分页")
		},
		Trunk: map[string]*spider.Rule{
			"企业分页": {
				AidFunc: func(ctx *spider.Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] <= loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  fmt.Sprintf("http://www.cnwood.cn/company/search.php?page=%d", loop[0]),
							Rule: aid["rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					dom.Find(".m .list").Each(func(i int, selection *goquery.Selection) {
						td := selection.Find("table tr").First().Find("td")
						li := td.Eq(2).Find("ul li")
						a := li.Eq(0).Find("a")
						url, _ := a.Attr("href")
						name := a.Find("strong").Text()
						major := li.Eq(1).Text()
						major = strings.ReplaceAll(major, "主营：", "")
						//e.g. [河北/邢台市]
						region := td.Eq(4).Text()
						region = reg.ReplaceAllString(region, "")
						arr := strings.Split(region, "/")
						province := ""
						if len(arr) > 0 {
							province = arr[0]
						}
						city := ""
						if len(arr) > 1 {
							city = arr[1]
						}
						ctx.AddQueue(&request.Request{
							Url:  url,
							Rule: "企业主页",
							Temp: map[string]interface{}{
								"name":     name,
								"major":    major,
								"province": province,
								"city":     city,
							},
						})
					})
				},
			},
			"企业主页": {
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					contact := ""
					tel := ""
					address := ""
					mobile := ""
					net := ""
					email := ""
					dl := dom.Find(".contact_way dl")
					if dl.Size() > 0 {
						dd := dl.Find("dd")
						dt := dl.Find("dt")
						dt.Each(func(i int, selection *goquery.Selection) {
							text := selection.Text()
							val := dd.Eq(i).Text()
							switch text {
							case "联系人":
								contact = val
							case "电话":
								tel = val
							case "地址":
								address = val
							case "手机":
								mobile = val
							case "网址":
								net = val
							case "Email", "邮件":
								email = val
							default:
								//
							}
						})
					} else {
						dom.Find(".m .side_body").Eq(2).Find("ul li").Each(func(i int, selection *goquery.Selection) {
							//联系人：周洋  邮件：zhouyang1000@21cn.com   等等
							text := strings.TrimSpace(selection.Text())
							arr := strings.Split(text, "：")
							if arr == nil || len(arr) != 2 {
								return
							}
							switch arr[0] {
							case "联系人":
								contact = arr[1]
							case "电话":
								tel = arr[1]
							case "地址":
								address = arr[1]
							case "手机":
								mobile = arr[1]
							case "网址":
								net = arr[1]
							case "Email", "邮件":
								email = arr[1]
							default:
								//
							}
						})
					}
					ctx.Output(map[int]interface{}{
						0: ctx.GetTemp("name", ""),
						1: ctx.GetTemp("major", ""),
						2: ctx.GetTemp("province", ""),
						3: ctx.GetTemp("city", ""),
						4: contact,
						5: tel,
						6: mobile,
						7: email,
					})
				},
				ItemFields: []string{
					"名称",
					"主营",
					"省份",
					"市区",
					"联系人",
					"电话",
					"手机",
					"邮件",
				},
			},
		},
	},
}
