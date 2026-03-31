/**
 * 系统常量：国家/地区代码、证件类型、年级、班级等
 */

// ISO 3166-1 alpha-3 国家/地区代码（精简版含中文名）
export const COUNTRIES = [
  { code: 'ABW', name: '阿鲁巴' },
  { code: 'AFG', name: '阿富汗' },
  { code: 'AGO', name: '安哥拉' },
  { code: 'AIA', name: '安圭拉' },
  { code: 'ALA', name: '奥兰群岛' },
  { code: 'ALB', name: '阿尔巴尼亚' },
  { code: 'AND', name: '安道尔' },
  { code: 'ARE', name: '阿拉伯联合酋长国' },
  { code: 'ARG', name: '阿根廷' },
  { code: 'ARM', name: '亚美尼亚' },
  { code: 'ASM', name: '美属萨摩亚' },
  { code: 'ATA', name: '南极洲' },
  { code: 'ATF', name: '法属南部领地' },
  { code: 'ATG', name: '安提瓜和巴布达' },
  { code: 'AUS', name: '澳大利亚' },
  { code: 'AUT', name: '奥地利' },
  { code: 'AZE', name: '阿塞拜疆' },
  { code: 'BDI', name: '布隆迪' },
  { code: 'BEL', name: '比利时' },
  { code: 'BEN', name: '贝宁' },
  { code: 'BES', name: '博纳尔、圣尤斯特歇斯和萨巴' },
  { code: 'BFA', name: '布基纳法索' },
  { code: 'BGD', name: '孟加拉国' },
  { code: 'BGR', name: '保加利亚' },
  { code: 'BHR', name: '巴林' },
  { code: 'BHS', name: '巴哈马' },
  { code: 'BIH', name: '波斯尼亚和黑塞哥维那' },
  { code: 'BLM', name: '圣巴泰勒米' },
  { code: 'BLR', name: '白俄罗斯' },
  { code: 'BLZ', name: '伯利兹' },
  { code: 'BMU', name: '百慕大' },
  { code: 'BOL', name: '玻利维亚' },
  { code: 'BRA', name: '巴西' },
  { code: 'BRB', name: '巴巴多斯' },
  { code: 'BRN', name: '文莱' },
  { code: 'BTN', name: '不丹' },
  { code: 'BVT', name: '布韦岛' },
  { code: 'BWA', name: '博茨瓦纳' },
  { code: 'CAF', name: '中非共和国' },
  { code: 'CAN', name: '加拿大' },
  { code: 'CCK', name: '科科斯（基林）群岛' },
  { code: 'CHE', name: '瑞士' },
  { code: 'CHL', name: '智利' },
  { code: 'CHN', name: '中国' },
  { code: 'CIV', name: '科特迪瓦' },
  { code: 'CMR', name: '喀麦隆' },
  { code: 'COD', name: '刚果（金）' },
  { code: 'COG', name: '刚果（布）' },
  { code: 'COK', name: '库克群岛' },
  { code: 'COL', name: '哥伦比亚' },
  { code: 'COM', name: '科摩罗' },
  { code: 'CPV', name: '佛得角' },
  { code: 'CRI', name: '哥斯达黎加' },
  { code: 'CUB', name: '古巴' },
  { code: 'CUW', name: '库拉索' },
  { code: 'CXR', name: '圣诞岛' },
  { code: 'CYM', name: '开曼群岛' },
  { code: 'CYP', name: '塞浦路斯' },
  { code: 'CZE', name: '捷克' },
  { code: 'DEU', name: '德国' },
  { code: 'DJI', name: '吉布提' },
  { code: 'DMA', name: '多米尼克' },
  { code: 'DNK', name: '丹麦' },
  { code: 'DOM', name: '多米尼加共和国' },
  { code: 'DZA', name: '阿尔及利亚' },
  { code: 'ECU', name: '厄瓜多尔' },
  { code: 'EGY', name: '埃及' },
  { code: 'ERI', name: '厄立特里亚' },
  { code: 'ESH', name: '西撒哈拉' },
  { code: 'ESP', name: '西班牙' },
  { code: 'EST', name: '爱沙尼亚' },
  { code: 'ETH', name: '埃塞俄比亚' },
  { code: 'FIN', name: '芬兰' },
  { code: 'FJI', name: '斐济' },
  { code: 'FLK', name: '福克兰群岛' },
  { code: 'FRA', name: '法国' },
  { code: 'FRO', name: '法罗群岛' },
  { code: 'FSM', name: '密克罗尼西亚联邦' },
  { code: 'GAB', name: '加蓬' },
  { code: 'GBR', name: '英国' },
  { code: 'GEO', name: '格鲁吉亚' },
  { code: 'GGY', name: '根西岛' },
  { code: 'GHA', name: '加纳' },
  { code: 'GIB', name: '直布罗陀' },
  { code: 'GIN', name: '几内亚' },
  { code: 'GLP', name: '瓜德罗普' },
  { code: 'GMB', name: '冈比亚' },
  { code: 'GNB', name: '几内亚比绍' },
  { code: 'GNQ', name: '赤道几内亚' },
  { code: 'GRC', name: '希腊' },
  { code: 'GRD', name: '格林纳达' },
  { code: 'GRL', name: '格陵兰' },
  { code: 'GTM', name: '危地马拉' },
  { code: 'GUF', name: '法属圭亚那' },
  { code: 'GUM', name: '关岛' },
  { code: 'GUY', name: '圭亚那' },
  { code: 'HKG', name: '中国香港' },
  { code: 'HMD', name: '赫德岛和麦克唐纳群岛' },
  { code: 'HND', name: '洪都拉斯' },
  { code: 'HRV', name: '克罗地亚' },
  { code: 'HTI', name: '海地' },
  { code: 'HUN', name: '匈牙利' },
  { code: 'IDN', name: '印度尼西亚' },
  { code: 'IMN', name: '马恩岛' },
  { code: 'IND', name: '印度' },
  { code: 'IOT', name: '英属印度洋领地' },
  { code: 'IRL', name: '爱尔兰' },
  { code: 'IRN', name: '伊朗' },
  { code: 'IRQ', name: '伊拉克' },
  { code: 'ISL', name: '冰岛' },
  { code: 'ISR', name: '以色列' },
  { code: 'ITA', name: '意大利' },
  { code: 'JAM', name: '牙买加' },
  { code: 'JEY', name: '泽西岛' },
  { code: 'JOR', name: '约旦' },
  { code: 'JPN', name: '日本' },
  { code: 'KAZ', name: '哈萨克斯坦' },
  { code: 'KEN', name: '肯尼亚' },
  { code: 'KGZ', name: '吉尔吉斯斯坦' },
  { code: 'KHM', name: '柬埔寨' },
  { code: 'KIR', name: '基里巴斯' },
  { code: 'KNA', name: '圣基茨和尼维斯' },
  { code: 'KOR', name: '韩国' },
  { code: 'KWT', name: '科威特' },
  { code: 'LAO', name: '老挝' },
  { code: 'LBN', name: '黎巴嫩' },
  { code: 'LBR', name: '利比里亚' },
  { code: 'LBY', name: '利比亚' },
  { code: 'LCA', name: '圣卢西亚' },
  { code: 'LIE', name: '列支敦士登' },
  { code: 'LKA', name: '斯里兰卡' },
  { code: 'LSO', name: '莱索托' },
  { code: 'LTU', name: '立陶宛' },
  { code: 'LUX', name: '卢森堡' },
  { code: 'LVA', name: '拉脱维亚' },
  { code: 'MAC', name: '中国澳门' },
  { code: 'MAF', name: '法属圣马丁' },
  { code: 'MAR', name: '摩洛哥' },
  { code: 'MCO', name: '摩纳哥' },
  { code: 'MDA', name: '摩尔多瓦' },
  { code: 'MDG', name: '马达加斯加' },
  { code: 'MDV', name: '马尔代夫' },
  { code: 'MEX', name: '墨西哥' },
  { code: 'MHL', name: '马绍尔群岛' },
  { code: 'MKD', name: '北马其顿' },
  { code: 'MLI', name: '马里' },
  { code: 'MLT', name: '马耳他' },
  { code: 'MMR', name: '缅甸' },
  { code: 'MNE', name: '黑山' },
  { code: 'MNG', name: '蒙古' },
  { code: 'MNP', name: '北马里亚纳群岛' },
  { code: 'MOZ', name: '莫桑比克' },
  { code: 'MRT', name: '毛里塔尼亚' },
  { code: 'MSR', name: '蒙特塞拉特' },
  { code: 'MTQ', name: '马提尼克' },
  { code: 'MUS', name: '毛里求斯' },
  { code: 'MWI', name: '马拉维' },
  { code: 'MYS', name: '马来西亚' },
  { code: 'MYT', name: '马约特' },
  { code: 'NAM', name: '纳米比亚' },
  { code: 'NCL', name: '新喀里多尼亚' },
  { code: 'NER', name: '尼日尔' },
  { code: 'NFK', name: '诺福克岛' },
  { code: 'NGA', name: '尼日利亚' },
  { code: 'NIC', name: '尼加拉瓜' },
  { code: 'NIU', name: '纽埃' },
  { code: 'NLD', name: '荷兰' },
  { code: 'NOR', name: '挪威' },
  { code: 'NPL', name: '尼泊尔' },
  { code: 'NRU', name: '瑙鲁' },
  { code: 'NZL', name: '新西兰' },
  { code: 'OMN', name: '阿曼' },
  { code: 'PAK', name: '巴基斯坦' },
  { code: 'PAN', name: '巴拿马' },
  { code: 'PCN', name: '皮特凯恩群岛' },
  { code: 'PER', name: '秘鲁' },
  { code: 'PHL', name: '菲律宾' },
  { code: 'PLW', name: '帕劳' },
  { code: 'PNG', name: '巴布亚新几内亚' },
  { code: 'POL', name: '波兰' },
  { code: 'PRI', name: '波多黎各' },
  { code: 'PRK', name: '朝鲜' },
  { code: 'PRT', name: '葡萄牙' },
  { code: 'PRY', name: '巴拉圭' },
  { code: 'PSE', name: '巴勒斯坦' },
  { code: 'PYF', name: '法属波利尼西亚' },
  { code: 'QAT', name: '卡塔尔' },
  { code: 'REU', name: '留尼汪' },
  { code: 'ROU', name: '罗马尼亚' },
  { code: 'RUS', name: '俄罗斯' },
  { code: 'RWA', name: '卢旺达' },
  { code: 'SAU', name: '沙特阿拉伯' },
  { code: 'SDN', name: '苏丹' },
  { code: 'SEN', name: '塞内加尔' },
  { code: 'SGP', name: '新加坡' },
  { code: 'SGS', name: '南乔治亚岛和南桑威奇群岛' },
  { code: 'SHN', name: '圣赫勒拿、阿森松和特里斯坦-达库尼亚' },
  { code: 'SJM', name: '斯瓦尔巴和扬马延' },
  { code: 'SLB', name: '所罗门群岛' },
  { code: 'SLE', name: '塞拉利昂' },
  { code: 'SLV', name: '萨尔瓦多' },
  { code: 'SMR', name: '圣马力诺' },
  { code: 'SOM', name: '索马里' },
  { code: 'SPM', name: '圣皮埃尔和密克隆' },
  { code: 'SRB', name: '塞尔维亚' },
  { code: 'SSD', name: '南苏丹' },
  { code: 'STP', name: '圣多美和普林西比' },
  { code: 'SUR', name: '苏里南' },
  { code: 'SVK', name: '斯洛伐克' },
  { code: 'SVN', name: '斯洛文尼亚' },
  { code: 'SWE', name: '瑞典' },
  { code: 'SWZ', name: '斯威士兰' },
  { code: 'SXM', name: '荷属圣马丁' },
  { code: 'SYC', name: '塞舌尔' },
  { code: 'SYR', name: '叙利亚' },
  { code: 'TCA', name: '特克斯和凯科斯群岛' },
  { code: 'TCD', name: '乍得' },
  { code: 'TGO', name: '多哥' },
  { code: 'THA', name: '泰国' },
  { code: 'TJK', name: '塔吉克斯坦' },
  { code: 'TKL', name: '托克劳' },
  { code: 'TKM', name: '土库曼斯坦' },
  { code: 'TLS', name: '东帝汶' },
  { code: 'TON', name: '汤加' },
  { code: 'TTO', name: '特立尼达和多巴哥' },
  { code: 'TUN', name: '突尼斯' },
  { code: 'TUR', name: '土耳其' },
  { code: 'TUV', name: '图瓦卢' },
  { code: 'TWN', name: '中国台湾' },
  { code: 'TZA', name: '坦桑尼亚' },
  { code: 'UGA', name: '乌干达' },
  { code: 'UKR', name: '乌克兰' },
  { code: 'UMI', name: '美国本土外小岛屿' },
  { code: 'URY', name: '乌拉圭' },
  { code: 'USA', name: '美国' },
  { code: 'UZB', name: '乌兹别克斯坦' },
  { code: 'VAT', name: '梵蒂冈' },
  { code: 'VCT', name: '圣文森特和格林纳丁斯' },
  { code: 'VEN', name: '委内瑞拉' },
  { code: 'VGB', name: '英属维尔京群岛' },
  { code: 'VIR', name: '美属维尔京群岛' },
  { code: 'VNM', name: '越南' },
  { code: 'VUT', name: '瓦努阿图' },
  { code: 'WLF', name: '瓦利斯和富图纳' },
  { code: 'WSM', name: '萨摩亚' },
  { code: 'YEM', name: '也门' },
  { code: 'ZAF', name: '南非' },
  { code: 'ZMB', name: '赞比亚' },
  { code: 'ZWE', name: '津巴布韦' }
]

