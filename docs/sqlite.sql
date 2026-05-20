PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA temp_store = MEMORY;
PRAGMA cache_size = -20000;

-- ------------- USERS MODULE -------------
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  email TEXT UNIQUE,
  phone TEXT UNIQUE,
  password_hash TEXT NOT NULL,

  first_name TEXT,
  last_name TEXT,
  avatar_url TEXT,

  email_confirmed_at DATETIME,
  phone_confirmed_at DATETIME,
  confirmed_at DATETIME,

  confirmation_token TEXT,
  recovery_token TEXT,

  email_change_token TEXT,
  phone_change_token TEXT,

  otp_code TEXT,

  confirmation_sent_at DATETIME,
  recovery_sent_at DATETIME,

  email_change_sent_at DATETIME,
  phone_change_sent_at DATETIME,
  otp_sent_at DATETIME,

  email_change TEXT,
  phone_change TEXT,

  banned_until DATETIME,
  last_sign_in_at DATETIME,

  metadata JSON DEFAULT '{}',
  settings JSON DEFAULT '{}',
  attributes JSON DEFAULT '{}',

  is_super_admin INTEGER NOT NULL DEFAULT 0,
  is_sso_user INTEGER NOT NULL DEFAULT 0,

  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME,
  deleted_at DATETIME
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

