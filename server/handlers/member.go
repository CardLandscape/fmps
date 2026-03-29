package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"fmps/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MemberHandler struct {
	DB *gorm.DB
}

// taiwanRegionCodes defines Taiwan ID prefix codes with optional stop dates (zero = active).
var taiwanRegionCodes = map[byte]time.Time{
	'A': {},           // 台北市 – active
	'B': {},           // 台中市 – active
	'C': {},           // 基隆市 – active
	'D': {},           // 台南市 – active
	'E': {},           // 高雄市 – active
	'F': {},           // 新北市 – active
	'G': {},           // 宜兰县 – active
	'H': {},           // 桃园市 – active
	'I': {},           // 嘉义市 – active
	'J': {},           // 新竹县 – active
	'K': {},           // 苗栗县 – active
	'L': time.Date(2010, time.December, 25, 0, 0, 0, 0, time.UTC), // 台中县 – stopped
	'M': {},           // 南投县 – active
	'N': {},           // 彰化县 – active
	'O': {},           // 新竹市 – active
	'P': {},           // 云林县 – active
	'Q': {},           // 嘉义县 – active
	'R': time.Date(2010, time.December, 25, 0, 0, 0, 0, time.UTC), // 台南县 – stopped
	'S': time.Date(2010, time.December, 25, 0, 0, 0, 0, time.UTC), // 高雄县 – stopped
	'T': {},           // 屏东县 – active
	'U': {},           // 花莲县 – active
	'V': {},           // 台东县 – active
	'W': {},           // 金门县 – active
	'X': {},           // 澎湖县 – active
	'Y': time.Date(1974, time.January, 1, 0, 0, 0, 0, time.UTC),   // 阳明山管理局 – stopped
	'Z': {},           // 连江县 – active
}

// validateIDNumber 根据证件类型和国籍校验证件号码
func validateIDNumber(docType, number, nationality string) error {
	if number == "" {
		return nil // 允许空值
	}
	switch docType {
	case "01", "91":
		if err := validateChineseID(number); err != nil {
			return err
		}
	case "11":
		var prefix string
		if nationality == "HKG" {
			prefix = "810000"
		} else if nationality == "MAC" {
			prefix = "820000"
		}
		if prefix != "" && !strings.HasPrefix(number, prefix) {
			return fmt.Errorf("港澳居民居住证号码必须以 %s 开头", prefix)
		}
		if err := validateChineseID(number); err != nil {
			return err
		}
	case "21":
		if !strings.HasPrefix(number, "830000") {
			return fmt.Errorf("台湾居民居住证号码必须以 830000 开头")
		}
		if err := validateChineseID(number); err != nil {
			return err
		}
	case "31":
		// 9 + 2位行政区划 + 3位数字国籍代码 + 剩余按18位身份证规则
		if len(number) != 18 || number[0] != '9' {
			return fmt.Errorf("外国人永久居留身份证号码必须以 9 开头且为 18 位")
		}
		if err := validateChineseID(number); err != nil {
			return err
		}
	case "02":
		// H/M + 8位数字，国籍HKG必须H开头，MAC必须M开头
		matched, _ := regexp.MatchString(`^[HM]\d{8}$`, number)
		if !matched {
			return fmt.Errorf("港澳通行证号码格式错误（H/M + 8位数字）")
		}
		if nationality == "HKG" && number[0] != 'H' {
			return fmt.Errorf("国籍为HKG时，港澳通行证号码必须以H开头")
		}
		if nationality == "MAC" && number[0] != 'M' {
			return fmt.Errorf("国籍为MAC时，港澳通行证号码必须以M开头")
		}
	case "03":
		// 8位数字
		matched, _ := regexp.MatchString(`^\d{8}$`, number)
		if !matched {
			return fmt.Errorf("台湾居民来往大陆通行证号码必须为 8 位数字")
		}
	case "04":
		// E + 8位数字 或 E + 1字母 + 7位数字
		matched1, _ := regexp.MatchString(`^E\d{8}$`, number)
		matched2, _ := regexp.MatchString(`^E[A-Za-z]\d{7}$`, number)
		if !matched1 && !matched2 {
			return fmt.Errorf("中国护照号码格式错误（E+8位数字 或 E+1字母+7位数字）")
		}
	case "05":
		l := len(number)
		if l < 6 || l > 9 {
			return fmt.Errorf("外国护照号码长度必须为 6-9 位")
		}
	case "52":
		// HA/MA + 7位数字
		matched, _ := regexp.MatchString(`^(HA|MA)\d{7}$`, number)
		if !matched {
			return fmt.Errorf("港澳通行证（非中国籍）号码格式错误（HA/MA + 7位数字）")
		}
	// 辅助证件类型
	case "90", "92", "95":
		// 1-2位字母 + 7位数字；禁止以 W 或 WX 开头
		if err := validateHKMOID(number); err != nil {
			return err
		}
	case "93":
		// 1位字母（台湾地区码）+ 9位数字；检查停发日期
		return nil // 93 的出生日期需要外层传入，此处仅做格式检查
	case "94":
		// 证明文件号码，按 proof_doc_type 区分。此处不校验，由 validateAux94 处理
		return nil
	case "96":
		// 8位数字，以1开头
		matched, _ := regexp.MatchString(`^1\d{7}$`, number)
		if !matched {
			return fmt.Errorf("澳门居民身份证号码格式错误（8位数字，以1开头）")
		}
	case "97":
		// 8位数字，以5开头
		matched, _ := regexp.MatchString(`^5\d{7}$`, number)
		if !matched {
			return fmt.Errorf("澳门永久性居民身份证号码格式错误（8位数字，以5开头）")
		}
	case "98":
		// 8位数字，以7开头
		matched, _ := regexp.MatchString(`^7\d{7}$`, number)
		if !matched {
			return fmt.Errorf("澳门永久性居民身份证（外国籍）号码格式错误（8位数字，以7开头）")
		}
	}
	return nil
}