// 证件类型
export const ID_DOC_TYPES = [
  { code: '01', name: '中国居民身份证' },
  { code: '91', name: '居民户口簿' },
  { code: '11', name: '港澳居民居住证' },
  { code: '21', name: '台湾居民居住证' },
  { code: '31', name: '2023版外国人永久居留身份证' },
  { code: '02', name: '港澳居民来往内地通行证' },
  { code: '03', name: '台湾居民来往大陆通行证' },
  { code: '04', name: '中国护照' },
  { code: '05', name: '外国护照' },
  { code: '52', name: '港澳居民来往内地通行证（非中国籍）' }
]

/**
 * 根据证件类型返回允许的国籍代码列表（空数组代表无限制）
 * @param {string} docType
 * @returns {string[]} 空数组=无限制，有值=白名单，特殊值'!'前缀=黑名单处理在外层
 */
export function getAllowedNationalities(docType) {
  switch (docType) {
    case '01':
    case '91':
    case '04':
      return ['CHN']
    case '11':
    case '02':
      return ['HKG', 'MAC']
    case '21':
    case '03':
      return ['TWN']
    case '31':
    case '05':
    case '52':
      return null // 特殊：排除 CHN/HKG/MAC/TWN
    default:
      return [] // 无限制
  }
}

const RESTRICTED_NATIONS = new Set(['CHN', 'HKG', 'MAC', 'TWN'])

