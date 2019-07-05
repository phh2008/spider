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
	zk71_com.Register()
}

var zk71_com = &spider.Spider{
	Name:        "中科商务网",
	Description: "http://www.zk71.com/company/13_0_0_0_0_0/page1/",
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.Aid(map[string]interface{}{
				"loop": [2]int{1, 115},
				"rule": "企业分页",
			}, "企业分页")
		},
		Trunk: map[string]*spider.Rule{
			"企业分页": {
				AidFunc: func(ctx *spider.Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] <= loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  fmt.Sprintf("http://www.zk71.com/company/13_0_0_0_0_0/page%d/", loop[0]),
							Rule: aid["rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *spider.Context) {
					ctx.GetDom().Find(".company-list li").Each(func(i int, selection *goquery.Selection) {
						p := selection.Find(".companyinfo p")
						ac := p.Filter(".coName").Find("a")
						url, _ := ac.Attr("href")
						name := strings.TrimSpace(ac.Text())
						major := strings.TrimSpace(p.Filter(".keywords").Text())
						major = strings.ReplaceAll(major, "主营产品：", "")
						ctx.AddQueue(&request.Request{
							Url:  url,
							Rule: "企业主页",
							Temp: map[string]interface{}{
								"name":  name,
								"major": major,
							},
						})
					})
				},
			},
			"企业主页": {
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					text := dom.Find(".lxwm").Text()
					var contact, tel, mobile, fax, address string
					text = regBlank2.ReplaceAllString(text, "|")
					arr := strings.Split(text, "|")
					for _, v := range arr {
						lineArr := strings.Split(v, "：")
						if len(lineArr) < 2 {
							continue
						}
						switch lineArr[0] {
						case "联系人":
							contact = lineArr[1]
						case "电话":
							tel = lineArr[1]
						case "手机":
							mobile = lineArr[1]
						case "传真":
							fax = lineArr[1]
						case "地址":
							address = lineArr[1]
						}
					}
					ctx.Output(map[int]interface{}{
						0: ctx.GetTemp("name", ""),
						1: ctx.GetTemp("major", ""),
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

var regBlank2 = regexp.MustCompile("\\s{2,}")
