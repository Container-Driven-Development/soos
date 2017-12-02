defmodule Config do

  def parse do
    {:ok, ini} = File.read ".soos"
    Ini.decode(ini)
  end

end
