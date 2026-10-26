package helper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/mozillazg/go-pinyin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// HandleMobile 处理手机号隐藏中间四位
func HandleMobile(mobile string) string {
	if len(mobile) <= 10 {
		return mobile
	}
	return mobile[:3] + "****" + mobile[len(mobile)-4:]
}

func FilterEmoji(content string) string {
	newContent := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			newContent += string(value)
		}
	}
	return newContent
}

// CleanUserInput removes invalid UTF-8 bytes and control characters while
// preserving tabs, newlines, carriage returns, and all other valid Unicode.
func CleanUserInput(text string) string {
	text = strings.ToValidUTF8(text, "")

	var builder strings.Builder
	builder.Grow(len(text))
	for _, r := range text {
		if unicode.IsControl(r) && r != '\t' && r != '\n' && r != '\r' {
			continue
		}
		builder.WriteRune(r)
	}

	return builder.String()
}

// RandStr 伪随机生成字符串
func RandStr(len uint) (string, error) {
	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%X", b), nil
}

// Md5V 生存MD5字符串
func Md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

const (
	GBK     string = "GBK"
	UTF8    string = "UTF8"
	UNKNOWN string = "UNKNOWN"
)

func ConvertToStringByCode(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func ConvertToStringGbkToUtf8(src string) string {
	if isUtf8([]byte(src)) {
		return src
	}
	var srcCode = "gbk"
	var tagCode = "utf-8"
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func ConvertToString(src string) string {
	if isUtf8([]byte(src)) {
		return src
	}
	if isGBK([]byte(src)) {
		utf8Data, err := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(src))
		if err != nil {
			return src
		}
		if isUtf8(utf8Data) {
			return string(utf8Data)
		} else {
			return ConvertToStringGbkToUtf8(string(utf8Data))
		}
	}
	sink, err := simplifiedchinese.GB18030.NewDecoder().Bytes([]byte(src))
	if err != nil {
		return src
	}
	return string(sink)
}

func isGBK(data []byte) bool {
	length := len(data)
	var i = 0
	for i < length {
		if data[i] <= 0x7f {
			// 编码0~127,只有一个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			// 大于127的使用双字节编码，落在gbk编码范围内的字符
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}

func preNUm(data byte) int {
	var mask byte = 0x80
	var num = 0
	// 8bit中首个0bit前有多少个1bits
	for i := 0; i < 8; i++ {
		if (data & mask) == mask {
			num++
			mask = mask >> 1
		} else {
			break
		}
	}
	return num
}

func isUtf8(data []byte) bool {
	i := 0
	for i < len(data) {
		if (data[i] & 0x80) == 0x00 {
			// 0XXX_XXXX
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			// 110X_XXXX 10XX_XXXX
			// 1110_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// preNUm() 返回首个字节的8个bits中首个0bit前面1bit的个数，该数量也是该字符所使用的字节数
			i++
			for j := 0; j < num-1; j++ {
				// 判断后面的 num - 1 个字节是不是都是10开头
				if (data[i] & 0xc0) != 0x80 {
					return false
				}
				i++
			}
		} else {
			// 其他情况说明不是utf-8
			return false
		}
	}
	return true
}

func GetStrCoding(data []byte) string {
	if isUtf8(data) == true {
		return UTF8
	} else if isGBK(data) == true {
		return GBK
	} else {
		return UNKNOWN
	}
}

// GetFileExt 根据文件路径获取文件后缀名（带.)
func GetFileExt(filename string) string {
	lastIndex := strings.LastIndex(filename, ".")
	if lastIndex == -1 {
		return ""
	}
	return `.` + filename[lastIndex+1:]
}

// CovertPinYin 将中文转为拼音
func CovertPinYin(a string) [][]string {
	args := pinyin.NewArgs()
	args.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}
	return pinyin.Pinyin(a, args)
}

// GetPinYinFirstLetter 获取拼音的首字母
func GetPinYinFirstLetter(pinyin [][]string) string {
	var firstLetter strings.Builder
	for i, s := range pinyin {
		if len(s) > 0 && i == 0 {
			firstLetter.WriteString(strings.ToLower(s[0][:1]))
		}
	}
	return firstLetter.String()
}

// GenerateRandomString 生成一个两位的随机字符串
func GenerateRandomString() string {
	// 定义字符集（字母和数字）
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())

	// 随机生成两位字符
	result := make([]byte, 2)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

// DesensitizeNickname 名称第一个字展示，后面为*
func DesensitizeNickname(fullName string) string {
	if fullName == "" {
		return ""
	}

	// Convert the string to a slice of runes to handle Unicode characters properly
	runes := []rune(fullName)

	// Check if there are any characters
	if len(runes) == 0 {
		return ""
	}

	// The first rune
	firstRune := runes[0]

	// Create a builder to construct the masked name
	var maskedName strings.Builder
	maskedName.WriteRune(firstRune) // Write the first character

	// Append the appropriate number of asterisks
	for i := 1; i < len(runes); i++ {
		maskedName.WriteRune('*')
	}

	return maskedName.String()
}

// DesensitizeAddress 地址的前11个展示，后面为*
func DesensitizeAddress(address string, sensitiveSize int) string {
	if address == "" {
		return ""
	}

	length := utf8.RuneCountInString(address)
	if sensitiveSize < 0 {
		// 处理敏感长度为负数的情况
		return address
	}

	if sensitiveSize >= length {
		// Mask the entire address if sensitiveSize is greater than or equal to the address length
		return strings.Repeat("*", length)
	}

	// Mask the part of the address
	// 获取每个字符的切割位置
	runes := []rune(address)
	maskedPart := strings.Repeat("*", sensitiveSize)
	return string(runes[:length-sensitiveSize]) + maskedPart
}

// TruncateString 支持中文的字符串截断函数
func TruncateString(s string, maxLen int) string {
	// 将字符串转换为rune切片，正确处理中文
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}

	// 截断并添加省略号
	return string(runes[:maxLen]) + "..."
}

// GetPinyin 获取拼音首字母（如 "北京" -> "bj"）
func GetPinyin(s string) string {
	args := pinyin.NewArgs()
	args.Style = pinyin.FirstLetter // 只取首字母
	pinyinSlice := pinyin.Pinyin(s, args)
	var result string
	for _, py := range pinyinSlice {
		if len(py) > 0 {
			result += py[0]
		}
	}
	return result
}

// GetPinyinAll  获取拼音全拼
func GetPinyinAll(s string) string {
	args := pinyin.NewArgs()
	args.Style = pinyin.Normal
	pinyinSlice := pinyin.Pinyin(s, args)
	var result string
	for _, py := range pinyinSlice {
		if len(py) > 0 {
			// 将第一个字符转换为大写
			first := unicode.ToUpper(rune(py[0][0]))
			// 剩余部分保持不变
			rest := py[0][1:]
			result += string(first) + rest
		}
	}
	return result
}

func RemoveTrailingAreaSuffix(s string) string {
	suffixes := []string{"特别行政区", "维吾尔自治区", "壮族自治区", "回族自治区", "自治区", "自治县", "新区", "市", "区", "县", "省"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return s[:len(s)-len(suffix)]
		}
	}
	return s
}

