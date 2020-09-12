--[[

  AvorionControl - data/scripts/lib/avocontrol.utils.lua
  -----------------------------

  Library file for the script mods that are provided by the AvorionControl
  bot. Includes a number of helper functions including (most importantly) the 
  function used to grab and verify data from the Server() object.

  License: WTFPL
  Info: https://en.wikipedia.org/wiki/WTFPL

]]

-- FileExists returns true if a file exists (and is a file)
--
-- Returns:
--  @1    Boolean
function FileExists(name)
  local f = io.open(name,"r")
  
  if f ~= nil then
      io.close(f)
      return true
  end

  return false
end


-- TblLen returns a shallow count of entries in a table (#table only returns
-- the highest index) or false if the variable provided is not a table.
--
-- Returns:
--  @1    Integer or Boolean
function TblLen(t)
  if type(t) ~= "table" then
    print("avocontrol-utils: TblLen: Input is not a table")
    return false
  end

  local c = 0
  for _ in ipairs(t) do
    c = c + 1
  end

  return c
end


-- FetchConfigData fetches data values from the server given a prefix, a table of
--  keys and their intended datatype. Data returned will always be a flat table
--  matching the same keynames as what was passed with one excpection: if the
--  datatype of a set value does not match the type provided, it will not be set.
--  Optionally, if the datatype provided is a function, that function will perform
--  it's own validation/transformation to the set value. In that case, the return
--  of the function is what is applied to the new table.
--
-- Returns:
--  @1    Prefix
--  @2    Table
function FetchConfigData(prefix, wants)
  if type(prefix) ~= "string" then
    return false
  end
  
  if type(wants) ~= "table" then
    return false
  end

  local server = Server()
  local serverkey, valuetype
  local config = {}

  for k, t in pairs(wants) do
    serverkey = "avorioncontrol:"..prefix..":"..k
    servervalue = server:getValue(serverkey)

    if type(t) == "function" then
      serverkey = t(serverkey)
    elseif type(servervalue) == t then
      config[k] = servervalue
    elseif type(t) ~= "string" then
      print("avocontrol-utilities: FetchConfigData: Invalid object type string: "..t)
    end
  end

  return config
end


-- SetConfigData sets data values on the server given a prefix and a table of
--  values that are to be set.
--
-- Returns:
--  @1    Boolean
function SetConfigData(prefix, sets)
  if type(prefix) ~= "string" then
    return false
  end

  if type(sets) ~= "table" then
    return false
  end

  local server = Server()
  local serverkey = ""

  for k,v in pairs(sets) do
    serverkey = "avorioncontrol:"..prefix..":"..k
    server:setValue(serverkey, v)
  end

  return true
end


-- FileSlurp returns all of the text from a file as a string
--
-- Returns:
--  @1    String
function FileSlurp(f)
  local FILE = assert(io.open(f, "r"))
  local d = FILE:read("*all")
  FILE:close()
  FILE = nil
  return d
end