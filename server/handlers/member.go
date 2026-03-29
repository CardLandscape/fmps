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

// validateIDNumber 根据证件类型和国籍校验证件号码（主证件及辅助证件通用）
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
		// H/M + 8位数字
		matched, _ := regexp.MatchString(`^[HM]\d{8}$`, number)
		if !matched {
			return fmt.Errorf("港澳通行证号码格式错误（H/M + 8位数字）")
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
	// ── 辅助证件类型 ──────────────────────────────────────────────────────────
	case "90", "92", "95":
		// 1-2位字母 + 7位数字；W/WX 开头直接拒绝
		upper := strings.ToUpper(number)
		if strings.HasPrefix(upper, "W") {
			return fmt.Errorf("此证件号码为根据补充劳工计划签发给来港就业的工人的身份证号码，该证持有人不具有香港居留权及不合资格申请回乡证")
		}
		matched, _ := regexp.MatchString(`^[A-Za-z]{1,2}\d{7}$`, number)
		if !matched {
			return fmt.Errorf("证件号码格式错误（1-2位字母 + 7位数字）")
		}
	case "93":
		// 1位字母 + 9位数字（停发代码校验需另行调用 validate93IDNumber）
		matched, _ := regexp.MatchString(`^[A-Za-z]\d{9}$`, number)
		if !matched {
			return fmt.Errorf("台湾居民身份证号码格式错误（1位字母 + 9位数字）")
		}
	case "96", "97", "98":
		// 8位数字，以1/5/7开头
		matched, _ := regexp.MatchString(`^[157]\d{7}$`, number)
		if !matched {
			return fmt.Errorf("澳门居民身份证号码格式错误（8位数字，以1/5/7开头）")
		}
	}
	return nil
}

// validate93IDNumber 校验台湾居民身份证（93）：包含停发代码与出生日期联动检查
func validate93IDNumber(number, birthDate string) error {
	if number == "" {
		return nil
	}
	if err := validateIDNumber("93", number, ""); err != nil {
		return err
	}
	code := strings.ToUpper(number)[0]
	stopDates := map[byte]string{
		'L': "2010-12-25",
		'R': "2010-12-25",
		'S': "2010-12-25",
		'Y': "1974-01-01",
	}
	if stopDateStr, stopped := stopDates[code]; stopped && birthDate != "" {
		stop, err1 := time.Parse("2006-01-02", stopDateStr)
		birth, err2 := time.Parse("2006-01-02", birthDate)
		if err1 == nil && err2 == nil && birth.After(stop) {
			return fmt.Errorf("代码 %c 已于 %s 停止赋配，出生日期 %s 晚于停发日期，不允许使用此代码", code, stopDateStr, birthDate)
		}
	}
	return nil
}

// containsString 判断字符串是否在切片中
func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
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
	// 保存时同步写回 name
	if m.NameCn != "" {
		m.Name = m.NameCn
	} else if m.NameEn != "" {
		m.Name = m.NameEn
	}
	// 迁移旧辅助证件字段到 aux1（旧数据兼容）
	if m.Aux1Type == "" && m.AuxDocType != "" {
		m.Aux1Type = m.AuxDocType
		m.Aux1Number = m.AuxDocNumber
	}
}

