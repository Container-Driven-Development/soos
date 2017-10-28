defmodule Soos do
  @moduledoc """
  Documentation for Soos.
  """

  @doc """
  Hello world.

  ## Examples

      iex> Soos.hello
      :world

  """
  def hello do
    :world
  end

  def main(args) do
    args |> parse_args |> process
  end

  def process([]) do
    IO.puts "No arguments given"
  end

  def process(options) do

    filename = options[:name]

    case File.read(filename) do
      {:ok, body} -> :crypto.hash(:sha, body) |> Base.encode16 |> IO.puts
      {:error, reason} -> IO.puts("Unable to open '#{filename}' because of following reason: #{reason}")
    end

  end

  defp parse_args(args) do
    {options, _, _} = OptionParser.parse(args,
      switches: [name: :string]
    )
    options
  end

end
