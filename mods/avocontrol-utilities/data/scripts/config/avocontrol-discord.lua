-- Bot controlled configuration file

Discord = {}

do
  local __discordUrl = "%INVLINK%"
  local __discordBot = "%BOTNAME%"

  -- Discord.Url() returns the configured Discord URL
  function Discord.Url()
    return __discordUrl
  end

  -- Discord.Bot() returns the current name of the bot managing this config file
  function Discord.Bot()
    return __discordBot
  end
end

-- Discord.IsLinked() checks the player index for a linked Discord account and
--  returns string if its valid.
function Discord.IsLinked(index)
  local l = Player(index):getValue("discorduserid")
  return (tonumber(l) and l or "")
end