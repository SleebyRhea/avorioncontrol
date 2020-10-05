package.path = package.path .. ";data/scripts/lib/?.lua"
include ("stringutility")

-- namespace DuelZone
DuelZone = {}

if onServer() then
  local isEternal = false

  function DuelZone.initialize()
      Sector().pvpDamage = 0
  end

  function DuelZone.secure()
    return {eternal = isEternal and 1 or 0}
  end

  function DuelZone.secure(data)
    if type(data) == "nil" then
      return DuelZone.MakeEphemeral()
    end

    local eternalType = type(data["eternal"])

    if eternalType == "nil" or eternalType == "string" then
      return DuelZone.MakeEphemeral()
    end

    if data.eternal > 0 then
      DuelZone.MakeEternal()
    end
  end

  function DuelZone.onRemove()
    DuelZone.DisablePVP()
  end

  function DuelZone.MakeEternal()
    isEternal = true
  end

  function DuelZone.MakeEphemeral()
    isEternal = false
  end

  function DuelZone.EnablePVP(is_eternal)
    local s = Sector()
    if s.pvpDamage ~= 1 then
      s.pvpDamage = 1
      local msg = "This sector has been marked as a duel zone! PVP is on!"
      s:broadcastChatMessage("", 0, msg)
      s:broadcastChatMessage("", 3, msg)
      if s:hasScript("sector/background/warzonecheck.lua") then
        s:removeScript("sector/background/warzonecheck.lua")
      end
      isEternal = (is_eternal or false)
    end
  end

  function DuelZone.DisablePVP(msg)
    local s = Sector()
    
    if type(msg) == nil then
      msg = "The fight has ended"
    end

    msg = "${m}. PVP is now off."%_T % {m = msg}

    if s.pvpDamage ~= 0 then
      s.pvpDamage = 0
      s:broadcastChatMessage("", 0, msg)
      s:broadcastChatMessage("", 3, msg)
      s:addScriptOnce("data/sector/background/warzonecheck.lua")
    end
  end

  function DuelZone.onPlayerLeft(index)
    local msg = "${p} has left"
    if not isEternal then
      DuelZone.DisablePVP(msg)
    end
  end

  function DuelZone.onPlayerEntered(index)
    local player = Player(index).name
    local msg = "${p} has arrived, and has interrupted the fight"%_T % {
      p = player}

    if not isEternal then
      DuelZone.DisablePVP(msg)
    end
  end
end