CREATE TABLE IF NOT EXISTS projects (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  owner_id INTEGER NOT NULL,
  settings JSON DEFAULT '{}',
  locale TEXT DEFAULT 'en',

  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME,

  FOREIGN KEY(owner_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_projects_owner ON projects(owner_id);

CREATE TABLE IF NOT EXISTS roles (
  id INTEGER PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  description TEXT,
  is_system INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS permissions (
  id INTEGER PRIMARY KEY,
  action TEXT NOT NULL,
  resource TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  UNIQUE(action, resource)
);

CREATE TABLE IF NOT EXISTS role_permissions (
  role_id INTEGER NOT NULL,
  permission_id INTEGER NOT NULL,

  PRIMARY KEY(role_id, permission_id),

  FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE,
  FOREIGN KEY(permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS project_members (
  project_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  role_id INTEGER NOT NULL,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  PRIMARY KEY(project_id, user_id, role_id),

  FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY(role_id) REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS user_roles (
  user_id INTEGER NOT NULL,
  role_id INTEGER NOT NULL,
  project_id INTEGER,

  PRIMARY KEY(user_id, role_id, project_id),

  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY(role_id) REFERENCES roles(id),
  FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_identities (
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL,

  provider TEXT NOT NULL,
  provider_user_id TEXT NOT NULL,

  data JSON DEFAULT '{}',

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME,

  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,

  UNIQUE(provider, provider_user_id)
);

CREATE TABLE IF NOT EXISTS sessions (
  id TEXT PRIMARY KEY,
  user_id INTEGER NOT NULL,
  project_id INTEGER,

  ip_address TEXT,
  user_agent TEXT,

  last_seen_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  expires_at DATETIME NOT NULL,
  revoked_at DATETIME,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE INDEX idx_sessions_user ON sessions(user_id);

CREATE TABLE IF NOT EXISTS refresh_tokens (
  id TEXT PRIMARY KEY,

  token_hash TEXT UNIQUE NOT NULL,
  user_id INTEGER NOT NULL,

  family_id TEXT NOT NULL,

  session_id TEXT,
  parent_id TEXT,

  used_at DATETIME,
  revoked_at DATETIME,
  revoked_reason TEXT,

  expires_at DATETIME NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY(session_id) REFERENCES sessions(id) ON DELETE CASCADE,
  FOREIGN KEY(parent_id) REFERENCES refresh_tokens(id) ON DELETE SET NULL
);

CREATE INDEX idx_refresh_user ON refresh_tokens(user_id);

-- ------------- COLLECTIONS MODULE -------------
CREATE TABLE IF NOT EXISTS collections (
  id INTEGER PRIMARY KEY,
  project_id INTEGER NOT NULL,

  name TEXT NOT NULL,
  slug TEXT NOT NULL,

  description TEXT,
  icon TEXT,

  is_system INTEGER DEFAULT 0,
  supports_versioning INTEGER DEFAULT 0,
  supports_localization INTEGER DEFAULT 0,

  settings JSON DEFAULT '{}',

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME,
  deleted_at DATETIME,

  FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,

  UNIQUE(project_id, slug)
);

CREATE TABLE IF NOT EXISTS  collection_fields (
  id INTEGER PRIMARY KEY,

  collection_id INTEGER NOT NULL,

  name TEXT NOT NULL,
  slug TEXT NOT NULL,

  field_type TEXT NOT NULL,

  is_required INTEGER DEFAULT 0,
  is_unique INTEGER DEFAULT 0,
  is_localized INTEGER DEFAULT 0,
  is_indexed INTEGER DEFAULT 0,

  default_value TEXT,

  validation_rules JSON DEFAULT '{}',
  options JSON DEFAULT '{}',

  display_order INTEGER DEFAULT 0,

  help_text TEXT,
  placeholder TEXT,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME,

  FOREIGN KEY(collection_id) REFERENCES collections(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS  entries (
  id INTEGER PRIMARY KEY,

  project_id INTEGER NOT NULL,
  collection_id INTEGER NOT NULL,
  author_id INTEGER,

  status TEXT NOT NULL DEFAULT 'draft',

  version INTEGER NOT NULL DEFAULT 1,
  parent_version_id INTEGER,

  locale TEXT DEFAULT 'en',

  title TEXT,
  slug TEXT,

  published_at DATETIME,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME,
  deleted_at DATETIME,

  FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
  FOREIGN KEY(collection_id) REFERENCES collections(id) ON DELETE CASCADE,
  FOREIGN KEY(author_id) REFERENCES users(id),
  FOREIGN KEY(parent_version_id) REFERENCES entries(id)
);

CREATE INDEX idx_entries_collection ON entries(collection_id);
CREATE INDEX idx_entries_project ON entries(project_id);

CREATE TABLE IF NOT EXISTS  entry_fields (
  id INTEGER PRIMARY KEY,

  entry_id INTEGER NOT NULL,
  field_id INTEGER NOT NULL,

  locale TEXT,

  value TEXT,
  value_json JSON,
  value_number REAL,
  value_boolean INTEGER,
  value_date DATETIME,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME,

  FOREIGN KEY(entry_id) REFERENCES entries(id) ON DELETE CASCADE,
  FOREIGN KEY(field_id) REFERENCES collection_fields(id) ON DELETE CASCADE
);

CREATE INDEX idx_entry_fields_entry ON entry_fields(entry_id);

CREATE TABLE IF NOT EXISTS  entry_relations (
  source_entry_id INTEGER NOT NULL,
  target_entry_id INTEGER NOT NULL,
  field_id INTEGER NOT NULL,

  relation_order INTEGER DEFAULT 0,

  PRIMARY KEY(source_entry_id, target_entry_id, field_id),

  FOREIGN KEY(source_entry_id) REFERENCES entries(id) ON DELETE CASCADE,
  FOREIGN KEY(target_entry_id) REFERENCES entries(id) ON DELETE CASCADE,
  FOREIGN KEY(field_id) REFERENCES collection_fields(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS  folder_objects (
  id TEXT PRIMARY KEY,

  owner_id INTEGER,
  project_id INTEGER,

  parent_id TEXT,

  name TEXT NOT NULL,
  path TEXT NOT NULL,

  depth INTEGER DEFAULT 0,

  file_count INTEGER DEFAULT 0,
  subfolder_count INTEGER DEFAULT 0,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME,

  FOREIGN KEY(parent_id) REFERENCES folder_objects(id) ON DELETE CASCADE,
  FOREIGN KEY(owner_id) REFERENCES users(id),
  FOREIGN KEY(project_id) REFERENCES projects(id)
);

CREATE TABLE IF NOT EXISTS  file_objects (
  id TEXT PRIMARY KEY,

  owner_id INTEGER,
  project_id INTEGER,

  folder_id TEXT,

  name TEXT NOT NULL,
  mime_type TEXT,
  category TEXT,

  size INTEGER,

  provider TEXT,
  provider_key TEXT,

  metadata JSON DEFAULT '{}',

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME,

  FOREIGN KEY(folder_id) REFERENCES folder_objects(id) ON DELETE SET NULL,
  FOREIGN KEY(owner_id) REFERENCES users(id),
  FOREIGN KEY(project_id) REFERENCES projects(id)
);

CREATE TABLE IF NOT EXISTS  file_usage (
  file_id TEXT NOT NULL,
  entry_id INTEGER NOT NULL,
  field_id INTEGER NOT NULL,

  usage_order INTEGER DEFAULT 0,

  PRIMARY KEY(file_id, entry_id, field_id),

  FOREIGN KEY(file_id) REFERENCES file_objects(id) ON DELETE CASCADE,
  FOREIGN KEY(entry_id) REFERENCES entries(id) ON DELETE CASCADE,
  FOREIGN KEY(field_id) REFERENCES collection_fields(id)
);

CREATE TABLE IF NOT EXISTS  plugins (
  id INTEGER PRIMARY KEY,

  name TEXT NOT NULL,
  version TEXT,

  description TEXT,

  author_name TEXT,
  author_email TEXT,
  author_homepage TEXT,

  entry_point TEXT,

  dependencies JSON DEFAULT '[]',
  manifest JSON DEFAULT '{}',

  checksum TEXT,

  is_active INTEGER DEFAULT 0,
  is_system INTEGER DEFAULT 0,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME
);

CREATE TABLE IF NOT EXISTS  plugin_hooks_log (
  id INTEGER PRIMARY KEY,

  plugin_id INTEGER,

  hook_name TEXT,

  execution_time_ms INTEGER,

  success INTEGER DEFAULT 1,
  error_message TEXT,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY(plugin_id) REFERENCES plugins(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS  notifications (
  id INTEGER PRIMARY KEY,

  title TEXT NOT NULL,
  body TEXT,

  type TEXT,

  status TEXT NOT NULL DEFAULT 'unread'
  CHECK(status IN ('sent','read','unread','dismissed','failed')),

  actions JSON DEFAULT '[]',

  sent_at DATETIME,
  read_at DATETIME,

  user_id INTEGER,
  project_id INTEGER,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME,
  deleted_at DATETIME,

  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS  notification_logs (
  id INTEGER PRIMARY KEY,

  type TEXT,
  sender TEXT,

  channel TEXT,
  destination TEXT,

  url TEXT,

  request_body TEXT,
  response_body TEXT,

  response_status INTEGER,

  error TEXT,

  notifications_id INTEGER,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY(notifications_id) REFERENCES notifications(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS  activity_logs (
  id INTEGER PRIMARY KEY,

  project_id INTEGER,
  actor_id INTEGER,

  action TEXT NOT NULL,

  resource_type TEXT,
  resource_id INTEGER,
  resource_slug TEXT,

  message TEXT,

  changes JSON DEFAULT '{}',
  payload JSON DEFAULT '{}',

  origin TEXT,
  severity TEXT,

  request_id TEXT,

  ip_address TEXT,
  user_agent TEXT,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY(project_id) REFERENCES projects(id),
  FOREIGN KEY(actor_id) REFERENCES users(id)
);

CREATE INDEX idx_activity_project ON activity_logs(project_id);