// validateHKMOID 校验90/92/95类证件：1-2位字母 + 7位数字，禁止以 W 开头（含 WX）
func validateHKMOID(number string) error {
	matched, _ := regexp.MatchString(`^[A-Za-z]{1,2}\d{7}$`, number)
	if !matched {
		return fmt.Errorf("证件号码格式错误（1-2位字母 + 7位数字）")
	}
	if number[0] == 'W' || number[0] == 'w' {
		return fmt.Errorf("此证件号码为根据补充劳工计划签发给来港就业的工人的身份证号码，该证持有人不具有香港居留权及不合资格申请回乡证")
	}
	return nil
}

// validateTaiwan93 校验93类台湾居民身份证号码，需传入出生日期用于停发代码检查
func validateTaiwan93(number string, birthDate string) error {
	if len(number) != 10 {
		return fmt.Errorf("台湾居民身份证号码必须为10位（1位字母 + 9位数字）")
	}
	prefix := strings.ToUpper(number[:1])[0]
	stopDate, ok := taiwanRegionCodes[prefix]
	if !ok {
		return fmt.Errorf("台湾居民身份证号码首位字母无效（%c）", number[0])
	}
	// 检查后9位是否为数字
	digits, _ := regexp.MatchString(`^\d{9}$`, number[1:])
	if !digits {
		return fmt.Errorf("台湾居民身份证号码格式错误（1位字母 + 9位数字）")
	}
	// 若停发日期非零，且出生日期晚于停发日期则拒绝
	if !stopDate.IsZero() && birthDate != "" {
		bd, err := time.Parse("2006-01-02", birthDate)
		if err == nil && bd.After(stopDate) {
			return fmt.Errorf("台湾居民身份证地区码 %c 已于 %s 停止赋配，出生日期晚于停发日期不允许使用此代码",
				prefix, stopDate.Format("2006-01-02"))
		}
	}
	return nil
}

