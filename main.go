package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	tt "github.com/huyoufu/go-timetracker"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func NewDocumentFromURL(url string) *goquery.Document {

	time.Sleep(time.Millisecond * 100)
	fmt.Printf("正在加载%s\n", url)
	client := &http.Client{}
	client.Timeout = time.Millisecond * 2000
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		fmt.Println("重定向过多!")
		return http.ErrUseLastResponse
	}
	var res *http.Response
	for {
		res, err = client.Do(request)
		if err != nil {
			log.Printf("重试吧~~~~~%s", err)
		} else {
			break
		}
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("url %s,status code error: %d %s", url, res.StatusCode, res.Status)
	}

	//reader := simplifiedchinese.GBK.NewDecoder().Reader(res.Body)
	//doc, err := goquery.NewDocumentFromReader(reader)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

//获取省份
func getProvinces(root *Region) []*Region {
	parent := root.Url
	if parent == "" {
		return nil
	}
	doc := NewDocumentFromURL(parent)

	parent = getRoot(parent)

	//解析出省份
	doc.Find(".provincetr ").Each(func(i int, s *goquery.Selection) {
		s.Find("td a").Each(func(i int, s *goquery.Selection) {
			provinceName := s.Text()
			//获取省份的连接地址!!
			provinceUrl, _ := s.Attr("href")
			provinceUrl = parent + "/" + provinceUrl

			root.add(NewRegion(RegionType_City, provinceName, "", "", provinceUrl))

		})
	})
	return root.Children
}
func getCities(root *Region) []*Region {
	parent := root.Url
	if parent == "" {
		return nil
	}
	doc := NewDocumentFromURL(parent)
	parent = getRoot(parent)
	//解析出省份
	doc.Find(".citytr ").Each(func(i int, s *goquery.Selection) {

		regionName := ""
		regionCode := ""
		regionUrl := ""

		s.Find("td a").Each(func(i int, s *goquery.Selection) {
			if i%2 == 0 {
				regionCode = s.Text()
				//获取省份的连接地址!!
				regionUrl, _ = s.Attr("href")
				regionUrl = parent + "/" + regionUrl
			} else {
				regionName = s.Text()
			}
		})
		root.add(NewRegion(RegionType_City, regionName, regionCode, "", regionUrl))
	})
	return root.Children
}
func getCounty(root *Region) []*Region {
	parent := root.Url
	if parent == "" {
		return nil
	}
	doc := NewDocumentFromURL(parent)
	parent = getRoot(parent)
	//解析县区
	doc.Find(".countytr ").Each(func(i int, s *goquery.Selection) {

		regionName := ""
		regionCode := ""
		regionUrl := ""

		s.Find("td ").Each(func(i int, s *goquery.Selection) {
			a := s.Find("a")
			if a.Length() > 0 {

				if i%2 == 0 {
					regionCode = a.Text()
					//获取县区的连接地址!!
					regionUrl, _ = a.Attr("href")
					regionUrl = parent + "/" + regionUrl
				} else {
					regionName = a.Text()
				}
			} else {
				//fmt.Println("没有找到a标签")
				//为空的话 说明没有下一级了
				if i%2 == 0 {
					regionCode = s.Text()
				} else {
					regionName = s.Text()
				}
			}

		})
		root.add(NewRegion(RegionType_County, regionName, regionCode, "", regionUrl))
	})
	return root.Children
}

//获取乡镇
func getTowns(root *Region) []*Region {
	parent := root.Url
	if parent == "" {
		return nil
	}
	doc := NewDocumentFromURL(parent)
	parent = getRoot(parent)
	//获取乡镇
	doc.Find(".towntr ").Each(func(i int, s *goquery.Selection) {

		regionName := ""
		regionCode := ""
		regionUrl := ""

		s.Find("td ").Each(func(i int, s *goquery.Selection) {
			a := s.Find("a")
			if a.Length() > 0 {

				if i%2 == 0 {
					regionCode = a.Text()
					//获取乡镇
					regionUrl, _ = a.Attr("href")
					regionUrl = parent + "/" + regionUrl
				} else {
					regionName = a.Text()
				}
			} else {
				//fmt.Println("没有找到a标签")
				//为空的话 说明没有下一级了
				if i%2 == 0 {
					regionCode = s.Text()
				} else {
					regionName = s.Text()
				}
			}

		})
		root.add(NewRegion(RegionType_Town, regionName, regionCode, "", regionUrl))
	})
	return root.Children
}

//获取村委会
func getVillages(root *Region) []*Region {
	parent := root.Url
	if parent == "" {
		return nil
	}
	doc := NewDocumentFromURL(parent)
	parent = getRoot(parent)
	//获取乡镇
	doc.Find(".villagetr ").Each(func(i int, s *goquery.Selection) {

		regionName := ""
		regionCode := ""
		regionClass := ""
		regionUrl := ""

		s.Find("td ").Each(func(i int, s *goquery.Selection) {

			if i%3 == 0 {
				regionCode = s.Text()
			} else if i%3 == 1 {
				regionClass = s.Text()
			} else {
				regionName = s.Text()
			}

		})
		root.add(NewRegion(RegionType_Village, regionName, regionCode, regionClass, regionUrl))
	})
	return root.Children
}

func main() {
	tracker := tt.NewTimeTracker("开始下载中国区域信息数据")

	url := "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2021/index.html"

	root := NewRegion(RegionType_Country, "中国", "", "", url)
	provinces := getProvinces(root)

	//fmt.Println(provinces)
	//遍历province  获取城市列表
	for _, province := range provinces {
		cities := getCities(province)
		//getCities(province)

		//获取县
		for _, city := range cities {

			counties := getCounty(city)
			//获取乡镇列表
			for _, county := range counties {
				towns := getTowns(county)

				//获取村庄列表
				for _, town := range towns {
					getVillages(town)
				}

			}
		}

	}

	//fmt.Println(root)

	bytes, _ := json.MarshalIndent(root, "", "\t")

	ioutil.WriteFile("data1.json", bytes, os.ModePerm)

	tracker.Close()
	tracker.PrintBeautiful()

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
