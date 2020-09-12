--[[

  AvorionControl - data/scripts/commands/allplayerdata.lua
  -----------------------------

  Returns data on all players and player alliances

  License: WTFPL
  Info: https://en.wikipedia.org/wiki/WTFPL

]]

local command
command             = include("avocontrol-command")
command.name        = "getplayerdata"
command.description = "Returns data on all players and player alliances"
command.usage       = command.name .. ""

command.SetExecute(func(user, cmnd, ...) {
  local alliances = {}
  local output = ""

  for _, player in ipairs({Server().getPlayers}) do
    if player.alliance then
      alliances[player.alliance] = player.alliance
    end

    output = output .. "${pi} ${ps} ${m} "%_T % {
      pi = player.index,
      ps = player.numships,
      m  = player.money}
    
    for _, resource in ipairs(player:getResources()) do
      output = output .. "${mn}:${ma}"%_T % {
        mn = resource.name,
        ma = resource.value}
    end

    return 0, output, ""
  end
})


function getHelp()
  return command:GetHelp()
end

function getDescription()
  return command:GetDescription()
end

function execute(...)
  return command:Execute(...)
end