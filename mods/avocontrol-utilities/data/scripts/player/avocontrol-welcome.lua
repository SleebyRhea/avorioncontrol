--[[

  AvorionControl - data/scripts/galaxy/server.lua
  -----------------------------------------------

  Sends off an email to a newly joined player using either the defaults
  provided below, or via a WelcomeEmail.txt file located in the Server
  root directory. Also optionally (and by default) adds one or more
  turrets of your spefication to the email as an attachment.

  License: BSD-3-Clause
  https://opensource.org/licenses/BSD-3-Clause

]]


--namespace WelcomeEmail
WelcomeEmail = {}

package.path = package.path .. ";data/scripts/lib/?.lua"
local mailfile = Server().folder .. "/WelcomeEmail.txt"
local player   = Player()

include("stringutility")
include("weapontype")
include("avocontrol-utils")

local function isValidTurret(string)
end

-- Make sure that the material we've configured is valid, and return a reference
--  to an object of it's type if it is. Otherwise, return a new reference to
--  a MaterialType.Iron object.
--
-- Return:
--  @1    Material
local function check_material(string)
  string = TransformServerData("string", string)
  if IsValidMaterialString(string) then
    return Material(MaterialType[string])
  end  
  return nil
end


-- Make sure that the type of turret that we've configured is valid. If it is,
--  return the enum for that WeaponType. Otherwise, return nil
--
-- Return:
--  @1    WeaponType
local function check_type(string)
  string = TransformServerData("string", string)
  if WeaponType[string] == "nil" then
    print("avocontrol (server.lua): Invalid turret type: " .. string)
    return WeaponType.MiningLaser
  end
  return WeaponType[string]
end


-- Make sure that the configured rarity is a valid rarity and return a 
--  new reference to a Rarity object of that type. Otherwise, return a
--  nil
--
-- Return:
--  @1    Rarity
local function check_rarity(number)
  if type(number) ~= "number" then
    return Rarity(2)
  end
  return Rarity(number)
end


-- Make sure that the coordinate passed is both a number and falls between
--  -999 and 999. Otherwise, return 999 (default to very weak)
--
-- Return:
--  @1    Number
local function check_coord(number)
  local n = math.abs(TransformServerData(number) or 500)
  return (n < 0 or n > 500) and 500 or number
end


function WelcomeEmail.initialize()
end







-- function onPlayerCreated (index)
--   vanilla_onPlayerCreated(index)
local msgfooter = "Used default text body."
local madeTurret, usedMailFile = false, false
local player = Player()
local server = Server()
local mail   = Mail()

-- Fetch our mail configuration from the server
local maildata = FetchConfigData("welcomeemail", {
  enabled  = "boolean",
  text     = "string",
  sender   = "string",
  header   = "string",
  money    = "number",
  Iron     = "number",
  Titanium = "number",
  Naonite  = "number",
  Trinium  = "number",
  Xanion   = "number",
  Ogonite  = "number",
  Avorion  = "number"})

if type(maildata.enabled) == "nil" then
  print("Skipping welcome email (disabled)")
  return
end

if maildata.enabled == false then
  print("Skipping welcome email (disabled)")
  return
end

-- Fetch our turret from the server. Anonymous functions used here so that 
--  we dont need to loop this table again.
local turretdata = FetchConfigData("welcometurret", {
  count    = "number",
  offset   = "number",
  material = __check_material,
  rarity   = __check_rarity,
  type     = __check_type,
  x        = __check_coord,
  y        = __check_coord})

-- Update our maildata with the contents of our email file if it exists
if FileExists(__mailfile) then
  maildata.text = FileSlurp(__mailfile)
  usedMailFile = true
end

-- Apply our configurations to the email. Failed/missing configurations are
-- replaced with the following defaults.
mail.sender = (maildata.sender or "Server")
mail.header = (maildata.header or "Welcome!")
mail.text   = (maildata.text   or "Welcome to our server!")
mail.money  = (maildata.money  or 0)

mail:setResources(
  maildata.Iron     or 0,
  maildata.Titanium or 0,
  maildata.Naonite  or 0,
  maildata.Trinium  or 0,
  maildata.Xanion   or 0,
  maildata.Ogonite  or 0,
  maildata.Avorion  or 0)

-- Generate and add our turrets, or break if that fails
for _, t in ipairs(__make_turrets(turretdata) or {}) do
  if type(t) == "userdata" then
    mail:addTurret(t)
    madeTurret = true
  else
    print("Skipped adding turret to player mail. " ..
      "Datatype is incorrect (${s})"%_T % {s=type(t)})
    break
  end
end

if usedMailFile then
  msgfooter = "Used "..__mailfile.." for mail text."
end

if madeTurret then
  msgfooter = msgfooter.." Attached "..turretdata.count.." turret[s]."
end

player:addMail(mail)
print("Sent welcome email to ${p}. ${f}"%_T % {
  p=player.name, f=msgfooter})