// validateAux94Number 校验94类证件号码（按 proofDocType 区分）
func validateAux94Number(number string, proofDocType string) error {
	if number == "" {
		return nil
	}
	if proofDocType == "94NP" {
		// H + 12位数字
		matched, _ := regexp.MatchString(`^H\d{12}$`, number)
		if !matched {
			return fmt.Errorf("94NP证件号码格式错误（H + 12位数字）")
		}
		return nil
	}
	// 其他94类：无强制格式要求
	return nil
}

// validateNationalityDocType 校验国籍与证件类型的组合
func validateNationalityDocType(docType, nationality string) error {
	if docType == "" || nationality == "" {
		return nil
	}
	switch docType {
	case "01", "91", "04":
		if nationality != "CHN" {
			return fmt.Errorf("证件类型 %s 要求国籍必须为 CHN（中国）", docType)
		}
	case "11", "02":
		if nationality != "HKG" && nationality != "MAC" {
			return fmt.Errorf("证件类型 %s 要求国籍为 HKG（香港）或 MAC（澳门）", docType)
		}
	case "21", "03":
		if nationality != "TWN" {
			return fmt.Errorf("证件类型 %s 要求国籍为 TWN（台湾）", docType)
		}
	case "31", "05", "52":
		restricted := map[string]bool{"CHN": true, "HKG": true, "MAC": true, "TWN": true}
		if restricted[nationality] {
			return fmt.Errorf("证件类型 %s 要求国籍不得为 CHN/HKG/MAC/TWN", docType)
		}
	}
	return nil
}

// validateChineseID 校验 18 位中国居民身份证
func validateChineseID(id string) error {
	if len(id) != 18 {
		return fmt.Errorf("身份证号码必须为 18 位")
	}
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checkCodes := "10X98765432"
	sum := 0
	for i := 0; i < 17; i++ {
		c := id[i]
		if c < '0' || c > '9' {
			return fmt.Errorf("身份证号码前17位必须为数字")
		}
		sum += int(c-'0') * weights[i]
	}
	expected := checkCodes[sum%11]
	last := id[17]
	// 大写处理，兼容小写 x
	if last >= 'a' && last <= 'z' {
		last = last - 32
	}
	if last != expected {
		return fmt.Errorf("身份证校验码错误（期望 %c，实际 %c）", expected, last)
	}
	// 校验日期
	year, _ := strconv.Atoi(id[6:10])
	month, _ := strconv.Atoi(id[10:12])
	day, _ := strconv.Atoi(id[12:14])
	if month < 1 || month > 12 || day < 1 || day > 31 {
		return fmt.Errorf("身份证日期部分无效")
	}
	if year < 1900 || year > time.Now().Year() {
		return fmt.Errorf("身份证年份无效")
	}
	return nil
}

// normalizeMember 补全 name 字段并处理兼容逻辑
func normalizeMember(m *models.Member) {
	// 若 name_cn 为空但 name 有值，回填 name_cn（兼容旧数据）
	if m.NameCn == "" && m.Name != "" {
		m.NameCn = m.Name
	}
	// 保存时同步写回 name = name_cn
	if m.NameCn != "" {
		m.Name = m.NameCn
	}
}

// validateIDGenderBirthdate checks that the gender digit and birthdate embedded in
// Chinese-ID-format numbers (types 01/91/11/21/31) match the member's stated values.
// Returns a generic "证件号码无效" on any mismatch to avoid exposing rule hints.
func validateIDGenderBirthdate(docType, number, gender, birthDate string) error {
	switch docType {
	case "01", "91", "11", "21", "31":
		if len(number) != 18 {
			return fmt.Errorf("证件号码无效")
		}
		// Digits 7-14 (index 6-13) encode birthdate as YYYYMMDD
		if birthDate != "" {
			bd, err := time.Parse("2006-01-02", birthDate)
			if err == nil {
				if number[6:14] != bd.Format("20060102") {
					return fmt.Errorf("证件号码无效")
				}
			}
		}
		// 17th digit (index 16): odd = male (男), even = female (女)
		if number[16] < '0' || number[16] > '9' {
			return fmt.Errorf("证件号码无效")
		}
		d := int(number[16] - '0')
		idGender := "男"
		if d%2 == 0 {
			idGender = "女"
		}
		if gender != "" && gender != idGender {
			return fmt.Errorf("证件号码无效")
		}
	}
	return nil
}

