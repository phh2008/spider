package rules

import (
	"fmt"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"github.com/henrylee2cn/pholcus/logs"
	"strings"
)

func init() {
	cailiao_com.Register()
}

var cailiao_com = &spider.Spider{
	Name:        "材料网",
	Description: "http://www.cailiao.com/frontend/Search/index?keyword=&search_type=company&page=1",
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.Aid(map[string]interface{}{
				"loop": [2]int{1, 1},
				"rule": "企业分页",
			}, "企业分页")
		},
		Trunk: map[string]*spider.Rule{
			"企业分页": {
				AidFunc: func(ctx *spider.Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] <= loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  fmt.Sprintf("http://www.cailiao.com/frontend/Search/index?keyword=&search_type=company&page=%d", loop[0]),
							Rule: aid["rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *spider.Context) {
					ctx.GetDom().Find(".product_main .product_list li").Each(func(i int, selection *goquery.Selection) {
						a := selection.Find(".detailed a")
						url, _ := a.Attr("href")
						name := strings.TrimSpace(a.Text())
						//e.g 广东,东莞,东莞周边
						region := strings.TrimSpace(selection.Find(".address-new .fl c66").Text())
						var province, city, street string
						arr := strings.Split(region, ",")
						if len(arr) > 0 {
							province = arr[0]
						}
						if len(arr) > 1 {
							city = arr[1]
						}
						if len(arr) > 2 {
							street = arr[2]
						}
						ctx.AddQueue(&request.Request{
							Url:  url,
							Rule: "企业主页",
							Temp: map[string]interface{}{
								"name":     name,
								"province": province,
								"city":     city,
								"street":   street,
							},
						})
					})
				},
			},
			"企业主页": {
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					mem := dom.Find(".pt_member li")
					if mem.Size() > 0 {
						var contact, tel, mobile, major, address, email string
						mem.Each(func(i int, selection *goquery.Selection) {
							span := selection.Find("span")
							if span.Size() < 2 {
								return
							}
							tmp := strings.TrimSpace(span.Eq(1).Text())
							switch strings.TrimSpace(span.Eq(0).Text()) {
							case "联系人：", "联系卖家":
								contact = tmp
							case "电话号码：", "电话":
								tel = tmp
							case "手机号码：", "手机":
								mobile = tmp
							case "所在地区：", "地址":
								address = tmp
							case "主营产品：", "主营":
								major = tmp
							case "邮件", "邮箱":
								email = tmp
							}
						})
						ctx.SetTemp("contact", contact)
						ctx.SetTemp("tel", tel)
						ctx.SetTemp("mobile", mobile)
						ctx.SetTemp("major", major)
						ctx.SetTemp("address", address)
						ctx.SetTemp("email", email)
						ctx.Parse("企业信息")
					} else {
						major := strings.TrimSpace(dom.Find(".para.m-t-30,.m-t-30.main").Text())
						if major == "" {
							dom.Find(".intro .register p").Each(func(i int, selection *goquery.Selection) {
								span := selection.Find("span")
								if span.Size() > 1 && span.Eq(0).Text() == "主营产品：" {
									major = span.Eq(1).Text()
								}
							})
						}
						temps := ctx.GetTemps()
						temps["major"] = major
						ctx.AddQueue(&request.Request{
							Url:  ctx.GetUrl() + "/contact",
							Rule: "企业联系",
							Temp: temps,
						})
					}
				},
			},
			"企业联系": {
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					var contact, tel, mobile, major, address, email string
					des := dom.Find(".des")
					logs.Log.Informational("des.size:>>>>>>>>>>>>>:%d", des.Size())
					des.Each(func(i int, selection *goquery.Selection) {
						text := strings.TrimSpace(selection.Text())
						logs.Log.Informational("企业联系:>>>>>>>>>>>>>:" + text)
						//e.g 联系人：王先生
						arr := strings.Split(text, "：")
						if len(arr) < 2 {
							return
						}
						switch arr[0] {
						case "联系人":
							contact = arr[1]
						case "电话":
							tel = arr[1]
						case "地址":
							address = arr[1]
						case "邮箱":
							email = arr[1]
						case "手机":
							mobile = arr[1]
						case "主营":
							major = arr[1]
						}
					})
					if major != "" {
						//前面一步也有获取主营
						ctx.SetTemp("major", major)
					}
					ctx.SetTemp("contact", contact)
					ctx.SetTemp("tel", tel)
					ctx.SetTemp("mobile", mobile)
					ctx.SetTemp("address", address)
					ctx.SetTemp("email", email)
					ctx.Parse("企业信息")
				},
			},
			"企业信息": {
				ParseFunc: func(ctx *spider.Context) {
					ctx.Output(map[int]interface{}{
						0: ctx.GetTemp("name", ""),
						1: ctx.GetTemp("major", ""),
						2: ctx.GetTemp("province", ""),
						3: ctx.GetTemp("city", ""),
						4: ctx.GetTemp("street", ""),
						5: ctx.GetTemp("contact", ""),
						6: ctx.GetTemp("tel", ""),
						7: ctx.GetTemp("mobile", ""),
						8: ctx.GetTemp("email", ""),
						9: ctx.GetTemp("address", ""),
					})
				},
				ItemFields: []string{
					"名称",
					"主营",
					"省",
					"市",
					"区",
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
