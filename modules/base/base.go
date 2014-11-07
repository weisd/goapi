// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	// "crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	// "html/template"
	"math"
	"regexp"
	"strings"
	"time"
	"unicode"

	// "github.com/Unknwon/com"
	// "github.com/Unknwon/i18n"

	// "github.com/gogits/gogs/modules/setting"
)

// Encode string to md5 hex value
func EncodeMd5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

// GetRandomString generate random string by specify chars.
func GetRandomString(n int, alphabets ...byte) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		if len(alphabets) == 0 {
			bytes[i] = alphanum[b%byte(len(alphanum))]
		} else {
			bytes[i] = alphabets[b%byte(len(alphabets))]
		}
	}
	return string(bytes)
}

// http://code.google.com/p/go/source/browse/pbkdf2/pbkdf2.go?repo=crypto
func PBKDF2(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}

// verify time limit code
// func VerifyTimeLimitCode(data string, minutes int, code string) bool {
// 	if len(code) <= 18 {
// 		return false
// 	}

// 	// split code
// 	start := code[:12]
// 	lives := code[12:18]
// 	if d, err := com.StrTo(lives).Int(); err == nil {
// 		minutes = d
// 	}

// 	// right active code
// 	retCode := CreateTimeLimitCode(data, minutes, start)
// 	if retCode == code && minutes > 0 {
// 		// check time is expired or not
// 		before, _ := DateParse(start, "YmdHi")
// 		now := time.Now()
// 		if before.Add(time.Minute*time.Duration(minutes)).Unix() > now.Unix() {
// 			return true
// 		}
// 	}

// 	return false
// }

const TimeLimitCodeLength = 12 + 6 + 40

// create a time limit code
// code format: 12 length date time string + 6 minutes string + 40 sha1 encoded string
// func CreateTimeLimitCode(data string, minutes int, startInf interface{}) string {
// 	format := "YmdHi"

// 	var start, end time.Time
// 	var startStr, endStr string

// 	if startInf == nil {
// 		// Use now time create code
// 		start = time.Now()
// 		startStr = DateFormat(start, format)
// 	} else {
// 		// use start string create code
// 		startStr = startInf.(string)
// 		start, _ = DateParse(startStr, format)
// 		startStr = DateFormat(start, format)
// 	}

// 	end = start.Add(time.Minute * time.Duration(minutes))
// 	endStr = DateFormat(end, format)

// 	// create sha1 encode string
// 	sh := sha1.New()
// 	sh.Write([]byte(data + setting.SecretKey + startStr + endStr + com.ToStr(minutes)))
// 	encoded := hex.EncodeToString(sh.Sum(nil))

// 	code := fmt.Sprintf("%s%06d%s", startStr, minutes, encoded)
// 	return code
// }

// Seconds-based time units
const (
	Minute = 60
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Week   = 7 * Day
	Month  = 30 * Day
	Year   = 12 * Month
)

func computeTimeDiff(diff int64) (int64, string) {
	diffStr := ""
	switch {
	case diff <= 0:
		diff = 0
		diffStr = "now"
	case diff < 2:
		diff = 0
		diffStr = "1 second"
	case diff < 1*Minute:
		diffStr = fmt.Sprintf("%d seconds", diff)
		diff = 0

	case diff < 2*Minute:
		diff -= 1 * Minute
		diffStr = "1 minute"
	case diff < 1*Hour:
		diffStr = fmt.Sprintf("%d minutes", diff/Minute)
		diff -= diff / Minute * Minute

	case diff < 2*Hour:
		diff -= 1 * Hour
		diffStr = "1 hour"
	case diff < 1*Day:
		diffStr = fmt.Sprintf("%d hours", diff/Hour)
		diff -= diff / Hour * Hour

	case diff < 2*Day:
		diff -= 1 * Day
		diffStr = "1 day"
	case diff < 1*Week:
		diffStr = fmt.Sprintf("%d days", diff/Day)
		diff -= diff / Day * Day

	case diff < 2*Week:
		diff -= 1 * Week
		diffStr = "1 week"
	case diff < 1*Month:
		diffStr = fmt.Sprintf("%d weeks", diff/Week)
		diff -= diff / Week * Week

	case diff < 2*Month:
		diff -= 1 * Month
		diffStr = "1 month"
	case diff < 1*Year:
		diffStr = fmt.Sprintf("%d months", diff/Month)
		diff -= diff / Month * Month

	case diff < 2*Year:
		diff -= 1 * Year
		diffStr = "1 year"
	default:
		diffStr = fmt.Sprintf("%d years", diff/Year)
		diff = 0
	}
	return diff, diffStr
}

