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
	}
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

// validateMember 验证成员数据合法性
func validateMember(m *models.Member) error {
	// 国籍为 CHN 时，中文姓名必填
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
	// 辅助证件校验
	if err := validateNationalityDocType(m.AuxDocType, m.Nationality); err != nil {
		return fmt.Errorf("辅助证件：%s", err)
	}
	if err := validateIDNumber(m.AuxDocType, m.AuxDocNumber, m.Nationality); err != nil {
		return fmt.Errorf("辅助证件号码：%s", err)
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
