-- Migration 001: initial schema
-- Applies to: members, cases, rules, records, penalty_points, settings
-- Note: the application uses GORM AutoMigrate; this file serves as a
--       human-readable reference and can be applied manually when needed.

-- ──────────────────────────────────────────────────────────────────────────────
-- Members table
-- ──────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS members (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    name                 TEXT,               -- legacy compatibility field (mirrors name_cn)
    name_cn              TEXT,               -- 中文姓名（CHN/HKG/MAC/TWN 国籍必填）
    name_en              TEXT,               -- 英文姓名（任何情况必填）
    role                 TEXT,               -- "parent" | "child"  (一经创建不可更改)
    avatar               TEXT,

    -- 基本信息
    gender               TEXT,               -- "男" | "女"
    nationality          TEXT,               -- ISO 3166-1 alpha-3（必填）
    birth_date           TEXT,               -- YYYY-MM-DD（必填）

    -- 主要证件（全部必填）
    id_doc_type          TEXT,               -- 01/91/11/21/31/02/03/04/05/52
    id_doc_number        TEXT,
    id_issue_date        TEXT,
    id_expiry_date       TEXT,
    id_issue_authority   TEXT,

    -- 辅助证件（旧字段，保留向下兼容）
    aux_doc_type         TEXT,
    aux_doc_number       TEXT,

    -- 辅助证件1
    aux1_doc_type        TEXT,
    aux1_doc_number      TEXT,

    -- 辅助证件2
    aux2_doc_type        TEXT,
    aux2_doc_number      TEXT,

    -- 主证件为04时的补充字段
    proof_doc_type       TEXT,               -- 94RV | 94PV | 94PC | 94PE | 94NP
    proof_issue_country  TEXT,               -- ISO alpha-3（94NP 时强制 CHN）

    -- 学籍信息
    school_name          TEXT,               -- 须含「小学」「中学」「大学」「学院」之一
    grade                TEXT,
    class_name           TEXT,
    class_teacher_name   TEXT,
    class_teacher_phone  TEXT,

    -- 外出权限
    outing_permission    TEXT,               -- "许可" | "不许可" | "受限"
    outing_dates         TEXT,               -- JSON array of YYYY-MM-DD strings
    outing_time_ranges   TEXT,               -- JSON array of {start,end} objects

    created_at           DATETIME,
    updated_at           DATETIME
);

-- ──────────────────────────────────────────────────────────────────────────────
-- Cases table
-- ──────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS cases (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    member_id            INTEGER,            -- legacy: kept for backward compat
    parent_member_id     INTEGER,            -- 家长成员 ID（必须为 role=parent 的成员）
    child_member_id      INTEGER,            -- 小孩成员 ID（必须为 role=child 的成员）
    title                TEXT,
    description          TEXT,
    punishment_process   TEXT,               -- 管道符分隔的惩罚步骤文本
    status               TEXT,               -- "pending" | "active" | "completed"
    start_time           DATETIME,
    created_at           DATETIME,
    updated_at           DATETIME,
    FOREIGN KEY (member_id)         REFERENCES members(id),
    FOREIGN KEY (parent_member_id)  REFERENCES members(id),
    FOREIGN KEY (child_member_id)   REFERENCES members(id)
);

-- ──────────────────────────────────────────────────────────────────────────────
-- Rules table
-- ──────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS rules (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    title       TEXT,
    description TEXT,
    points      INTEGER,
    created_at  DATETIME,
    updated_at  DATETIME
);

-- ──────────────────────────────────────────────────────────────────────────────
-- Records table
-- ──────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS records (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    member_id   INTEGER,
    rule_id     INTEGER,
    points      INTEGER,
    note        TEXT,
    created_at  DATETIME,
    updated_at  DATETIME,
    FOREIGN KEY (member_id) REFERENCES members(id),
    FOREIGN KEY (rule_id)   REFERENCES rules(id)
);

-- ──────────────────────────────────────────────────────────────────────────────
-- Penalty points table (case-level deductions)
-- ──────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS penalty_points (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    case_id       INTEGER,
    member_id     INTEGER,
    rule_text     TEXT,
    score_delta   INTEGER,
    reason        TEXT,
    revoked       BOOLEAN DEFAULT 0,
    revoked_at    DATETIME,
    revoke_reason TEXT,
    created_at    DATETIME,
    updated_at    DATETIME,
    FOREIGN KEY (case_id)   REFERENCES cases(id),
    FOREIGN KEY (member_id) REFERENCES members(id)
);

-- ──────────────────────────────────────────────────────────────────────────────
-- Settings table
-- ──────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS settings (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    key         TEXT UNIQUE,
    value       TEXT,
    created_at  DATETIME,
    updated_at  DATETIME
);

-- Seed default authorization password
INSERT OR IGNORE INTO settings (key, value) VALUES ('authorization_password', '123456');

-- ──────────────────────────────────────────────────────────────────────────────
-- Backfill helpers (run once after first migration on existing data)
-- ──────────────────────────────────────────────────────────────────────────────
-- Copy legacy name → name_cn where name_cn is still empty
UPDATE members SET name_cn = name WHERE name_cn IS NULL OR name_cn = '';