// validateMember 验证成员数据合法性
func validateMember(m *models.Member) error {
	// 英文姓名必填
	if strings.TrimSpace(m.NameEn) == "" {
		return fmt.Errorf("英文姓名（name_en）为必填项")
	}
	// 国籍必填
	if strings.TrimSpace(m.Nationality) == "" {
		return fmt.Errorf("国籍（nationality）为必填项")
	}
	// 出生日期必填
	if strings.TrimSpace(m.BirthDate) == "" {
		return fmt.Errorf("出生日期（birth_date）为必填项")
	}
	// CHN：中英文姓名均必填
	if m.Nationality == "CHN" && strings.TrimSpace(m.NameCn) == "" {
		return fmt.Errorf("国籍为 CHN 时，中文姓名（name_cn）为必填项")
	}
	// 主证件校验
	if err := validateNationalityDocType(m.IdDocType, m.Nationality); err != nil {
		return fmt.Errorf("主证件：%s", err)
	}
	if err := validateIDNumber(m.IdDocType, m.IdDocNumber, m.Nationality); err != nil {
		return fmt.Errorf("主证件号码：%s", err)
	}

	// 辅助证件校验（按主证件类型判断）
	switch m.IdDocType {
	case "01", "91":
		// 不允许录入辅助证件
		if m.Aux1Type != "" || m.Aux1Number != "" || m.Aux2Type != "" || m.Aux2Number != "" {
			return fmt.Errorf("证件类型 %s 不允许录入辅助证件", m.IdDocType)
		}

	case "11":
		// 辅助证件1必须是 02，辅助证件2必须是 90/92/96/97
		if m.Aux1Type != "02" {
			return fmt.Errorf("港澳居民居住证持有人辅助证件1类型必须为 02（港澳居民来往内地通行证）")
		}
		if strings.TrimSpace(m.Aux1Number) == "" {
			return fmt.Errorf("辅助证件1号码为必填项")
		}
		if err := validateIDNumber("02", m.Aux1Number, m.Nationality); err != nil {
			return fmt.Errorf("辅助证件1号码：%s", err)
		}
		if !containsString([]string{"90", "92", "96", "97"}, m.Aux2Type) {
			return fmt.Errorf("港澳居民居住证持有人辅助证件2类型必须为 90/92/96/97 之一")
		}
		if strings.TrimSpace(m.Aux2Number) == "" {
			return fmt.Errorf("辅助证件2号码为必填项")
		}
		if err := validateIDNumber(m.Aux2Type, m.Aux2Number, m.Nationality); err != nil {
			return fmt.Errorf("辅助证件2号码：%s", err)
		}

	case "21":
		// 辅助证件1必须是 03，辅助证件2必须是 93
		if m.Aux1Type != "03" {
			return fmt.Errorf("台湾居民居住证持有人辅助证件1类型必须为 03（台湾居民来往大陆通行证）")
		}
		if strings.TrimSpace(m.Aux1Number) == "" {
			return fmt.Errorf("辅助证件1号码为必填项")
		}
		if err := validateIDNumber("03", m.Aux1Number, m.Nationality); err != nil {
			return fmt.Errorf("辅助证件1号码：%s", err)
		}
		if m.Aux2Type != "93" {
			return fmt.Errorf("台湾居民居住证持有人辅助证件2类型必须为 93（台湾居民身份证）")
		}
		if strings.TrimSpace(m.Aux2Number) == "" {
			return fmt.Errorf("辅助证件2号码为必填项")
		}
		if err := validate93IDNumber(m.Aux2Number, m.BirthDate); err != nil {
			return fmt.Errorf("辅助证件2号码：%s", err)
		}

	case "04":
		// 辅助证件类型必须是 94，且 proof_doc_type / proof_issue_country 必填
		if m.Aux1Type != "94" {
			return fmt.Errorf("中国护照持有人辅助证件类型必须为 94（证明文件）")
		}
		if strings.TrimSpace(m.Aux1Number) == "" {
			return fmt.Errorf("辅助证件（证明文件）号码为必填项")
		}
		if m.ProofDocType == "" {
			return fmt.Errorf("主证件为中国护照时，证明文件类型（proof_doc_type）为必填项")
		}
		validProofTypes := map[string]bool{
			"94RV": true, "94PV": true, "94PC": true, "94PE": true, "94NP": true,
		}
		if !validProofTypes[m.ProofDocType] {
			return fmt.Errorf("证明文件类型无效（须为 94RV/94PV/94PC/94PE/94NP 之一）")
		}
		if m.ProofIssueCountry == "" {
			return fmt.Errorf("主证件为中国护照时，签发国家（proof_issue_country）为必填项")
		}
		if m.ProofDocType == "94NP" {
			if m.ProofIssueCountry != "CHN" {
				return fmt.Errorf("证明文件类型为 94NP 时，签发国家必须为 CHN")
			}
			// 校验 94NP 号码构成规则：H + 12位数字
			matched, _ := regexp.MatchString(`^H\d{12}$`, m.Aux1Number)
			if !matched {
				return fmt.Errorf("94NP证件号码格式错误（H + 12位数字）")
			}
		} else {
			restricted := map[string]bool{"CHN": true, "HKG": true, "MAC": true, "TWN": true}
			if restricted[m.ProofIssueCountry] {
				return fmt.Errorf("证明文件签发国家不得为 CHN/HKG/MAC/TWN")
			}
		}

	case "05":
		// 隐藏辅助证件，不允许录入
		if m.Aux1Type != "" || m.Aux2Type != "" {
			return fmt.Errorf("外国护照（05）不需要录入辅助证件")
		}

	case "52":
		// 辅助证件类型只限 95/98
		if m.Aux1Type != "" {
			if !containsString([]string{"95", "98"}, m.Aux1Type) {
				return fmt.Errorf("港澳居民来往内地通行证（非中国籍）持有人辅助证件类型只限 95 或 98")
			}
			if err := validateIDNumber(m.Aux1Type, m.Aux1Number, m.Nationality); err != nil {
				return fmt.Errorf("辅助证件号码：%s", err)
			}
		}

	default:
		// 其他主证件类型：辅助证件可选，有填则校验
		if m.Aux1Type != "" && m.Aux1Number != "" {
			if m.Aux1Type == "93" {
				if err := validate93IDNumber(m.Aux1Number, m.BirthDate); err != nil {
					return fmt.Errorf("辅助证件1号码：%s", err)
				}
			} else {
				if err := validateIDNumber(m.Aux1Type, m.Aux1Number, m.Nationality); err != nil {
					return fmt.Errorf("辅助证件1号码：%s", err)
				}
			}
		}
		if m.Aux2Type != "" && m.Aux2Number != "" {
			if m.Aux2Type == "93" {
				if err := validate93IDNumber(m.Aux2Number, m.BirthDate); err != nil {
					return fmt.Errorf("辅助证件2号码：%s", err)
				}
			} else {
				if err := validateIDNumber(m.Aux2Type, m.Aux2Number, m.Nationality); err != nil {
					return fmt.Errorf("辅助证件2号码：%s", err)
				}
			}
		}
	}

	return nil
}

func (h *MemberHandler) List(c *gin.Context) {
	var members []models.Member
	if err := h.DB.Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败"})
		return
	}
	// 兼容旧数据：name_cn 为空时从 name 回填
	for i := range members {
		if members[i].NameCn == "" && members[i].Name != "" {
			members[i].NameCn = members[i].Name
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

	var member models.Member
	if err := h.DB.First(&member, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "成员不存在"})
		return
	}

	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求格式错误"})
		return
	}
	member.ID = uint(id)
	normalizeMember(&member)
	if err := validateMember(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.DB.Save(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, member)
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