/**
 * 根据证件类型过滤国家列表
 */
export function getAvailableCountries(docType) {
  const allowed = getAllowedNationalities(docType)
  if (!docType || (Array.isArray(allowed) && allowed.length === 0)) {
    return COUNTRIES
  }
  if (Array.isArray(allowed)) {
    return COUNTRIES.filter(c => allowed.includes(c.code))
  }
  // null => 排除模式
  return COUNTRIES.filter(c => !RESTRICTED_NATIONS.has(c.code))
}

/**
 * 校验国籍与证件类型是否相容
 */
export function validateNationalityDocType(docType, nationality) {
  if (!docType || !nationality) return null
  const allowed = getAllowedNationalities(docType)
  if (Array.isArray(allowed) && allowed.length === 0) return null
  if (Array.isArray(allowed)) {
    if (!allowed.includes(nationality)) {
      return `证件类型"${ID_DOC_TYPES.find(t => t.code === docType)?.name || docType}"要求国籍为：${allowed.join('/')}，当前国籍不符合`
    }
    return null
  }
  // null => 排除模式
  if (RESTRICTED_NATIONS.has(nationality)) {
    return `证件类型"${ID_DOC_TYPES.find(t => t.code === docType)?.name || docType}"要求国籍不得为 CHN/HKG/MAC/TWN`
  }
  return null
}

// Generic error returned for all ID number format/content violations
const ID_NUMBER_INVALID = '证件号码无效'

/**
 * 18位中国居民身份证格式校验（仅校验格式和校验码，不校验性别/出生日期）
 * @returns {boolean} true = valid format
 */
function validateChineseID(id) {
  if (!id || id.length !== 18) return false
  const weights = [7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2]
  const checkCodes = '10X98765432'
  let sum = 0
  for (let i = 0; i < 17; i++) {
    const c = id.charCodeAt(i)
    if (c < 48 || c > 57) return false
    sum += (c - 48) * weights[i]
  }
  const expected = checkCodes[sum % 11]
  const last = id[17].toUpperCase()
  if (last !== expected) return false
  const year = parseInt(id.substring(6, 10))
  const month = parseInt(id.substring(10, 12))
  const day = parseInt(id.substring(12, 14))
  if (month < 1 || month > 12 || day < 1 || day > 31) return false
  if (year < 1900 || year > new Date().getFullYear()) return false
  return true
}

/**
 * 前端证件号码校验
 * All errors return a generic '证件号码无效' to avoid exposing validation rules.
 * @param {string} docType
 * @param {string} number
 * @param {string} nationality
 * @param {Object} [opts] - extra options
 * @param {string} [opts.birthDate] - YYYY-MM-DD, used for birthdate consistency and type-93 stop-date
 * @param {string} [opts.proofDocType] - proof doc type for type 94
 * @param {string} [opts.gender] - '男' or '女', used for gender consistency check
 * @returns {string|null} '证件号码无效' or null
 */
export function validateIDNumber(docType, number, nationality, opts = {}) {
  if (!number) return null
  const INVALID = ID_NUMBER_INVALID
  switch (docType) {
    case '01':
    case '91':
      if (!validateChineseID(number)) return INVALID
      break
    case '11': {
      let prefix = ''
      if (nationality === 'HKG') prefix = '810000'
      else if (nationality === 'MAC') prefix = '820000'
      if (prefix && !number.startsWith(prefix)) return INVALID
      if (!validateChineseID(number)) return INVALID
      break
    }
    case '21':
      if (!number.startsWith('830000')) return INVALID
      if (!validateChineseID(number)) return INVALID
      break
    case '31':
      if (number.length !== 18 || number[0] !== '9') return INVALID
      if (!validateChineseID(number)) return INVALID
      break
    case '02':
      if (!/^[HM]\d{8}$/.test(number)) return INVALID
      if (nationality === 'HKG' && number[0] !== 'H') return INVALID
      if (nationality === 'MAC' && number[0] !== 'M') return INVALID
      return null
    case '03':
      return /^\d{8}$/.test(number) ? null : INVALID
    case '04':
      return /^E\d{8}$/.test(number) || /^E[A-Za-z]\d{7}$/.test(number) ? null : INVALID
    case '05':
      return number.length >= 6 && number.length <= 9 ? null : INVALID
    case '52':
      return /^(HA|MA)\d{7}$/.test(number) ? null : INVALID
    // 辅助证件类型
    case '90':
    case '92':
    case '95':
      return validateHKMOID(number) === null ? null : INVALID
    case '93': {
      const fmtErr = validateTaiwan93(number, opts.birthDate || '')
      if (fmtErr) return INVALID
      // Second character must be 1 (male) or 2 (female) and must match gender
      if (number.length >= 2) {
        const c = number[1]
        if (c !== '1' && c !== '2') return INVALID
        if (opts.gender) {
          const numGender = c === '1' ? '男' : '女'
          if (opts.gender !== numGender) return INVALID
        }
      }
      return null
    }
    case '94':
      return validateAux94Number(number, opts.proofDocType || '') === null ? null : INVALID
    case '96':
      return /^1\d{7}$/.test(number) ? null : INVALID
    case '97':
      return /^5\d{7}$/.test(number) ? null : INVALID
    case '98':
      return /^7\d{7}$/.test(number) ? null : INVALID
    default:
      return null
  }

  // Gender & birthdate consistency for Chinese-ID-format types (01/91/11/21/31)
  if (['01', '91', '11', '21', '31'].includes(docType) && number.length === 18) {
    // Digits 7-14 (index 6-13) = YYYYMMDD birthdate
    // Parse YYYY-MM-DD directly (avoid timezone issues with new Date())
    if (opts.birthDate) {
      const parts = opts.birthDate.split('-')
      if (parts.length === 3) {
        const expectedBirth = parts[0] + parts[1].padStart(2, '0') + parts[2].padStart(2, '0')
        if (number.substring(6, 14) !== expectedBirth) return INVALID
      }
    }
    // 17th digit (index 16): odd = male (男), even = female (女)
    if (opts.gender) {
      const ch = number[16]
      if (ch < '0' || ch > '9') return INVALID
      const d = parseInt(ch, 10)
      const idGender = d % 2 !== 0 ? '男' : '女'
      if (opts.gender !== idGender) return INVALID
    }
  }

  return null
}

