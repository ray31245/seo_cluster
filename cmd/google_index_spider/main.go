package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	domains := []string{
		"www.example1.com",
		"www.example2.com",
	}

	set := []func() (string, string){
		func() (string, string) {
			return "192.168.1.2", "192.168.1.30"
		},
	}

	requester, err := NewMultiIpRequester(set)
	if err != nil {
		log.Fatal(err)
	}

	total := len(domains)
	indexedNum := 0
	table := make([][]string, 0, total)

	for _, domain := range domains {
		isIndexed, err := requester.IsGoogleIndexed(ctx, domain)
		if err != nil {
			log.Fatal(err)
		}

		if isIndexed {
			indexedNum++
		}

		table = append(table, []string{domain, fmt.Sprintf("%t", isIndexed)})

		log.Printf("%s: %v", domain, isIndexed)
		time.Sleep(10 * time.Millisecond)
	}
	log.Printf("Total: %d, Indexed: %d", total, indexedNum)

	csvFile, err := os.OpenFile("isIndexed.csv", os.O_CREATE|os.O_WRONLY, 0o755)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	csvWriter.WriteAll(table)

	// ip1, err := requester.MyIp(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// ip1 = strings.Trim(ip1, " ")

	// log.Println(ip1)

	// count := 1
	// ip2 := ""
	// for ip2 != ip1 {
	// 	ip, err := requester.MyIp(ctx)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	ip = strings.Trim(ip, " ")
	// 	ip2 = ip

	// 	log.Printf("%s, %s, %v", ip, ip1, ip1 != ip)
	// 	count++
	// }
	// log.Println(count)
}

type IPRange struct {
	IPStart string
	IPEnd   string
}

type MultiIpRequester struct {
	IPRanges     []IPRange
	currentIndex int
	currentIP    int64
}

func NewMultiIpRequester(setIpRange []func() (start, end string)) (*MultiIpRequester, error) {
	if len(setIpRange) == 0 {
		return nil, fmt.Errorf("NewMultiIpRequester: empty ip range")
	}

	ipRanges := make([]IPRange, 0, len(setIpRange))

	for _, ipRange := range setIpRange {
		start, end := ipRange()
		if net.ParseIP(start) == nil || net.ParseIP(end) == nil {
			return nil, fmt.Errorf("NewMultiIpRequester: invalid ip range: %s, %s", start, end)
		}

		if InetAtoN(start) >= InetAtoN(end) {
			return nil, fmt.Errorf("NewMultiIpRequester: invalid ip range: %s-%s", start, end)
		}

		ipRanges = append(ipRanges, IPRange{IPStart: start, IPEnd: end})
	}

	return &MultiIpRequester{
		IPRanges:     ipRanges,
		currentIndex: 0,
		currentIP:    InetAtoN(ipRanges[0].IPStart),
	}, nil
}

func (m *MultiIpRequester) NextIP() string {
	if m.currentIP > InetAtoN(m.IPRanges[m.currentIndex].IPEnd) {
		m.currentIndex++
		if m.currentIndex >= len(m.IPRanges) {
			m.currentIndex = 0
		}

		m.currentIP = InetAtoN(m.IPRanges[m.currentIndex].IPStart)
	}

	ip := InetNtoA(m.currentIP)
	m.currentIP++

	return ip
}

func (m *MultiIpRequester) getClient() *http.Client {
	ip := m.NextIP()
	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP: net.ParseIP(ip),
		},
	}

	transport := &http.Transport{
		DialContext: dialer.DialContext,
	}

	return &http.Client{
		Transport: transport,
	}
}

func (m *MultiIpRequester) IsGoogleIndexed(ctx context.Context, domain string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.google.com/search?q=site:"+domain, nil)
	if err != nil {
		return false, fmt.Errorf("IsGoogleIndexed: new request error: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	c := m.getClient()

	resp, err := c.Do(req)
	if err != nil {
		return false, fmt.Errorf("IsGoogleIndexed: do request error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("IsGoogleIndexed: status code error: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return false, fmt.Errorf("IsGoogleIndexed: new document error: %w", err)
	}

	findNums := doc.Find(".sCuL3").Length()
	if findNums > 0 {
		return true, nil
	}

	return false, nil
}

func (m *MultiIpRequester) MyIp(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.ipify.org?format=json", nil)
	if err != nil {
		return "", fmt.Errorf("MyIp: new request error: %w", err)
	}

	c := m.getClient()

	resp, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("MyIp: do request error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("MyIp: status code error: %d", resp.StatusCode)
	}

	var ip struct {
		Ip string `json:"ip"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ip); err != nil {
		return "", fmt.Errorf("MyIp: decode error: %w", err)
	}

	return ip.Ip, nil
}

func InetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())

	return ret.Int64()
}
