package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Municipality struct {
	Pref   string
	PrefEn string
	Name   string
	NameEn string
	URL    string
}

func main() {
	fileJp, err := os.Open("KEN_ALL_ROME.CSV")
	if err != nil {
		log.Fatal(err)
	}
	defer fileJp.Close()

	fileC4F, err := os.Open("localgovjp-sjis.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer fileC4F.Close()

	outputSjisCRLF, err := os.Create("localgov_sjis_crlf.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer outputSjisCRLF.Close()

	outputUtf8CRLF, err := os.Create("localgov_utf8_crlf.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer outputUtf8CRLF.Close()

	outputSjisLF, err := os.Create("localgov_sjis_lf.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer outputSjisLF.Close()

	outputUtf8LF, err := os.Create("localgov_utf8_lf.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer outputUtf8LF.Close()

	readerJpPost := csv.NewReader(transform.NewReader(fileJp, japanese.ShiftJIS.NewDecoder()))
	readerJpPost.LazyQuotes = true

	readerForC4F := csv.NewReader(transform.NewReader(fileC4F, japanese.ShiftJIS.NewDecoder()))
	readerForC4F.LazyQuotes = true

	writerSjisCRLF := csv.NewWriter(transform.NewWriter(outputSjisCRLF, japanese.ShiftJIS.NewEncoder()))
	writerSjisCRLF.UseCRLF = true

	writerUtf8CRLF := csv.NewWriter(outputUtf8CRLF)
	writerUtf8CRLF.UseCRLF = true

	writerSjisLF := csv.NewWriter(transform.NewWriter(outputSjisLF, japanese.ShiftJIS.NewEncoder()))
	writerSjisLF.UseCRLF = false

	writerUtf8LF := csv.NewWriter(outputUtf8LF)
	writerUtf8LF.UseCRLF = false

	cityToURL := make(map[string]string)
	for {
		rawRecord, err := readerForC4F.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		var record []string
		for i, v := range rawRecord {
			if i > 0 {
				record = append(record, v)
			}
		}

		key := record[0] + record[2]
		cityToURL[key] = record[6]
	}

	// code4fukui/localgovjp のデータと市区町村名が完全にマッチしない自治体のURLを追加
	// 宮城県七ヶ宿町
	// https://town.shichikashuku.miyagi.jp/
	// 兵庫県丹波篠山市
	// https://www.city.tambasasayama.lg.jp/
	// 福岡県那珂川市
	// https://www.city.nakagawa.lg.jp/
	// 福岡県須惠町
	// https://www.town.sue.fukuoka.jp/
	cityToURL["宮城県七ヶ宿町"] = "https://town.shichikashuku.miyagi.jp/"
	cityToURL["兵庫県丹波篠山市"] = "https://www.city.tambasasayama.lg.jp/"
	cityToURL["福岡県那珂川市"] = "https://www.city.nakagawa.lg.jp/"
	cityToURL["福岡県須惠町"] = "https://www.town.sue.fukuoka.jp/"

	var municipalities []Municipality

	for {
		rawRecord, err := readerJpPost.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		var record []string
		for i, v := range rawRecord {
			if i > 0 {
				record = append(record, v)
			}
		}

		if strings.Contains(record[1], "　") {
			s := strings.Split(record[1], "　")
			if strings.Contains(record[1], "区") {
				record[1] = s[len(s)-2]
			} else {
				record[1] = s[len(s)-1]
			}
		}

		if strings.Contains(record[4], " ") {
			s := strings.Split(record[4], " ")
			if len(s) == 2 {
				// 〜市/町/村/区(東京都)
				record[4] = strings.Title(strings.ToLower(s[0])) + strings.ToLower("-"+s[1])
			} else if len(s) == 4 {
				// 〜郡 〜市/町/村 or 〜市/町/村 〜区
				if s[1] == "GUN" {
					record[4] = strings.Title(strings.ToLower(s[2])) + strings.ToLower("-"+s[3])
				} else {
					record[4] = strings.Title(strings.ToLower(s[0])) + strings.ToLower("-"+s[1])
				}
			} else {
				// 例外パターン
				if record[4] == "MIYAKEJIMA MIYAKE MURA" {
					record[4] = "Miyake-mura"
				} else if record[4] == "HACHIJOJIMA HACHIJO MACHI" {
					record[4] = "Hachijo-machi"
				} else if record[4] == "NISHIYATSUSHIRO GUN ICHIKAWAMISATO" {
					record[4] = "Ichikawamisato-cho"
				}
			}
		}

		if record[3] == "HOKKAIDO" {
			record[3] = "Hokkaido"
		} else {
			s := strings.Split(record[3], " ")
			record[3] = strings.Title(strings.ToLower(s[0]))
			// record[3] = strings.Title(strings.ToLower(s[0])) + "-" + strings.ToLower(s[1])
		}

		dup := false

		for _, m := range municipalities {
			if m.Pref == record[0] && m.Name == record[1] {
				dup = true
				break
			}
		}

		if !dup {
			municipalities = append(municipalities, Municipality{
				Pref:   record[0],
				PrefEn: record[3],
				Name:   record[1],
				NameEn: record[4],
				URL:    cityToURL[record[0]+record[1]],
			})
			if err := writerSjisCRLF.Write([]string{record[0], record[3], record[1], record[4], cityToURL[record[0]+record[1]]}); err != nil {
				fmt.Println(err)
			}
			if err := writerUtf8CRLF.Write([]string{record[0], record[3], record[1], record[4], cityToURL[record[0]+record[1]]}); err != nil {
				fmt.Println(err)
			}
			if err := writerSjisLF.Write([]string{record[0], record[3], record[1], record[4], cityToURL[record[0]+record[1]]}); err != nil {
				fmt.Println(err)
			}
			if err := writerUtf8LF.Write([]string{record[0], record[3], record[1], record[4], cityToURL[record[0]+record[1]]}); err != nil {
				fmt.Println(err)
			}
		}
	}

	writerSjisCRLF.Flush()
	writerUtf8CRLF.Flush()
	writerSjisLF.Flush()
	writerUtf8LF.Flush()
}
