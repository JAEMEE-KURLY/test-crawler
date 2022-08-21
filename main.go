package main

import (
    "context"
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "github.com/chromedp/chromedp"
    "log"
    _ "log"
    "strconv"
    "strings"
    "time"
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

    url := "https://emart.ssg.com/search.ssg?target=all&query=%EA%B9%80%EC%B9%98"

    var data string

    if err := chromedp.Run(ctx,

        chromedp.Navigate(url),
        chromedp.Sleep(time.Second*1),
        chromedp.OuterHTML("html", &data, chromedp.ByQuery),
    ); err != nil {

        log.Fatal(err)
    }

    //fmt.Printf("%s", data)

    var doc, _ = goquery.NewDocumentFromReader(strings.NewReader(data))
    temp := doc.Find("a")
    temp.Each(func(i int, s *goquery.Selection) {
        class, _ := s.Attr("class")
        if class == "clickable" {
            temp2 := s.Find("em")
            temp2.Each(func(i int, s2 *goquery.Selection) {
                class2, _ := s2.Attr("class")
                if class2 == "tx_ko" {
                    fmt.Printf("%s\n", strings.TrimSpace(s2.Text()))
                }
            })
        }
    })

    temp3 := doc.Find("em")
    temp3.Each(func(i int, s3 *goquery.Selection) {
        class3, _ := s3.Attr("class")
        if class3 == "ssg_price" {
            //fmt.Printf("%s", strings.Trim(s3.Text(), ","))
            intVar, _ := strconv.Atoi(strings.Replace(s3.Text(), ",", "", -1))
            if intVar <= 100 {
                return
            }
            fmt.Printf("%s\n", strings.TrimSpace(s3.Text()))
        }
    })

    return nil
}

func GetTypeHomeplusCrawlingInfo() error {
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    url := "https://front.homeplus.co.kr/search?entry=recent&keyword=%EC%82%AC%EA%B3%BC"

    var data [5]string
    //var example string
    //var nodes []*cdp.Node
    chromedp.Run(ctx,

        chromedp.Navigate(url),
        chromedp.Sleep(time.Second*1),
        chromedp.OuterHTML("html", &data[0], chromedp.ByQuery),
        chromedp.ActionFunc(func(ctx context.Context) error {
            for i := 1; i < 4; i++ {
                // todo search button and  click
                fmt.Printf("%d", i)
                chromedp.Click("button.btnNext.css-1qz8j5i-buttonPagination", chromedp.ByQueryAll).Do(ctx)
                chromedp.Sleep(time.Second * 2).Do(ctx)
                chromedp.OuterHTML("html", &data[i], chromedp.ByQuery).Do(ctx)
                fmt.Printf("%s", data[i])
            }
            return nil
        }),
        //chromedp.Click("#root > div > div.css-1di1x1r-container > div.css-oiwa5q-defaultStyle-gridRow-IntegratedSearch > div.mainWrap > div:nth-child(2) > div > div.pagination-js.css-dpcmyw-defaultStyle > button:nth-child(11)", chromedp.ByQueryAll),
        //chromedp.Sleep(time.Second*5),
        //chromedp.OuterHTML("html", &data2, chromedp.ByQuery),
    )
    for i := 0; i < len(data); i++ {
        fmt.Printf("%d", i)
        fmt.Printf("=======================================================")
        var doc, _ = goquery.NewDocumentFromReader(strings.NewReader(data[i]))
        temp := doc.Find("a")
        temp.Each(func(i int, s *goquery.Selection) {
            class, _ := s.Attr("class")
            if class == "productTitle css-y9z3ts-defaultStyle-Linked" {
                temp2 := s.Find("p")
                temp2.Each(func(i int, s2 *goquery.Selection) {
                    class2, _ := s2.Attr("class")
                    if class2 == "css-12cdo53-defaultStyle-Typography-ellips" {
                        fmt.Printf("%s\n", strings.TrimSpace(s2.Text()))
                    }
                })
            }
        })

        temp3 := doc.Find("strong")
        temp3.Each(func(i int, s3 *goquery.Selection) {
            class3, _ := s3.Attr("class")
            if class3 == "priceValue" {
                //fmt.Printf("%s", strings.Trim(s3.Text(), ","))
                intVar, _ := strconv.Atoi(strings.Replace(s3.Text(), ",", "", -1))
                if intVar <= 100 {
                    return
                }
                fmt.Printf("%s\n", strings.TrimSpace(s3.Text()))
            }
        })
    }

    return nil
}

func GetTypeLotteCrawlingInfo() error {
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    url := "https://www.lotteon.com/search/search/search.ecn?render=search&platform=pc&q=%EC%82%AC%EA%B3%BC&mallId=1"

    var data [10]string

    if err := chromedp.Run(ctx,

        chromedp.Navigate(url),
        chromedp.Sleep(time.Second*1),
        chromedp.OuterHTML("html", &data[0], chromedp.ByQuery),
    ); err != nil {

        log.Fatal(err)
    }

    var doc, _ = goquery.NewDocumentFromReader(strings.NewReader(data[1]))
    temp := doc.Find("a")
    temp.Each(func(i int, s *goquery.Selection) {
        class, _ := s.Attr("class")
        if class == "srchGridProductUnitLink" {
            temp2 := s.Find("div")
            temp2.Each(func(i int, s2 *goquery.Selection) {
                class2, _ := s2.Attr("class")
                if class2 == "srchProductUnitTitle" {
                    fmt.Printf("%s : ", strings.TrimSpace(s2.Text()))
                }
            })
            temp3 := s.Find("span")
            temp3.Each(func(i int, s3 *goquery.Selection) {
                class3, _ := s3.Attr("class")
                if class3 == "s-product-price__number" {
                    //fmt.Printf("%s", strings.Trim(s3.Text(), ","))
                    intVar, _ := strconv.Atoi(strings.Replace(s3.Text(), ",", "", -1))
                    if intVar <= 100 {
                        return
                    }
                    fmt.Printf("%s\n", strings.TrimSpace(s3.Text()))
                }
            })
        }
    })

    //// HTML 읽기
    //html, err := goquery.NewDocumentFromReader(res.Body)
    //if err != nil {
    //    log.Fatal(err)
    //}

    return nil
}

func main() {
    //GetTypeSsgCrawlingInfo()
    GetTypeHomeplusCrawlingInfo()
    //GetTypeLotteCrawlingInfo()
}