/**
 * 校验香港/澳门类身份证（90/92/95）：1-2位字母 + 7位数字，禁止W或WX开头
 * @returns {string|null}
 */
function validateHKMOID(number) {
  if (!/^[A-Za-z]{1,2}\d{7}$/.test(number)) {
    return '证件号码格式错误（1-2位字母 + 7位数字）'
  }
  if (number[0].toUpperCase() === 'W') {
    return '此证件号码为根据补充劳工计划签发给来港就业的工人的身份证号码，该证持有人不具有香港居留权及不合资格申请回乡证'
  }
  return null
}

/**
 * 校验台湾居民身份证（93）：1位字母（台湾地区码）+ 9位数字
 * @param {string} number
 * @param {string} birthDate YYYY-MM-DD
 * @returns {string|null}
 */
function validateTaiwan93(number, birthDate) {
  if (number.length !== 10) {
    return '台湾居民身份证号码必须为10位（1位字母 + 9位数字）'
  }
  const prefix = number[0].toUpperCase()
  if (!TAIWAN_REGION_CODES.has(prefix)) {
    return `台湾居民身份证号码首位字母无效（${number[0]}）`
  }
  if (!/^\d{9}$/.test(number.substring(1))) {
    return '台湾居民身份证号码格式错误（1位字母 + 9位数字）'
  }
  const stopDateStr = TAIWAN_REGION_STOP_DATES[prefix]
  if (stopDateStr && birthDate) {
    const stopDate = new Date(stopDateStr)
    const bd = new Date(birthDate)
    if (!isNaN(bd.getTime()) && bd > stopDate) {
      return `台湾居民身份证地区码 ${prefix} 已于 ${stopDateStr} 停止赋配，出生日期晚于停发日期不允许使用此代码`
    }
  }
  return null
}

/**
 * 校验94类辅助证件号码（按 proofDocType）
 */
function validateAux94Number(number, proofDocType) {
  if (!number) return null
  if (proofDocType === '94NP') {
    return /^H\d{12}$/.test(number) ? null : '94NP证件号码格式错误（H + 12位数字）'
  }
  return null
}

// 辅助证件类型
export const AUX_DOC_TYPES = [
  { code: '02', name: '02-港澳居民来往内地通行证' },
  { code: '03', name: '03-台湾居民来往大陆通行证' },
  { code: '05', name: '05-外国护照' },
  { code: '90', name: '90-香港永久性居民身份证' },
  { code: '92', name: '92-香港居民身份证' },
  { code: '93', name: '93-台湾居民身份证' },
  { code: '94', name: '94-中国公民定居国外的证明文件' },
  { code: '95', name: '95-香港永久性居民身份证（外国籍）' },
  { code: '96', name: '96-澳门居民身份证' },
  { code: '97', name: '97-澳门永久性居民身份证' },
  { code: '98', name: '98-澳门永久性居民身份证（外国籍）' }
]

// 证明文件类型（主证件为04时的辅助证件94的补充字段）
export const PROOF_DOC_TYPES = [
  { code: '94RV', name: '94RV-定居签证' },
  { code: '94PV', name: '94PV-永久居留签证' },
  { code: '94PC', name: '94PC-永久居留卡' },
  { code: '94PE', name: '94PE-永久居留电子签证' },
  { code: '94NP', name: '94NP-国家移民管理局护照查询结果' }
]

// 台湾地区代码停发日期（零值表示仍在使用）
const TAIWAN_REGION_STOP_DATES = {
  L: '2010-12-25', // 台中县
  R: '2010-12-25', // 台南县
  S: '2010-12-25', // 高雄县
  Y: '1974-01-01'  // 阳明山管理局
}
// 所有有效台湾地区代码（含已停发）
const TAIWAN_REGION_CODES = new Set([
  'A','B','C','D','E','F','G','H','I','J','K','L','M','N','O',
  'P','Q','R','S','T','U','V','W','X','Y','Z'
])

// 年级选项
export const GRADES = [
  '一年级', '二年级', '三年级', '四年级', '五年级', '六年级',
  '初一', '初二', '初三',
  '高一', '高二', '高三'
]

// 班级选项（1-40）
export const CLASSES = Array.from({ length: 40 }, (_, i) => String(i + 1))

// ─── Date validation helpers ─────────────────────────────────────────────────

/**
 * Parse a YYYY-MM-DD string into { y, m, d }.  Returns null on failure.
 */
function parseYMD(dateStr) {
  if (!dateStr) return null
  const parts = String(dateStr).split('-')
  if (parts.length !== 3) return null
  const y = parseInt(parts[0], 10)
  const m = parseInt(parts[1], 10)
  const d = parseInt(parts[2], 10)
  if (isNaN(y) || isNaN(m) || isNaN(d)) return null
  if (m < 1 || m > 12 || d < 1 || d > 31) return null
  return { y, m, d }
}

/** Return today as { y, m, d } in local time. */
function getTodayYMD() {
  const now = new Date()
  return { y: now.getFullYear(), m: now.getMonth() + 1, d: now.getDate() }
}

/**
 * Compare two { y, m, d } objects.
 * Returns negative if a < b, 0 if equal, positive if a > b.
 */
function compareYMD(a, b) {
  if (a.y !== b.y) return a.y - b.y
  if (a.m !== b.m) return a.m - b.m
  return a.d - b.d
}

/** Return true if year is a leap year. */
function isLeapYear(year) {
  return (year % 4 === 0 && year % 100 !== 0) || year % 400 === 0
}

/** Return the number of days in a given month. */
function daysInMonth(year, month) {
  if (month === 2) return isLeapYear(year) ? 29 : 28
  return [4, 6, 9, 11].includes(month) ? 30 : 31
}

/**
 * Calculate the age at a given date.
 *
 * Feb 29 birthday handling (feb29IsMar1 = true, default):
 *   The age milestone falls on Mar 1 in non-leap years.
 *   This means: if asOf (the date being evaluated) is before Mar 1 of the current
 *   age-year for a Feb-29-born person, the age has not yet been reached.
 *   When asOf itself is Feb 29 of a leap year, Feb 29 < Mar 1, so the age is
 *   still counted as not yet reached — matching the regulation for type 01.
 *
 * @param {{ y, m, d }} birth
 * @param {{ y, m, d }} asOf
 * @param {boolean} [feb29IsMar1=true]
 * @returns {number}
 */
