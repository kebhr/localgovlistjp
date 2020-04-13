# localgovlistjp
## About
全国1,741市区町村（市町村・特別区）のCSVファイルです。  
以下の項目を含んでいます。
- 都道府県名（日本語）
- 都道府県名（ローマ字）
- 自治体名（日本語）
- 自治体名（ローマ字）
- WebサイトURL

ファイルは文字コード・改行コード別に4つに分かれています。

|ファイル名|文字コード|改行コード|
|---|---|---|
|`localgov_sjis_crlf.csv`|Shift-JIS|CRLF|
|`localgov_sjis_lf.csv`|Shift-JIS|LF|
|`localgov_utf8_crlf.csv`|UTF-8|CRLF|
|`localgov_utf8_lf.csv`|UTF-8|LF|

## Copyright
`KEN_ALL_ROME.CSV` は日本郵便が配布している郵便番号データです。  
都道府県名（日本語・英語）・自治体名（日本語/英語）のデータ元として使用しています。  
再配布が許可されているファイルです。[出典](https://www.post.japanpost.jp/zipcode/dl/readme_ro.html)  
[ダウンロード先](https://www.post.japanpost.jp/zipcode/dl/roman-zip.html)

`localgovjp-sjis.csv` は [code4fukui/localgovjp](https://github.com/code4fukui/localgovjp) にてCC0で配布されています。  
WebサイトURLのデータ元として使用しています。

## License
`KEN_ALL_ROME.CSV`と`localgovjp-sjis.csv`は配布元のライセンスに従います。  
これらを除くファイルのライセンスは[CC0](https://github.com/kebhr/localgovlistjp/blob/master/LICENSE)です。

## File Generate
CSVファイルの生成には `main.go` を使用しています。  
`$ go run main.go`
