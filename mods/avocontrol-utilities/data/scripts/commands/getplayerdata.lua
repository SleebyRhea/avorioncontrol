--[[

  AvorionControl - data/scripts/commands/getplayerdata.lua
  -----------------------------

  IMPORTANT: Bug reports with cases where this script was modified will be closed.

  A reimplementation of the playerinfo command. While this does not provide data
  on the players Steam64ID nor their IP address, this command does have the ability
  to output data on every player and alliance that currently exists in the game.

  More importantly than that though, this command exists primarily for use by the
  manager object to have a command with a controllable regex for matching against.

  For users of this bot, do NOT modify this command unless you are 100% sure that
  the output will not differ. If you do so, and this either breaks or it's return
  output differs from the bots regex definitions, player and alliance tracking
  will no longer function and will likely break other functionality.

  License: WTFPL
  Info: https://en.wikipedia.org/wiki/WTFPL

]]

package.path = package.path .. ";data/scripts/lib/?.lua"
include("stringutility")

local command       = include("avocontrol-command")
command.name        = "getplayerdata"
command.description = "Returns data on all players and player alliances"

local restypes = {}
table.insert(restypes, "iron")
table.insert(restypes, "titanium")
table.insert(restypes, "naonite")
table.insert(restypes, "trinium")
table.insert(restypes, "xanian")
table.insert(restypes, "ogonite")
table.insert(restypes, "avorion")

command:AddFlag({
  short = "p",
  long  = "player",
  usage = "[-p|--player] playerindex",
  help  = "Return data on specified player index",
  func  = function(arg)
    if type(command.data.players) == "nil" then
      command.data.players = {}
    end
    table.insert(command.data.players, arg)
  end})

command:AddFlag({
  short = "a",
  long  = "alliance",
  usage = "[-a|--alliance] allianceindex",
  help  = "Return data on specified Alliance index",
  func  = function(arg)
    if type(command.data.alliances) == "nil" then
      command.data.alliances = {}
    end
    table.insert(command.data.alliances, arg)
  end})

command:SetExecute(function ()
  local doEveryPlayer   = true
  local doEveryAlliance = true

  local alliances  = {}
  local playerlist = {}
  local output     = ""

  -- Process our provided players if they were given. We also disable the default
  --  behaviour of processing all data here if this is processed
  if type(command.data.players) ~= "nil" and #command.data.players > 0 then
    for _, index in ipairs(command.data.players) do
      if type(tonumber(index)) ~= "number" then
        return 1, "Index must be a number", ""
      end

      local p = Player(index)
      
      if p == nil then
        return 1, "Failed to acquire data for index: "..index, ""
      end

      table.insert(playerlist, p)
    end
    doEveryPlayer = false
    doEveryAlliance = false
  end

  -- Process our provided alliances if they were given. We also disable the default
  --  behaviour of processing all data here if this is processed
  if type(command.data.alliances) ~= "nil" and #command.data.alliances > 0 then
    for _, index in ipairs(command.data.alliances) do
      if type(tonumber(index)) ~= "number" then
        return 1, "Index must be a number", ""
      end

      local a = Alliance(index)

      if a == nil then
        return 1, "Failed to acquire data for index: "..index, ""
      end

      alliances[index] = a
    end
    doEveryPlayer = false
    doEveryAlliance = false
  end

  -- Default to processing every player
  if doEveryPlayer then
    playerlist = {Server():getPlayers()}
  end

  for _, player in ipairs(playerlist) do
    if doEveryAlliance then
      if player.alliance then
        alliances[player.alliance] = player.alliance
      end
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