// TimeSincePro calculates the time interval and generate full user-friendly string.
func TimeSincePro(then time.Time) string {
	now := time.Now()
	diff := now.Unix() - then.Unix()

	if then.After(now) {
		return "future"
	}

	var timeStr, diffStr string
	for {
		if diff == 0 {
			break
		}

		diff, diffStr = computeTimeDiff(diff)
		timeStr += ", " + diffStr
	}
	return strings.TrimPrefix(timeStr, ", ")
}

// func timeSince(then time.Time, lang string) string {
// 	now := time.Now()

// 	lbl := i18n.Tr(lang, "tool.ago")
// 	diff := now.Unix() - then.Unix()
// 	if then.After(now) {
// 		lbl = i18n.Tr(lang, "tool.from_now")
// 		diff = then.Unix() - now.Unix()
// 	}

// 	switch {
// 	case diff <= 0:
// 		return i18n.Tr(lang, "tool.now")
// 	case diff <= 2:
// 		return i18n.Tr(lang, "tool.1s", lbl)
// 	case diff < 1*Minute:
// 		return i18n.Tr(lang, "tool.seconds", diff, lbl)

// 	case diff < 2*Minute:
// 		return i18n.Tr(lang, "tool.1m", lbl)
// 	case diff < 1*Hour:
// 		return i18n.Tr(lang, "tool.minutes", diff/Minute, lbl)

// 	case diff < 2*Hour:
// 		return i18n.Tr(lang, "tool.1h", lbl)
// 	case diff < 1*Day:
// 		return i18n.Tr(lang, "tool.hours", diff/Hour, lbl)

// 	case diff < 2*Day:
// 		return i18n.Tr(lang, "tool.1d", lbl)
// 	case diff < 1*Week:
// 		return i18n.Tr(lang, "tool.days", diff/Day, lbl)

// 	case diff < 2*Week:
// 		return i18n.Tr(lang, "tool.1w", lbl)
// 	case diff < 1*Month:
// 		return i18n.Tr(lang, "tool.weeks", diff/Week, lbl)

// 	case diff < 2*Month:
// 		return i18n.Tr(lang, "tool.1mon", lbl)
// 	case diff < 1*Year:
// 		return i18n.Tr(lang, "tool.months", diff/Month, lbl)

// 	case diff < 2*Year:
// 		return i18n.Tr(lang, "tool.1y", lbl)
// 	default:
// 		return i18n.Tr(lang, "tool.years", diff/Year, lbl)
// 	}
// }

// TimeSince calculates the time interval and generate user-friendly string.
// func TimeSince(t time.Time, lang string) template.HTML {
// 	return template.HTML(fmt.Sprintf(`<span class="time-since" title="%s">%s</span>`, t.Format(setting.TimeFormat), timeSince(t, lang)))
// }

const (
	Byte  = 1
	KByte = Byte * 1024
	MByte = KByte * 1024
	GByte = MByte * 1024
	TByte = GByte * 1024
	PByte = TByte * 1024
	EByte = PByte * 1024
)

var bytesSizeTable = map[string]uint64{
	"b":  Byte,
	"kb": KByte,
	"mb": MByte,
	"gb": GByte,
	"tb": TByte,
	"pb": PByte,
	"eb": EByte,
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%dB", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f+"%s", val, suffix)
}

// FileSize calculates the file size and generate user-friendly string.
func FileSize(s int64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(uint64(s), 1024, sizes)
}

