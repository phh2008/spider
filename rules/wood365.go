package rules

import (
	"com.phh/spider/utils"
	"fmt"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"strings"
)

func init() {
	Wood365.Register()
}

var Wood365 = &spider.Spider{
	Name:         "木业网",
	Description:  "http://www.wood365.cn/Corp/",
	EnableCookie: false,
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.Aid(map[string]interface{}{
				"loop": [2]int{1, 30},
				"rule": "企业分页",
			}, "企业分页")
		},
		Trunk: map[string]*spider.Rule{
			"企业分页": {
				AidFunc: func(ctx *spider.Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] <= loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  fmt.Sprintf("http://www.wood365.cn/Corp/corp_%d.html", loop[0]),
							Rule: aid["rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					dom.Find("#htmlCorpList li").Each(func(i int, li *goquery.Selection) {
						ctx.SetTemp("li", li)
						ctx.Parse("企业条目")
					})
				},
			},
			"企业条目": {
				ParseFunc: func(ctx *spider.Context) {
					li := ctx.GetTemp("li", &goquery.Selection{}).(*goquery.Selection)
					supply := li.Find(".supply-table-wrap")
					corpA := supply.Find("h2 a")
					corpName := corpA.Text()
					corpUrl, _ := corpA.Attr("href")
					corpUrl = "http://www.wood365.cn" + corpUrl
					p := supply.Find("p")
					//主营
					//主营产品：单板烘干机,单板干燥机,踏式单板剪切机,旋切机木轴压辊,木材单板烘干机,木材单板烘干机,滚筒式单板烘干机,网带式单板烘干机,网带式单板干燥机,滚筒式单板干燥机
					major := p.Eq(0).Text()
					major = strings.ReplaceAll(major, "\"", "")
					majorArr := strings.Split(major, "：")
					major = utils.If(len(majorArr) > 1, majorArr[1], "").(string)
					//联系人：李先生
					contact := p.Eq(1).Text()
					contactArr := strings.Split(contact, "：")
					contact = utils.If(len(contactArr) > 1, contactArr[1], "").(string)
					//联系电话：13089848211  13304837119
					tel := p.Eq(2).Text()
					telArr := strings.Split(tel, "：")
					tel = utils.If(len(telArr) > 1, telArr[1], "").(string)
					ctx.AddQueue(&request.Request{
						Url:  corpUrl,
						Rule: "企业信息",
						Temp: map[string]interface{}{
							"corpName": corpName,
							"major":    major,
							"contact":  contact,
							"tel":      tel,
						},
					})
				},
			},
			"企业信息": {
				ItemFields: []string{
					"企业名称",
					"简介",
					"主营",
					"联系人",
					"联系电话",
					"邮箱",
					"地址",
				},
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					var address, email, about string
					p := dom.Find("#Contact_list p")
					if p.Size() > 0 {
						address = p.Eq(1).Text()
						email = p.Eq(4).Text()
						about = dom.Find(".banner-para").Text()
					} else {
						p = dom.Find(".lianxi-fl .hot-read-lis p")
						address = p.Eq(1).Text()
						email = p.Eq(4).Text()
						about = dom.Find(".company-intro-content p").First().Text()
					}
					about = strings.TrimSpace(about)
					emailArr := strings.Split(email, "：")
					email = utils.If(len(emailArr) > 1, emailArr[1], "").(string)
					addressArr := strings.Split(address, "：")
					address = utils.If(len(addressArr) > 1, addressArr[1], "").(string)
					ctx.Output(map[int]interface{}{
						0: ctx.GetTemp("corpName", ""),
						1: about,
						2: ctx.GetTemp("major", ""),
						3: ctx.GetTemp("contact", ""),
						4: ctx.GetTemp("tel", ""),
						5: email,
						6: address,
					})
				},
			},
		},
	},
}
