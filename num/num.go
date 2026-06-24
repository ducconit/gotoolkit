package num

import (
	"errors"
	"maps"
	"math"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// ErrUnsupportedLanguage trả về khi ngôn ngữ yêu cầu chưa được đăng ký.
var ErrUnsupportedLanguage = errors.New("unsupported language")

// Translator định nghĩa hàm chuyển đổi số thành chữ cho một ngôn ngữ cụ thể.
type Translator func(val int64) string

var (
	mu       sync.RWMutex
	registry atomic.Value // Lưu trữ map[string]Translator
)

// Register đăng ký một Translator cho mã ngôn ngữ (ví dụ: "vi", "en").
func Register(lang string, t Translator) {
	mu.Lock()
	defer mu.Unlock()

	var oldMap map[string]Translator
	if val := registry.Load(); val != nil {
		oldMap = val.(map[string]Translator)
	} else {
		oldMap = make(map[string]Translator)
	}

	newMap := make(map[string]Translator, len(oldMap)+1)
	maps.Copy(newMap, oldMap)
	newMap[strings.ToLower(lang)] = t
	registry.Store(newMap)
}

// ToWords chuyển đổi một số nguyên thành chữ dựa trên ngôn ngữ được chọn.
// Mặc định hỗ trợ "en" (tiếng Anh) và "vi" (tiếng Việt).
func ToWords(val int64, lang string) (string, error) {
	m, ok := registry.Load().(map[string]Translator)
	if !ok {
		return "", ErrUnsupportedLanguage
	}
	t, found := m[strings.ToLower(lang)]
	if !found {
		return "", ErrUnsupportedLanguage
	}
	return t(val), nil
}

// ToShorthand chuyển đổi số nguyên thành dạng viết tắt ngắn gọn (ví dụ: 1.5K, 2M hoặc 1.5 nghìn, 2 triệu).
// Hỗ trợ mã ngôn ngữ "en" (K, M, B, T, Qa, Qi) và "vi" (nghìn, triệu, tỷ, nghìn tỷ, triệu tỷ, tỷ tỷ).
// Mặc định nếu ngôn ngữ không phải "vi" sẽ hiển thị theo định dạng quốc tế ("en").
func ToShorthand(val int64, lang string) string {
	var absVal float64
	var isNegative bool

	if val == math.MinInt64 {
		isNegative = true
		absVal = 9223372036854775808
	} else if val < 0 {
		isNegative = true
		absVal = float64(-val)
	} else {
		absVal = float64(val)
	}

	// Kiểm tra ngôn ngữ không allocate
	isVI := len(lang) == 2 && (lang[0] == 'v' || lang[0] == 'V') && (lang[1] == 'i' || lang[1] == 'I')

	var buf [64]byte
	b := buf[:0]
	if isNegative {
		b = append(b, '-')
	}

	if isVI {
		if absVal < 1000 {
			b = strconv.AppendFloat(b, absVal, 'f', -1, 64)
		} else if absVal < 1000000 {
			b = appendFormattedFloat(b, absVal/1000)
			b = append(b, " nghìn"...)
		} else if absVal < 1000000000 {
			b = appendFormattedFloat(b, absVal/1000000)
			b = append(b, " triệu"...)
		} else if absVal < 1000000000000 {
			b = appendFormattedFloat(b, absVal/1000000000)
			b = append(b, " tỷ"...)
		} else if absVal < 1000000000000000 {
			b = appendFormattedFloat(b, absVal/1000000000000)
			b = append(b, " nghìn tỷ"...)
		} else if absVal < 1000000000000000000 {
			b = appendFormattedFloat(b, absVal/1000000000000000)
			b = append(b, " triệu tỷ"...)
		} else {
			b = appendFormattedFloat(b, absVal/1000000000000000000)
			b = append(b, " tỷ tỷ"...)
		}
	} else {
		if absVal < 1000 {
			b = strconv.AppendFloat(b, absVal, 'f', -1, 64)
		} else if absVal < 1000000 {
			b = appendFormattedFloat(b, absVal/1000)
			b = append(b, 'K')
		} else if absVal < 1000000000 {
			b = appendFormattedFloat(b, absVal/1000000)
			b = append(b, 'M')
		} else if absVal < 1000000000000 {
			b = appendFormattedFloat(b, absVal/1000000000)
			b = append(b, 'B')
		} else if absVal < 1000000000000000 {
			b = appendFormattedFloat(b, absVal/1000000000000)
			b = append(b, 'T')
		} else if absVal < 1000000000000000000 {
			b = appendFormattedFloat(b, absVal/1000000000000000)
			b = append(b, "Qa"...)
		} else {
			b = appendFormattedFloat(b, absVal/1000000000000000000)
			b = append(b, "Qi"...)
		}
	}

	return string(b)
}

func appendFormattedFloat(b []byte, val float64) []byte {
	val = math.Round(val*100) / 100
	return strconv.AppendFloat(b, val, 'f', -1, 64)
}

// Format định dạng số thực thành chuỗi có dấu phân tách hàng nghìn và phần thập phân tùy chọn.
// Tham số precision xác định số chữ số sau dấu thập phân (precision >= 0 hoặc âm để tự động).
// thousandSep là ký tự phân tách hàng nghìn (ví dụ: ",", ".", " " hoặc "" nếu không muốn dùng).
// decimalSep là ký tự phân tách phần thập phân (ví dụ: ".", ",").
func Format(val float64, precision int, thousandSep string, decimalSep string) string {
	var prefix string
	absVal := val
	if math.Signbit(val) {
		prefix = "-"
		absVal = -val
	}

	var rawBuf [128]byte
	rawBytes := strconv.AppendFloat(rawBuf[:0], absVal, 'f', precision, 64)

	dotIdx := -1
	for i, c := range rawBytes {
		if c == '.' {
			dotIdx = i
			break
		}
	}

	var intPart, decPart []byte
	if dotIdx != -1 {
		intPart = rawBytes[:dotIdx]
		decPart = rawBytes[dotIdx+1:]
	} else {
		intPart = rawBytes
	}

	n := len(intPart)

	sepCount := 0
	if thousandSep != "" && n > 3 {
		sepCount = (n - 1) / 3
	}

	totalLen := len(prefix) + n + sepCount*len(thousandSep)
	if precision != 0 && len(decPart) > 0 {
		totalLen += len(decimalSep) + len(decPart)
	}

	var resBuf [128]byte
	var resBytes []byte
	if totalLen <= 128 {
		resBytes = resBuf[:0]
	} else {
		resBytes = make([]byte, 0, totalLen)
	}

	if prefix != "" {
		resBytes = append(resBytes, prefix...)
	}

	for i := 0; i < n; i++ {
		if i > 0 && (n-i)%3 == 0 && thousandSep != "" {
			resBytes = append(resBytes, thousandSep...)
		}
		resBytes = append(resBytes, intPart[i])
	}

	if precision != 0 && len(decPart) > 0 {
		resBytes = append(resBytes, decimalSep...)
		resBytes = append(resBytes, decPart...)
	}

	return string(resBytes)
}


func init() {
	Register("en", toWordsEN)
	Register("vi", toWordsVI)
}

// --- THUẬT TOÁN TIẾNG ANH (en) ---

var (
	enLessThan20 = []string{
		"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
		"ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen",
	}
	enTens = []string{
		"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety",
	}
	enGroups = []string{
		"", "thousand", "million", "billion", "trillion", "quadrillion", "quintillion",
	}
)

func toWordsEN(val int64) string {
	if val == 0 {
		return "zero"
	}

	var sb strings.Builder
	var absVal uint64

	if val == math.MinInt64 {
		sb.WriteString("minus ")
		absVal = 9223372036854775808
	} else if val < 0 {
		sb.WriteString("minus ")
		absVal = uint64(-val)
	} else {
		absVal = uint64(val)
	}

	// Tách thành các nhóm 3 chữ số bằng array trên stack để tránh heap allocation
	var parts [7]uint64
	nParts := 0
	for absVal > 0 {
		parts[nParts] = absVal % 1000
		nParts++
		absVal /= 1000
	}

	var wordsBuf [32]string
	words := wordsBuf[:0]
	for i := nParts - 1; i >= 0; i-- {
		part := parts[i]
		if part == 0 {
			continue
		}

		partStr := readThreeDigitsEN(part)
		groupName := enGroups[i]

		if groupName != "" {
			words = append(words, partStr, groupName)
		} else {
			words = append(words, partStr)
		}
	}

	sb.WriteString(strings.Join(words, " "))
	return sb.String()
}

func readThreeDigitsEN(n uint64) string {
	var words []string

	hundreds := n / 100
	rem := n % 100

	if hundreds > 0 {
		words = append(words, enLessThan20[hundreds], "hundred")
	}

	if rem > 0 {
		if rem < 20 {
			words = append(words, enLessThan20[rem])
		} else {
			tens := rem / 10
			ones := rem % 10
			if ones > 0 {
				words = append(words, enTens[tens]+"-"+enLessThan20[ones])
			} else {
				words = append(words, enTens[tens])
			}
		}
	}

	return strings.Join(words, " ")
}

// --- THUẬT TOÁN TIẾNG VIỆT (vi) ---

var (
	viDigits = []string{
		"không", "một", "hai", "ba", "bốn", "năm", "sáu", "bảy", "tám", "chín",
	}
	viGroups = []string{
		"", "nghìn", "triệu", "tỷ", "nghìn tỷ", "triệu tỷ", "tỷ tỷ",
	}
)

func toWordsVI(val int64) string {
	if val == 0 {
		return "không"
	}

	var sb strings.Builder
	var absVal uint64

	if val == math.MinInt64 {
		sb.WriteString("âm ")
		absVal = 9223372036854775808
	} else if val < 0 {
		sb.WriteString("âm ")
		absVal = uint64(-val)
	} else {
		absVal = uint64(val)
	}

	// Tách thành các nhóm 3 chữ số từ hàng đơn vị trở lên bằng array trên stack
	var parts [7]uint64
	nParts := 0
	for absVal > 0 {
		parts[nParts] = absVal % 1000
		nParts++
		absVal /= 1000
	}

	var wordsBuf [32]string
	words := wordsBuf[:0]
	hasHigher := false

	for i := nParts - 1; i >= 0; i-- {
		part := parts[i]
		if part == 0 {
			// Nếu nhóm hiện tại bằng 0 nhưng là nhóm cuối cùng (đơn vị) và trước đó đã có đọc số lớn hơn
			// thì ta không làm gì (ví dụ: 1,000,000 đọc là "một triệu" chứ không đọc nghìn hay đơn vị).
			continue
		}

		partStr := readThreeDigitsVI(part, hasHigher)
		groupName := viGroups[i]

		if groupName != "" {
			words = append(words, partStr, groupName)
		} else {
			words = append(words, partStr)
		}

		hasHigher = true
	}

	sb.WriteString(strings.Join(words, " "))
	return sb.String()
}

func readThreeDigitsVI(n uint64, hasHigher bool) string {
	var words []string

	hundreds := n / 100
	tens := (n % 100) / 10
	ones := n % 10

	// Đọc hàng trăm
	if hundreds > 0 {
		words = append(words, viDigits[hundreds], "trăm")
	} else if hasHigher {
		words = append(words, "không", "trăm")
	}

	// Đọc hàng chục và đơn vị
	if tens == 0 {
		if ones > 0 {
			if hundreds > 0 || hasHigher {
				words = append(words, "lẻ", viDigits[ones])
			} else {
				words = append(words, viDigits[ones])
			}
		}
	} else if tens == 1 {
		words = append(words, "mười")
		if ones == 5 {
			words = append(words, "lăm")
		} else if ones > 0 {
			words = append(words, viDigits[ones])
		}
	} else {
		words = append(words, viDigits[tens], "mươi")
		if ones == 1 {
			words = append(words, "mốt")
		} else if ones == 4 {
			words = append(words, "tư")
		} else if ones == 5 {
			words = append(words, "lăm")
		} else if ones > 0 {
			words = append(words, viDigits[ones])
		}
	}

	return strings.Join(words, " ")
}
