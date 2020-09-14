package.path = package.path .. ";data/scripts/lib/?.lua"
local command = include("avocontrol-command")
command.name        = "testcommand"
command.description = "Does nothing, only for testing functions of the Command library"

command:AddFlag({
  help  = "Testing flag",
  long  = "test",
  short = "t",
  usage = "[-t|--test]",
  func  = function(arg)
    print("Got: "..arg)
  end})

command:SetExecute(function(user, ...)
  print("Ran by: "..user)
  for i, v in ipairs({...}) do
    print(i.. ": "..v)
  end
end)