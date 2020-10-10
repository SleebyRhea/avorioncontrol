--[[

  DuelZoneSector - data/scripts/sector/background/duelzone.lua
  ------------------------------------------------------------

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

DuelZoneSector.isEternal = false
DuelZoneSector.isEnabled = false

-- stripZone strips a zone from the list of zones
local function stripZone(zones, zone)
  zones = string.gsub(zones, "^"..zone.."$", "")
  zones = string.gsub(zones, "%:"..zone.."$", "")
  zones = string.gsub(zones, "^"..zone.."%:", "")
  zones = string.gsub(zones, "%:"..zone.."%:", "")
  return zones
end

local function pvpon()
  local s = Sector()
  local x, y   = s:getCoordinates()
  local zone   = x.."_"..y
  local zones = (Server():getValue("duelzones") or "") 
  s.pvpDamage = true
  s:setValue("duelzone", true)
  DuelZoneSector.isEnabled = true
  print(zone..": PvP enabled")
  Server():setValue("duelzones", 
    string.gsub(stripZone(zones, zone) .. ":" .. zone, "^:*(.-):*$", "%1"))
end

local function pvpoff()
  local s = Sector()
  local x, y   = s:getCoordinates()
  local zone   = x.."_"..y
  local zones = (Server():getValue("duelzones") or "") 
  s.pvpDamage = false
  s:setValue("duelzone", false)
  DuelZoneSector.isEnabled = false
  print(zone..": PvP disabled")
  Server():setValue("duelzones", string.gsub(stripZone(zones, zone),
    "^:*(.-):*$", "%1"))
end

local function eternalon()
  DuelZoneSector.isEternal = true
end

local function eternaloff()
  DuelZoneSector.isEternal = false
end


-- DuelZoneSector.initialize sets pvp to 0 when this script is first
--  loaded.
--
-- Returns:
--  None
function DuelZoneSector.initialize()
  local s = Sector()
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
    eternal = (DuelZoneSector.isEternal and 1 or 0),
    enabled = (DuelZoneSector.isEnabled and 1 or 0)}
end


-- DuelZoneSector.restore restores the table returned in DuelZoneSector.secure
--  on script reload. In practice, this is really only for making sure
--  that we are properly setting DuelZoneSector.isEternal and the pvpState.
--
-- Returns:
--  None
function DuelZoneSector.restore(data)
  local s    = Sector()
  local x, y = s:getCoordinates()
  local name = x.."_"..y

  if type(data) == "nil" then
    eternaloff()
    pvpoff()
    return
  end

  local eternalType = type(data["eternal"])
  local enabledType = type(data["enabled"])

  if enabledType ~= "number" then
    pvpoff()
  else
    if data.enabled > 0 then
      pvpon()
    else
      pvpoff()
    end
  end


  if eternalType ~= "number" and DuelZoneSector.isEnabled then
    eternaloff()
  else
    if data.eternal > 0 then
      eternalon()
    else
      eternaloff()
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

  if s:callbacksRegistered("onPlayerLeft", "onPlayerLeft") > 0 then
    s:unregisterCallback("onPlayerLeft", "onPlayerLeft")
  end

  if s:callbacksRegistered("onPlayerEntered", "onPlayerEntered") then
    s:unregisterCallback("onPlayerEntered", "onPlayerEntered")
  end

  pvpoff()
end


-- DuelZoneSector.MakeEternal sets the zone to be a have pvp enabled
--  permanently when its activated. When DuelZoneSector.isEternal is true, the 
--  onPlayerEntered and onPlayerLeft callbacks do nothing.
--
-- Returns:
--  None
function DuelZoneSector.MakeEternal()
  eternaloff()
end


-- DuelZoneSector.MakeEphemeral configures the zone to disable pvp when
--  either of the onPlayerEntered or onPlayerLeft callbacks are run.
--
-- Returns:
--  None
function DuelZoneSector.MakeEphemeral()
  eternaloff()
end


-- DuelZoneSector.EnablePVP enables pvp in the sector and broadcasts a
--  message to all players in said sector.
--
-- Returns:
--  None
function DuelZoneSector.EnablePVP(is_eternal, no_broadcast)
  local s   = Sector()
  local msg = "This sector has been marked as a duelzone. PVP damage is enabled!"

  if not no_broadcast then
    s:broadcastChatMessage("", 3, msg)
  end

  if is_eternal then
    eternalon()
  end

  pvpon()
end


-- DuelZoneSector.DisablePVP disables pvp in the sector and broadcasts
--  a message to all players in said sector.
--
-- Returns:
--  None
function DuelZoneSector.DisablePVP(msg, no_broadcast)
  eternaloff()
  pvpoff()
  
  if not no_broadcast then
    Sector():broadcastChatMessage("", 2, (
      msg and msg..", PVP is disabled." or "PVP is disabled."))
  end
end


-- DuelZoneSector.onPlayerLeft disables pvp when a player leaves, and
--  when the zone is ephemeral
--
-- Returns:
--  None
function DuelZoneSector.onPlayerLeft(index)
  local msg = "${p} has left"%_T % {p = Player(index).name}
  if not DuelZoneSector.isEternal then
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
  local s = Sector()

  local msg = "${p} has arrived, and has interrupted the fight"%_T % {
    p = player.name}

  -- Nested so that we don't run a len everytime this callback is run
  if not DuelZoneSector.isEternal then
    if DuelZoneSector.isEnabled then
      if len(#{s:getPlayers()}) > 0 then
        DuelZoneSector.DisablePVP(msg)
        return
      end
    end
  end

  if  DuelZoneSector.isEternal
  and DuelZoneSector.isEnabled then
    msg = "You have entered a designated PVP area."
    player:sendChatMessage("", 2, msg)
    return
  end

  pvpoff()
end