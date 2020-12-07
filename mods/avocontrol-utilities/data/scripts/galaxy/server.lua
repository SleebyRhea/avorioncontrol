--[[

  AvorionControl - data/scripts/galaxy/server.lua
  -----------------------------------------------

  Add player LogIn/Off event output and bot related scripts

  License: BSD-3-Clause
  https://opensource.org/licenses/BSD-3-Clause

]]

-- Create our own login event output for more reliable tracking
function onPlayerLogIn_AvoControl(playerIndex)
  local p = Player(playerIndex)
  print("playerJoinEvent: ${i} ${n}"%_T % {i=p.index, n=p.name})
end

function onPlayerLogOff_AvoControl(playerIndex)
  local p = Player(playerIndex)
  print("playerLeftEvent: ${i} ${n}"%_T % {i=p.index, n=p.name})
end

Server():registerCallback("onPlayerLogIn", "onPlayerLogIn_AvoControl")
Server():registerCallback("onPlayerLogOff", "onPlayerLogOff_AvoControl")