function calcAge(birth, asOf, feb29IsMar1 = true) {
  let age = asOf.y - birth.y
  const isBirthFeb29 = birth.m === 2 && birth.d === 29

  if (isBirthFeb29 && feb29IsMar1) {
    // Milestone is Mar 1; anything before Mar 1 means age not yet reached
    if (asOf.m < 3) age--
  } else {
    const bm = isBirthFeb29 ? 3 : birth.m
    const bd = isBirthFeb29 ? 1 : birth.d
    if (asOf.m < bm || (asOf.m === bm && asOf.d < bd)) age--
  }
  return age
}

/**
 * Add `years` years to a { y, m, d }, keeping the same month/day.
 * If the original date is Feb 29 and the target year has no Feb 29, use Mar 1.
 */
function addYearsSameDay(d, years) {
  const targetYear = d.y + years
  if (d.m === 2 && d.d === 29 && !isLeapYear(targetYear)) {
    return { y: targetYear, m: 3, d: 1 }
  }
  return { y: targetYear, m: d.m, d: d.d }
}

/**
 * Subtract one calendar day from { y, m, d }.
 */
function subtractOneDay(d) {
  let { y, m } = d
  let day = d.d - 1
  if (day < 1) {
    m--
    if (m < 1) { m = 12; y-- }
    day = daysInMonth(y, m)
  }
  return { y, m, d: day }
}

/**
 * Add `years` years then subtract 1 day (used by types 31/02/03/52/04).
 * Feb 29 issue date: addYearsSameDay may map to Mar 1, then -1 day = Feb 28.
 */
function addYearsMinusOneDay(d, years) {
  return subtractOneDay(addYearsSameDay(d, years))
}

/** Format { y, m, d } as YYYY-MM-DD string. */
function ymdToStr(d) {
  return `${d.y}-${String(d.m).padStart(2, '0')}-${String(d.d).padStart(2, '0')}`
}

// ─── Public date validation functions ────────────────────────────────────────

/**
 * Validate a birthdate string (YYYY-MM-DD).
 * Rules: must be strictly before today; must not be more than 100 years ago.
 * Returns null if valid, or an i18n error key string on failure.
 *
 * @param {string} dateStr
 * @returns {string|null}
 */
export function validateBirthDate(dateStr) {
  if (!dateStr) return null
  const bd = parseYMD(dateStr)
  if (!bd) return 'invalidBirthDate'
  const today = getTodayYMD()
  // Must be strictly before today
  if (compareYMD(bd, today) >= 0) return 'invalidBirthDate'
  // Must not be older than 100 years (person must be at most 100 years old today)
  const centuryAgo = { y: today.y - 100, m: today.m, d: today.d }
  if (compareYMD(bd, centuryAgo) < 0) return 'invalidBirthDate'
  return null
}

/**
 * Validate a document issue date string (YYYY-MM-DD).
 * Rules: must be today or earlier; must be within the past 20 years.
 * Returns null if valid, or an i18n error key string on failure.
 *
 * @param {string} dateStr
 * @returns {string|null}
 */
export function validateIssueDate(dateStr) {
  if (!dateStr) return null
  const id = parseYMD(dateStr)
  if (!id) return 'invalidIssueDate'
  const today = getTodayYMD()
  // Must not be in the future
  if (compareYMD(id, today) > 0) return 'invalidIssueDate'
  // Must be within the past 20 years
  const twentyYearsAgo = { y: today.y - 20, m: today.m, d: today.d }
  if (compareYMD(id, twentyYearsAgo) < 0) return 'invalidIssueDate'
  return null
}

/**
 * Validate the expiry date for a document.
 * Returns null if valid, or an i18n error key string on failure.
 *
 * @param {string} expiryDateStr   - YYYY-MM-DD
 * @param {string} docType         - e.g. '01', '91', '11', …
 * @param {string} birthDateStr    - YYYY-MM-DD
 * @param {string} issueDateStr    - YYYY-MM-DD
 * @param {string} role            - 'parent' | 'child'
 * @returns {string|null}
 */
export function validateExpiryDate(expiryDateStr, docType, birthDateStr, issueDateStr, role) {
  if (!expiryDateStr || !docType || !birthDateStr || !issueDateStr) return null

  const expiry = parseYMD(expiryDateStr)
  const birth = parseYMD(birthDateStr)
  const issue = parseYMD(issueDateStr)
  const today = getTodayYMD()

  if (!expiry || !birth || !issue) return 'invalidExpiryDate'

  // ── Special handling for type 91 ───────────────────────────────────────────
  // Expiry must equal the member's 16th birthday exactly.
  // Per the domain rule, a person born on Feb 29 will always have a leap-year
  // 16th birthday (birth.y + 16 is leap for all valid Feb-29 birth years in scope),
  // so their 16th birthday date is Feb 29 of that year — handled in _sixteenthBirthday.
  // If the member is already 16+ today, type 91 may not be used as primary ID.
  if (docType === '91') {
    const exp16 = _sixteenthBirthday(birth)
    if (compareYMD(today, exp16) >= 0) return 'doc91Expired'
    if (compareYMD(expiry, exp16) !== 0) return 'invalidExpiryDate'
    return null
  }

  // ── Generic check: expiry must be strictly after today ─────────────────────
  if (compareYMD(expiry, today) <= 0) return 'expiryBeforeToday'

  switch (docType) {
    case '01':  return _validateExpiry01(expiry, birth, issue, role)
    case '11':
    case '21':  return _validateExpiry1121(expiry, issue)
    case '31':
    case '02':  return _validateExpiry3102(expiry, birth, issue)
    case '03':
    case '52':  return _validateExpiry0352(expiry, issue)
    case '04':  return _validateExpiry04(expiry, birth, issue)
    case '05':  return _validateExpiry05(expiry, issue)
    default:    return null
  }
}

// ─── Per-type expiry validators ───────────────────────────────────────────────

/**
 * Type 01 — Resident ID.
 * Age tiers (at issue date, feb29 = not-yet-reached):
 *   <16  → 5 y  |  16-25 → 10 y  |  26-45 → 20 y  |  ≥46 → 2099-12-31 (long-term)
 * Expiry = issue + term, same month/day; Feb 29 issue → Mar 1 if target non-leap.
 */
function _validateExpiry01(expiry, birth, issue, role) {
  const age = calcAge(birth, issue, true) // feb29 milestone = Mar 1

  if (age >= 46) {
    // Long-term: the only accepted value is 2099-12-31
    if (expiry.y === 2099 && expiry.m === 12 && expiry.d === 31) return null
    return 'invalidExpiryDate'
  }

  let years
  if (age < 16)       years = 5
  else if (age < 26)  years = 10
  else                years = 20

  const expected = addYearsSameDay(issue, years)
  return compareYMD(expiry, expected) === 0 ? null : 'invalidExpiryDate'
}