// validateAux93Gender checks that the second character of a type-93 number is 1 or 2
// and matches the member's stated gender (1=male/男, 2=female/女).
// Returns a generic "证件号码无效" on any violation.
func validateAux93Gender(number, gender string) error {
	if len(number) < 2 {
		return fmt.Errorf("证件号码无效")
	}
	c := number[1]
	if c != '1' && c != '2' {
		return fmt.Errorf("证件号码无效")
	}
	numGender := "男"
	if c == '2' {
		numGender = "女"
	}
	if gender != "" && gender != numGender {
		return fmt.Errorf("证件号码无效")
	}
	return nil
}

// validateMember 验证成员数据合法性
func validateMember(m *models.Member) error {
	// 英文姓名必填（任何情况）
	if strings.TrimSpace(m.NameEn) == "" {
		return fmt.Errorf("英文姓名（name_en）为必填项")
	}
	// 国籍为 CHN/HKG/MAC/TWN 时，中文姓名必填
	chineseNats := map[string]bool{"CHN": true, "HKG": true, "MAC": true, "TWN": true}
	if chineseNats[m.Nationality] && strings.TrimSpace(m.NameCn) == "" {
		return fmt.Errorf("国籍为 CHN/HKG/MAC/TWN 时，中文姓名（name_cn）为必填项")
	}
	// 国籍必填
	if strings.TrimSpace(m.Nationality) == "" {
		return fmt.Errorf("国籍（nationality）为必填项")
	}
	// 出生日期必填
	if strings.TrimSpace(m.BirthDate) == "" {
		return fmt.Errorf("出生日期（birth_date）为必填项")
	}
	// 主证件类型必填
	if strings.TrimSpace(m.IdDocType) == "" {
		return fmt.Errorf("主证件类型（id_doc_type）为必填项")
	}
	// 主证件号码必填
	if strings.TrimSpace(m.IdDocNumber) == "" {
		return fmt.Errorf("主证件号码（id_doc_number）为必填项")
	}
	// 主证件签发日期必填
	if strings.TrimSpace(m.IdIssueDate) == "" {
		return fmt.Errorf("主证件签发日期（id_issue_date）为必填项")
	}
	// 主证件有效期必填
	if strings.TrimSpace(m.IdExpiryDate) == "" {
		return fmt.Errorf("主证件有效期（id_expiry_date）为必填项")
	}
	// 主证件签发机关必填
	if strings.TrimSpace(m.IdIssueAuthority) == "" {
		return fmt.Errorf("主证件签发机关（id_issue_authority）为必填项")
	}
	// 主证件校验
	if err := validateNationalityDocType(m.IdDocType, m.Nationality); err != nil {
		return fmt.Errorf("主证件：%s", err)
	}
	if err := validateIDNumber(m.IdDocType, m.IdDocNumber, m.Nationality); err != nil {
		return fmt.Errorf("证件号码无效")
	}
	// 主证件性别与出生日期一致性校验（类型 01/91/11/21/31）
	if err := validateIDGenderBirthdate(m.IdDocType, m.IdDocNumber, m.Gender, m.BirthDate); err != nil {
		return err
	}

	// 辅助证件约束
	if err := validateAuxDocs(m); err != nil {
		return err
	}

	// 就读学校名称需包含关键词
	schoolName := strings.TrimSpace(m.SchoolName)
	if schoolName != "" {
		keywords := []string{"小学", "中学", "大学", "学院"}
		hasKeyword := false
		for _, kw := range keywords {
			if strings.Contains(schoolName, kw) {
				hasKeyword = true
				break
			}
		}
		if !hasKeyword {
			return fmt.Errorf("就读学校名称须包含「小学」、「中学」、「大学」或「学院」之一")
		}
	}
	return nil
}

