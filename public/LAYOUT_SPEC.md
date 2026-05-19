# Layout Spec Guidelines

## 1. Dashboard

**Route:** `/{project-slug}/dashboard`
**Purpose:** Overview of the project, KPIs, recent activity, and important notifications.

**Components:**

1. **Weekly Traffic Chart**

   * **Type:** Chart
   * **Data Source:** `activity_logs`
   * **Filters:** Last 7 days, Device type
   * **Description:** Shows the number of visits per day for the last week.

2. **Recent System Logs**

   * **Type:** Table
   * **Columns:** `timestamp`, `actor`, `action`, `resource_type`, `severity`
   * **Filters:** Actor, Resource type, Severity
   * **Pagination:** Cursor-based
   * **Order:** `timestamp` descending
   * **Description:** Quick view of the most recent system events.

3. **Important Notifications**

   * **Type:** Table / Summary
   * **Columns:** `title`, `type`, `status`, `created_at`, `read_at`
   * **Filters:** Status, Type
   * **Pagination:** Cursor-based
   * **Description:** Shows critical notifications for the project.

4. **Recent User Actions**

   * **Type:** Table
   * **Columns:** `actor`, `action`, `resource_type`, `resource_slug`, `timestamp`
   * **Filters:** Actor, Resource type, Date range
   * **Pagination:** Cursor-based
   * **Description:** Highlights the last actions performed by users.

---

## 2. Content Types

**Route:** `/{project-slug}/content-types/:contentType?`
**Purpose:** Manage content types (collections) and their fields.

**Components:**

1. **Content Type Cards**

   * **Type:** Card list
   * **Fields Displayed:** `name`, `slug`, `num_fields`, `type` (single/collection)
   * **Interactions:**

     * Click → Opens field types table
     * Delete → Only allowed if no entries exist
   * **Description:** Shows a summary of all content types in the project.

2. **Field Types Table**

   * **Type:** Table
   * **Columns:** `name`, `slug`, `field_type`, `required`, `unique`, `localized`, `default_value`, `help_text`, `display_order`
   * **Filters:** Field type, Required, Unique, Localized
   * **Pagination:** Offset-based
   * **Interactions:** Add new field, Edit field, Delete field
   * **Description:** Detailed view of all fields in a content type, with full customization.

---

## 3. Entries

**Route:** `/{project-slug}/entries/:contentType?/:entrySlug?`
**Purpose:** Manage and edit content entries.

**Components:**

1. **Entries Table**

   * **Columns:** `title`, `slug`, `status`, `author`, `locale`, `created_at`, `updated_at`, `published_at`
   * **Filters:** Status, Author, Locale, Collection
   * **Grouping:** By collection
   * **Pagination:** Cursor-based
   * **Interactions:**

     * Click → Open entry detail
     * New → Create new entry
     * Preview → If plugin available
   * **Description:** Full list of entries across the project, filterable by type.

2. **Entry Detail / Editor**

   * **Fields Source:** `entry_fields`
   * **Supports:** Multi-locale, versioning, relations (`entry_relations`)
   * **Interactions:**

     * Update fields
     * Preview
     * Publish / Unpublish
     * Track versions
   * **Description:** Single entry view with full editing capabilities and plugin support.

---

## 4. Media / Storage

**Route:** `/{project-slug}/media`
**Purpose:** File and media management.

**Components:**

1. **Media Browser**

   * **Views:** Cards, List
   * **Features:**

     * Infinite scroll
     * Folder tree navigation
     * Drag & drop
     * Select file or folder
   * **Filters:** MIME type, Media type, Folder, Search
   * **Sorting:** Size, Name, Type (only in List view)
   * **Columns (List View):** `name`, `mime_type`, `category`, `size`, `created_at`, `updated_at`
   * **Pagination:** Infinite scroll
   * **Description:** File explorer style view for media files with OS-like navigation and sorting.

---

## 5. Users

**Route:** `/{project-slug}/users`
**Purpose:** Manage project users and their access.

**Components:**

* **Users Table**

  * **Columns:** `username`, `email`, `phone`, `role`, `status`, `last_sign_in_at`, `created_at`
  * **Filters:** Role, Status, Email/Phone search
  * **Pagination:** Cursor-based
  * **Interactions:** Add, Edit, Delete users, Assign roles
  * **Description:** Full user management interface.

---

## 6. Roles

**Route:** `/{project-slug}/roles`
**Purpose:** Manage roles and permissions for the project.

**Components:**

* **Roles Table**

  * **Columns:** `name`, `description`, `is_system`, `created_at`
  * **Filters:** System roles only
  * **Pagination:** Offset-based
  * **Interactions:** Add, Edit, Delete roles, Assign permissions
  * **Description:** Allows defining RBAC for the project.

---

## 7. Audit Logs

**Route:** `/{project-slug}/audit`
**Purpose:** Review activity logs and system events.

**Components:**

* **Audit Table**

  * **Columns:** `timestamp`, `actor`, `action`, `resource_type`, `resource_slug`, `changes`, `severity`, `origin`
  * **Filters:** Actor, Resource type, Severity, Origin, Date range
  * **Pagination:** Cursor-based
  * **Description:** Full audit trail for project activities.

---

## 8. Settings

**Route:** `/{project-slug}/settings`
**Purpose:** Configure project and CMS settings.

**Components:**

* **Tabs:**

  1. **Domain**

     * Configure domain if supported.
  2. **LLM / AI Agent**

     * AI agent settings.
  3. **General**

     * Project-specific options, feature toggles.
* **Description:** Central hub for all project configuration options.

---

## 9. Integrations / Plugins

**Route:** `/{project-slug}/integrations`
**Purpose:** Manage plugins and external integrations.

**Components:**

* **Plugins Table**

  * **Columns:** `name`, `version`, `author_name`, `status`, `is_system`, `created_at`
  * **Filters:** Active / Inactive, System / Custom
  * **Pagination:** Offset-based
  * **Interactions:** Activate / Deactivate, Configure, View logs
* **Plugin Logs**

  * **Columns:** `hook_name`, `execution_time_ms`, `success`, `error_message`, `created_at`
  * **Filters:** Success / Failure, Hook name
  * **Pagination:** Offset-based
  * **Description:** Detailed logs for plugin executions and hook events.

