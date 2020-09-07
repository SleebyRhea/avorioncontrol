if not onServer() then
  return
end

mod = {
  name        = "linkdiscord",
  description = "Generate a PIN to link your discord account with Avorion"
}

-- getDescription returns this commands description. For use with /help
function getDescription()
  return mod.description
end

-- getHelp returns this commands help syntax. For use with /help
function getHelp(cmnd)
  return "Usage: " .. (cmnd or mod.name)
end

-- execute is the main function that is executed when this command is run
function execute(user, cmnd, ...)
  if type(user) == "nil" then
    return 1, "Cannot generate an integration request for RCON"
  end

  package.path = package.path .. ";data/scripts/lib/?.lua"
  package.path = package.path .. ";data/scripts/config/?.lua"
  include("avocontrol-discord")
  include("stringutility")
  include("randomext")

  local plr = Player(user)
  local pin = ""
  local du  = ""

  for i=1, 5, 1 do
    pin = pin .. getInt(0,9)
  end

  if type(plr) == "nil" then
    return 0, "Failed to get playerdata"
  end

  du = Discord.IsLinked(plr.index)
  if du ~= "" then
    return 0, "You've already linked your Discord account to ${d}"%_T % {d=du}, ""
  end

  local msg = "Hello, ${p}!\n\n"%_T % {p=plr.name}
    .. "Here is your Discord integration code: ${i}:${c}"%_T % {i=plr.index, c=pin}
    .. "\n\n"
    .. "To use that code, all you need to do is direct message it to our Discord "
    .. "bot. Make sure to do this *after* you have joined the Discord server ("
    .. "link is provided below). If you have any questions, feel free to reach "
    .. "out to us!\n\n"
    .. "Take care!\n\n"
    .. "Discord Link: ${l}\nDiscord Bot:  ${b}"%_T % {
      l=Discord.Url(), b=Discord.Bot()}

  m = Mail()
  m.text = msg
  m.sender = "Server"
  m.header = "Discord Integration Request"
  plr:addMail(m)

  print("DiscordIntegrationRequest: ${i} ${c}"%_T % {i=plr.index, c=pin})
  return 1, "Code sent (check your mail for instructions)", ""
end