defmodule Exec do
  def run(command, attrs) do
    IO.puts command, attrs
    System.cmd command, attrs
  end
end
