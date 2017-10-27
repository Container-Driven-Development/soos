defmodule SoosTest do
  use ExUnit.Case
  doctest Soos

  test "greets the world" do
    assert Soos.hello() == :world
  end
end
