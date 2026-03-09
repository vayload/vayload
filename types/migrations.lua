---@meta

-- =========================
-- Schema
-- =========================

---@class Schema
Schema = {}

---@param tableName string
---@param callback fun(table: Blueprint)
---@return boolean
function Schema.create(tableName, callback) end

---@param tableName string
---@param callback fun(table: Blueprint)
---@return boolean
function Schema.table(tableName, callback) end

---@param tableName string
---@return boolean
function Schema.drop(tableName) end

---@param tableName string
---@return boolean
function Schema.dropIfExists(tableName) end

---@param tableName string
---@param columns string[]
---@param name? string
---@return boolean
function Schema.index(tableName, columns, name) end

---@param tableName string
---@param name string
---@return boolean
function Schema.dropIndex(tableName, name) end

---@param tableName string
---@param columns string[]
---@param name? string
---@return boolean
function Schema.unique(tableName, columns, name) end

---@param tableName string
---@param name string
---@return boolean
function Schema.dropUnique(tableName, name) end

---@param tableName string
---@param column string
---@param type string
---@return boolean
function Schema.addColumn(tableName, column, type) end

---@param tableName string
---@param column string
---@return boolean
function Schema.dropColumn(tableName, column) end

---@param tableName string
---@param from string
---@param to string
---@return boolean
function Schema.renameColumn(tableName, from, to) end

---@param tableName string
---@return boolean
function Schema.hasTable(tableName) end

---@param tableName string
---@param column string
---@return boolean
function Schema.hasColumn(tableName, column) end


-- =========================
-- Blueprint
-- =========================

---@class Blueprint
local Blueprint = {}

-- Column creators
---@return Column
function Blueprint:id() end

---@param name string
---@param length? integer
---@return Column
function Blueprint:string(name, length) end

---@param name string
---@return Column
function Blueprint:text(name) end

---@param name string
---@return Column
function Blueprint:integer(name) end

---@param name string
---@return Column
function Blueprint:bigInteger(name) end

---@param name string
---@return Column
function Blueprint:boolean(name) end

---@param name string
---@return Column
function Blueprint:timestamp(name) end

---@param name string
---@param precision? integer
---@param scale? integer
---@return Column
function Blueprint:decimal(name, precision, scale) end

---@param name string
---@return Column
function Blueprint:float(name) end

---@param name string
---@return Column
function Blueprint:json(name) end

-- Helpers
---@return nil
function Blueprint:timestamps() end

---@return nil
function Blueprint:softDeletes() end

-- Foreign keys
---@param name string
---@return ForeignColumn
function Blueprint:foreign(name) end

-- Indexes
---@param columns string[]
---@param name? string
---@return nil
function Blueprint:index(columns, name) end

---@param columns string[]
---@param name? string
---@return nil
function Blueprint:unique(columns, name) end

-- Alter table
---@param column string
---@return nil
function Blueprint:dropColumn(column) end

---@param from string
---@param to string
---@return nil
function Blueprint:renameColumn(from, to) end

---@param name string
---@return nil
function Blueprint:dropIndex(name) end


-- =========================
-- Column (chainable)
-- =========================

---@class Column
local Column = {}

---@return Column
function Column:nullable() end

---@param value any
---@return Column
function Column:default(value) end

---@return Column
function Column:unique() end

---@return Column
function Column:unsigned() end

---@return Column
function Column:autoIncrement() end

---@return Column
function Column:primary() end


-- =========================
-- ForeignColumn
-- =========================

---@class ForeignColumn
local ForeignColumn = {}

---@param column string
---@return ForeignColumn
function ForeignColumn:references(column) end

---@param table string
---@return ForeignColumn
function ForeignColumn:on(table) end

---@param action string
---@return ForeignColumn
function ForeignColumn:onDelete(action) end

---@return ForeignColumn
function ForeignColumn:nullable() end
