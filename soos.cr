require "admiral"

class Soos < Admiral::Command
  define_flag planet

  def run
    puts "Hello #{flags.planet || "World"}"
  end
end

Soos.run
