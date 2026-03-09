# Schema Builder API Reference

The Schema Builder provides a database-agnostic way to manipulate tables. It features a clean, expressive API that works across all supported database systems (MySQL, PostgreSQL, SQLite, etc.), allowing you to define your database structure using simple code rather than raw SQL.

## Table of Contents

1. [Creating Tables](#creating-tables)
2. [Column Types & Modifiers](#column-types--modifiers)
3. [Indexes & Keys](#indexes--keys)
4. [Foreign Key Constraints](#foreign-key-constraints)
5. [Modifying Tables](#modifying-tables)
6. [Checking Existence](#checking-existence)
7. [Dropping Tables](#dropping-tables)
8. [Full Example: Blog System](#full-example-blog-system)

---

## Creating Tables

To create a new database table, use the `Schema.create` method. The callback function receives a `table` object (blueprint) used to define the columns.

```lua
Schema.create("users", function(table)
    table.id()                     -- primary key
    table.string("name", 100)      -- VARCHAR(100)
    table.string("email").unique() -- Unique constraint
    table.boolean("is_active").default(true)
    table.timestamps()             -- Adds created_at and updated_at
end)

```

### Common Table Options

| Method                | Description                                                                            |
| --------------------- | -------------------------------------------------------------------------------------- |
| `table.id()`          | Alias for an integer (real: sqlite, mysql: BIGINT, postgresql: BIGSERIAL) primary key. |
| `table.timestamps()`  | Adds nullable `created_at` and `updated_at` timestamps.                                |
| `table.softDeletes()` | Adds a `deleted_at` column for soft deletion support.                                  |

---

## Column Types & Modifiers

The builder supports a variety of column types that map to the underlying database driver.

### Available Types

```lua
Schema.create("products", function(table)
    table.string("name", 200)      -- String with length
    table.text("description")      -- Long text field
    table.integer("stock")         -- Integer
    table.decimal("price", 10, 2)  -- Decimal with precision/scale
    table.float("weight")          -- Float
    table.boolean("available")     -- Boolean (TINYINT/BOOL)
    table.json("metadata")         -- JSON field
    table.timestamp("published_at") -- Timestamp
end)

```

### Column Modifiers

You can chain modifiers to column definitions to add constraints or default values.

| Modifier           | Description                                                                                              | Example                                 |
| ------------------ | -------------------------------------------------------------------------------------------------------- | --------------------------------------- |
| `.autoIncrement()` | Alias for an auto-incrementing integer (real: sqlite, mysql: BIGINT, postgresql: BIGSERIAL) primary key. | `table.integer("id").autoIncrement()`   |
| `.nullable()`      | Allows NULL values to be inserted.                                                                       | `table.integer("age").nullable()`       |
| `.default(value)`  | Sets a default value for the column.                                                                     | `table.boolean("active").default(true)` |
| `.unsigned()`      | Sets integer columns to UNSIGNED.                                                                        | `table.integer("qty").unsigned()`       |
| `.unique()`        | Adds a unique index constraint.                                                                          | `table.string("sku").unique()`          |

---

## Indexes & Keys

The schema builder supports several types of indexes. You can define them fluently on the column or separately.

### Defining Indexes

You can pass a single column name or an array of columns for compound indexes. Optionally, you can specify a custom index name as the second argument.

```lua
Schema.create("products", function(table)
    -- Simple Index
    table.index({"name"})

    -- Compound Index with custom name
    table.index({"name", "available"}, "idx_products_name_available")

    -- Unique Index
    table.unique({"name"}, "uniq_products_name")
end)

```

### Standalone Index Operations

If you need to add an index without modifying the table structure (e.g., inside a migration file but outside a `create` block):

```lua
Schema.index("users", {"email", "is_active"}, "idx_users_email_active")
Schema.unique("products", {"sku"}, "uniq_products_sku")
Schema.dropIndex("users", "idx_old_index_name")

```

---

## Foreign Key Constraints

The builder provides a fluent syntax for defining foreign key constraints, including cascading actions.

```lua
Schema.create("orders", function(table)
    table.id()
    table.integer("user_id").unsigned()

    -- Define the constraint
    table.foreign("user_id")
         .references("id")
         .on("users")
         .onDelete("cascade") -- Options: cascade, set null, restrict
end)

```

### Many-to-Many (Pivot Table) Example

Creating a pivot table usually involves two foreign keys and a composite unique key.

```lua
Schema.create("order_products", function(table)
    table.integer("order_id").unsigned()
    table.integer("product_id").unsigned()

    table.foreign("order_id").references("id").on("orders").onDelete("cascade")
    table.foreign("product_id").references("id").on("products").onDelete("cascade")

    -- Prevent duplicate pairs
    table.unique({"order_id", "product_id"})
end)

```

---

## Modifying Tables

The `Schema.table` method is used to update existing tables. You can add new columns, rename them, or drop them.

### Adding & Modifying Columns

```lua
Schema.table("users", function(table)
    -- Add new columns
    table.string("phone", 20).nullable()
    table.timestamp("last_login").nullable()

    -- Add new indexes to existing columns
    table.index({"phone"})
end)

```

### Renaming & Dropping

_Note: Make sure to check database driver compatibility for renaming columns._

```lua
Schema.table("users", function(table)
    -- Rename a column
    table.renameColumn("name", "full_name")

    -- Drop a column
    table.dropColumn("old_field")

    -- Drop an index
    table.dropIndex("idx_old_index")
end)

```

### Direct Operations

For simple one-off changes, you can use direct Schema methods:

```lua
Schema.addColumn("users", "nickname", "VARCHAR(50)")
Schema.dropColumn("orders", "old_status")
Schema.renameColumn("products", "desc", "description")

```

---

## Checking Existence

You can easily check for the existence of tables or columns before running operations. This is useful for idempotent migrations.

```lua
if Schema.hasTable("users") then
    -- Table exists logic
end

if Schema.hasColumn("users", "email") then
    -- Column exists logic
end

```

**Example: Conditional Creation**

```lua
if not Schema.hasTable("logs") then
    Schema.create("logs", function(table)
        table.id()
        table.text("message")
    end)
end

```

---

## Dropping Tables

Methods to remove tables from the database.

```lua
-- Drop a table if it exists (prevents errors)
Schema.dropIfExists("temp_table")

-- Drop a table (throws error if table not found)
Schema.drop("old_table")

```

---

## Full Example: Blog System

Below is a comprehensive example of a database schema for a blogging platform, demonstrating relationships, constraints, and migrations.

### 1. Categories & Users

```lua
Schema.create("categories", function(table)
    table.id()
    table.string("name").unique()
    table.timestamps()
end)

-- (Users table assumed created previously)

```

### 2. Posts (With Relations)

```lua
Schema.create("posts", function(table)
    table.id()
    table.integer("user_id").unsigned()
    table.integer("category_id").unsigned().nullable()
    table.string("title")
    table.string("slug").unique()
    table.text("content")
    table.string("status").default("draft")
    table.timestamps()
    table.softDeletes()

    -- Constraints
    table.foreign("user_id").references("id").on("users").onDelete("cascade")
    table.foreign("category_id").references("id").on("categories").onDelete("set null")
end)

```

### 3. Comments

```lua
Schema.create("comments", function(table)
    table.id()
    table.integer("post_id").unsigned()
    table.text("content")
    table.timestamps()

    table.foreign("post_id").references("id").on("posts").onDelete("cascade")
end)

```

### 4. Incremental Migration (Settings)

This demonstrates how to evolve your schema safely.

```lua
-- Step 1: Create base table
if not Schema.hasTable("settings") then
    Schema.create("settings", function(table)
        table.id()
        table.string("key").unique()
        table.text("value").nullable()
        table.timestamps()
    end)
end

-- Step 2: Add 'group' column later
if Schema.hasTable("settings") and not Schema.hasColumn("settings", "group") then
    Schema.table("settings", function(table)
        table.string("group").default("general")
        table.index({"group", "key"})
    end)
end

```
