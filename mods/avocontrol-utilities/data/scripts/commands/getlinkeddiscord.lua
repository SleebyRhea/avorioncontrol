package.path = package.path .. ";data/scripts/config/?.lua"
include("avocontrol-discord")

local command = {
  name = "getlinkeddiscord",
  desc = "Returns a given players Discord userid, or nothing if its unlinked"
}

function getDescription()
  return command.desc
end

function getHelp()
  return "Usage: "..command.name
end

function execute(_, cmdn, index)
  return 0, Discord.IsLinked(index), ""
end