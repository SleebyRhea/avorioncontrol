package.path = package.path .. ";data/scripts/lib/?.lua"

function execute(runner, cmnd, color, user, message)
  if type(runner) ~= "nil" then
    return 1, "This command is only intended for bot use", ""
  end

  include("avocontrol-discord").Say(color or "default", user or "Server",
    message or "Bad message (file a bug report please)")
  return 1, "", ""
end

function getHelp()
end

function getDescription()
  return "(Bot Only) Display a message from a Discord user"
end