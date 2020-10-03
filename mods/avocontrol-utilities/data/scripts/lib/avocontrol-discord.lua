--[[

  AvorionControl - data/scripts/lib/avocontrol-discord.lua
  --------------------------------------------------------

  Discord helper functions for Avorion

  License: BSD-3-Clause
  https://opensource.org/licenses/BSD-3-Clause

]]

include("avocontrol-utils")

do
  local data = FetchConfigData("Discord", {
    discordUrl = "string",
    discordBot = "string"})

  -- Discord.Url() returns the configured Discord URL
  local function __discordUrl()
    return data.discordUrl
  end

  -- Discord.Bot() returns the current name of the bot managing this config file
  local function __discordBot()
    return data.discordBot
  end

  -- Discord.IsLinked() checks the player index for a linked Discord account and
  --  returns string if its valid.
  function __discordIsLinked(index)
    local l = Player(index):getValue("discorduserid")
    return (tonumber(l) and l or "")
  end

  function __discordSay(user, message)
    local server = Server()
    server:broadcastChatMessage("D", ChatMessageType.Normal,
      "\\c(77c)<"..user.."> \\c(ddd)"..message)
  end

  return {
    Url      = __discordUrl,
    Bot      = __discordBot,
    Say      = __discordSay,
    IsLinked = __discordIsLinked}
end

