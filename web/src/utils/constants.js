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

/**
 * 18位中国居民身份证校验
 */
function validateChineseID(id) {
  if (!id || id.length !== 18) return '身份证号码必须为18位'
  const weights = [7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2]
  const checkCodes = '10X98765432'
  let sum = 0
  for (let i = 0; i < 17; i++) {
    const c = id.charCodeAt(i)
    if (c < 48 || c > 57) return '身份证号码前17位必须为数字'
    sum += (c - 48) * weights[i]
  }
  const expected = checkCodes[sum % 11]
  const last = id[17].toUpperCase()
  if (last !== expected) return `身份证校验码错误（期望 ${expected}，实际 ${last}）`
  const year = parseInt(id.substring(6, 10))
  const month = parseInt(id.substring(10, 12))
  const day = parseInt(id.substring(12, 14))
  if (month < 1 || month > 12 || day < 1 || day > 31) return '身份证日期部分无效'
  if (year < 1900 || year > new Date().getFullYear()) return '身份证年份无效'
  return null
}

/**
 * 前端证件号码校验
 * @param {string} docType
 * @param {string} number
 * @param {string} nationality
 * @param {Object} [opts] - extra options
 * @param {string} [opts.birthDate] - birth date for type 93 stop-date check
 * @param {string} [opts.proofDocType] - proof doc type for type 94
 * @returns {string|null} error message or null
 */
export function validateIDNumber(docType, number, nationality, opts = {}) {
  if (!number) return null
  switch (docType) {
    case '01':
    case '91':
      return validateChineseID(number)
    case '11': {
      let prefix = ''
      if (nationality === 'HKG') prefix = '810000'
      else if (nationality === 'MAC') prefix = '820000'
      if (prefix && !number.startsWith(prefix)) {
        return `港澳居民居住证号码必须以 ${prefix} 开头`
      }
      return validateChineseID(number)
    }
    case '21':
      if (!number.startsWith('830000')) return '台湾居民居住证号码必须以 830000 开头'
      return validateChineseID(number)
    case '31':
      if (number.length !== 18 || number[0] !== '9') return '外国人永久居留身份证必须以 9 开头且为18位'
      return validateChineseID(number)
    case '02':
      return /^[HM]\d{8}$/.test(number) ? null : '港澳通行证格式错误（H/M + 8位数字）'
    case '03':
      return /^\d{8}$/.test(number) ? null : '台湾居民来往大陆通行证必须为8位数字'
    case '04':
      return /^E\d{8}$/.test(number) || /^E[A-Za-z]\d{7}$/.test(number) ? null : '中国护照格式错误（E+8位数字 或 E+1字母+7位数字）'
    case '05':
      return number.length >= 6 && number.length <= 9 ? null : '外国护照号码长度必须为 6-9 位'
    case '52':
      return /^(HA|MA)\d{7}$/.test(number) ? null : '港澳通行证（非中国籍）格式错误（HA/MA + 7位数字）'
    // 辅助证件类型
    case '90':
    case '92':
    case '95':
      return validateHKMOID(number)
    case '93':
      return validateTaiwan93(number, opts.birthDate || '')
    case '94':
      return validateAux94Number(number, opts.proofDocType || '')
    case '96':
      return /^1\d{7}$/.test(number) ? null : '澳门居民身份证号码格式错误（8位数字，以1开头）'
    case '97':
      return /^5\d{7}$/.test(number) ? null : '澳门永久性居民身份证号码格式错误（8位数字，以5开头）'
    case '98':
      return /^7\d{7}$/.test(number) ? null : '澳门永久性居民身份证（外国籍）号码格式错误（8位数字，以7开头）'
    default:
      return null
  }
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
