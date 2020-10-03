--[[

  AvorionControl - data/scripts/commands/sendmail.lua
  ---------------------------------------------------

  Sends a player, a list of players, or all players an email.

  License: BSD-3-Clause
  https://opensource.org/licenses/BSD-3-Clause

]]

package.path = package.path .. ";data/scripts/lib/?.lua"
include("stringutility")
include("avocontrol-utils")

local command       = include("avocontrol-command")
command.name        = "sendmail"
command.description = "Send mail to a player, or set of players"
command.help        = "Use -- to signify the beginning of a message"
  .. " body (or just use -f)"

local restypes       = {}
restypes["iron"]     = true
restypes["credits"]  = true
restypes["titanium"] = true
restypes["naonite"]  = true
restypes["trinium"]  = true
restypes["xanian"]   = true
restypes["ogonite"]  = true
restypes["avorion"]  = true

local maildef = {
  sender    = "Server",
  header    = "",
  ircpt     = {},
  nrcpt     = {},
  text      = "",
  resources = {}}

command:AddFlag({
  usage = "sender",
  short = "s",
  long  = "sender",
  help  = "Set the sender name for the mail",
  func  = function(arg)
    maildef.sender = arg
  end})


command:AddFlag({
  usage = "header",
  short = "h",
  long  = "header",
  help  = "Set the header to use for the mail",
  func  = function(arg)
    maildef.header = arg
  end})


command:AddFlag({
  usage = "name",
  short = "p",
  long  = "player-name",
  help  = "Add a player to the list of recipients based on the name",
  func  = function(...)
    for _, n in ipairs({...}) do
      table.insert(maildef.nrcpt, n)
    end
  end})


command:AddFlag({
  usage = "index1,index2,...",
  short = "i",
  long  = "player-index",    
  help  = "Add a player to the list of recipients based on the index",
  func  = function(...)
    local arg = table.concat({...}, ",")
    for m in string.gmatch(arg, "[^, ]+") do
      table.insert(maildef.ircpt, m)
    end
  end})


command:AddFlag({
  usage = "",
  short = "b",
  long  = "broadcast",
  help  = "Send the email to all players",
  func  = function(arg)
    if arg then
      return "broadcast does not take inputs"
    end
  end})


command:AddFlag({
  usage = "filename",
  short = "f",
  long  = "file",
  help  = "Specify the file to use for the mail",
  func  = function(arg)
    local file = Server().folder.."/messages/"..arg
    if FileExists(file) then
      maildef.text = FileSlurp(file)
      return
    end
    return file.." does not exist"
  end})


command:AddFlag({
  usage = "name count",
  short = "r",
  long  = "resource",
  help  = "Set a resource to add to the mail",
  func  = function(res, cnt, ...)
    if ... then
      return "Too many inputs given: ${a}, ${b}, ${c}"%_T % {
        a=tostring(res), b=tostring(cnt), c=table.concat({...},", ")}
    end

    if not type(res) == "string" then
      return "Resource name ${r} is invalid"%_T % {r=tostring(res)}
    end

    res = string.lower(res)

    if not restypes[res] then
      return "Resource name ${r} is invalid"%_T % {r=tostring(res)}
    end

    if type(cnt) == "string" then
      cnt = tonumber(cnt) or nil
    end

    if not cnt then
      return "Please provide a valid amount for resource "..res
    end

    maildef.resources[res] = cnt
  end})


command:SetExecute(function(user, ...)
  if not command:FlagPassed("file") then
    maildef.text = table.concat({...}, " ")
  end

  if maildef.sender == "" or type(maildef.sender) == "nil" then
    return 1, "", "Please supply a sender (or leave out the argument)"
  end

  if maildef.header == "" or type(maildef.header) == "nil" then
    return 1, "", "Please supply a header"
  end

  if maildef.text == "" or type(maildef.text) == "nil" then
    return 1, "", "Please supply a message body"
  end

  local out   = ""
  local sent  = 0
  local fail  = 0
  local mail  = Mail()
  local res   = maildef.resources

  mail.text    = maildef.text
  mail.sender  = maildef.sender
  mail.header  = maildef.header
  mail.money   = (res.credits or 0)

  mail:setResources(
    (res.iron     or 0),
    (res.titanium or 0),
    (res.naonite  or 0),
    (res.trinium  or 0),
    (res.xanion   or 0),
    (res.ogonite  or 0),
    (res.avorion  or 0))

  -- If we're broadcasting, then we don't need to do anything else here
  --  just send the email to all players.
  if command:FlagPassed("broadcast") then
    maildef.rcpt = {}
    for _, p in ipairs({Server():getPlayers()}) do
      p:addMail(mail)
      sent = sent + 1
    end
    return 0, "", "Sent email to ${n} players."%_T % {n=sent}
  end

  -- If we aren't broadcasting, then we need to process the recipients
  --  and map them to player objects
  for _, p in ipairs(maildef.ircpt) do
    p = Player(p)
    if type(p) ~= "nil" then
      sent = sent + 1
    else
      fail = fail + 1
    end
  end

  for _, p in ipairs(maildef.nrcpt) do
    p = FindPlayerByName(p)
    if type(p) ~= "nil" then
      p:addMail(mail)
      sent = sent + 1
    else
      fail = fail + 1
    end
  end

  if sent > 0 then
    out = "${o}Sent email ${n} players."%_T % {o=out, n=sent}
  end

  if fail > 0 then
    out = "${o} Failed to send ${n} emails."%_T % {o=out, n=fail}
  end

  return 0, out , ""
end)