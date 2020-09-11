--[[

  AvorionControl - data/scripts/commands/sendmail.lua
  -----------------------------

  Sends a player, a list of players, or all players an email.

  License: WTFPL
  Info: https://en.wikipedia.org/wiki/WTFPL

]]

-- We use unpack, so this is present in case Avorion ever moves towards lua 5.3
_G.unpack = (type(table.unpack) == "function" and table.unpack or _G.unpack)

-- GetTblLen returns the actual length of a table
--  @1    Table
function GetTblLen(__tbl)
  local c = 0
  for _, _ in pairs(__tbl) do c=c+1 end
  return c
end


-- FindPlayerByName returns a reference to a player given their username
--
-- Returns:
--  @1    Player
function FindPlayerByName(request_name, index)
  for _, p in ipairs({Server():getPlayers()}) do
    print("Found "..p.baseName)
    if p.baseName == request_name then
      return p
    end
  end
  return nil
end


-- FileSlurp returns all of the text from a file as a string
--
-- Returns:
--  @1    String
function FileSlurp(f)
  local FILE = assert(io.open(f, "r"))
  local d = FILE:read("*all")
  FILE:close()
  FILE = nil
  return d
end


-- FileExists returns true if a file exists (and is a file)
--
-- Returns:
--  @1    Boolean
function FileExists(name)
  local f = io.open(name,"r")
  if f ~= nil then
      io.close(f)
      return true
  end
  return false
end

  
do
  local __oldpath = package.path
  local err = nil
  package.path = package.path .. ";data/scripts/lib/?.lua"

  include("stringutility")

  local __err = {
    no_arg = "Please provide an argument",
    bad_arg = "Please provide a valid argument"
  }

  local __modinfo = {
    i_modname = "ds9utils-commandpack",
    i_commandname = "sendmail",
    i_description = "Sends mail via command (for admin/bot use only)",
    i_message_dir = Server().folder .. "/messages/"
  }

  -- __sendEmail is the takes the data that is gathered as part of the "execute"
  --  function and sends an email to the specified players.
  --
  -- Returns:
  --  @1    int (sucess)
  --  @2    string (output)
  --  @3    string (empty, fulfills Avorions command interface)
  local function __sendEmail(__m)
    local __mail   = Mail()
    local __sent   = 0
    local __failed = 0

    __mail.text   = __m.m_text
    __mail.sender = __m.m_sender
    __mail.header = __m.m_header
    __mail.money  = __m.r_credits
    
    __mail:setResources(
      (__m.r_iron     or 0),
      (__m.r_titanium or 0),
      (__m.r_naonite  or 0),
      (__m.r_trinium  or 0),
      (__m.r_xanion   or 0),
      (__m.r_ogonite  or 0),
      (__m.r_avorion  or 0))

  for _, __n in ipairs(__m.m_rcpt_n) do
    local p = FindPlayerByName(__n)
    if type(p) ~= "nil" then
      p:addMail(__mail)
      __sent = __sent + 1
    else
      __failed = __failed + 1
    end
  end

  for _, __i in ipairs(__m.m_rcpt_i) do
    local p = Player(__i)
    if type(p) ~= "nil" then
      p:addMail(__mail)
      __sent = __sent + 1
    else
      __failed = __failed + 1
    end
  end

  out = ""

  if __sent > 0 then
    out = "Sent ${n} players email."%_T % {n=__sent}
  end

  if __failed > 0 then
    out = out .. " Failed to send ${n} emails."%_T % {n=__failed}
  end

  return 0, out , ""
end


-- __addAllPlayersData
--  Append the index for all players to the playerindex DB
local __addAllPlayersData = {
  description = "Add all players the recipients list",
  usage = false,
  func = function (__data)
    __data.m_rcpt_i = {}
    __data.m_rcpt_n = {}
    
    -- TODO: Fix this cludge. This function should just set a boolean value which
    --  is checked later.
    for i, n in ipairs({Galaxy():getPlayerNames()}) do
      table.insert(__data.m_rcpt_n, n)
    end
  end,
}


-- __addPlayerName adds a player name to the list of intended recipients.
local __addPlayerName = {
  description = "Adds a player to the recipients list",
  usage = "playername",
  func = function (__data, ...)
    for _, __id in ipairs({...}) do
      table.insert(__data.m_rcpt_n, __id)
    end
  end,
}

-- __addPlayerIndex adds a player index to the list of intended recipients
local __addPlayerIndex = {
  description = "Adds a player to the recipients list using their index",
  usage = "playerindex",
  func = function (__data, ...)
    for _, __id in ipairs({...}) do
      table.insert(__data.m_rcpt_i, __id)
    end
  end,
}

