package main

import (
    "context"
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "github.com/chromedp/chromedp"
    "log"
    _ "log"
    "math"
    libUrl "net/url"
    "strconv"
    "strings"
    "time"

    "regexp"
)

//func findComment(n *html.Node) *html.Node {
//    if n == nil {
//        return nil
//    }
//    if n.Type == html.CommentNode {
//        return n
//    }
//    if res := findComment(n.FirstChild); res != nil {
//        return res
//    }
//    if res := findComment(n.NextSibling); res != nil {
//        return res
//    }
//    return nil
//}

type ItemInfo struct {
    Date     string
    Name     string
    Price    string
    SiteName string
    Weight   string
    Cnt      string
    Unit     string
    Bundle   string
    ImageSrc string
    Size     string
}

func GetTypeSsgCrawlingInfo() error {
    //c := colly.NewCollector()
    ////var comment *html.Node
    //// Find and visit all links
    //c.OnHTML("a[class]", func(e *colly.HTMLElement) {
    //    //e.Request.Visit(e.Attr("class"))
    //    if e.Attr("class") == "clickable" {
    //        e.ForEach("em", func(_ int, elem *colly.HTMLElement) {
    //            if elem.Attr("class") == "tx_ko" {
    //                fmt.Println(elem.Text)
    //            }
    //        })
    //    }
    //})
    //
    //c.OnHTML("em[class]", func(e *colly.HTMLElement) {
    //    //e.Request.Visit(e.Attr("class"))
    //    if e.Attr("class") == "ssg_price" {
    //        fmt.Println(e.Text)
    //    }
    //})
    //
    //c.OnRequest(func(r *colly.Request) {
    //    fmt.Println("Visiting", r.URL)
    //})
    //
    //c.Visit("https://emart.ssg.com/search.ssg?target=all&query=%EA%B9%80%EC%B9%98")
    //
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    var name []string
    var price []string
    var imageSrc []string

    var data [5]string

    url := "https://www.ssg.com/search.ssg?target=all&query=%EA%B9%80%EC%B9%98"

    if err := chromedp.Run(ctx,

        chromedp.Navigate(url),
        chromedp.Sleep(time.Second*1),
        chromedp.OuterHTML("html", &data[0], chromedp.ByQuery),
        chromedp.ActionFunc(func(ctx context.Context) error {
            for i := 1; i <= 4; i++ {
                // todo search button and  click
                chromedp.Click("//a[text() = '"+strconv.Itoa(i+1)+"']", chromedp.BySearch).Do(ctx)
                chromedp.Sleep(time.Second * 2).Do(ctx)
                chromedp.OuterHTML("html", &data[i], chromedp.ByQuery).Do(ctx)
                //fmt.Printf("%s", data[i])
            }
            return nil
        }),
    ); err != nil {

        log.Fatal(err)
    }

    //fmt.Printf("%s", data)
    for i := 0; i < len(data); i++ {
        var doc, _ = goquery.NewDocumentFromReader(strings.NewReader(data[i]))
        temp0 := doc.Find("div")
        temp0.Each(func(i int, s0 *goquery.Selection) {
            class, _ := s0.Attr("class")
            if class == "tmpl_itemlist" {
                temp := s0.Find("li")
                temp.Each(func(i int, s1 *goquery.Selection) {
                    class2, _ := s1.Attr("class")
                    if class2 == "cunit_t232" {
                        temp1 := s1.Find("div")
                        temp1.Each(func(i int, sImage *goquery.Selection) {
                            classImage, _ := sImage.Attr("class")
                            if classImage == "thmb" {
                                tempImage := sImage.Find("img")
                                tempImage.Each(func(i int, sImage2 *goquery.Selection) {
                                    imgSrc, exists := sImage2.Attr("src")
                                    if exists {
                                        //fmt.Println(imgSrc)
                                        imageSrc = append(imageSrc, imgSrc)
                                    }
                                })
                            }
                        })

                        temp2 := s1.Find("a")
                        temp2.Each(func(i int, s2 *goquery.Selection) {
                            class3, _ := s2.Attr("class")
                            if class3 == "clickable" {
                                temp3 := s2.Find("em")
                                temp3.Each(func(i int, s3 *goquery.Selection) {
                                    class4, _ := s3.Attr("class")
                                    if class4 == "tx_ko" {
                                        //fmt.Printf("%s\n", strings.TrimSpace(s3.Text()))
                                        name = append(name, s3.Text())
                                    }
                                })
                            }
                        })

                        temp3 := s1.Find("em")
                        tempPrice := math.MaxInt
                        var tempPriceString string
                        temp3.Each(func(i int, s3 *goquery.Selection) {
                            class3, _ := s3.Attr("class")
                            if class3 == "ssg_price" {
                                //fmt.Printf("%s\n", strings.Trim(s3.Text(), ","))
                                intVar, _ := strconv.Atoi(strings.Replace(s3.Text(), ",", "", -1))
                                if intVar <= 100 {
                                    return
                                }
                                if intVar < tempPrice {
                                    tempPrice = intVar
                                    tempPriceString = s3.Text()
                                }
                            }
                        })
                        //fmt.Printf("%s\n", tempPrice)
                        price = append(price, tempPriceString)
                    }
                })
            }
        })
    }

    tempUrl, err := libUrl.Parse(url)
    if err != nil {
        log.Fatal(err)
    }
    parts := strings.Split(tempUrl.Hostname(), ".")
    //domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
    //fmt.Println(parts[len(parts)-2])

    currentTime := time.Now()
    currentDate := currentTime.Format("2006-01-02")

    var itemInfo []ItemInfo

    for i := 0; i < len(name); i++ {
        r, _ := regexp.Compile("(([0-9]*[.])?[0-9]+(g|kg))")

        weight := r.FindString(name[i])

        r, _ = regexp.Compile("(([0-9]?[-])?[0-9]+(개|봉|입|과))")

        cnt := r.FindString(name[i])

        r, _ = regexp.Compile("(g|kg)")

        unit := r.FindString(name[i])

        r, _ = regexp.Compile("(([0-9]+[+][0-9]+))")

        bundle := r.FindString(name[i])

        temp := ItemInfo{
            Date:     currentDate,
            Name:     name[i],
            Price:    price[i],
            SiteName: parts[1],
            Weight:   weight,
            Cnt:      cnt,
            Unit:     unit,
            Bundle:   bundle,
            ImageSrc: imageSrc[i],
        }
        itemInfo = append(itemInfo, temp)
        fmt.Printf("%v\n", temp)
    }

    return nil
}