// Subtract deals with subtraction of all types of number.
func Subtract(left interface{}, right interface{}) interface{} {
	var rleft, rright int64
	var fleft, fright float64
	var isInt bool = true
	switch left.(type) {
	case int:
		rleft = int64(left.(int))
	case int8:
		rleft = int64(left.(int8))
	case int16:
		rleft = int64(left.(int16))
	case int32:
		rleft = int64(left.(int32))
	case int64:
		rleft = left.(int64)
	case float32:
		fleft = float64(left.(float32))
		isInt = false
	case float64:
		fleft = left.(float64)
		isInt = false
	}

	switch right.(type) {
	case int:
		rright = int64(right.(int))
	case int8:
		rright = int64(right.(int8))
	case int16:
		rright = int64(right.(int16))
	case int32:
		rright = int64(right.(int32))
	case int64:
		rright = right.(int64)
	case float32:
		fright = float64(left.(float32))
		isInt = false
	case float64:
		fleft = left.(float64)
		isInt = false
	}

	if isInt {
		return rleft - rright
	} else {
		return fleft + float64(rleft) - (fright + float64(rright))
	}
}

// DateFormat pattern rules.
var datePatterns = []string{
	// year
	"Y", "2006", // A full numeric representation of a year, 4 digits   Examples: 1999 or 2003
	"y", "06", //A two digit representation of a year   Examples: 99 or 03

	// month
	"m", "01", // Numeric representation of a month, with leading zeros 01 through 12
	"n", "1", // Numeric representation of a month, without leading zeros   1 through 12
	"M", "Jan", // A short textual representation of a month, three letters Jan through Dec
	"F", "January", // A full textual representation of a month, such as January or March   January through December

	// day
	"d", "02", // Day of the month, 2 digits with leading zeros 01 to 31
	"j", "2", // Day of the month without leading zeros 1 to 31

	// week
	"D", "Mon", // A textual representation of a day, three letters Mon through Sun
	"l", "Monday", // A full textual representation of the day of the week  Sunday through Saturday

	// time
	"g", "3", // 12-hour format of an hour without leading zeros    1 through 12
	"G", "15", // 24-hour format of an hour without leading zeros   0 through 23
	"h", "03", // 12-hour format of an hour with leading zeros  01 through 12
	"H", "15", // 24-hour format of an hour with leading zeros  00 through 23

	"a", "pm", // Lowercase Ante meridiem and Post meridiem am or pm
	"A", "PM", // Uppercase Ante meridiem and Post meridiem AM or PM

	"i", "04", // Minutes with leading zeros    00 to 59
	"s", "05", // Seconds, with leading zeros   00 through 59

	// time zone
	"T", "MST",
	"P", "-07:00",
	"O", "-0700",

	// RFC 2822
	"r", time.RFC1123Z,
}

// Parse Date use PHP time format.
func DateParse(dateString, format string) (time.Time, error) {
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	return time.ParseInLocation(format, dateString, time.Local)
}

// Date takes a PHP like date func to Go's time format.
func DateFormat(t time.Time, format string) string {
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	return t.Format(format)
}

type xssFilter struct {
	reg  *regexp.Regexp
	repl []byte
}

var (
	whiteSpace = []byte(" ")
	xssFilters = []xssFilter{
		{regexp.MustCompile(`\ [ONon]\w*=["]*`), whiteSpace},
		{regexp.MustCompile(`<[SCRIPTscript]{6}`), whiteSpace},
		{regexp.MustCompile(`=[` + "`" + `'"]*[JAVASCRIPTjavascript \t\0&#x0D;]*:`), whiteSpace},
	}
)

// XSS goes through all the XSS filters to make user input content as safe as possible.
func XSS(in []byte) []byte {
	for _, filter := range xssFilters {
		in = filter.reg.ReplaceAll(in, filter.repl)
	}
	return in
}

func XSSString(in string) string {
	return string(XSS([]byte(in)))
}

func IsNumber(str string) bool {
	for _, r := range str {
		if !unicode.IsNumber(r) {
			return false
		}
	}

	return true
}
