package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func NewDocumentFromURL(url string) *goquery.Document {
	fmt.Printf("正在加载%s\n", url)
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	res, err := client.Do(request)

	//res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	reader := simplifiedchinese.GBK.NewDecoder().Reader(res.Body)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func getProvinces(parent string, root *Region) []*Region {

	doc := NewDocumentFromURL(parent)

	parent = getRoot(parent)

	//解析出省份
	doc.Find(".provincetr ").Each(func(i int, s *goquery.Selection) {
		s.Find("td a").Each(func(i int, s *goquery.Selection) {
			provinceName := s.Text()
			//获取省份的连接地址!!
			provinceUrl, _ := s.Attr("href")
			provinceUrl = parent + "/" + provinceUrl

			root.add(NewRegion(RegionType_City, provinceName, "", provinceUrl))

		})
	})
	return root.Children
}
func getCities(parent string, root *Region) []*Region {

	doc := NewDocumentFromURL(parent)
	parent = getRoot(parent)
	//解析出省份
	doc.Find(".citytr ").Each(func(i int, s *goquery.Selection) {

		cityName := ""
		cityCode := ""
		cityUrl := ""

		s.Find("td a").Each(func(i int, s *goquery.Selection) {
			if i%2 == 0 {
				cityCode = s.Text()
				//获取省份的连接地址!!
				cityUrl, _ = s.Attr("href")
				cityUrl = parent + "/" + cityUrl
			} else {
				cityName = s.Text()
			}
		})
		root.add(NewRegion(RegionType_City, cityName, cityCode, cityUrl))
	})
	return root.Children
}
func getCounty(parent string, root *Region) []*Region {

	doc := NewDocumentFromURL(parent)
	parent = getRoot(parent)
	if root.Name == "郑州市" {
		fmt.Println("~~~~~~~")
	}
	//解析出省份
	doc.Find(".countytr ").Each(func(i int, s *goquery.Selection) {

		countyName := ""
		countyCode := ""
		countyUrl := ""

		s.Find("td ").Each(func(i int, s *goquery.Selection) {
			a := s.Find("a")
			if a.Length() > 0 {

				if i%2 == 0 {
					countyCode = a.Text()
					//获取省份的连接地址!!
					countyUrl, _ = a.Attr("href")
					countyUrl = parent + "/" + countyUrl
				} else {
					countyName = a.Text()
				}
			} else {
				//fmt.Println("没有找到a标签")
				//为空的话 说明没有下一级了
				if i%2 == 0 {
					countyCode = s.Text()
				} else {
					countyName = s.Text()
				}
			}

		})
		root.add(NewRegion(RegionType_County, countyName, countyCode, countyUrl))
	})
	return root.Children
}

func main() {

	url := "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2020/index.html"

	root := NewRegion(RegionType_Country, "中国", "", url)
	provinces := getProvinces(url, root)

	//fmt.Println(provinces)
	//遍历province  获取城市列表
	for _, province := range provinces {
		cities := getCities(province.Url, province)

		//获取县
		for _, city := range cities {
			getCounty(city.Url, city)
		}

	}

	//fmt.Println(root)

	bytes, _ := json.MarshalIndent(root, "", "\t")

	ioutil.WriteFile("data.json", bytes, os.ModePerm)

}

func getRoot(source string) string {
	return substring(source, 0, strings.LastIndex(source, "/"))
}

func substring(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	var substring = ""
	for i := start; i < end; i++ {
		substring += string(r[i])
	}

	return substring
}