func GetTypeHomeplusCrawlingInfo() error {
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    var name []string
    var price []string
    var imageSrc []string

    url := "https://front.homeplus.co.kr/search?entry=recent&keyword=%EC%82%AC%EA%B3%BC"

    var data [5]string
    //var example string
    //var nodes []*cdp.Node
    if err := chromedp.Run(ctx,

        chromedp.Navigate(url),
        chromedp.Sleep(time.Second*1),
        chromedp.OuterHTML("html", &data[0], chromedp.ByQuery),
        chromedp.ActionFunc(func(ctx context.Context) error {
            for i := 1; i <= 4; i++ {
                // todo search button and  click
                chromedp.Click("button.btnNext.css-1qz8j5i-buttonPagination", chromedp.ByQueryAll).Do(ctx)
                //chromedp.Click("//button[text() = '"+strconv.Itoa(i+1)+"']", chromedp.BySearch).Do(ctx)
                chromedp.Sleep(time.Second * 3).Do(ctx)
                chromedp.OuterHTML("html", &data[i], chromedp.ByQuery).Do(ctx)
            }
            return nil
        }),
        //chromedp.Click("#root > div > div.css-1di1x1r-container > div.css-oiwa5q-defaultStyle-gridRow-IntegratedSearch > div.mainWrap > div:nth-child(2) > div > div.pagination-js.css-dpcmyw-defaultStyle > button:nth-child(11)", chromedp.ByQueryAll),
        //chromedp.Sleep(time.Second*5),
        //chromedp.OuterHTML("html", &data2, chromedp.ByQuery),
    ); err != nil {

        log.Fatal(err)
    }
    for i := 0; i < len(data); i++ {
        //fmt.Printf("%d\n", i)
        var doc, _ = goquery.NewDocumentFromReader(strings.NewReader(data[i]))
        temp0 := doc.Find("div")
        temp0.Each(func(i int, s0 *goquery.Selection) {
            class, _ := s0.Attr("class")
            if class == "itemDisplayList" {
                temp := s0.Find("div")
                temp.Each(func(i int, s1 *goquery.Selection) {
                    class2, _ := s1.Attr("class")
                    if class2 == "unitItemBox mid cardType css-t6werr-itemDisplayStyle" {
                        temp1 := s1.Find("div")
                        temp1.Each(func(i int, sImage *goquery.Selection) {
                            classImage, _ := sImage.Attr("class")
                            if classImage == "thumbWrap " {
                                tempImage := sImage.Find("img")
                                tempImage.Each(func(i int, sImage2 *goquery.Selection) {
                                    imgSrc, exists := sImage2.Attr("src")
                                    if exists {
                                        //fmt.Println(imgSrc)
                                        imageSrc = append(imageSrc, imgSrc)
                                    }
                                })
                            }
                        })

                        temp2 := s1.Find("a")
                        temp2.Each(func(i int, s2 *goquery.Selection) {
                            class3, _ := s2.Attr("class")
                            if class3 == "productTitle css-y9z3ts-defaultStyle-Linked" {
                                temp3 := s2.Find("p")
                                temp3.Each(func(i int, s3 *goquery.Selection) {
                                    class4, _ := s3.Attr("class")
                                    if class4 == "css-12cdo53-defaultStyle-Typography-ellips" {
                                        //fmt.Printf("%s\n", strings.TrimSpace(s3.Text()))
                                        name = append(name, s3.Text())
                                    }
                                })
                            }
                        })
                        temp3 := s1.Find("strong")
                        tempPrice := math.MaxInt
                        var tempPriceString string
                        temp3.Each(func(i int, s3 *goquery.Selection) {
                            class3, _ := s3.Attr("class")
                            if class3 == "priceValue" {
                                //fmt.Printf("%s\n", strings.Trim(s3.Text(), ","))
                                intVar, _ := strconv.Atoi(strings.Replace(s3.Text(), ",", "", -1))
                                if intVar <= 100 {
                                    return
                                }
                                if intVar < tempPrice {
                                    tempPrice = intVar
                                    tempPriceString = s3.Text()
                                }
                            }
                        })
                        //fmt.Printf("%s\n", tempPrice)
                        price = append(price, tempPriceString)
                    }
                })
            }
        })
    }
    //for i := 0; i < len(data); i++ {
    //	var doc, _ = goquery.NewDocumentFromReader(strings.NewReader(data[i]))
    //	temp0 := doc.Find("div")
    //	temp0.Each(func(i int, s0 *goquery.Selection) {
    //		class, _ := s0.Attr("class")
    //		if class == "itemDisplayList" {
    //			temp := s0.Find("a")
    //			temp.Each(func(i int, s *goquery.Selection) {
    //				class, _ := s.Attr("class")
    //				if class == "productTitle css-y9z3ts-defaultStyle-Linked" {
    //					temp2 := s.Find("p")
    //					temp2.Each(func(i int, s2 *goquery.Selection) {
    //						class2, _ := s2.Attr("class")
    //						if class2 == "css-12cdo53-defaultStyle-Typography-ellips" {
    //							//fmt.Printf("%s\n", strings.TrimSpace(s2.Text()))
    //							name = append(name, s2.Text())
    //						}
    //					})
    //				}
    //			})
    //
    //			temp3 := s0.Find("strong")
    //			temp3.Each(func(i int, s3 *goquery.Selection) {
    //				class3, _ := s3.Attr("class")
    //				if class3 == "priceValue" && s3.Parent().Children().Length() == 1 {
    //					//fmt.Printf("%s", strings.Trim(s3.Text(), ","))
    //					intVar, _ := strconv.Atoi(strings.Replace(s3.Text(), ",", "", -1))
    //					if intVar <= 100 {
    //						return
    //					}
    //					//fmt.Printf("%s\n", strings.TrimSpace(s3.Text()))
    //					price = append(price, s3.Text())
    //				}
    //			})
    //		}
    //	})
    //}

    tempUrl, err := libUrl.Parse(url)
    if err != nil {
        log.Fatal(err)
    }
    parts := strings.Split(tempUrl.Hostname(), ".")
    //domain := parts[len(parts)-2] + "." + parts[len(parts)-1]

    currentTime := time.Now()
    currentDate := currentTime.Format("2006-01-02")

    var itemInfo []ItemInfo

    for i := 0; i < len(name); i++ {
        r, _ := regexp.Compile("(([0-9]*[.])?[0-9]+(g|kg))")

        weight := r.FindString(name[i])

        r, _ = regexp.Compile("(([0-9]*(-|~))?[0-9]+(개|봉|입|과))")

        cnt := r.FindString(name[i])

        r, _ = regexp.Compile("(g|kg)")

        unit := r.FindString(name[i])

        r, _ = regexp.Compile("(([0-9]+[+][0-9]+))")

        bundle := r.FindString(name[i])

        temp := ItemInfo{
            Date:     currentDate,
            Name:     name[i],
            Price:    price[i],
            SiteName: parts[1],
            Weight:   weight,
            Cnt:      cnt,
            Unit:     unit,
            Bundle:   bundle,
            ImageSrc: imageSrc[i],
        }
        itemInfo = append(itemInfo, temp)
        fmt.Printf("%v\n", temp)
    }

    return nil
}

