/**
 * Minimal reactive i18n store for FMPS.
 * Supports Chinese (zh) and English (en).
 * Language selection is persisted in localStorage under 'fmps_lang'.
 */
import { ref, computed } from 'vue'

const STORAGE_KEY = 'fmps_lang'

export const LANGUAGES = [
  { code: 'zh', label: '中文' },
  { code: 'en', label: 'English' }
]

const lang = ref(localStorage.getItem(STORAGE_KEY) || 'zh')

export function setLang(code) {
  lang.value = code
  localStorage.setItem(STORAGE_KEY, code)
}

const messages = {
  zh: {
    // App-wide
    appTitle: '家庭惩戒管理系统',
    appSubtitle: 'FMPS',

    // Login page
    loginUsername: '用户名',
    loginPassword: '密码',
    loginUsernamePlaceholder: '请输入用户名',
    loginPasswordPlaceholder: '请输入密码',
    loginButton: '登 录',
    loginDefaultHint: '默认账号：admin　/　默认密码：123456',
    loginFailed: '登录失败，请检查用户名和密码',
    loginUsernameRequired: '请输入用户名',
    loginPasswordRequired: '请输入密码',
    switchLang: 'English',

    // Layout / nav
    navDashboard: '仪表盘',
    navMembers: '家庭成员',
    navRules: '惩戒规则',
    navRecords: '惩戒记录',
    navCases: '案件管理',
    navSettings: '系统设置',
    logout: '退出登录',

    // Dashboard
    statMembers: '家庭成员数',
    statRules: '规则数量',
    statRecords: '记录总数',
    statPoints: '累计分值',
    recentRecords: '最近惩戒记录',
    colMemberName: '成员姓名',
    colViolation: '违规项目',
    colPoints: '分值',
    colNote: '备注',
    colOccurredAt: '发生时间',

    // Member form / validation messages
    invalidBirthDate: '出生日期无效（须早于今日且不超过100年前）',
    invalidIssueDate: '签发日期无效（须为今日及今日前20年内）',
    invalidExpiryDate: '证件有效期不符合规定',
    expiryBeforeToday: '有效期必须晚于当日',
    doc91Expired: '该成员已年满16周岁，不得录入主证件为91'
  },

  en: {
    // App-wide
    appTitle: 'Family Management & Penalty System',
    appSubtitle: 'FMPS',

    // Login page
    loginUsername: 'Username',
    loginPassword: 'Password',
    loginUsernamePlaceholder: 'Enter your username',
    loginPasswordPlaceholder: 'Enter your password',
    loginButton: 'Sign In',
    loginDefaultHint: 'Default account: admin　/　Default password: 123456',
    loginFailed: 'Login failed. Please check your username and password.',
    loginUsernameRequired: 'Username is required',
    loginPasswordRequired: 'Password is required',
    switchLang: '中文',

    // Layout / nav
    navDashboard: 'Dashboard',
    navMembers: 'Family Members',
    navRules: 'Penalty Rules',
    navRecords: 'Penalty Records',
    navCases: 'Cases',
    navSettings: 'Settings',
    logout: 'Sign Out',

    // Dashboard
    statMembers: 'Members',
    statRules: 'Rules',
    statRecords: 'Total Records',
    statPoints: 'Total Points',
    recentRecords: 'Recent Penalty Records',
    colMemberName: 'Member',
    colViolation: 'Violation',
    colPoints: 'Points',
    colNote: 'Notes',
    colOccurredAt: 'Date & Time',

    // Member form / validation messages
    invalidBirthDate: 'Invalid date of birth (must be before today and within the past 100 years)',
    invalidIssueDate: 'Invalid issue date (must be within the past 20 years up to today)',
    invalidExpiryDate: 'The document expiry date is invalid',
    expiryBeforeToday: 'Expiry date must be after today',
    doc91Expired: 'This member has already turned 16 and cannot use document type 91 as the primary ID'
  }
}

/**
 * Translate a message key using the current language.
 * Falls back to Chinese, then returns the key itself.
 */
export function t(key) {
  return messages[lang.value]?.[key] ?? messages['zh'][key] ?? key
}

/**
 * Vue composable — returns reactive helpers.
 */
export function useI18n() {
  const currentLang = lang
  const translate = (key) => messages[currentLang.value]?.[key] ?? messages['zh'][key] ?? key
  return { lang: currentLang, t: translate, setLang, LANGUAGES }
}