func StripHTML(input string) string {
	re := regexp.MustCompile("<[^>]*>")
	return re.ReplaceAllString(input, "")
}

// ConvertChineseToPinyin 将字符串中的中文转为拼音，其他字符保留原样
// 拼音之间用下划线连接
func ConvertChineseToPinyin(s string) string {
	if s == "" {
		return ""
	}
	var result strings.Builder
	var chineseBuffer strings.Builder

	args := pinyin.NewArgs()
	args.Style = pinyin.Normal // 不带声调的风格

	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			// 是中文字符，先缓存起来
			chineseBuffer.WriteRune(r)
		} else {
			// 不是中文，先处理缓存的中文
			if chineseBuffer.Len() > 0 {
				pinyinSlice := pinyin.Pinyin(chineseBuffer.String(), args)
				for i, words := range pinyinSlice {
					for _, word := range words {
						result.WriteString(word)
						if i < len(pinyinSlice)-1 {
							result.WriteString("_")
						}
					}
				}
				chineseBuffer.Reset()
			}
			// 添加当前非中文字符
			result.WriteRune(r)
		}
	}

	// 处理最后可能剩余的中文
	if chineseBuffer.Len() > 0 {
		pinyinSlice := pinyin.Pinyin(chineseBuffer.String(), args)
		for i, words := range pinyinSlice {
			for _, word := range words {
				result.WriteString(word)
				if i < len(pinyinSlice)-1 {
					result.WriteString("_")
				}
			}
		}
	}

	return result.String()
}
func FormatStoreVisitTimeSpan(language string, TimeSpan string) string {
	var timeSpanMap = map[string]string{
		"上午": "AM",
		"下午": "PM",
	}
	if language == "en" && (TimeSpan == "上午" || TimeSpan == "下午") {
		return timeSpanMap[TimeSpan]
	}
	return TimeSpan
}