func GetTypeLotteCrawlingInfo() error {
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    var name []string
    var price []string
    var imageSrc []string

    url := "https://www.lotteon.com/search/search/search.ecn?render=search&platform=pc&q=%EC%82%AC%EA%B3%BC&mallId=1"

    var data [5]string
    //var example string
    //var nodes []*cdp.Node
    if err := chromedp.Run(ctx,

        chromedp.Navigate(url),
        chromedp.Sleep(time.Second*1),
        chromedp.OuterHTML("html", &data[0], chromedp.ByQuery),
        chromedp.ActionFunc(func(ctx context.Context) error {
            for i := 1; i <= 4; i++ {
                // todo search button and  click
                //chromedp.Click("button.btnNext.css-1qz8j5i-buttonPagination", chromedp.ByQueryAll).Do(ctx)
                chromedp.Click("//a[text() = '"+strconv.Itoa(i+1)+"']", chromedp.BySearch).Do(ctx)
                chromedp.Sleep(time.Second * 3).Do(ctx)
                chromedp.OuterHTML("html", &data[i], chromedp.ByQuery).Do(ctx)
            }
            return nil
        }),
        //chromedp.Click("#root > div > div.css-1di1x1r-container > div.css-oiwa5q-defaultStyle-gridRow-IntegratedSearch > div.mainWrap > div:nth-child(2) > div > div.pagination-js.css-dpcmyw-defaultStyle > button:nth-child(11)", chromedp.ByQueryAll),
        //chromedp.Sleep(time.Second*5),
        //chromedp.OuterHTML("html", &data2, chromedp.ByQuery),
    ); err != nil {

        log.Fatal(err)
    }
    for i := 0; i < len(data); i++ {
        fmt.Printf("%d\n", i)
        var doc, _ = goquery.NewDocumentFromReader(strings.NewReader(data[i]))
        temp0 := doc.Find("div")
        temp0.Each(func(i int, s0 *goquery.Selection) {
            class, _ := s0.Attr("class")
            if class == "srchResultProductArea" {
                temp := s0.Find("li")
                temp.Each(func(i int, s1 *goquery.Selection) {
                    class2, _ := s1.Attr("class")
                    if class2 == "srchProductItem" {
                        temp1 := s1.Find("div")
                        temp1.Each(func(i int, sImage *goquery.Selection) {
                            classImage, _ := sImage.Attr("class")
                            if classImage == "srchThumbImageWrap" {
                                tempImage := sImage.Find("img")
                                tempImage.Each(func(i int, sImage2 *goquery.Selection) {
                                    imgSrc, exists := sImage2.Attr("src")
                                    if exists {
                                        //fmt.Println(imgSrc)
                                        imageSrc = append(imageSrc, imgSrc)
                                    }
                                })
                            }
                        })

                        temp2 := s1.Find("a")
                        temp2.Each(func(i int, s2 *goquery.Selection) {
                            class3, _ := s2.Attr("class")
                            if class3 == "srchGridProductUnitLink" {
                                temp3 := s2.Find("div")
                                temp3.Each(func(i int, s3 *goquery.Selection) {
                                    class4, _ := s3.Attr("class")
                                    if class4 == "srchProductUnitTitle" {
                                        //fmt.Printf("%s\n", strings.TrimSpace(s3.Text()))
                                        name = append(name, s3.Text())
                                    }
                                })
                            }
                        })
                        temp3 := s1.Find("span")
                        tempPrice := math.MaxInt
                        var tempPriceString string
                        temp3.Each(func(i int, s3 *goquery.Selection) {
                            class3, _ := s3.Attr("class")
                            if class3 == "s-product-price__number" {
                                //fmt.Printf("%s\n", strings.Trim(s3.Text(), ","))
                                intVar, _ := strconv.Atoi(strings.Replace(s3.Text(), ",", "", -1))
                                if intVar <= 100 {
                                    return
                                }
                                if intVar < tempPrice {
                                    tempPrice = intVar
                                    tempPriceString = s3.Text()
                                }
                            }
                        })
                        //fmt.Printf("%s\n", tempPrice)
                        price = append(price, tempPriceString)
                    }
                })
            }
        })
    }

    tempUrl, err := libUrl.Parse(url)
    if err != nil {
        log.Fatal(err)
    }
    parts := strings.Split(tempUrl.Hostname(), ".")
    //domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
    //fmt.Println(parts[len(parts)-2])

    currentTime := time.Now()
    currentDate := currentTime.Format("2006-01-02")

    var itemInfo []ItemInfo

    for i := 0; i < len(name); i++ {
        r, _ := regexp.Compile("(([0-9]*[.])?[0-9]+(g|kg))")

        weight := r.FindString(name[i])

        r, _ = regexp.Compile("(([0-9]?[-])?[0-9]+(개|봉|입|과))")

        cnt := r.FindString(name[i])

        r, _ = regexp.Compile("(g|kg)")

        unit := r.FindString(name[i])

        r, _ = regexp.Compile("(([0-9]+[+][0-9]+))")

        bundle := r.FindString(name[i])

        temp := ItemInfo{
            Date:     currentDate,
            Name:     name[i],
            Price:    price[i],
            SiteName: parts[1],
            Weight:   weight,
            Cnt:      cnt,
            Unit:     unit,
            Bundle:   bundle,
            ImageSrc: imageSrc[i],
        }
        itemInfo = append(itemInfo, temp)
        fmt.Printf("%v\n", temp)
    }

    return nil
}

