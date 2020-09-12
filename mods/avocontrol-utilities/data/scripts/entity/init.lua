if onServer() then
  local entity = Entity()
  if not entity.aiOwned then 
    if entity.isShip or entity.isStation or entity.isDrone then
        -- entity:removeScript("data/scripts/entity/avocontrol-shiptracker.lua")      
        entity:addScriptOnce("data/scripts/entity/avocontrol-shiptracker.lua")
    end
  end
end