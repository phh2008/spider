package rules

import (
	"fmt"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"regexp"
	"strings"
)

func init() {
	jj718.Register()
}

var reg = regexp.MustCompile("\\[|\\]")

var jj718 = &spider.Spider{
	Name:        "家具材料网",
	Description: "http://www.jj718.com/company/search.php?page=1",
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.Aid(map[string]interface{}{
				"loop": [2]int{1, 347},
				"rule": "企业分页",
			}, "企业分页")
		},
		Trunk: map[string]*spider.Rule{
			"企业分页": {
				AidFunc: func(ctx *spider.Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] <= loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  fmt.Sprintf("http://www.jj718.com/company/search.php?page=%d", loop[0]),
							Rule: aid["rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					dom.Find(".list").Each(func(i int, selection *goquery.Selection) {
						td := selection.Find("table").Find("tr").First().Find("td")
						li := td.Eq(2).Find("ul li")
						a := li.First().Find("a")
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
							Rule: "企业信息",
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
			"企业信息": {
				ItemFields: []string{
					"名称",
					"主营",
					"省份",
					"市区",
					"联系人",
					"电话",
					"邮件",
					"传真",
				},
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					contact := ""
					tel := ""
					email := ""
					fax := ""
					dom.Find("#side .side_body").Eq(1).Find("ul li").Each(func(i int, selection *goquery.Selection) {
						text := strings.TrimSpace(selection.Text())
						//联系人：罗慧  电话：0578-2761888
						arr := strings.Split(text, "：")
						if arr == nil || len(arr) != 2 {
							return
						}
						switch arr[0] {
						case "联系人":
							contact = arr[1]
						case "电话":
							tel = arr[1]
						case "邮件":
							email = arr[1]
						case "传真":
							fax = arr[1]
						default:
							//
						}
					})
					ctx.Output(map[int]interface{}{
						0: ctx.GetTemp("name", ""),
						1: ctx.GetTemp("major", ""),
						2: ctx.GetTemp("province", ""),
						3: ctx.GetTemp("city", ""),
						4: contact,
						5: tel,
						6: email,
						7: fax,
					})
				},
			},
		},
	},
}
