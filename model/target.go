package model

import (
	"afkl/fumofinger/utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

type ResponseForCapture struct {
	Url           string
	Body          []byte
	HeadersString string
	HeadersMap    http.Header
	Cert          []byte
}

type TargetList struct {
	List []string
	Len  int
}

func isTarget(target string) bool {
	u, err := url.Parse(target)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (list *TargetList) LoadTargetFromString(target string) {
	list.Len += 1
	list.List = append(list.List, target)
}

func (list *TargetList) LoadTargetFromFile(filename string) (len int, err error) {
	fp, err := utils.Open(filename)
	if err != nil {
		return 0, err
	}
	defer fp.Close()

	for {
		target, isNotEOF := fp.ReadLine()
		if isNotEOF {
			if isTarget(target) {
				list.LoadTargetFromString(target)
			}
		} else {
			break
		}
	}

	return list.Len, nil
}

func (list TargetList) Length() int {
	return list.Len
}

func (list TargetList) SplitTargetList(blockNum int) [][]string {
	max := list.Len
	if max < blockNum {
		blockNum = max
	}
	var segmens = make([][]string, 0)
	quantity := max / blockNum
	end := 0
	for i := 1; i <= blockNum; i++ {
		qu := i * quantity
		if i != blockNum {
			segmens = append(segmens, list.List[i-1+end:qu])
		} else {
			segmens = append(segmens, list.List[i-1+end:])
		}
		end = qu - i
	}
	return segmens
}

func RequestTarget(url string) (*ResponseForCapture, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	headers := ""
	for key, values := range resp.Header {
		tmp := ""
		for _, value := range values {
			tmp += value
		}
		headers += fmt.Sprintf("%s: %s\r\n", key, tmp)
	}

	return &ResponseForCapture{
		Url:           url,
		Body:          body,
		HeadersString: headers,
		HeadersMap:    resp.Header,
		Cert:          utils.GetCerts(resp),
	}, nil
}

func ReMatchBody(resp ResponseForCapture, re string) []string {
	compileRE := regexp.MustCompile(re)
	return compileRE.FindAllString(string(resp.Body), -1)
}
