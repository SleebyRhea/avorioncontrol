function execute(user, cmnd, ...)
  out = "Nil Command"

  if type(user) ~= "nil" then
    out = Player(user).name .. ": Ran nil command"
  end
  
  if type(...) ~= "nil" then
    out = out .. ": " .. table.concat({...}, " ")
  end

  print("nilCommandEvent: "..out)
  return 0, "", ""
end

function getHelp()
  return "None"
end

function getDescription()
  return "Report nil commands"
end
