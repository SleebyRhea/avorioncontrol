--[[

  AvorionControl - data/scripts/commands/getlinkeddiscord.lua
  -----------------------------

  Given a player index, this command returns the players Discord UID if it's
  been linked. Otherwise, an empty string is returned. This command is for use
  by Admins and (more importantly) RCON.

  License: WTFPL
  Info: https://en.wikipedia.org/wiki/WTFPL

]]

package.path = package.path .. ";data/scripts/lib/?.lua"

-- getDescription returns this commands description. For use with /help
function getDescription()
  return "Returns a given players Discord userid, or nothing if its unlinked"
end

-- getHelp returns this commands help syntax. For use with /help
function getHelp()
  return "Usage: getlinkeddiscord"
end

-- execute is the main function that is run when this command is run
function execute(_, _, index)
  return 0, include("avocontrol-discord").IsLinked(index), ""
end