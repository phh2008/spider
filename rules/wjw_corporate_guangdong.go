package rules

import (
	"com.phh/spider/utils"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"github.com/henrylee2cn/pholcus/logs"
	"strconv"
	"strings"
)

func init() {
	WjwCorporateGuangdong.Register()
}

var WjwCorporateGuangdong = &spider.Spider{
	Name:         "全球五金A",
	Description:  "全球五金[广东1-35] http://www.wjw.cn/gongying/quyu/guangdong",
	EnableCookie: false,
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.AddQueue(&request.Request{
				Url:  "http://www.wjw.cn/gongying/quyu/guangdong",
				Rule: "供应搜索",
			})
		},
		Trunk: map[string]*spider.Rule{
			"供应搜索": {
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					lastPage, err := dom.Find(".page a").Last().Attr("href")
					if err {
						logs.Log.Error("************未找到尾页信息*************")
					}
					//截取总页数
					idx := strings.LastIndex(lastPage, "/")
					total, _ := strconv.Atoi(lastPage[idx:])
					logs.Log.Informational(">>>>>>总页数：%d", total)
					for i := 1; i <= 35; i++ {
						ctx.AddQueue(&request.Request{
							Url:  "http://www.wjw.cn/gongying/quyu/guangdong/" + strconv.Itoa(i),
							Rule: "供应列表",
						})
					}
				},
			},
			"供应列表": {
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					dom.Find(".list-hen li").Each(func(i int, li *goquery.Selection) {
						gsA := li.Find(".lhitem-gs a").First()
						if url, ok := gsA.Attr("href"); ok {
							//公司名称
							corpName := gsA.Text()
							//地区 格式：广东 东莞 || 广东 佛山 辖区
							region := li.Find(".lhitem-dz p").First().Text()

							logs.Log.Informational(">>>>>>公司名称：%s", corpName)
							logs.Log.Informational(">>>>>>地   区：%s", region)
							ctx.AddQueue(&request.Request{
								Url:  url,
								Rule: "公司主页",
								Temp: map[string]interface{}{
									"corpName": corpName,
									"region":   region,
								},
							})
						}
					})
				},
			},
			"公司主页": {
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					var major []string
					dom.Find(".SaleKeywords a").Each(func(i int, a *goquery.Selection) {
						major = append(major, a.Text())
					})
					info := dom.Find(".cft tbody tr").First().Find("td").First().Text()
					//公司名称
					corp := dom.Find(".HeadCompany h1").Text()
					//公司主页
					url := ctx.GetUrl()
					//诚信档案
					profileUrl, _ := dom.Find(".Nav li").Eq(2).Find("a").Attr("href")
					//联系我们
					contactUsUrl, _ := dom.Find(".Nav li").Last().Find("a").Attr("href")
					//公司地址 格式：协兴螺丝工业(深圳)有限公司 版权所有 公司地址：中国 广东 广东省东莞市厚街镇白濠工业区源泉路8号
					var address = dom.Find(".Footer980 span").Last().Text()
					//截取
					address = utils.SubString1(address, utils.UnicodeLastIndex(address, "：")+1)

					logs.Log.Informational(">>>>>>corp：%s", corp)
					logs.Log.Informational(">>>>>>url：%s", url)
					logs.Log.Informational(">>>>>>address：%s", address)
					ctx.AddQueue(&request.Request{
						Url:  url + profileUrl,
						Rule: "诚信档案",
						Temp: map[string]interface{}{
							"corpName":     ctx.GetTemp("corpName", corp),
							"major":        strings.Join(major, ","),
							"info":         info,
							"home":         url,
							"region":       ctx.GetTemp("region", ""),
							"address":      address,
							"contactUsUrl": url + contactUsUrl,
						},
					})
				},
			},
			"诚信档案": {
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					var registrationNo string
					var trList = dom.Find("#ctl00_ShopBody_divQiye .Cxtr1 table tr")
					if trList.Length() > 0 {
						//注册号
						registrationNo = trList.Eq(1).Find("td").Eq(3).Text()
					} else {
						trList = dom.Find("#ctl00_ShopBody_SincerityDetail tbody tr")
						registrationNo = trList.Find("#ctl00_ShopBody_Label_regNumber").Text()
					}
					temp := ctx.GetTemps()
					temp["registrationNo"] = strings.TrimSpace(registrationNo)
					logs.Log.Informational(">>>>>>registNo：%s", registrationNo)
					ctx.AddQueue(&request.Request{
						Url:  ctx.GetTemp("contactUsUrl", "").(string),
						Rule: "联系我们",
						Temp: temp,
					})
				},
			},
			"联系我们": {
				ItemFields: []string{
					"公司名称",
					"主营",
					"注册号",
					"联系人",
					"性别",
					"职位",
					"手机",
					"电话",
					"简介",
					"主页",
					"地区",
					"地址",
				},
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					//电话手机用js设置的，无法获取，解析js试下
					//tel := dom.Find("#ctl00_ShopBody_divLianxi1 .P18 span").First().Text()
					//mobile := dom.Find("#ctl00_ShopBody_divLianxi1 .P18 span").Last().Text()
					//js文件如下
					//var itm = '13510831860';
					//var itt = '86-0755-88876123';
					var tel, mobile string
					dom.Find("body").Text()
					dom.Find("script").Each(func(i int, js *goquery.Selection) {
						match := regTelMobile.FindAllStringSubmatch(js.Text(), 2)
						for _, v := range match {
							for _, m := range v {
								if strings.Index(m, "var itm =") >= 0 {
									tmp := strings.ReplaceAll(m, "var itm = '", "")
									tmp = strings.ReplaceAll(tmp, "';", "")
									mobile = tmp
								} else if strings.Index(m, "var itt = '") >= 0 {
									tmp := strings.ReplaceAll(m, "var itt = '", "")
									tmp = strings.ReplaceAll(tmp, "';", "")
									tel = tmp
								}
							}
						}
					})
					//格式：郑 先生（销售 经理）|| 郑晶 女士（销售总部 销售顾问）|| 李小华 先生（销售代表）
					contactStr := dom.Find("#ctl00_ShopBody_divLianxi1 .P17").Text()
					contactStr = strings.ReplaceAll(contactStr, "(", " ")
					contactStr = strings.ReplaceAll(contactStr, "（", " ")
					contactStr = strings.ReplaceAll(contactStr, ")", "")
					contactStr = strings.ReplaceAll(contactStr, "）", "")
					strList := strings.Split(contactStr, " ")
					contact := utils.If(len(strList) > 0, strList[0], "").(string)
					sex := utils.If(len(strList) > 1, strList[1], "").(string)
					job := utils.If(len(strList) > 2, strings.Join(strList[2:], ","), "").(string)

					logs.Log.Informational(">>>>>>contact：%s", contact)
					logs.Log.Informational(">>>>>>sex：%s", sex)
					logs.Log.Informational(">>>>>>job：%s", job)
					logs.Log.Informational(">>>>>>tel：%s", tel)
					logs.Log.Informational(">>>>>>mobile：%s", mobile)
					//输出内容
					ctx.Output(map[int]interface{}{
						0:  ctx.GetTemp("corpName", ""),
						1:  ctx.GetTemp("major", ""),
						2:  ctx.GetTemp("registrationNo", ""),
						3:  contact,
						4:  sex,
						5:  job,
						6:  mobile,
						7:  tel,
						8:  ctx.GetTemp("info", ""),
						9:  ctx.GetTemp("home", ""),
						10: ctx.GetTemp("region", ""),
						11: ctx.GetTemp("address", ""),
					})
				},
			},
		},
	},
}