func GetScrawlingInfo(buttonElem string, buttonClass string, divContainerClass string, splitElem string,
    splitElemClass string, divImageclass string, aTitleclass string, titleElem string,
    titleClass string, priceElem string, priceClass string, url string, item string) error {
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    var name []string
    var price []string
    var imageSrc []string

    realUrl := url + item

    var data [5]string
    if err := chromedp.Run(ctx,

        chromedp.Navigate(realUrl),
        chromedp.Sleep(time.Second*1),
        chromedp.OuterHTML("html", &data[0], chromedp.ByQuery),
        chromedp.ActionFunc(func(ctx context.Context) error {
            for i := 1; i <= 4; i++ {
                // todo search button and  click
                if buttonElem == "button" {
                    chromedp.Click(buttonElem+"."+strings.Join(strings.Split(buttonClass, " "), "."), chromedp.ByQueryAll).Do(ctx)
                } else {
                    chromedp.Click("//"+buttonElem+"[text() = '"+strconv.Itoa(i+1)+"']", chromedp.BySearch).Do(ctx)
                }
                chromedp.Sleep(time.Second * 3).Do(ctx)
                chromedp.OuterHTML("html", &data[i], chromedp.ByQuery).Do(ctx)
            }
            return nil
        }),
        //chromedp.Click("#root > div > div.css-1di1x1r-container > div.css-oiwa5q-defaultStyle-gridRow-IntegratedSearch > div.mainWrap > div:nth-child(2) > div > div.pagination-js.css-dpcmyw-defaultStyle > button:nth-child(11)", chromedp.ByQueryAll),
        //chromedp.Sleep(time.Second*5),
        //chromedp.OuterHTML("html", &data2, chromedp.ByQuery),
    ); err != nil {

        log.Fatal(err)
    }
    for i := 0; i < len(data); i++ {
        var doc, _ = goquery.NewDocumentFromReader(strings.NewReader(data[i]))
        temp0 := doc.Find("div")
        temp0.Each(func(i int, s0 *goquery.Selection) {
            class, _ := s0.Attr("class")
            if class == divContainerClass {
                temp := s0.Find(splitElem)
                temp.Each(func(i int, s1 *goquery.Selection) {
                    class2, _ := s1.Attr("class")
                    if class2 == splitElemClass {
                        temp1 := s1.Find("div")
                        temp1.Each(func(i int, sImage *goquery.Selection) {
                            classImage, _ := sImage.Attr("class")
                            if classImage == divImageclass {
                                tempImage := sImage.Find("img")
                                tempImage.Each(func(i int, sImage2 *goquery.Selection) {
                                    imgSrc, exists := sImage2.Attr("src")
                                    if exists {
                                        //fmt.Println(imgSrc)
                                        imageSrc = append(imageSrc, imgSrc)
                                    }
                                })
                            }
                        })

                        temp2 := s1.Find("a")
                        temp2.Each(func(i int, s2 *goquery.Selection) {
                            class3, _ := s2.Attr("class")
                            if class3 == aTitleclass {
                                temp3 := s2.Find(titleElem)
                                temp3.Each(func(i int, s3 *goquery.Selection) {
                                    class4, _ := s3.Attr("class")
                                    if class4 == titleClass {
                                        //fmt.Printf("%s\n", strings.TrimSpace(s3.Text()))
                                        name = append(name, s3.Text())
                                    }
                                })
                            }
                        })
                        temp3 := s1.Find(priceElem)
                        tempPrice := math.MaxInt
                        var tempPriceString string
                        temp3.Each(func(i int, s3 *goquery.Selection) {
                            class3, _ := s3.Attr("class")
                            if class3 == priceClass {
                                //fmt.Printf("%s\n", strings.Trim(s3.Text(), ","))
                                intVar, _ := strconv.Atoi(strings.Replace(s3.Text(), ",", "", -1))
                                if intVar <= 100 {
                                    return
                                }
                                if intVar < tempPrice {
                                    tempPrice = intVar
                                    tempPriceString = s3.Text()
                                }
                            }
                        })
                        //fmt.Printf("%s\n", tempPrice)
                        price = append(price, tempPriceString)
                    }
                })
            }
        })
    }

    tempUrl, err := libUrl.Parse(url)
    if err != nil {
        log.Fatal(err)
    }
    parts := strings.Split(tempUrl.Hostname(), ".")
    //domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
    //fmt.Println(parts[len(parts)-2])

    currentTime := time.Now()
    currentDate := currentTime.Format("2006-01-02")

    var itemInfo []ItemInfo

    for i := 0; i < len(name); i++ {
        r, _ := regexp.Compile("(([0-9]*[.])?[0-9]+(g|kg))")

        weight := r.FindString(name[i])

        r, _ = regexp.Compile("(([0-9]?[-])?[0-9]+(개|봉|입|과))")

        cnt := r.FindString(name[i])

        r, _ = regexp.Compile("(g|kg)")

        unit := r.FindString(name[i])

        r, _ = regexp.Compile("(([0-9]+[+][0-9]+))")

        bundle := r.FindString(name[i])

        temp := ItemInfo{
            Date:     currentDate,
            Name:     name[i],
            Price:    price[i],
            SiteName: parts[1],
            Weight:   weight,
            Cnt:      cnt,
            Unit:     unit,
            Bundle:   bundle,
            ImageSrc: imageSrc[i],
        }
        itemInfo = append(itemInfo, temp)
        fmt.Printf("%v\n", temp)
    }

    return nil
}

