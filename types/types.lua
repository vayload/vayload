--[[
  Type definitions for Lua migrations
  This file serves as reference documentation for available types
  in the Lua migration environment. It is not executed.

  Version: 2.0
  Supported: MySQL, PostgreSQL, SQLite
]] --

---@class HTTPResponse
---@field status number HTTP status code (200, 404, 500, etc.)
---@field body string Response body as a string
---@field headers table<string, string> HTTP response headers

---@class HTTPHeaders
---@field headers table<string, string> HTTP headers as key-value pairs

---@class Json Global JSON API similar to JavaScript
---@field encode fun(value: any, space?: string|integer, replacer?: fun(key: string|integer|nil, value: any): any): string, string? Converts a Lua table/value into a JSON string. Returns JSON string or nil + error.
---@field decode fun(text: string, reviver?: fun(key: string|integer|nil, value: any): any): any, string? Parses a JSON string into Lua tables/values. Optionally transforms values using reviver.
Json = {}

---@class HttpClient
---@field get fun(url: string, headers?: HTTPHeaders): HTTPResponse?, string? Performs a GET request
---@field post fun(url: string, body: string, headers?: HTTPHeaders): HTTPResponse?, string? Performs a POST request
---@field put fun(url: string, body: string, headers?: HTTPHeaders): HTTPResponse?, string? Performs a PUT request
---@field patch fun(url: string, body: string, headers?: HTTPHeaders): HTTPResponse?, string? Performs a PATCH request
---@field delete fun(url: string, headers?: HTTPHeaders): HTTPResponse?, string? Performs a DELETE request

-- Register the httpclient for module
---@type HttpClient
require("vayload:http")

-- Available global utility functions:
-- print, log, toJSON, fromJSON, sleep, now

--[[
  -- HTTP request
  local response, err = Http.get("https://api.example.com/users")
  if not err then
    local data, parseErr = Json.decode(response.body)
    if not parseErr then
      print("Users retrieved: " .. #data)
    end
  end

  -- POST with JSON
  local jsonData = Json.encode({name = "John", email = "john@test.com"})
  local response, err = Http.post(
    "https://api.example.com/users",
    jsonData,
    "application/json"
  )

]] --

--[[
  MIGRATION EXAMPLES:

  -- Simple up migration
  Migration.up(function()
    DB.exec("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY, name VARCHAR(255))")
  end)

  -- Down migration with error handling
  Migration.down(function()
    local _, err = DB.exec("DROP TABLE IF EXISTS users")
    if err then
      log("ERROR", "Error dropping table: " .. err)
      return "Error dropping table: " .. err
    end
    return nil -- No error, success
  end)
]] --

---@class Migration
---@field up fun(closure: fun()) Executes the up migration
---@field down fun(closure: fun()): string? Executes the down migration, returns error message or nil if successful
Migration = {}