/**
 * Types 11, 21 — fixed 5-year validity, same month/day.
 * Feb 29 issue → expiry is Mar 1 of the 5th year.
 */
function _validateExpiry1121(expiry, issue) {
  const expected = addYearsSameDay(issue, 5)
  return compareYMD(expiry, expected) === 0 ? null : 'invalidExpiryDate'
}

/**
 * Types 31, 02 — age-based 5/10 years, expiry = issue + term − 1 day.
 * Age < 18 at issue → 5 y; age ≥ 18 → 10 y.
 */
function _validateExpiry3102(expiry, birth, issue) {
  const age = calcAge(birth, issue, true)
  const years = age < 18 ? 5 : 10
  const expected = addYearsMinusOneDay(issue, years)
  return compareYMD(expiry, expected) === 0 ? null : 'invalidExpiryDate'
}

/**
 * Types 03, 52 — fixed 5-year validity, expiry = issue + 5 y − 1 day.
 */
function _validateExpiry0352(expiry, issue) {
  const expected = addYearsMinusOneDay(issue, 5)
  return compareYMD(expiry, expected) === 0 ? null : 'invalidExpiryDate'
}

/**
 * Type 04 — age-based 5/10 years, expiry = issue + term − 1 day.
 * Age < 16 at issue → 5 y; age ≥ 16 → 10 y.
 */
function _validateExpiry04(expiry, birth, issue) {
  const age = calcAge(birth, issue, true)
  const years = age < 16 ? 5 : 10
  const expected = addYearsMinusOneDay(issue, years)
  return compareYMD(expiry, expected) === 0 ? null : 'invalidExpiryDate'
}

/**
 * Type 05 — country-defined; expiry must not exceed issue + 10 years (same day).
 */
function _validateExpiry05(expiry, issue) {
  const maxExpiry = addYearsSameDay(issue, 10)
  if (compareYMD(expiry, maxExpiry) > 0) return 'invalidExpiryDate'
  return null
}

/**
 * Compute the 16th-birthday date for doc type 91.
 * Per the domain specification, a person born on Feb 29 will have their 16th
 * birthday on Feb 29 of (birth.y + 16).  The specification states this date
 * always exists for valid Feb-29 birth years used in practice (e.g. 2000→2016,
 * 2004→2020, 2008→2024 are all leap years).  The function trusts this domain
 * invariant rather than re-deriving it.
 */
function _sixteenthBirthday(birth) {
  if (birth.m === 2 && birth.d === 29) {
    return { y: birth.y + 16, m: 2, d: 29 }
  }
  return { y: birth.y + 16, m: birth.m, d: birth.d }
}

/**
 * Validate whether the gender and birth date embedded in a Chinese-format ID
 * (types 01/91/11/21/31) are consistent with the separately entered form fields.
 *
 * This is intentionally separate from validateIDNumber so that field-level
 * errors can be shown on the correct form field rather than on the ID number field.
 *
 * @param {string} docType   - e.g. '01', '91', '11', '21', '31'
 * @param {string} number    - the raw ID string
 * @param {string} gender    - '男' | '女' (can be empty)
 * @param {string} birthDate - YYYY-MM-DD (can be empty)
 * @returns {{ field: 'birth_date'|'gender', key: string }|null}
 */
export function validateIDConsistency(docType, number, gender, birthDate) {
  if (!number || !['01', '91', '11', '21', '31'].includes(docType)) return null
  if (number.length !== 18) return null

  // Birth date embedded in digits 7-14 (index 6-13) as YYYYMMDD
  if (birthDate) {
    const parts = String(birthDate).split('-')
    if (parts.length === 3) {
      const expected = parts[0] + parts[1].padStart(2, '0') + parts[2].padStart(2, '0')
      if (number.substring(6, 14) !== expected) {
        return { field: 'birth_date', key: 'birthMismatchId' }
      }
    }
  }

  // 17th digit (index 16): odd = male (男), even = female (女)
  if (gender) {
    const ch = number[16]
    if (ch >= '0' && ch <= '9') {
      const d = parseInt(ch, 10)
      const idGender = d % 2 !== 0 ? '男' : '女'
      if (gender !== idGender) {
        return { field: 'gender', key: 'genderMismatchId' }
      }
    }
  }

  return null
}

// ─── Document Issuer Authority ────────────────────────────────────────────────

/** Fixed issuer for doc type 31 (NIA) */
export const ISSUER_NIA = '中华人民共和国国家移民管理局'

/** Fixed issuer for doc types 02/03/52 (MPS EEA) */
export const ISSUER_MPS_EEA = '中华人民共和国出入境管理局'

/** Issuer for type 04 when issue date ≥ 2019-03-04 */
export const ISSUER_NIA_04 = '中华人民共和国国家移民管理局'

/** Issuer for type 04 when issue date ≤ 2019-03-03 */
export const ISSUER_MPS_EEA_04 = '公安部出入境管理局'

/** MFA Commissioner's Office in Hong Kong */
export const ISSUER_MFA_HK = '中华人民共和国外交部驻香港特别行政区特派员公署'

/** MFA Commissioner's Office in Macao */
export const ISSUER_MFA_MO = '中华人民共和国外交部驻澳门特别行政区特派员公署'

/**
 * Built-in offline seed list of PRC embassies and consulates (excluding consular agencies).
 * Used as a fallback when the app is offline.
 * Source: Ministry of Foreign Affairs of the People's Republic of China.
 */
