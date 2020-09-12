package.path = package.path .. ";data/scripts/lib/?.lua"

function execute(runner, cmnd, user, message)
  if type(runner) ~= "nil" then
    return 1, "\\c(f00) Do not run this please", ""
  end

  discord = include("avocontrol-discord")
  discord.Say(user or "Server", message or "Bad message (file a bug report please)")
  return 1, "", ""
end

function getHelp()
end

function getDescription()
  return "(Bot Only) Display a message from a Discord user"
end