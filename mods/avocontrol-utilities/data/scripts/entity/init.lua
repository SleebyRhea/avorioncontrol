do
  if onServer() then
    package.path = package.path .. ";data/scripts/lib/?.lua"
    include("stringutility")

    function onJump(ship, x, y)
      print("Performing jump")
      -- ship    = Entity(ship)
      -- player  = Player(Faction(ship.factionindex).index)
      -- print("${e}: Player ${i} moved ship ${i}:${s} into (${x}:${y})"%_T % {
      --   e="PlayerJumpEvent", i=player.index, s=ship.index, x=x, y=y})
    end

    local entity = Entity()

    if entity.isShip or entity.isStation or entity.isDrone then
      if not entity.aiOwned then
        print("Added jump callback to "..entity.name)
        entity.registerCallback("onJump", "onJump")
      end
    end
  end
end