--[[

  AvorionControl - data/scripts/lib/avocontrol-command.lua
  -----------------------------

  Library file for the command mods that AvorionControl utilizes. Adds a simple
  interface for creating a command, and parsing it's arguments.

  License: WTFPL
  Info: https://en.wikipedia.org/wiki/WTFPL

]]

do
  package.path = package.path .. ";data/scripts/lib/?.lua"
  include("avocontrol-utils")

  local unpack = (type(table.unpack) == "function" and table.unpack or _G.unpack)

  local debug = (FetchConfigData("AvoDebug", {debug = "boolean"}).debug
    or false)

  local Command = {
    name        = "UnsetName",
    flags       = {},
    description = "No description defined",
    execute     = function ()
      print(trace.." Attempted to run command without running SetExecute")
    end}
  
  Command.__index = Command

  local function validFlag(self, arg)
    if type(arg) ~= "string" then
      return nil, nil
    end

    -- Check the argument against all assigned flags and if one fits, return the
    --  flag index for referencing but not the argument
    for index, flag in ipairs(self.flags) do
      if "-"..flag.short == arg or "--"..flag.long == arg then
        return index, nil
      end
    end

    -- If no matches were found, return the argument but not the index
    return nil, arg
  end

  -- Command.AddFlag adds an argument definition to the argument definition
  --  list for later processing.
  --
  -- Returns:
  --  @1    Boolean
  function Command.AddFlag(self, d)
    for k, v in pairs({short=d.short, long=d.long, usage=d.usage, help=d.help}) do
      if type(v) ~= "string" then
        print(self:Trace().."AddArgument: Invalid argument for "..k)
        return false, "Script error: Command argument definition is invalid"
      end
    end

    if type(d.func) ~= "function" then
      print(self:Trace().."AddArgument: Invalid function passed (not a function)")
        return false, "Script error: Command argument definition is invalid"
    end

    table.insert(self.flags, {
      execute = d.func,
      help    = d.help,
      long    = d.long,
      usage   = d.usage,
      short   = d.short})
    
    print("Added flag: "..d.long)
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

    if #self.flags < 1 then
      return true
    end
   
    local cur     = false
    local extra   = {}
    local handled = {}

    for _, v in ipairs(input) do
      local flag, arg = validFlag(self, v)

      -- Debug output. TODO: Remove this when this is done
      print((type(flag)~="nil" and flag or "nil") .. ":"
        .. (type(arg) ~= "nil" and arg or "nil"))

      -- If flag is set, and its data is present, then it's been handled before
      --  and we should specify this.
      if flag then
        if type(self.flags[flag].data) == "table" then
          handled[flag] = true
        else
          handled[flag] = false
        end

        cur = flag
        goto continue
      end

      -- Assign any arguments that do not have a given flag to the extra table.
      --  These will be unpacked into the command.execute function
      if not cur and arg then
        print("Adding argument to extra: "..arg)
        table.insert(extra, arg)
        goto continue
      end

      -- Catch bad inputs. Execution ends here.
      if not cur and not arg then
        return false, "Invalid argument supplied"
      end

      -- If the current argument has already input, process its data and set its
      --  handled value to false and reset the flag data table
      if handled[cur] then
        print("Running flag: "..self.flags[cur].long)
        local err = self.flags[cur].execute(unpack(self.flags[cur].data))
        self.flags[cur].data = nil
        if err then
          return false, err
        end
      end

      -- Add our argument data and set the to false to complete the input
      print("Adding \""..arg.."\" to flag: "..self.flags[cur].long)
      self.flags[cur].data = {}
      table.insert(self.flags[cur].data, arg)
      cur = false

      ::continue::
    end

    for i, _ in ipairs(self.flags) do
      if type(self.flags[i].data) == "table" then
        self.flags[i].execute(unpack(self.flags[i].data))
      end
    end

    self.data.extra = extra
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
    local ok, err = self:ParseFlags(...)
    
    if not ok then
      return 1, err, ""
    end
    
    return self.execute(user, unpack(self.data.extra))
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
    if #self.flags > 0 then
      -- local output = "Example: /"..self.name.." [--parameter|-p] argument\n"
      -- for i, v in self.flags do
      --   output = 
      -- end
    else
      return "Example: /"..self.name
    end
  end


  -- Command.Trace produces tracing information for debug output
  --
  -- Return:
  --  @1    String
  function Command.Trace(self)
    return "avocontrol: command: "..self.name..": "
  end

  local command = setmetatable({data = {}}, Command)

  -- Set the global functions that Avorion looks for. Doing this here means that
  --  simply sourcing our library creates a usable command.
  function _G.getHelp()
    return command:GetHelp()
  end
  
  function _G.getDescription()
    return command:GetDescription()
  end
  
  function _G.execute(...)
    return command:Execute(...)
  end

  return command
end