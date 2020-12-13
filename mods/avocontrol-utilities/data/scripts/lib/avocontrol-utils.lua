--[[

  AvorionControl - data/scripts/lib/avocontrol.utils.lua
  ------------------------------------------------------

  Library file for the script mods that are provided by the AvorionControl
  bot. Includes a number of helper functions including (most importantly) the 
  function used to grab and verify data from the Server() object.

  License: BSD-3-Clause
  https://opensource.org/licenses/BSD-3-Clause

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

do
  local toboolean = {["true"] = true, ["false"] = false}
  function toboolean.__toindex (data)
    print("avocontrol-utils: TransformServerData: invalid string: "
      .. data)
    return false
  end

  -- TransformServerData converts an object to the given datatype if it
  --  can. This is really only useful for cases where FetchConfigData
  --  is used.
  -- 
  -- Returns:
  --  @1    Var or Nil
  function TransformServerData(dataType, data)
    if dataType == "boolean" then
      if type(servervalue) == "number" then
        return (servervalue > 0 and true or false)
      end

      if dataType == "string" then
        return toboolean[string.lower(data)]
      end
    end

    if dataType == "number" and type(data) == "string" then
      return tonumber(servervalue)
    end

    return data
  end
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
    return false, "Invalid prefix type: "..type(prefix)
  end
  
  if type(wants) ~= "table" then
    return falsem "Invalid wants type: "..type(wants)
  end

  local server = Server()
  local config = {}

  for k, t in pairs(wants) do
    local serverkey = "avorioncontrol:"..prefix..":"..k
    local servervalue = server:getValue(serverkey)
    if type(t) == "function" then
      servervalue = t(servervalue)
    elseif type(servervalue) == t then
      config[k] = TransformServerData(t, servervalue)
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


-- FindPlayerByName returns a reference to a player given their username
--
-- Returns:
--  @1    Player
function FindPlayerByName(request_name, index)
  for _, p in ipairs({Server():getPlayers()}) do
    if p.baseName == request_name then
      return p
    end
  end
  return nil
end

local materials = {
  iron     = MaterialType.Iron,
  titanium = MaterialType.Titanium,
  naonite  = MaterialType.Naonite,
  trinium  = MaterialType.Trinium,
  xanion   = MaterialType.Xanion,
  ogonite  = MaterialType.Ogonite,
  avorion  = MaterialType.Avorion}

local rarities = {
  petty       = RarityType.Petty,
  common      = RarityType.Common,
  uncommon    = RarityType.Uncommon,
  rare        = RarityType.Rare,
  exceptional = RarityType.Exceptional,
  exotic      = RarityType.Exotic,
  legendary   = RarityType.Legendary}

-- IsValidMaterialString returns either the enum for the given 
--  material string (ie: iron) or false. This is for easy Material
--  checking since there isn't actually a way to do this natively
--  in the game (that I could find) outside of indexing MaterialType
--  directly, which throws an error if it fails. This also allows you
--  to use lowercase material strings as well.
--
-- Returns:
--  @1    Boolean or Number
function IsValidMaterialString(s)
  return materials[string.lower(tostring(s))]
end

-- IsValidRarityString returns either the enum for the given 
--  rarity string (ie: common) or false. This is for easy Rarity
--  checking since there isn't actually a way to do this natively
--  in the game (that I could find) outside of indexing RarityType
--  directly, which throws an error if it fails. This also allows you
--  to index lowercase rarity strings as well.
--
-- Returns:
--  @1    Boolean or Number
function IsValidRarityString(s)
  return rarities[string.lower(tostring(s))]
end
