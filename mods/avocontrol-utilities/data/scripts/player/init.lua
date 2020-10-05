--[[

  AvorionControl - data/scripts/player/init.lua
  ---------------------------------------------

  This lua script outputs information on the player and their
  current sector on login

  License: BSD-3-Clause
  https://opensource.org/licenses/BSD-3-Clause

]]

do
  if onServer() then
    local player   = Player()
    local sector   = Sector()
    local alliance = (player.alliance and player.alliance or nil)

    local old_path = package.path
    package.path = package.path .. ";data/scripts/lib/?.lua"
    include("stringutility")
    include("avocontrol-utils")

    -- Get ownership counts for the entities in the given table
    local function getOwnerCount(t)
      local p, a, o = 0, 0, 0
      for _, v in ipairs(t) do
        if v.playerOwned then
          p = p + 1
        elseif v.allianceOwned then
          a = a + 1
        else
          o = o + 1
        end
      end
      return p, a, o
    end

    local prefix = "LoginInfo: "..player.name.." "
    local drone  = player.name .. "'s Drone"
    local x, y   = sector:getCoordinates()
    local plan

    if player.craft then
      if player.craft.type ~= EntityType.Drone and player.craft.name then
        if player.craft.playerOwned then
          plan = player:getShipPlan(player.craft.name)
        elseif player.craft.allianceOwned then
          plan = alliance:getShipPlan(player.craft.name)
        end
      end
    end

    local stations_player, stations_alliance = getOwnerCount({
      sector:getEntitiesByType(EntityType.Station)})

    local ships_player, ships_alliance = getOwnerCount({
      sector:getEntitiesByType(EntityType.Ship)})

    print("${p}Ship: ${n}, ${b} Blocks, ${v}k m3"%_T % {
      p = prefix,
      n = (player.craft and player.craft.name or drone),
      b = (plan and plan.numBlocks or "Error"),
      v = (plan and plan.volume or "Error")})

    print("${p}System: ${x}:${y} ${n1} objects, ${n2} players"%_T % {
      p  = prefix,
      x  = x,
      y  = y,
      n1 = sector.numEntities,
      n2 = sector.numPlayers})

    print("${p}System: ${n1} drones, ${n2} player ships, ${n3} alliance ships"%_T % {
      p  = prefix,
      n1 = #{sector:getEntitiesByType(EntityType.Drone)},
      n2 = ships_player,
      n3 = ships_alliance})

    print("${p}System: ${n1} player stations, ${n2} alliance stations"%_T % {
      p  = prefix,
      n1 = stations_player,
      n2 = stations_alliance})
    
    package.path = old_path
  end
end