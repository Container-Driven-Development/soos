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
  def main(args) do
    # args |> parse_args |> process

    imageNameTagged = "#{Config.parse()[:imageName]}:#{Tokenizer.depToken()}"

    IO.puts(imageNameTagged)

    # {stdOut, returnCode} = System.cmd "docker", ["pull", imageNameTagged]
    {stdOut, returnCode} = System.cmd "docker", ["image", "ls", "-q", "--filter=reference=#{imageNameTagged}"]

    IO.puts "O: #{stdOut}; S: #{returnCode}"

    if stdOut == "" do
      Path.expand('./Dockerfile') |> Path.absname |> File.write(Npm.getDockerFile, [:write])
      {stdOut, returnCode} = System.cmd "docker", ["build", "-t", imageNameTagged, "."]
    end

    {stdOut, returnCode} = System.cmd "docker", ["run", "--rm", "-v", "#{System.cwd}:/build/app/", imageNameTagged]

    IO.puts stdOut
  end

  defp parse_args(args) do
    {options, _, _} = OptionParser.parse(args,
      switches: [name: :string]
    )
    options
  end

end