// validateAuxDocs 根据主证件类型校验辅助证件规则
func validateAuxDocs(m *models.Member) error {
	mainType := m.IdDocType
	aux1Type := strings.TrimSpace(m.Aux1DocType)
	aux1Num := strings.TrimSpace(m.Aux1DocNumber)
	aux2Type := strings.TrimSpace(m.Aux2DocType)
	aux2Num := strings.TrimSpace(m.Aux2DocNumber)

	switch mainType {
	case "01", "91":
		// 不允许辅助证件
		if aux1Type != "" || aux2Type != "" {
			return fmt.Errorf("主证件类型 %s 不允许录入辅助证件", mainType)
		}

	case "11":
		// 辅助1必须为02，辅助2必须为90/92/96/97，两者均必填
		if aux1Type != "02" {
			return fmt.Errorf("主证件为11时，辅助证件1类型必须为02（港澳居民来往内地通行证）")
		}
		if aux1Num == "" {
			return fmt.Errorf("辅助证件1号码为必填项")
		}
		if err := validateIDNumber("02", aux1Num, m.Nationality); err != nil {
			return fmt.Errorf("证件号码无效")
		}
		allowed2 := map[string]bool{"90": true, "92": true, "96": true, "97": true}
		if !allowed2[aux2Type] {
			return fmt.Errorf("主证件为11时，辅助证件2类型只能为90/92/96/97")
		}
		if aux2Num == "" {
			return fmt.Errorf("辅助证件2号码为必填项")
		}
		if err := validateIDNumber(aux2Type, aux2Num, m.Nationality); err != nil {
			return fmt.Errorf("证件号码无效")
		}

	case "21":
		// 辅助1必须为03，辅助2必须为93，两者均必填
		if aux1Type != "03" {
			return fmt.Errorf("主证件为21时，辅助证件1类型必须为03（台湾居民来往大陆通行证）")
		}
		if aux1Num == "" {
			return fmt.Errorf("辅助证件1号码为必填项")
		}
		if err := validateIDNumber("03", aux1Num, m.Nationality); err != nil {
			return fmt.Errorf("证件号码无效")
		}
		if aux2Type != "93" {
			return fmt.Errorf("主证件为21时，辅助证件2类型必须为93（台湾居民身份证）")
		}
		if aux2Num == "" {
			return fmt.Errorf("辅助证件2号码为必填项")
		}
		if err := validateTaiwan93(aux2Num, m.BirthDate); err != nil {
			return fmt.Errorf("证件号码无效")
		}
		if err := validateAux93Gender(aux2Num, m.Gender); err != nil {
			return err
		}

	case "04":
		// 辅助证件必须为94，proof_doc_type必填，proof_issue_country 必填（94NP时须为CHN）
		if aux1Type != "94" {
			return fmt.Errorf("主证件为04时，辅助证件类型必须为94")
		}
		if aux1Num == "" {
			return fmt.Errorf("辅助证件号码为必填项")
		}
		proofType := strings.TrimSpace(m.ProofDocType)
		if proofType == "" {
			return fmt.Errorf("主证件为04时，证明文件类型（proof_doc_type）为必填项")
		}
		validProofTypes := map[string]bool{"94RV": true, "94PV": true, "94PC": true, "94PE": true, "94NP": true}
		if !validProofTypes[proofType] {
			return fmt.Errorf("证明文件类型无效（%s）", proofType)
		}
		proofCountry := strings.TrimSpace(m.ProofIssueCountry)
		if proofCountry == "" {
			return fmt.Errorf("主证件为04时，签发国家（proof_issue_country）为必填项")
		}
		if proofType == "94NP" && proofCountry != "CHN" {
			return fmt.Errorf("证明文件类型为94NP时，签发国家必须为CHN")
		}
		if proofType != "94NP" {
			restricted := map[string]bool{"CHN": true, "HKG": true, "MAC": true, "TWN": true}
			if restricted[proofCountry] {
				return fmt.Errorf("签发国家不得为CHN/HKG/MAC/TWN（94NP除外）")
			}
		}
		// 校验94号码
		if err := validateAux94Number(aux1Num, proofType); err != nil {
			return fmt.Errorf("证件号码无效")
		}

	case "05":
		// 隐藏辅助证件
		if aux1Type != "" || aux2Type != "" {
			return fmt.Errorf("主证件类型05时不允许录入辅助证件")
		}

	case "31":
		// 辅助证件类型只能为05（外国护照），可不填
		if aux1Type != "" {
			if aux1Type != "05" {
				return fmt.Errorf("主证件为31时，辅助证件类型只能为05（外国护照）")
			}
			if aux1Num != "" {
				if err := validateIDNumber("05", aux1Num, m.Nationality); err != nil {
					return fmt.Errorf("证件号码无效")
				}
			}
		}

	case "02":
		// 辅助类型限90/92/96/97；HKG仅90/92，MAC仅96/97
		if aux1Type != "" {
			allowedAux := map[string]bool{"90": true, "92": true, "96": true, "97": true}
			if m.Nationality == "HKG" {
				allowedAux = map[string]bool{"90": true, "92": true}
			} else if m.Nationality == "MAC" {
				allowedAux = map[string]bool{"96": true, "97": true}
			}
			if !allowedAux[aux1Type] {
				if m.Nationality == "HKG" {
					return fmt.Errorf("国籍为HKG时，辅助证件类型只能为90/92")
				} else if m.Nationality == "MAC" {
					return fmt.Errorf("国籍为MAC时，辅助证件类型只能为96/97")
				}
				return fmt.Errorf("主证件为02时，辅助证件类型只能为90/92/96/97")
			}
			if aux1Num != "" {
				if err := validateIDNumber(aux1Type, aux1Num, m.Nationality); err != nil {
					return fmt.Errorf("证件号码无效")
				}
			}
		}

	case "03":
		// 辅助证件类型只能为93（台湾居民身份证），可不填
		if aux1Type != "" {
			if aux1Type != "93" {
				return fmt.Errorf("主证件为03时，辅助证件类型只能为93（台湾居民身份证）")
			}
			if aux1Num != "" {
				if err := validateTaiwan93(aux1Num, m.BirthDate); err != nil {
					return fmt.Errorf("证件号码无效")
				}
				if err := validateAux93Gender(aux1Num, m.Gender); err != nil {
					return err
				}
			}
		}

	case "52":
		// 辅助类型只能为95或98
		if aux1Type != "" {
			allowed1 := map[string]bool{"95": true, "98": true}
			if !allowed1[aux1Type] {
				return fmt.Errorf("主证件为52时，辅助证件类型只能为95或98")
			}
			if aux1Num != "" {
				if err := validateIDNumber(aux1Type, aux1Num, m.Nationality); err != nil {
					return fmt.Errorf("证件号码无效")
				}
			}
		}
	}
	return nil
}

