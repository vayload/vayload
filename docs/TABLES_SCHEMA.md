```sql
// Project tables
Table projects {
  id bigint [pk, note: "Primary key del usuario"]
  name varchar(125) [not null]
  slug varchar(125) [not null, unique]
  owner_id bigint [ref: > users.id]
  settings json
  locale varchar(4)
  created_at datetime
  updated_at datetime
}

Table project_members {
  project_id bigint [ref: > projects.id]
  user_id bigint [ref: > users.id]
  role_id bigint [ref: > roles.id]

  indexes {
    (project_id, user_id, role_id) [pk]
  }
}

// Users module
Table users {
  id bigint [pk, note: "Primary key del usuario"]
  first_name varchar(125)
  last_name varchar(125)
  username varchar(125)
  email varchar(255)
  password_hash varchar(255) [null]
  email_confirmed_at datetime
  invited_at datetime
  confirmation_token varchar(255)
  confirmation_sent_at datetime
  recovery_token varchar(255)
  recovery_sent_at datetime
  email_change_token_new varchar(255)
  email_change varchar(255)
  email_change_sent_at datetime
  phone varchar(255) [unique, note: "Número de teléfono"]
  phone_confirmed_at datetime
  phone_change varchar(255) [default: "''"]
  phone_change_token varchar(255) [default: "''"]
  phone_change_sent_at datetime
  confirmed_at datetime [note: "GENERATED ALWAYS AS LEAST(email_confirmed_at, phone_confirmed_at)"]
  email_change_token_current varchar(255) [default: "''"]
  email_change_confirm_status smallint [default: 0, note: "0 = pendiente, 1 = confirmado, 2 = cancelado"]
  otp_code VARCHAR(200)
  otp_sent_at DATETIME
  banned_until datetime
  reauthentication_token varchar(255) [default: "''"]
  reauthentication_sent_at datetime
  avatar_url VARCHAR(255) [null]
  metadata json
  settings json
  attributes json // for ABAC
  last_sign_in_at datetime
  is_super_admin boolean
  is_sso_user boolean
  created_at datetime [default: `CURRENT_TIMESTAMP`]
  updated_at datetime [default: `CURRENT_TIMESTAMP`]
  deleted_at datetime [null]
}

Table user_identities {
  id bigint [pk, note: "Primary key del usuario"]
  user_id bigint [ref: > users.id]
  identity_id varchar(125) [not null]
  data json
  created_at datetime
  updated_at datetime
}

Table sessions {
  id uuid [pk]
  user_id bigint [ref: > users.id]
  project_id bigint [null] // opcional: último proyecto activo
  ip_address varchar(64)
  user_agent text
  last_seen_at datetime
  expires_at datetime
  revoked_at datetime
  created_at datetime
}

Table refresh_tokens {
  id uuid [pk]
  user_id bigint [ref: > users.id]
  session_id uuid [ref: > sessions.id]
  token_hash varchar(255)
  expires_at datetime
  revoked_at datetime
  created_at datetime
}

// Authorization RBAC
Table roles {
  id bigint [pk]
  name varchar [unique, not null]
  description text
}

Table permissions {
  id bigint [pk]
  action varchar [not null]     // create, read, update, delete, publish
  resource varchar [not null]   // article, page, media
}

Table role_permissions {
  role_id bigint [ref: > roles.id]
  permission_id bigint [ref: > permissions.id]

  indexes {
    (role_id, permission_id) [pk, unique]
  }
}

Table user_roles {
  user_id bigint [ref: > users.id]
  role_id bigint [ref: > roles.id]

  indexes {
    (user_id, role_id) [pk, unique]
  }
}

// User module logs
Table audit_log_entries {
  id bigint [pk]
  payload json
  created_at timestamp
  ip_address varchar(64)
  user_agent text
  actor_id bigint
  action varchar(255)
}

// Notifications
Table notifications {
  id bigint [pk]
  title varchar(255)
  body text
  datetime datetime
  status enum("sent", "read", "unread", "dismissed", "failed")
  type varchar(255)
  actions json
  sent_at datetime
  read_at datetime
  user_id bigint [ref: > users.id]
  project_id bigint [ref: > projects.id]
  created_at datetime
  updated_at datetime
  deleted_at datetime [null]
}

// Logs
Table notification_logs {
  id	bigint
  type	varchar(100)
  sender	varchar(255)
  channel varchar(20)
  destination	varchar(255)
  url	varchar(255)
  request_body	text
  response_body	text
  response_status	int
  error	text
  notifications_id	varchar(50)
  created_at	datetime
}

// Collections for posts, products, etc
Table collections {
  id bigint
  project_id bigint [ref: > projects.id]
  name varchar(125) [not null]
  slug varchar(125) [not null]
  fields_schema json // {title: string, body: richtext}
  settings json
  created_at datetime
  deleted_at datetime
  Indexes {
    (project_id, slug) [unique]
  }
}

Table entries {
  id bigint
}

// Manage blobs (files)
Table file_objects {
  id bigint [pk]
  owner_id bigint [ref: > users.id]
  project_id bigint [ref: > projects.id]
  name varchar(255) [not null]
  mime_type varchar(100) [not null]
  category file_category [not null]
  size int [not null]
  sha256 varchar(32) [not null]
  key varchar(255) [not null]
  folder varchar(255)
  created_at datetime [not null]
  updated_at datetime [not null]

  Indexes {
    (project_id)
    (owner_id)
    (sha256)
  }
}
```