-- __addResourceData parses and adds the resources provided to the list of 
--  resources to add.
--
-- TODO: This needs to be modified to not parse the resource directly into a
--  table key. It should utilize a lookup table instead.
local __addResourceData = {
  description = "Add resources to an email",
  usage = "resourcename amount",
  func = function (__data, __res, __cnt, ...)
    if ... then
      return "Too many inputs given"
    end

    if not type(__res) == "string" or not __data["r_"..__res] then
      return "Resource name \"${res}\"is invalid"%_T % {res=tostring(__res)}
    end

    if type(__cnt) == "string" then
      __cnt = ( tonumber(__cnt) or nil )
    end

    if not __cnt then
      return "Please provide a valid amount for resource \"${res}\""%_T % {res=__res}
    end

    print("Adding ${c} ${r}"%_T % {
      c=__cnt, r=__res})

    __data["r_"..__res] = __cnt
  end,
}

-- __setHeaderData sets the emails header line
local __setHeaderData = {
  description = "Set the subject line of the message",
  usage = "string",
  func = function(__d, header)
    __d.m_header = header
  end,
}

-- Determines whether or not the file that was provided exists and read it into
--  a string if it does.
local __setFileData = {
  description = "Specify the file (stored in "..
  __modinfo.i_message_dir..
  ") to use for the email",
  usage = "filename",
  func = function (__d, fileName)
    local file = __modinfo.i_message_dir .. fileName
    if FileExists(file) then
      __d.m_text = FileSlurp(file)
      return
    end
    return "${f} does not exist"%_T % {f=file}
  end,
}

-- __setSender sets the sending "address" of the message
local __setSender = {
  description = "Specify the sender",
  usage = "sender",
  func = function (__data, sender)
    __data.m_sender = sender
  end,
}

local function __getAvailableArguments()
  local __data = {}
  __data["-p"] = __addPlayerName
  __data["-i"] = __addPlayerIndex
  __data["-f"] = __setFileData
  __data["-h"] = __setHeaderData
  __data["-b"] = __addAllPlayersData
  __data["-r"] = __addResourceData
  __data["-s"] = __setSender
  return __data
end

function getDescription(sender)
  return __modinfo.i_description
end

function getHelp(sender)
  local __usage = ""
  local __valid_arguments, err = __getAvailableArguments(sender)

  for k,v in pairs(__valid_arguments) do
    __usage = __usage .. "  " .. k ..
    ": <" .. (v.usage and v.usage or "none") .. ">\n        " ..
    v.description .. "\n"
  end

  local out = "Usage: ${c} <option> [parameter] ... -- Message Text "%_T % {
    c = __modinfo.i_commandname}

  out = out .. "\n  Message text starts at \"--\", with all arguments afterwards being ignored\n" ..
    "\nOptions:\n${u}"%_T % {u = __usage}

  return out
end

function execute(sender, commandName, ...)
  if not ... then
    return 1, "", __err.no_arg
  end

  -- If an error is received, cancel the command and return it
  local __valid_arguments, err = __getAvailableArguments()
  if err then
    return 1, "", err
  end

  -- Initilize our message data now that we know the user is authorized
  local __arg_map = {}

  -- Message data
  local __msg = {
    m_sender    = "Server",
    m_header    = "",
    m_text      = "",
    m_rcpt_n    = {}, --RCPT by name
    m_rcpt_i    = {}, --RCPT by index (faster)

    -- Resources
    r_credits   = 0,
    r_iron      = 0,
    r_titanium  = 0,
    r_naonite   = 0,
    r_trinium   = 0,
    r_xanion    = 0,
    r_ogonite   = 0,
    r_avorion   = 0,

    -- Execution Flags
    f_do_resources_send = false,
  }

  -- Locate valid arguments from the ones provided to the command
  do
    local __in_cmd = false
    local __command_data = {...}

    repeat
      local v = table.remove(__command_data, 1)
      if v == "--" then
        __msg.m_text = table.concat(__command_data, " ")
        break
      end

      if type(__valid_arguments[v]) == "nil" and not __in_cmd then
        return 1, "", __err.bad_arg
      elseif type(__valid_arguments[v]) == "table" then
        -- If the last command process was the same, we don't want to overwrite
        -- that entry. So, we perform its run now instead.
        if v == __in_cmd then
          err, _ = __valid_arguments[v].func(__msg, unpack(__arg_map[v]))
          __arg_map[v] = nil
          if err then
            return 1, "", err
          end
        end

        __in_cmd = v
        __arg_map[__in_cmd] = {}

      else
        table.insert(__arg_map[__in_cmd], v)
      end
    until GetTblLen(__command_data) < 1
  end

  -- Process the argument data
  for k, v in pairs(__arg_map) do
    err, _ = __valid_arguments[k].func(__msg, unpack(v))
    if err then
      return 1, "", err
    end
  end

  if GetTblLen(__msg.m_rcpt_n) + GetTblLen(__msg.m_rcpt_i) < 1 then
    return 1, "", "Please supply a recipient"
  end

  if __msg.m_sender == "" or type(__msg.m_sender) == "nil" then
    return 1, "", "Please supply a sender (or leave out the argument)"
  end

  if __msg.m_header == "" or type(__msg.m_header) == "nil" then
    return 1, "", "Please supply a header"
  end

  if __msg.m_text == "" or type(__msg.m_text) == "nil" then
    return 1, "", "Please supply a message body"
  end

  return __sendEmail(__msg)
end

package.path = __oldpath
end