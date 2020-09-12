--[[

  AvorionControl - data/scripts/entity/avocontrol-shiptracker.lua
  -----------------------------

  Emit ship jump and deletion events to stdout for players and alliances

  License: WTFPL
  Info: https://en.wikipedia.org/wiki/WTFPL

]]

-- namespace AvorionControlShipTracker
AvorionControlShipTracker = {}
local index = ""

package.path = package.path .. ";data/scripts/lib/?.lua"
include("stringutility")

function AvorionControlShipTracker.initialize()
  local ship = Entity()
  local x, y = Sector():getCoordinates()
  index = Uuid(ship.index).number
  ship:registerCallback("onDestroyed", "onDestroyed")
  print("shipTrackInitEvent: ${oi} ${si} ${x}:${y} ${sn}"%_T % {
    oi=ship.factionIndex, si=index, x=x, y=y, sn=ship.name})
end

function AvorionControlShipTracker.onSectorChanged()
  local ship  = Entity()
  local x, y  = Sector():getCoordinates()
  print("shipJumpEvent: ${oi} ${si} ${x}:${y} ${sn}"%_T % {
    oi=ship.factionIndex, si=index, x=x, y=y, sn=ship.name})
end

function AvorionControlShipTracker.onDestroyed(destroyedId, destroyerId)
  entity = Entity()
  print("shipDestroyedEvent: "..index)
end

function AvorionControlShipTracker.onDelete()
  print("shipDeletedEvent: "..index)
end