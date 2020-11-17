--[[

  AvorionControl - data/scripts/commands/linkdiscordacct.lua
  ----------------------------------------------------------

  Generates a pin that is output both on STDOUT with an identifying
  string for bot processing, and sends the player an email with
  instructions on how to link their discord user to their Avorion user.
  If used by RCON, this command can be used to generate an Integration
  request for a given user.

  License: BSD-3-Clause
  https://opensource.org/licenses/BSD-3-Clause

]]

mod = {
  name        = "linkdiscordacct",
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

-- execute is the main function that is run when this command is run
function execute(user, cmnd, ...)
  if type(user) == "nil" then
    args    = {...}
    id      = args[1]
    discord = args[2]
    if type(id) == "nil" or type(Player(id)) == "nil" then
      print("Invalid: "..tostring(id)..tostring(Player(id)))
      return 1, "Please supply a valid user index", ""
    end

    if type(discord) == "nil" then
      return 1, "Please supply a valid discord ID", ""
    end

    Player(id):setValue("discorduserid", discord)
    return 0, "Set user integration", ""
  end

  package.path = package.path .. ";data/scripts/lib/?.lua"
  discord = include("avocontrol-discord")
  include("stringutility")
  include("randomext")

  local player = Player(user)
  local pin = ""
  local du  = ""

  if type(player) == "nil" then
    return 0, "Failed to get playerdata"
  end

  for i=1, 10, 1 do
    pin = pin .. getInt(0,9)
  end

  du = discord.IsLinked(player.index)
  if du ~= "" then
    return 0, "You've already linked your Discord account to ${d}"%_T % {d=du}, ""
  end

  local msg = "Hello, ${p}!\n\n"%_T % {p=player.name}
    .. "Here is your Discord integration code: ${i}:${c}"%_T % {i=player.index, c=pin}
    .. "\n\n"
    .. "To use that code, all you need to do is direct message it to our Discord "
    .. "bot. Make sure to do this *after* you have joined the Discord server ("
    .. "link is provided below). If you have any questions, feel free to reach "
    .. "out to us!\n\n"
    .. "Take care!\n\n"
    .. "Discord Link: ${l}\nDiscord Bot:  ${b}"%_T % {
          l=discord.Url(), b=discord.Bot()}

  mail = Mail()
  mail.text = msg
  mail.sender = "Server"
  mail.header = "Discord Integration Request"
  player:addMail(mail)

  print("discordIntegrationRequestEvent: ${i} ${c}"%_T % {i=player.index, c=pin})
  return 1, "Code sent (check your mail for instructions)", ""
end