--[[

  AvorionControl - data/scripts/commands/rules.lua
  ------------------------------------------------

  When invoked, this command will read from a file named "RulesList.txt"
  placed within the galaxy root (the output of Server().folder). The
  syntax said file is as follows:

  ";": A comment. Lines beginning with this character are ignored

  "#": Disable formatting. Lines beginning with this character will
       simply be output as is.

  "-" or "+": These lines are formatted as a bullet point for a rule

  All other non-empty lines present are treated as new rules and are
  output with numbered formatting. Here is an source, and its output:

  Note: Leading and trailing whitespaces are trimmed.

  Example Source:
    # =|Server Rules|=
    # NOTE: Simple note
    ; This is completely ignored
    Rule about something
      - More information on that rule
    Rule about something else
      + Another bit of info, but on this rule instead

  Example Output:
    =|Server Rules|=
    NOTE: Simple note
    1. Rule about something
      - More information on that rule
    2. Rule about something else
      - Another bit of info, but on this rule instead
  
  
  License: BSD-3-Clause
  https://opensource.org/licenses/BSD-3-Clause

]]

package.path = package.path .. ";data/scripts/lib/?.lua"
include("stringutility")
include("avocontrol-utils")

local command       = include("avocontrol-command")
command.name        = "rules"
command.description = "Output the server rules"

command:SetExecute(function(sender)
  local file = Server().folder .. "/RulesList.txt"
  local resp = ""
  
  if not FileExists(file) then
    return 1, "", "No rules defined!"
  end
  
  local c = 0
  for l in io.lines(file) do
    local parsed = false
    
    l = string.gsub(l, "^%s*(.-)%s*$", "%1")
    
    if type(l) == "nil" or l == "" or string.sub(l,1,1) == ";" then
      goto continue
      
    elseif string.sub(l,1,1) == "#" then
      l = string.gsub(l, "^%s*#*%s*(.-)%s*$", "%1")
      parsed = "${r}\n"%_T % {c=c,r=l}
      
    elseif string.sub(l,1,1) == "-" or string.sub(l,1,1) == "+" then
      l = string.gsub(string.sub(l,2,-1), "^%s*[-+]*%s*(.-)%s*$", "%1")
      parsed = "    - ${r}\n"%_T % {c=c, r=l}
      
    else
      c = c + 1
      parsed = "${c}: ${r}\n"%_T % {c=c, r=l}
    end
    
    ::continue::    
    resp = (parsed and resp..parsed or resp)
  end
  
  if type(sender) ~= "nil" then
    print("Player <${p}> has read the server rules"%_T % {
      p=Player(sender).name})
  end

  return 0, resp, ""
end)