package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()

	uri := "https://www.binance.com/bapi/composite/v3/friendly/pgc/content/article/list?pageIndex=1&pageSize=20&type=1"

	// reqBody := []byte(`{"pageIndex":1,"pageSize":20,"scene":"web-homepage","contentIds":[]}`)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	// req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewReader(reqBody))
	if err != nil {
		log.Fatalln(err)
	}

	// req.Header.Set("User-Agent", "Mozilla/5.0")
	// req.Header.Set("accept", "*/*")
	// req.Header.Set("accept-language", "zh-TW,zh;q=0.9")
	// req.Header.Set("bnc-location", "")
	// req.Header.Set("bnc-uuid", "5b474df5-8e38-42e1-bf8c-48ab09851d0b")
	req.Header.Set("clienttype", "web")
	req.Header.Set("content-type", "application/json")
	// req.Header.Set("cookie", "theme=dark; bnc-uuid=5b474df5-8e38-42e1-bf8c-48ab09851d0b; source=organic; campaign=www.google.com; changeBasisTimeZone=; BNC_FV_KEY=33d6767e1ca9a981a7d584cb10d54694a6b6cde5; userPreferredCurrency=USD_USD; BNC_FV_KEY_T=101-DLaTWrOTv%2FUJFLDTB2Ca3RLAa8m1FQ5LtPM9gJaaDfJu1zpT7KH6pkoFOi8ikyxS0u6s72NaYaQlXSmcTTyMLQ%3D%3D-uPEonGDrH3syXnwon656IA%3D%3D-60; BNC_FV_KEY_EXPIRE=1727702860438; _gid=GA1.2.1220371624.1727681262; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%221912681d64a581-08252742a6bbcd8-19525637-2073600-1912681d64b1676%22%2C%22first_id%22%3A%22%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%2C%22%24latest_referrer%22%3A%22%22%7D%2C%22identities%22%3A%22eyIkaWRlbnRpdHlfY29va2llX2lkIjoiMTkxMjY4MWQ2NGE1ODEtMDgyNTI3NDJhNmJiY2Q4LTE5NTI1NjM3LTIwNzM2MDAtMTkxMjY4MWQ2NGIxNjc2In0%3D%22%2C%22history_login_id%22%3A%7B%22name%22%3A%22%22%2C%22value%22%3A%22%22%7D%2C%22%24device_id%22%3A%22192315802701fa-010e77d15d7afa7-19525637-1484784-192315802712dfd%22%7D; _gat_UA-162512367-1=1; _gat=1; _gcl_au=1.1.1628456763.1727681699; _ga=GA1.2.67312340.1722927928; OptanonConsent=isGpcEnabled=0&datestamp=Mon+Sep+30+2024+15%3A34%3A59+GMT%2B0800+(%E5%8F%B0%E5%8C%97%E6%A8%99%E6%BA%96%E6%99%82%E9%96%93)&version=202407.2.0&browserGpcFlag=0&isIABGlobal=false&hosts=&consentId=de6a662c-aedc-4b1a-92e0-6b2af2e7ec1e&interactionCount=0&isAnonUser=1&landingPath=NotLandingPage&groups=C0001%3A1%2CC0003%3A1%2CC0004%3A1%2CC0002%3A1&AwaitingReconsent=false; _uetsid=7f9108807efe11ef8bf47798310bae6e; _uetvid=7f9104e07efe11efaf1df3225961fdee; lang=zh-cn; _ga_3WP50LGEEC=GS1.1.1727681262.5.1.1727681700.3.0.0")
	// req.Header.Set("csrftoken", "d41d8cd98f00b204e9800998ecf8427e")
	// req.Header.Set("device-info", "eyJzY3JlZW5fcmVzb2x1dGlvbiI6IjE5MjAsMTA4MCIsImF2YWlsYWJsZV9zY3JlZW5fcmVzb2x1dGlvbiI6IjE5MjAsMTA4MCIsInN5c3RlbV92ZXJzaW9uIjoiTWFjIE9TIDEwLjE1LjciLCJicmFuZF9tb2RlbCI6InVua25vd24iLCJzeXN0ZW1fbGFuZyI6InpoLVRXIiwidGltZXpvbmUiOiJHTVQrMDg6MDAiLCJ0aW1lem9uZU9mZnNldCI6LTQ4MCwidXNlcl9hZ2VudCI6Ik1vemlsbGEvNS4wIChNYWNpbnRvc2g7IEludGVsIE1hYyBPUyBYIDEwXzE1XzcpIEFwcGxlV2ViS2l0LzUzNy4zNiAoS0hUTUwsIGxpa2UgR2Vja28pIENocm9tZS8xMjYuMC4wLjAgU2FmYXJpLzUzNy4zNiIsImxpc3RfcGx1Z2luIjoiUERGIFZpZXdlcixDaHJvbWUgUERGIFZpZXdlcixDaHJvbWl1bSBQREYgVmlld2VyLE1pY3Jvc29mdCBFZGdlIFBERiBWaWV3ZXIsV2ViS2l0IGJ1aWx0LWluIFBERiIsImNhbnZhc19jb2RlIjoiMzNjMmE3ODYiLCJ3ZWJnbF92ZW5kb3IiOiJHb29nbGUgSW5jLiAoQXBwbGUpIiwid2ViZ2xfcmVuZGVyZXIiOiJBTkdMRSAoQXBwbGUsIEFOR0xFIE1ldGFsIFJlbmRlcmVyOiBBcHBsZSBNMyBQcm8sIFVuc3BlY2lmaWVkIFZlcnNpb24pIiwiYXVkaW8iOiIxMjQuMDQzNDY2MDcxMTQ3MTIiLCJwbGF0Zm9ybSI6Ik1hY0ludGVsIiwid2ViX3RpbWV6b25lIjoiQXNpYS9UYWlwZWkiLCJkZXZpY2VfbmFtZSI6IkNocm9tZSBWMTI2LjAuMC4wIChNYWMgT1MpIiwiZmluZ2VycHJpbnQiOiIxZmExNmE2ZjA5YzFmODcxYWZkNjVkNzhhMTJiNzhjOCIsImRldmljZV9pZCI6IiIsInJlbGF0ZWRfZGV2aWNlX2lkcyI6IiJ9")
	// req.Header.Set("fvideo-id", "33d6767e1ca9a981a7d584cb10d54694a6b6cde5")
	// req.Header.Set("fvideo-token", "l6klbfZhDRdQAmZ9wMHfQh1YFTWqzcfEIo5qg3tKQAXexyjkLlwCHUac2aeMR5Z4vJ5ojbHqrL+nkXR+I0x5454N4wqdnI8jQqlhA7wb9CeXCCFK/47TSZQP7VBiRbiBBLzzREprJ/9OC9JpWxupQ6lyM8hxCzj6RUW4wmnx1YkNggp75L+EO8A6mv5asCId8=05")

	// req.Header.Set("lang", "zh-CN")
	req.Header.Set("lang", "en")

	// req.Header.Set("origin", "https://www.binance.com")
	// req.Header.Set("priority", "u=1, i")
	// req.Header.Set("referer", "https://www.binance.com/zh-CN/square")
	// req.Header.Set("sec-ch-ua", "Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\"")
	// req.Header.Set("sec-ch-ua-mobile", "?0")
	// req.Header.Set("sec-ch-ua-platform", "macOS")
	// req.Header.Set("sec-fetch-dest", "empty")
	// req.Header.Set("sec-fetch-mode", "cors")
	// req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("versioncode", "web")
	// req.Header.Set("x-passthrough-token", "")
	// req.Header.Set("x-trace-id", "02536a65-48cb-4d86-8036-a5f94151fcfb")
	// req.Header.Set("x-ui-request-trace", "02536a65-48cb-4d86-8036-a5f94151fcfb")

	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	// do something with the response
	log.Println(resp.Status)

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(resBody))
	f, err := os.Create("test.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	err = os.WriteFile("test.json", resBody, 0o644)
	if err != nil {
		log.Fatalln(err)
	}
}