func main() {
    //GetTypeSsgCrawlingInfo()
    //GetTypeHomeplusCrawlingInfo()
    //GetTypeLotteCrawlingInfo()
    //SSG
    //GetScrawlingInfo("a", "", "tmpl_itemlist", "li",
    //    "cunit_t232", "thmb", "clickable", "em", "tx_ko",
    //    "em", "ssg_price",
    //    "https://www.ssg.com/search.ssg?target=all&query=", "%EA%B9%80%EC%B9%98")
    // HOMEPLUS
    //GetScrawlingInfo("button", "btnNext css-1qz8j5i-buttonPagination",
    //    "itemDisplayList", "div",
    //    "unitItemBox mid cardType css-t6werr-itemDisplayStyle", "thumbWrap ",
    //    "productTitle css-y9z3ts-defaultStyle-Linked", "p",
    //    "css-12cdo53-defaultStyle-Typography-ellips", "strong", "priceValue",
    //    "https://front.homeplus.co.kr/search?entry=recent&keyword=", "%EC%82%AC%EA%B3%BC")
    GetScrawlingInfo("a", "", "srchResultProductArea", "li",
        "srchProductItem", "srchThumbImageWrap", "srchGridProductUnitLink",
        "div", "srchProductUnitTitle", "span", "s-product-price__number",
        "https://www.lotteon.com/search/search/search.ecn?render=search&platform=pc&q=", "%EC%82%AC%EA%B3%BC")
}
