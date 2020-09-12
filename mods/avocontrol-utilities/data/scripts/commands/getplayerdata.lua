--[[

  AvorionControl - data/scripts/commands/getplayerdata.lua
  -----------------------------

  Returns data on all players and player alliances. This command exists solely
  for easy parsing by the bot, and to remove the need for playerinfo rcon loops

  License: WTFPL
  Info: https://en.wikipedia.org/wiki/WTFPL

]]

package.path = package.path .. ";data/scripts/lib/?.lua"
include("stringutility")

local command       = include("avocontrol-command")
command.name        = "getplayerdata"
command.description = "Returns data on all players and player alliances"

local restypes  = {
  [1] = "iron",    [2] = "titanium", 
  [3] = "naonite", [4] = "trinium",
  [5] = "xanian",  [6] = "ogonite", 
  [7] = "avorion"}

command:SetExecute(function ()
  local alliances = {}
  local output    = ""

  for _, player in ipairs({Server():getPlayers()}) do
    if player.alliance then
      alliances[player.alliance] = player.alliance
    end

    local x, y = player:getSectorCoordinates()

    output = output .. "player: ${pi} ${x}:${y} ${ps} ${pS} credits:${m}"%_T % {
      pi = player.index,
      ps = player.numShips,
      pS = player.numStations,
      m  = player.money,
      x  = x,
      y  = y}

    local i = 1
    for _, v in ipairs({player:getResources()}) do
      output = output .. " ${mn}:${ma}"%_T % {
        mn = restypes[i],
        ma = v}
      i = i + 1
    end

    output = output.." "..player.name.."\n"
  end

  for _, alliance in pairs(alliances) do
    output = output .. "alliance: ${pi} ${ps} ${pS} credits:${m}"%_T % {
      pi = alliance.index,
      ps = alliance.numShips,
      pS = alliance.numStations,
      m  = alliance.money}

      local i = 1
      for _, v in ipairs({alliance:getResources()}) do
        output = output .. " ${mn}:${ma}"%_T % {
          mn = restypes[i],
          ma = v}
        i = i + 1
      end
      output = output .. " "..alliance.name .. "\n"
  end

  return 0, output, ""
end)


function getHelp()
  return command:GetHelp()
end

function getDescription()
  return command:GetDescription()
end

function execute(...)
  return command:Execute(...)
end