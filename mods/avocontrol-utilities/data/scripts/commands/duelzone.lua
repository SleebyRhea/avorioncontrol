--[[

  DuelZones - data/scripts/commands/duelzone.lua
  ---------

  A simple sector mod that forces PvP to be disabled by default, with
  a command that allows for the sector from which the command is run
  to be temporarily designated as a duelzone.

  This script is the command section of this mod, and makes use of the
  AvorionControl command library.

  License: BSD-3-Clause
  https://opensource.org/licenses/BSD-3-Clause

]]

package.path = package.path .. ";data/scripts/lib/?.lua"
include("stringutility")

local command       = include("avocontrol-command")
command.name        = "duelzone"
command.description = "Enable or disable PvP for the sector in which its run"
local scriptname    = "data/sector/background/duelzone.lua"

command:AddFlag({
  short = "e",
  long  = "enable",
  usage = "none",
  help  = "Set the current sector to be a dueling zone",
  func  = function() end})

command:AddFlag({
  short = "p",
  long  = "permanent",
  usage = "none",
  help  = "Permanently set a dueling zone (ignore jumpin/jumpout)",
  func  = function() end})

command:AddFlag({
  short = "t",
  long  = "temporary",
  usage = "none",
  help  = "Unset permanency for a dueling zone (disable pvp on jumpin/jumpout)",
  func  = function() end})

command:AddFlag({
  short = "d",
  long  = "disable",
  usage = "none",
  help  = "Disable a dueling zone",
  func  = function() end})

command:SetExecute(function(user, ...)
  if type(user) == "nil" then
    return 1, "Please run this from in-game", ""
  end

  local doEnable    = command:FlagPassed("enable")
  local doDisable   = command:FlagPassed("disable")
  local isEternal   = command:FlagPassed("permanent")
  local isEphemeral = command:FlagPassed("temporary")

  if doEnable and doDisable then
    return 1, "", "--enable (-e) and --disable (-d) are not compatible"
  end

  if doDisable and (isEternal or isEphemeral) then
    return 1, "", "Permanency cannot be set with --disable (-d)"
  end

  if isEternal and isEphemeral then
    return 1, "", "--permanent (-p) and --temporary (-t) are not compatible"
  end

  local s = Sector()

  if not s:hasScript(scriptname) then
    s:addScriptOnce(scriptname)
  end

  if doDisable then
    s:invokeFunction(scriptname, "DisablePVP")
    return 0, "", "Disabled PVP in this sector"
  end

  if doEnable then
    s:invokeFunction(scriptname, "EnablePVP", isEternal)
    return 0, "", "Enabled PVP in this sector"
  end

  if isEternal then
    s:invokeFunction(scriptname, "MakeEternal")
    return 0, "", "Sector is now a permanent PVP zone"
  end

  if isEphemeral then
    s:invokeFunction(scriptname, "MakeEphemeral")
    return 0, "", "Sector is now a temporary PVP zone"
  end

  return 1, "Please supply an argument", ""
end)