export const EMBASSY_SEED_LIST = [
  // Africa — Embassies
  '中国驻阿尔及利亚大使馆',
  '中国驻安哥拉大使馆',
  '中国驻贝宁大使馆',
  '中国驻博茨瓦纳大使馆',
  '中国驻布隆迪大使馆',
  '中国驻喀麦隆大使馆',
  '中国驻佛得角大使馆',
  '中国驻中非共和国大使馆',
  '中国驻乍得大使馆',
  '中国驻科摩罗大使馆',
  '中国驻刚果共和国大使馆',
  '中国驻刚果民主共和国大使馆',
  '中国驻科特迪瓦大使馆',
  '中国驻吉布提大使馆',
  '中国驻埃及大使馆',
  '中国驻赤道几内亚大使馆',
  '中国驻厄立特里亚大使馆',
  '中国驻埃塞俄比亚大使馆',
  '中国驻加蓬大使馆',
  '中国驻冈比亚大使馆',
  '中国驻加纳大使馆',
  '中国驻几内亚大使馆',
  '中国驻几内亚比绍大使馆',
  '中国驻肯尼亚大使馆',
  '中国驻莱索托大使馆',
  '中国驻利比里亚大使馆',
  '中国驻利比亚大使馆',
  '中国驻马达加斯加大使馆',
  '中国驻马拉维大使馆',
  '中国驻马里大使馆',
  '中国驻毛里塔尼亚大使馆',
  '中国驻毛里求斯大使馆',
  '中国驻摩洛哥大使馆',
  '中国驻莫桑比克大使馆',
  '中国驻纳米比亚大使馆',
  '中国驻尼日尔大使馆',
  '中国驻尼日利亚大使馆',
  '中国驻卢旺达大使馆',
  '中国驻圣多美和普林西比大使馆',
  '中国驻塞内加尔大使馆',
  '中国驻塞舌尔大使馆',
  '中国驻塞拉利昂大使馆',
  '中国驻索马里大使馆',
  '中国驻南非大使馆',
  '中国驻南苏丹大使馆',
  '中国驻苏丹大使馆',
  '中国驻坦桑尼亚大使馆',
  '中国驻多哥大使馆',
  '中国驻突尼斯大使馆',
  '中国驻乌干达大使馆',
  '中国驻赞比亚大使馆',
  '中国驻津巴布韦大使馆',
  // Africa — Consulates-General
  '中国驻开普敦总领事馆',
  '中国驻德班总领事馆',
  '中国驻拉各斯总领事馆',
  // Americas — Embassies
  '中国驻安提瓜和巴布达大使馆',
  '中国驻阿根廷大使馆',
  '中国驻巴哈马大使馆',
  '中国驻巴巴多斯大使馆',
  '中国驻玻利维亚大使馆',
  '中国驻巴西大使馆',
  '中国驻加拿大大使馆',
  '中国驻智利大使馆',
  '中国驻哥伦比亚大使馆',
  '中国驻哥斯达黎加大使馆',
  '中国驻古巴大使馆',
  '中国驻多米尼克大使馆',
  '中国驻多米尼加共和国大使馆',
  '中国驻厄瓜多尔大使馆',
  '中国驻萨尔瓦多大使馆',
  '中国驻格林纳达大使馆',
  '中国驻危地马拉大使馆',
  '中国驻圭亚那大使馆',
  '中国驻海地大使馆',
  '中国驻洪都拉斯大使馆',
  '中国驻牙买加大使馆',
  '中国驻墨西哥大使馆',
  '中国驻尼加拉瓜大使馆',
  '中国驻巴拿马大使馆',
  '中国驻秘鲁大使馆',
  '中国驻苏里南大使馆',
  '中国驻特立尼达和多巴哥大使馆',
  '中国驻乌拉圭大使馆',
  '中国驻美国大使馆',
  '中国驻委内瑞拉大使馆',
  // Americas — Consulates-General
  '中国驻卡尔加里总领事馆',
  '中国驻蒙特利尔总领事馆',
  '中国驻多伦多总领事馆',
  '中国驻温哥华总领事馆',
  '中国驻芝加哥总领事馆',
  '中国驻休斯顿总领事馆',
  '中国驻洛杉矶总领事馆',
  '中国驻迈阿密总领事馆',
  '中国驻纽约总领事馆',
  '中国驻旧金山总领事馆',
  '中国驻圣保罗总领事馆',
  '中国驻里约热内卢总领事馆',
  '中国驻蒙特雷总领事馆',
  '中国驻瓜达拉哈拉总领事馆',
  // Asia & Middle East — Embassies
  '中国驻阿富汗大使馆',
  '中国驻亚美尼亚大使馆',
  '中国驻阿塞拜疆大使馆',
  '中国驻巴林大使馆',
  '中国驻孟加拉国大使馆',
  '中国驻文莱大使馆',
  '中国驻柬埔寨大使馆',
  '中国驻东帝汶大使馆',
  '中国驻格鲁吉亚大使馆',
  '中国驻印度大使馆',
  '中国驻印度尼西亚大使馆',
  '中国驻伊朗大使馆',
  '中国驻伊拉克大使馆',
  '中国驻以色列大使馆',
  '中国驻日本大使馆',
  '中国驻约旦大使馆',
  '中国驻哈萨克斯坦大使馆',
  '中国驻科威特大使馆',
  '中国驻吉尔吉斯斯坦大使馆',
  '中国驻老挝大使馆',
  '中国驻黎巴嫩大使馆',
  '中国驻马来西亚大使馆',
  '中国驻马尔代夫大使馆',
  '中国驻蒙古国大使馆',
  '中国驻缅甸大使馆',
  '中国驻尼泊尔大使馆',
  '中国驻朝鲜大使馆',
  '中国驻阿曼大使馆',
  '中国驻巴基斯坦大使馆',
  '中国驻巴勒斯坦大使馆',
  '中国驻菲律宾大使馆',
  '中国驻卡塔尔大使馆',
  '中国驻沙特阿拉伯大使馆',
  '中国驻新加坡大使馆',
  '中国驻韩国大使馆',
  '中国驻斯里兰卡大使馆',
  '中国驻叙利亚大使馆',
  '中国驻塔吉克斯坦大使馆',
  '中国驻泰国大使馆',
  '中国驻土耳其大使馆',
  '中国驻土库曼斯坦大使馆',
  '中国驻阿联酋大使馆',
  '中国驻乌兹别克斯坦大使馆',
  '中国驻越南大使馆',
  '中国驻也门大使馆',
  // Asia & Middle East — Consulates-General
  '中国驻阿拉木图总领事馆',
  '中国驻吉大港总领事馆',
  '中国驻登巴萨总领事馆',
  '中国驻棉兰总领事馆',
  '中国驻泗水总领事馆',
  '中国驻孟买总领事馆',
  '中国驻加尔各答总领事馆',
  '中国驻金奈总领事馆',
  '中国驻班加罗尔总领事馆',
  '中国驻伊斯坦布尔总领事馆',
  '中国驻福冈总领事馆',
  '中国驻名古屋总领事馆',
  '中国驻大阪总领事馆',
  '中国驻札幌总领事馆',
  '中国驻长崎总领事馆',
  '中国驻哥打基纳巴卢总领事馆',
  '中国驻古晋总领事馆',
  '中国驻槟城总领事馆',
  '中国驻曼德勒总领事馆',
  '中国驻卡拉奇总领事馆',
  '中国驻拉合尔总领事馆',
  '中国驻白沙瓦总领事馆',
  '中国驻宿务总领事馆',
  '中国驻釜山总领事馆',
  '中国驻光州总领事馆',
  '中国驻济州总领事馆',
  '中国驻宋卡总领事馆',
  '中国驻清迈总领事馆',
  '中国驻岘港总领事馆',
  '中国驻胡志明市总领事馆',
  // Europe — Embassies
  '中国驻阿尔巴尼亚大使馆',
  '中国驻奥地利大使馆',
  '中国驻白俄罗斯大使馆',
  '中国驻比利时大使馆',
  '中国驻波斯尼亚和黑塞哥维那大使馆',
  '中国驻保加利亚大使馆',
  '中国驻克罗地亚大使馆',
  '中国驻塞浦路斯大使馆',
  '中国驻捷克大使馆',
  '中国驻丹麦大使馆',
  '中国驻爱沙尼亚大使馆',
  '中国驻芬兰大使馆',
  '中国驻法国大使馆',
  '中国驻德国大使馆',
  '中国驻希腊大使馆',
  '中国驻匈牙利大使馆',
  '中国驻冰岛大使馆',
  '中国驻爱尔兰大使馆',
  '中国驻意大利大使馆',
  '中国驻拉脱维亚大使馆',
  '中国驻立陶宛大使馆',
  '中国驻卢森堡大使馆',
  '中国驻马耳他大使馆',
  '中国驻摩尔多瓦大使馆',
  '中国驻摩纳哥大使馆',
  '中国驻黑山大使馆',
  '中国驻荷兰大使馆',
  '中国驻北马其顿大使馆',
  '中国驻挪威大使馆',
  '中国驻波兰大使馆',
  '中国驻葡萄牙大使馆',
  '中国驻罗马尼亚大使馆',
  '中国驻俄罗斯大使馆',
  '中国驻圣马力诺大使馆',
  '中国驻塞尔维亚大使馆',
  '中国驻斯洛伐克大使馆',
  '中国驻斯洛文尼亚大使馆',
  '中国驻西班牙大使馆',
  '中国驻瑞典大使馆',
  '中国驻瑞士大使馆',
  '中国驻乌克兰大使馆',
  '中国驻英国大使馆',
  // Europe — Consulates-General
  '中国驻巴塞罗那总领事馆',
  '中国驻爱丁堡总领事馆',
  '中国驻曼彻斯特总领事馆',
  '中国驻法兰克福总领事馆',
  '中国驻汉堡总领事馆',
  '中国驻慕尼黑总领事馆',
  '中国驻杜塞尔多夫总领事馆',
  '中国驻日内瓦总领事馆',
  '中国驻佛罗伦萨总领事馆',
  '中国驻米兰总领事馆',
  '中国驻里昂总领事馆',
  '中国驻马赛总领事馆',
  '中国驻斯特拉斯堡总领事馆',
  '中国驻哥德堡总领事馆',
  '中国驻圣彼得堡总领事馆',
  '中国驻叶卡捷琳堡总领事馆',
  '中国驻哈巴罗夫斯克总领事馆',
  '中国驻海参崴总领事馆',
  '中国驻伊尔库茨克总领事馆',
  '中国驻克拉科夫总领事馆',
  '中国驻格但斯克总领事馆',
  '中国驻弗罗茨瓦夫总领事馆',
  '中国驻克卢日总领事馆',
  '中国驻阿姆斯特丹总领事馆',
  // Oceania — Embassies
  '中国驻澳大利亚大使馆',
  '中国驻斐济大使馆',
  '中国驻基里巴斯大使馆',
  '中国驻密克罗尼西亚大使馆',
  '中国驻新西兰大使馆',
  '中国驻巴布亚新几内亚大使馆',
  '中国驻萨摩亚大使馆',
  '中国驻所罗门群岛大使馆',
  '中国驻汤加大使馆',
  '中国驻瓦努阿图大使馆',
  // Oceania — Consulates-General
  '中国驻阿德莱德总领事馆',
  '中国驻布里斯班总领事馆',
  '中国驻墨尔本总领事馆',
  '中国驻珀斯总领事馆',
  '中国驻悉尼总领事馆',
  '中国驻奥克兰总领事馆',
  '中国驻基督城总领事馆'
]

