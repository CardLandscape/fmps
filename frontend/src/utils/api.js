// API wrapper for Wails Go backend calls
const go = () => window.go?.main?.App

export const api = {
  // Members
  getAllMembers: () => go().GetAllMembers(),
  getMember: (id) => go().GetMember(id),
  createMember: (member) => go().CreateMember(member),
  updateMember: (member) => go().UpdateMember(member),
  deleteMember: (id) => go().DeleteMember(id),
  setProtection: (memberId, hours) => go().SetProtection(memberId, hours),
  checkProtection: (memberId) => go().CheckProtection(memberId),

  // Clauses
  getAllClauses: () => go().GetAllClauses(),
  getClause: (id) => go().GetClause(id),
  createClause: (clause) => go().CreateClause(clause),
  updateClause: (clause) => go().UpdateClause(clause),
  deleteClause: (id) => go().DeleteClause(id),
  getClausesByCategory: (cat) => go().GetClausesByCategory(cat),

  // Templates
  getAllTemplates: () => go().GetAllTemplates(),
  getTemplate: (id) => go().GetTemplate(id),
  createTemplate: (tpl) => go().CreateTemplate(tpl),
  updateTemplate: (tpl) => go().UpdateTemplate(tpl),
  deleteTemplate: (id) => go().DeleteTemplate(id),

  // Cases
  getAllCases: () => go().GetAllCases(),
  getCase: (id) => go().GetCase(id),
  createCase: (c) => go().CreateCase(c),
  updateCaseStatus: (id, status) => go().UpdateCaseStatus(id, status),
  getCasesByMember: (memberId) => go().GetCasesByMember(memberId),
  getCaseStats: () => go().GetCaseStats(),
  startPunishment: (id) => go().StartPunishment(id),
  completePunishment: (id) => go().CompletePunishment(id),
  getCaseComments: (caseId) => go().GetCaseComments(caseId),
  addCaseComment: (comment) => go().AddCaseComment(comment),
  getRecentCases: (limit) => go().GetRecentCases(limit),

  // Appeals
  getAppealsByCase: (caseId) => go().GetAppealsByCase(caseId),
  createAppeal: (appeal) => go().CreateAppeal(appeal),
  reviewAppeal: (id, reviewerId, approved, comment) => go().ReviewAppeal(id, reviewerId, approved, comment),

  // Settings
  getSetting: (key) => go().GetSetting(key),
  setSetting: (key, value) => go().SetSetting(key, value),
  getAllSettings: () => go().GetAllSettings(),
}
