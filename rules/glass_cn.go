package rules

import (
	"com.phh/spider/utils"
	"fmt"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"regexp"
	"strings"
)

func init() {
	glassCn.Register()
}

var glassCn = &spider.Spider{
	Name:        "全球玻璃网",
	Description: "https://www.glass.cn/gCompany/cqtppicapg1ky.html",
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.Aid(map[string]interface{}{
				"loop": [2]int{1, 1349},
				"rule": "企业分页",
			}, "企业分页")
		},
		Trunk: map[string]*spider.Rule{
			"企业分页": {
				AidFunc: func(ctx *spider.Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] <= loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  fmt.Sprintf("https://www.glass.cn/gCompany/cqtppicapg%dky.html", loop[0]),
							Rule: aid["rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *spider.Context) {
					ctx.GetDom().Find(".selllist .itemBox").Each(func(i int, selection *goquery.Selection) {
						cbox := selection.Find(".box3")
						ca := cbox.Find("h2 a")
						url, _ := ca.Attr("href")
						name := ca.Text()
						major := cbox.Find(".description").Clone().RemoveFiltered("strong").End().Text()
						major = reg_glass_cn.ReplaceAllString(major, "")
						//e.g.  重庆-重庆
						region := strings.TrimSpace(selection.Find(".box4").Text())
						arr := strings.Split(region, "-")
						var province, city string
						if len(arr) > 0 {
							province = arr[0]
						}
						if len(arr) > 1 {
							city = arr[1]
						}
						ctx.AddQueue(&request.Request{
							Url:  url + "contact.html",
							Rule: "公司信息",
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
			"公司信息": {
				ParseFunc: func(ctx *spider.Context) {
					box := ctx.GetDom().Find(".contactbox")
					var contact, tel, mobile, fax, address string
					if box.Size() > 0 {
						box.Find("li").Each(func(i int, selection *goquery.Selection) {
							text := strings.TrimSpace(selection.Clone().RemoveFiltered("a").End().Text())
							arr := strings.Split(text, "：")
							switch i {
							case 1:
								contact = arr[1]
							case 2:
								tel = arr[1]
							case 3:
								mobile = arr[1]
							case 4:
								fax = arr[1]
							case 5:
								address = arr[1]
							default:
								//
							}
						})
					} else {
						ctx.GetDom().Find(".contact p").Each(func(i int, selection *goquery.Selection) {
							text := strings.TrimSpace(selection.Text())
							arr := strings.Split(text, "：")
							label := utils.RegBlank.ReplaceAllString(arr[0], "")
							label = utils.Nbsp.ReplaceAllString(label, "")
							switch label {
							case "联系人":
								contact = arr[1]
							case "电话":
								tel = arr[1]
							case "手机", "移动电话":
								mobile = arr[1]
							case "传真":
								fax = arr[1]
							case "地址":
								address = arr[1]
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
						7: fax,
						8: address,
					})
				},
				ItemFields: []string{
					"名称",
					"主营",
					"省份",
					"城市",
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

var reg_glass_cn = regexp.MustCompile("：|\\s+|\u3002")
