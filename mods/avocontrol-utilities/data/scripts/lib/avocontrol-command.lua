--[[

  AvorionControl - data/scripts/lib/avocontrol-command.lua
  -----------------------------

  Library file for the command mods that AvorionControl utilizes. Adds a simple
  interface for creating a command, and parsing it's arguments.

  License: WTFPL
  Info: https://en.wikipedia.org/wiki/WTFPL

]]

do
  local unpack = (type(table.unpack) == "function" and table.unpack or _G.unpack)

  local debug = (FetchConfigData("AvoDebug", {debug = "boolean"}).debug
    or false)

  local Command = {
    name        = "UnsetName",
    usage       = "No help text defined",
    arguments   = {},
    description = "No description defined",
    execute     = function ()
      print(trace.." Attempted to run command without running SetExecute")
    end}
  
  Command.__index = Command

  local function __valid_arg(in)
  end


  -- Command.AddArgument adds an argument definition to the argument definition
  --  list for later processing.
  --
  -- Returns:
  --  @1    Boolean
  function Command.AddArgument(self, kind, short, long, usage, help, func)
    for k, v in pairs({kind=kind, short=short, long=long, usage=usage, help=help}) do
      if type(v) ~= "string" then
        print(self:Trace.."AddArgument: Invalid argument for "..k)
        return false, "Script error: Command argument definition is invalid"
      end
    end

    if type(func) ~= "function" then
      print(self:Trace.."AddArgument: Invalid function passed (not a function)")
        return false, "Script error: Command argument definition is invalid"
    end

    table.insert(self.arguments, {
      exec  = exec,
      data  = {}
      help  = help,
      long  = long,
      usage = usage,
      short = short})
    
    return true
  end


  -- Command.ParseFlags parses the arguments provided to the command using the
  --  functions defined in the Command.argumnts list via Command.AddArgument
  --
  -- Returns:
  --  @1    Boolean (return status)
  --  @2    Error
  function Command.ParseFlags(self, ...)
    local input = {...}

    if #self.arguments < 1 then
      return true
    end
   
    local cur  = 0
    local last = 0

    repeat
      local arg, data = __valid_arg(table.remove(input, 1))

      -- If arg is set, then update the last variable to hold the old argument
      --  index and update cur to hold the new argument index
      if arg then
        last, cur = cur, arg
        goto continue
      end

      if not cur and not data then
        return false, "Invalid argument supplied"
      end

      -- If both the current argument and the previously processed argument
      --  were the same, then we run that arguments execution function
      --  and reset that table
      if last == cur then
        local err = self.arguments[cur].execute(unpack(self.arguments[cur].data))
        self.arguments[cur].data = {}

        if err
          return false, err
        end

        goto continue
      end

      table.insert(self.arguments.data, data)
      ::continue::
    until GetTblLen(data) < 1

    return true
  end


  -- Command.SetExecute sets the primary execution context for the command being
  --  created
  --
  -- Returns:
  --  @1    Boolean
  function Command.SetExecute(self, func)
    if type(func) ~= "function" then
      print(self:Trace().."SetExecute: Bad type (SetExecute expects a function)")
      return false
    end

    self.execute = func
    return true
  end


  -- Command.Execute runs the execution function that we have defined
  --
  -- Returns:
  --  @1    Int (return status)
  --  @2    String Output
  --  @3    String (Avorion uses this for something but its undocumented)
  function Command.Execute(self, user, cmnd, ...)
    if debug then
      print(self:Trace().."Execute: Running self.execute")
    end
    return self.execute(user, cmnd, ...)
  end


  -- Command.GetDescription returns a string containing the commands description
  --
  -- Returns:
  --  @1    String
  function Command.GetDescription(self)
    return self.description
  end


  -- Command.GetHelp generates help text using the set argument definitions and
  --  returns that as a string for output
  --
  -- Return:
  --  @1    String
  function Command.GetHelp(self)
  end


  -- Command.Trace produces tracing information for debug output
  --
  -- Return:
  --  @1    String
  function Command.Trace(self)
    return "avocontrol: command: "..self.name..": "
  end

  return setmetatable({}, Command)
end