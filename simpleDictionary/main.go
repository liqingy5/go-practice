package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type DictRequest struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
	UserID    string `json:"user_id"`
}

type DictResponse struct {
	Rc   int `json:"rc"`
	Wiki struct {
		KnownInLaguages int `json:"known_in_laguages"`
		Description     struct {
			Source string      `json:"source"`
			Target interface{} `json:"target"`
		} `json:"description"`
		ID   string `json:"id"`
		Item struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"item"`
		ImageURL  string `json:"image_url"`
		IsSubject string `json:"is_subject"`
		Sitelink  string `json:"sitelink"`
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

type BaiduResponse struct {
	TransResult struct {
		Data []struct {
			Dst        string          `json:"dst"`
			PrefixWrap int             `json:"prefixWrap"`
			Result     [][]interface{} `json:"result"`
			Src        string          `json:"src"`
		} `json:"data"`
		From     string `json:"from"`
		Status   int    `json:"status"`
		To       string `json:"to"`
		Type     int    `json:"type"`
		Phonetic []struct {
			SrcStr string `json:"src_str"`
			TrgStr string `json:"trg_str"`
		} `json:"phonetic"`
	} `json:"trans_result"`
	DictResult struct {
		Edict struct {
			Item []struct {
				TrGroup []struct {
					Tr          []string `json:"tr"`
					Example     []string `json:"example"`
					SimilarWord []string `json:"similar_word"`
				} `json:"tr_group"`
				Pos string `json:"pos"`
			} `json:"item"`
			Word string `json:"word"`
		} `json:"edict"`
		Collins struct {
			Entry []struct {
				EntryID string `json:"entry_id"`
				Type    string `json:"type"`
				Value   []struct {
					MeanType []struct {
						InfoType string `json:"info_type"`
						InfoID   string `json:"info_id"`
						Example  []struct {
							ExampleID string `json:"example_id"`
							TtsSize   string `json:"tts_size"`
							Tran      string `json:"tran"`
							Ex        string `json:"ex"`
							TtsMp3    string `json:"tts_mp3"`
						} `json:"example,omitempty"`
						Posc []struct {
							Tran    string `json:"tran"`
							PoscID  string `json:"posc_id"`
							Example []struct {
								ExampleID string `json:"example_id"`
								Tran      string `json:"tran"`
								Ex        string `json:"ex"`
								TtsMp3    string `json:"tts_mp3"`
							} `json:"example"`
							Def string `json:"def"`
						} `json:"posc,omitempty"`
					} `json:"mean_type"`
					Gramarinfo []struct {
						Tran  string `json:"tran"`
						Type  string `json:"type"`
						Label string `json:"label"`
					} `json:"gramarinfo"`
					Tran   string `json:"tran"`
					Def    string `json:"def"`
					MeanID string `json:"mean_id"`
					Posp   []struct {
						Label string `json:"label"`
					} `json:"posp"`
				} `json:"value"`
			} `json:"entry"`
			WordName      string `json:"word_name"`
			Frequence     string `json:"frequence"`
			WordEmphasize string `json:"word_emphasize"`
			WordID        string `json:"word_id"`
		} `json:"collins"`
		From        string `json:"from"`
		SimpleMeans struct {
			WordName  string   `json:"word_name"`
			From      string   `json:"from"`
			WordMeans []string `json:"word_means"`
			Exchange  struct {
				WordPl []string `json:"word_pl"`
			} `json:"exchange"`
			Tags struct {
				Core  []string `json:"core"`
				Other []string `json:"other"`
			} `json:"tags"`
			Symbols []struct {
				PhEn  string `json:"ph_en"`
				PhAm  string `json:"ph_am"`
				Parts []struct {
					Part  string   `json:"part"`
					Means []string `json:"means"`
				} `json:"parts"`
				PhOther string `json:"ph_other"`
			} `json:"symbols"`
		} `json:"simple_means"`
		Lang   string `json:"lang"`
		Oxford struct {
			Entry []struct {
				Tag  string `json:"tag"`
				Name string `json:"name"`
				Data []struct {
					Tag  string `json:"tag"`
					Data []struct {
						Tag  string `json:"tag"`
						Data []struct {
							Tag  string `json:"tag"`
							Data []struct {
								Tag  string `json:"tag"`
								Data []struct {
									Tag    string `json:"tag"`
									EnText string `json:"enText,omitempty"`
									ChText string `json:"chText,omitempty"`
									G      string `json:"g,omitempty"`
									Data   []struct {
										Text      string `json:"text"`
										HoverText string `json:"hoverText"`
									} `json:"data,omitempty"`
								} `json:"data"`
							} `json:"data"`
						} `json:"data,omitempty"`
						P     string `json:"p,omitempty"`
						PText string `json:"p_text,omitempty"`
						N     string `json:"n,omitempty"`
						Xt    string `json:"xt,omitempty"`
					} `json:"data"`
				} `json:"data"`
			} `json:"entry"`
			Unbox []struct {
				Tag  string `json:"tag"`
				Type string `json:"type"`
				Name string `json:"name"`
				Data []struct {
					Tag     string `json:"tag"`
					Text    string `json:"text,omitempty"`
					Words   string `json:"words,omitempty"`
					Outdent string `json:"outdent,omitempty"`
					Data    []struct {
						Tag    string `json:"tag"`
						EnText string `json:"enText"`
						ChText string `json:"chText"`
					} `json:"data,omitempty"`
				} `json:"data"`
			} `json:"unbox"`
		} `json:"oxford"`
		BaiduPhrase []struct {
			Tit   []string `json:"tit"`
			Trans []string `json:"trans"`
		} `json:"baidu_phrase"`
		QueryExplainVideo struct {
			ID           int    `json:"id"`
			UserID       string `json:"user_id"`
			UserName     string `json:"user_name"`
			UserPic      string `json:"user_pic"`
			Query        string `json:"query"`
			Direction    string `json:"direction"`
			Type         string `json:"type"`
			Tag          string `json:"tag"`
			Detail       string `json:"detail"`
			Status       string `json:"status"`
			SearchType   string `json:"search_type"`
			FeedURL      string `json:"feed_url"`
			Likes        string `json:"likes"`
			Plays        string `json:"plays"`
			CreatedAt    string `json:"created_at"`
			UpdatedAt    string `json:"updated_at"`
			DuplicateID  string `json:"duplicate_id"`
			RejectReason string `json:"reject_reason"`
			CoverURL     string `json:"coverUrl"`
			VideoURL     string `json:"videoUrl"`
			ThumbURL     string `json:"thumbUrl"`
			VideoTime    string `json:"videoTime"`
			VideoType    string `json:"videoType"`
		} `json:"queryExplainVideo"`
	} `json:"dict_result"`
	LijuResult struct {
		Double string   `json:"double"`
		Tag    []string `json:"tag"`
		Single string   `json:"single"`
	} `json:"liju_result"`
	Logid int `json:"logid"`
}

func queryCaiyun(word string, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{}
	request := DictRequest{TransType: "en2zh", Source: word}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("os-version", "")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	req.Header.Set("app-name", "xy")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("device-id", "")
	req.Header.Set("os-type", "web")
	req.Header.Set("X-Authorization", "token:qgemv4jr1y38jyq6vhvi")
	req.Header.Set("Origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "_ym_uid=16456948721020430059; _ym_d=1645694872")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}
	var dictResponse DictResponse
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("caiyun")
	fmt.Println(word, "UK:", dictResponse.Dictionary.Prons.En, "US:", dictResponse.Dictionary.Prons.EnUs)
	for _, item := range dictResponse.Dictionary.Explanations {
		fmt.Println(item)
	}
}

func queryBaidu(word string, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{}
	var data = strings.NewReader(`from=en&to=zh&query=` + word + `&transtype=realtime&simple_means_flag=3&sign=54706.276099&token=f03f3ff61fc27157bcf6b3b7050bd6cf&domain=common`)
	req, err := http.NewRequest("POST", "https://fanyi.baidu.com/v2transapi?from=en&to=zh", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	req.Header.Set("Acs-Token", "1673769896714_1673829283662_/jNTmMdug05eaCQYNoVx12H//AUDvT1OVFIc1YEyP35x5FqBhdI1kuYKANBm/la3C+OXMEJUUb6VTnUeapUl+WXq9MjdV6J9/DlVJ0judQVWZCq/D7FDCatiFVp6520RnLlIa4nXdEvGkwCuIqItTVzTJnVP5+6c2abzI8OzRlKHemqIHawSr95ZGTxoy3/jozpLcJpOzYa6EDlwts6xjGpWgbxUtlCteWModF5JCIA4AoGrRSooZhZYnrDwjjigU++aUEq4BKpJmGatXYCCYZAik2njsFDPSR4KnuMfigM=")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", "PSTM=1671158464; BIDUPSID=9DE284EE5EA6C5CD5DECB0DAA5563A52; BAIDUID=2F957799E79AEE928D3E809E6D1EE39C:FG=1; BAIDUID_BFESS=2F957799E79AEE928D3E809E6D1EE39C:FG=1; ZFY=sXQlb4DhbjDUWDBjcc3eSCAnTZpCsgRlOUfU:A8zkNZg:C; BDUSS=N1YTJIczE5Nmx3RXhYRERZeTZ4b0x1WDRGdjRZbkdyWktHNGtmZHJZUks5ZUpqSVFBQUFBJCQAAAAAAAAAAAEAAAC0wNEXTFFZTklDSwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEpou2NKaLtjUF; BDUSS_BFESS=N1YTJIczE5Nmx3RXhYRERZeTZ4b0x1WDRGdjRZbkdyWktHNGtmZHJZUks5ZUpqSVFBQUFBJCQAAAAAAAAAAAEAAAC0wNEXTFFZTklDSwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEpou2NKaLtjUF; BDRCVFR[feWj1Vr5u3D]=I67x6TjHwwYf0; PSINO=7; BA_HECTOR=2g010g20250h2g05848l2gmr1hs8ccl1k; delPer=0; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; H_PS_PSSID=36552_37971_37647_37552_37512_38018_36920_38034_37990_37935_37900_26350_37881; ZD_ENTRY=google; __bid_n=1843ba4bb98af652e54207; FPTOKEN=wk5mCfUe8IGYSFTc73kb4m6fD/+vPsUudAJDqnD/1gWnWzNd99f1xJPgasKG9KPDa2jfNZiYVKRz/kfXm1MWknI0kq3YZ8OfCPlBzRFVsXvqzMs1sgSq8yZTiNZ/wNe3dqBI/C8m9GONvew4EUEGnWpxEWWAdJwDfFtMlL8IPqxiI0LGW6cw5BJL7MXfhqu8S56gTCIH/sd/bTzJUAyJNgs5XNhDNTbQyBfUvbEwKAxk6+JiG5TfvcqIlQO3l+4Tx2MSePwVzTbUJEq9j5d5RHFn0Lf1nBk8clGF23YtXGPJJ4btjSmSUUd49x9UW1d/5LfqBCLINwhnfTjRD62tBvcvQPrcsQ2O1LDhtLaPoAT7wa8Eaf3t+BoKvb1Lg+xDOQbuGqkhRbXLyFnoclglNw==|nNnTDKo31nfW33fjkn2tPkGoR/daDBVvuL41C5ICxoE=|10|aa0cc876f4e34ae85c70d1fcd84bfc6c; APPGUIDE_10_0_2=1; REALTIME_TRANS_SWITCH=1; FANYI_WORD_SWITCH=1; HISTORY_SWITCH=1; SOUND_SPD_SWITCH=1; SOUND_PREFER_SWITCH=1; Hm_lvt_64ecd82404c51e03dc91cb9e8c025574=1673829065; Hm_lpvt_64ecd82404c51e03dc91cb9e8c025574=1673829263; ab_sr=1.0.1_ZDgzZjYyMGM0M2YzY2VkYzQ3NTMxZjk5ZWU1NThjMGZiN2MyNzFjYmQyOGE3OTRlNjgxZGY3NWVhNjE4MGJhZDk1ZDE5NjUxZWEwNzRjMzRhMTEyYTcwYzVhZmNkOTg2OTc0ZDEyNzNiMTc5NGZlNGU5MzdjMGE5ZDM3YjRlMWU3ZGY1MzY5OTZmODAxZjA4OTE0NWY1NTQ4NGVhNDY4NjEyMzY5ZWYxNTBjZWMwNDkwYzUyZDIwM2I3YWE1MTIx")
	req.Header.Set("Origin", "https://fanyi.baidu.com")
	req.Header.Set("Referer", "https://fanyi.baidu.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua", `"Not?A_Brand";v="8", "Chromium";v="108", "Google Chrome";v="108"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var baidurespon BaiduResponse
	err = json.Unmarshal(bodyText, &baidurespon)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Baidu:")
	fmt.Println(word, "UK:", baidurespon.DictResult.SimpleMeans.Symbols[0].PhEn, "US:", baidurespon.DictResult.SimpleMeans.Symbols[0].PhAm)
	for _, part := range baidurespon.DictResult.SimpleMeans.WordMeans {
		fmt.Println(part)
	}

}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage: simpleDict WORD
example: simpleDict hello
		`)
		os.Exit(1)
	}
	word := os.Args[1]
	wg := sync.WaitGroup{}
	wg.Add(2)
	go queryBaidu(word, &wg)
	go queryCaiyun(word, &wg)
	wg.Wait()
}
