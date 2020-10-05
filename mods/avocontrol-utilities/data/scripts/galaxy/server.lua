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

do
  local vanilla_onPlayerLogIn = onPlayerLogIn
  function onPlayerLogIn (index)
    vanilla_onPlayerLogIn(index)
    Player(index):addScriptOnce("data/scripts/player/avorioncontrol-welcome.lua")
  end
end