// getAuthPassword retrieves the authorization password from Settings.
func (h *MemberHandler) getAuthPassword() (string, error) {
	var setting models.Setting
	if err := h.DB.Where("key = ?", "authorization_password").First(&setting).Error; err != nil {
		return "", err
	}
	return setting.Value, nil
}

// checkAuthPassword validates the provided password against the stored authorization password.
func (h *MemberHandler) checkAuthPassword(provided string) error {
	stored, err := h.getAuthPassword()
	if err != nil {
		return fmt.Errorf("系统配置错误，无法验证授权密码")
	}
	if provided != stored {
		return fmt.Errorf("授权密码错误")
	}
	return nil
}

// countActiveParents returns the number of parent-role members in the database.
func (h *MemberHandler) countActiveParents(excludeID uint) (int64, error) {
	var count int64
	q := h.DB.Model(&models.Member{}).Where("role = ?", "parent")
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (h *MemberHandler) List(c *gin.Context) {
	var members []models.Member
	if err := h.DB.Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败"})
		return
	}
	// 兼容旧数据：name_cn 为空时从 name 回填；role="adult" 视为 "parent"
	for i := range members {
		if members[i].NameCn == "" && members[i].Name != "" {
			members[i].NameCn = members[i].Name
		}
		if members[i].Role == "adult" {
			members[i].Role = "parent"
		}
	}
	c.JSON(http.StatusOK, members)
}