const EMBASSY_STORAGE_KEY = 'fmps_embassy_list'

/**
 * Get the current embassy/consulate list.
 * Returns the cached list from localStorage if available, otherwise the seed list.
 * @returns {string[]}
 */
export function getEmbassyList() {
  try {
    const stored = localStorage.getItem(EMBASSY_STORAGE_KEY)
    if (stored) {
      const parsed = JSON.parse(stored)
      if (Array.isArray(parsed) && parsed.length > 0) {
        return parsed
      }
    }
  } catch {
    // ignore parse errors
  }
  return EMBASSY_SEED_LIST
}

/**
 * Attempt to refresh the embassy/consulate list from a remote source.
 * On success, updates localStorage cache and returns the new list.
 * On failure, silently falls back to the current list (seed or cached).
 * @returns {Promise<string[]>}
 */
export async function refreshEmbassyList() {
  try {
    const resp = await fetch(
      'https://en.wikipedia.org/w/api.php?action=parse&page=List_of_diplomatic_missions_of_China&prop=wikitext&format=json&origin=*',
      { signal: AbortSignal.timeout(10000) }
    )
    if (!resp.ok) throw new Error('HTTP ' + resp.status)
    const data = await resp.json()
    const wikitext = data?.parse?.wikitext?.['*'] || ''
    const missions = _parseEmbassiesFromWikitext(wikitext)
    if (missions.length > 10) {
      localStorage.setItem(EMBASSY_STORAGE_KEY, JSON.stringify(missions))
      return missions
    }
  } catch {
    // Network error or parse failure — fall through to seed
  }
  return getEmbassyList()
}

/**
 * Parse embassy/consulate names from Wikipedia wikitext.
 * Returns an array of Chinese name strings.
 * @param {string} wikitext
 * @returns {string[]}
 */
function _parseEmbassiesFromWikitext(wikitext) {
  const results = []
  const seen = new Set()
  const lines = wikitext.split('\n')
  for (const line of lines) {
    const m = line.match(/中国驻.+?(?:大使馆|总领事馆|领事馆)/)
    if (m) {
      const name = m[0].replace(/[[\]{}|]/g, '').trim()
      if (name && !seen.has(name)) {
        seen.add(name)
        results.push(name)
      }
    }
  }
  return results
}

/**
 * Get the issuer authority options for document type 04 (Chinese Passport).
 * The list varies based on issue date.
 *
 * @param {string} issueDate  - YYYY-MM-DD; may be empty while user is selecting
 * @returns {string[]} ordered array of issuer name strings
 */
export function getType04IssuerOptions(issueDate) {
  const embassies = getEmbassyList()
  const fixed = []
  if (issueDate) {
    const d = parseYMD(issueDate)
    // 2019-03-04: the date when the National Immigration Administration (NIA) was
    // established and took over passport issuance from the MPS Exit-Entry Administration.
    const cutoff = { y: 2019, m: 3, d: 4 }
    if (d && compareYMD(d, cutoff) >= 0) {
      fixed.push(ISSUER_NIA_04)
    } else if (d) {
      fixed.push(ISSUER_MPS_EEA_04)
    } else {
      fixed.push(ISSUER_NIA_04, ISSUER_MPS_EEA_04)
    }
  } else {
    // No date yet — offer both authority options
    fixed.push(ISSUER_NIA_04, ISSUER_MPS_EEA_04)
  }
  return [...fixed, ...embassies, ISSUER_MFA_HK, ISSUER_MFA_MO]
}

/**
 * Return the fixed issuer authority for a given doc type, or null if not fixed.
 * @param {string} docType
 * @returns {string|null}
 */
export function getFixedIssuer(docType) {
  if (docType === '31') return ISSUER_NIA
  if (['02', '03', '52'].includes(docType)) return ISSUER_MPS_EEA
  return null
}
