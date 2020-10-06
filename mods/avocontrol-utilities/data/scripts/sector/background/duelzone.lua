--[[

  DuelZoneSector - data/scripts/sector/background/duelzone.lua
  -------------------------------------------------------

  A simple sector mod that forces PvP to be disabled by default, with
  a command that allows for the sector from which the command is run
  to be temporarily designated as a DuelZone.

  This script is the sector script section of this mod.

  Please see the SectorFunctions.html documentation page for further
  information.

  License: BSD-3-Clause
  https://opensource.org/licenses/BSD-3-Clause

]]

package.path = package.path .. ";data/scripts/lib/?.lua"
include ("stringutility")

-- namespace DuelZoneSector
DuelZoneSector = {}
if not onServer() then
  return
end

local isEternal = false


-- stripZone strips a zone from the list of zones
local function stripZone(zones, zone)
  local patterns = {}
  table.insert(patterns, "^"..zone.."$")
  table.insert(patterns, "%:"..zone.."$")
  table.insert(patterns, "^"..zone.."%:")
  table.insert(patterns, "%:"..zone.."%:")

  for _, p in ipairs(patterns) do
    zones = string.gsub(zones, p, "")
  end

  return zones
end


-- DuelZoneSector.initialize sets pvp to 0 when this script is first
--  loaded.
--
-- Returns:
--  None
function DuelZoneSector.initialize()
  local s = Sector()
  DuelZoneSector:DisablePVP()
  s:setValue("duelzone", false)
  s:registerCallback("onPlayerLeft", "onPlayerLeft")
  s:registerCallback("onPlayerEntered", "onPlayerEntered")
end

-- DuelZoneSector.secure stores a table value containing the current
--  state of this sectors namespace (DuelZoneSector)
--
-- Returns:
--  @0    Table
function DuelZoneSector.secure()
  return {
    eternal = (isEternal and 1 or 0),
    enabled = (Sector().pvpDamage and 1 or 0)}
end


-- DuelZoneSector.restore restores the table returned in DuelZoneSector.secure
--  on script reload. In practice, this is really only for making sure
--  that we are properly setting isEternal and the pvpState.
--
-- Returns:
--  None
function DuelZoneSector.restore(data)
  if type(data) == "nil" then
    return DuelZoneSector.MakeEphemeral()
  end

  local eternalType = type(data["eternal"])
  local enabledType = type(data["enabled"])

  if eternalType ~= "number" then
    DuelZoneSector.MakeEphemeral()
  else
    if data.eternal > 0 then
      DuelZoneSector.MakeEternal()
    end
  end

  if enabledType ~= "number" then
    DuelZoneSector.DisablePVP()
  else
    if data.enabled > 0 then
      DuelZoneSector.EnablePVP()
    end
  end
end


-- DuelZoneSector.onRemove ensures that we disable pvp before this
--  script is unloaded. In addition, remove our registered callBacks
--  from the sector.
--
-- Returns:
--  None
function DuelZoneSector.onRemove()
  local s = Sector()

  DuelZoneSector.DisablePVP()

  if s:callbacksRegistered("onPlayerLeft", "onPlayerLeft") > 0 then
    s:unregisterCallback("onPlayerLeft", "onPlayerLeft")
  end

  if s:callbacksRegistered("onPlayerEntered", "onPlayerEntered") then
    s:unregisterCallback("onPlayerEntered", "onPlayerEntered")
  end
end


-- DuelZoneSector.MakeEternal sets the zone to be a have pvp enabled
--  permanently when its activated. When isEternal is true, the 
--  onPlayerEntered and onPlayerLeft callbacks do nothing.
--
-- Returns:
--  None
function DuelZoneSector.MakeEternal()
  isEternal = true
end


-- DuelZoneSector.MakeEphemeral configures the zone to disable pvp when
--  either of the onPlayerEntered or onPlayerLeft callbacks are run.
--
-- Returns:
--  None
function DuelZoneSector.MakeEphemeral()
  isEternal = false
end


-- DuelZoneSector.EnablePVP enables pvp in the sector and broadcasts a
--  message to all players in said sector.
--
-- Returns:
--  None
function DuelZoneSector.EnablePVP(is_eternal)
  local s = Sector()
  local msg = "This sector has been marked as a duelzone. PVP damage is on!"

  if not s.pvpDamage then
    s.pvpDamage = true
    s:broadcastChatMessage("", 3, msg)
    isEternal = (is_eternal or false)
    s:setValue("duelzone", true)

    local x, y   = s:getCoordinates()
    local zone   = x.."_"..y
    local zones  = (Server():getValue("duelzones") or "")

    Server():setValue("duelzones",
      string.gsub(stripZone(zones, zone) .. ":" .. zone,
      "^:*(.-):*$", "%1"))
  end
end


-- DuelZoneSector.DisablePVP disables pvp in the sector and broadcasts
--  a message to all players in said sector.
--
-- Returns:
--  None
function DuelZoneSector.DisablePVP(msg)
  local s = Sector()
  
  if type(msg) ~= "string" then
    msg = "The fight has ended"
  end

  msg = "${m}. PVP is now off."%_T % {m = msg}

  if s.pvpDamage then
    s.pvpDamage = false
    s:broadcastChatMessage("", 3, msg)
    s:addScriptOnce("data/scripts/sector/background/warzonecheck.lua")
    s:setValue("duelzone", false)

    local x, y   = s:getCoordinates()
    local zone   = x.."_"..y
    local zones  = (Server():getValue("duelzones") or "")

    Server():setValue("duelzones", string.gsub(stripZone(zones, zone),
      "^:*(.-):*$", "%1"))
  end
end


-- DuelZoneSector.onPlayerLeft disables pvp when a player leaves, and
--  when the zone is ephemeral
--
-- Returns:
--  None
function DuelZoneSector.onPlayerLeft(index)
  local msg = "${p} has left"
  if not isEternal then
    DuelZoneSector.DisablePVP(msg)
  end
end


-- DuelZoneSector.onPlayerLeft disables pvp when a player enters, and
--  when the zone is not ephemeral. When not ephemeral, it also sends
--  a warning to said player stating that pvp is enabled in the sector.
--
-- Returns:
--  None
function DuelZoneSector.onPlayerEntered(index)
  local player = Player(index)
  local msg = "${p} has arrived, and has interrupted the fight"%_T % {
    p = player.name}

    if not isEternal then
      DuelZoneSector.DisablePVP(msg)
    else
      if Sector().pvpDamage then
        msg = "You have entered a designated PVP area."
        player:sendChatMessage("", 2, msg)
      end
    end
end