func (h *MemberHandler) Create(c *gin.Context) {
	var member models.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求格式错误"})
		return
	}

	// Normalize legacy "adult" → "parent"
	if member.Role == "adult" {
		member.Role = "parent"
	}
	// Role must be "parent" or "child"
	if member.Role != "parent" && member.Role != "child" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "成员类型只能为「家长」或「小孩」"})
		return
	}

	// If creating a parent record, check if a child record with same name+birthdate already exists.
	// If so, require authorization password.
	if member.Role == "parent" {
		nameCn := strings.TrimSpace(member.NameCn)
		birthDate := strings.TrimSpace(member.BirthDate)
		if nameCn != "" && birthDate != "" {
			var existingCount int64
			h.DB.Model(&models.Member{}).
				Where("role = ? AND name_cn = ? AND birth_date = ?",
					"child", nameCn, birthDate).
				Count(&existingCount)
			if existingCount > 0 {
				if err := h.checkAuthPassword(member.AuthPassword); err != nil {
					c.JSON(http.StatusForbidden, gin.H{"message": "该成员已有「小孩」记录，添加「家长」记录需要授权密码"})
					return
				}
			}
		}
	}

	normalizeMember(&member)
	if err := validateMember(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := h.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建失败"})
		return
	}
	c.JSON(http.StatusCreated, member)
}

func (h *MemberHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的ID"})
		return
	}

	var existing models.Member
	if err := h.DB.First(&existing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "成员不存在"})
		return
	}

	var incoming models.Member
	if err := c.ShouldBindJSON(&incoming); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求格式错误"})
		return
	}

	// Require authorization password for editing any member
	if err := h.checkAuthPassword(incoming.AuthPassword); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "编辑成员资料需要授权密码"})
		return
	}

	// Normalize legacy "adult" → "parent"
	if incoming.Role == "adult" {
		incoming.Role = "parent"
	}
	existingRole := existing.Role
	if existingRole == "adult" {
		existingRole = "parent"
	}

	// Prevent role change after creation
	if incoming.Role != "" && incoming.Role != existingRole {
		c.JSON(http.StatusBadRequest, gin.H{"message": "成员类型一经确定不可更改"})
		return
	}
	// Keep original role if not provided
	if incoming.Role == "" {
		incoming.Role = existingRole
	}

	incoming.ID = uint(id)
	normalizeMember(&incoming)
	if err := validateMember(&incoming); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.DB.Save(&incoming).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, incoming)
}

// DeleteWithAuth handles member deletion with authorization password in request body.
// Used by the frontend POST /api/members/:id/delete endpoint.
func (h *MemberHandler) DeleteWithAuth(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的ID"})
		return
	}

	var req struct {
		AuthPassword string `json:"auth_password"`
	}
	// Body is optional (empty body means no auth_password provided).
	// Only return an error for genuinely malformed JSON, not for an empty body.
	if err := c.ShouldBindJSON(&req); err != nil {
		// io.EOF and io.ErrUnexpectedEOF mean empty body – treat as missing password
		errMsg := err.Error()
		if errMsg != "EOF" && errMsg != "unexpected EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "请求格式错误"})
			return
		}
	}

	var member models.Member
	if err := h.DB.First(&member, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "成员不存在"})
		return
	}

	// Require authorization password for both roles
	if err := h.checkAuthPassword(req.AuthPassword); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "删除成员需要授权密码"})
		return
	}

	role := member.Role
	if role == "adult" {
		role = "parent"
	}

	// For parent role: ensure at least one other parent will remain
	if role == "parent" {
		remaining, err := h.countActiveParents(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败"})
			return
		}
		if remaining < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "系统中必须至少保留一名有效的家长，无法删除最后一名家长"})
			return
		}
	}

	if err := h.DB.Delete(&models.Member{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (h *MemberHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的ID"})
		return
	}

	if err := h.DB.Delete(&models.Member{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
