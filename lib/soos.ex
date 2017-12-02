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

    {stdOut, returnCode} = System.cmd "docker", ["pull", imageNameTagged]

    IO.puts "O: #{stdOut}; S: #{returnCode}"

    if returnCode == 1 do
      Path.expand('./Dockerfile') |> Path.absname |> File.write("""
FROM kkarczmarczyk/node-yarn:8.0-wheezy

WORKDIR /build/app

ENV PATH=/build/node_modules/.bin:$PATH

ADD package.json /build/

RUN cd /build && \
  ([[ -n ${https_proxy} ]] && yarn config set proxy ${https_proxy} -g || true) && \
  yarn && \
  chmod -R 777 /build

RUN mkdir /.config /.cache && \
  chmod -R 777 /.config /.cache

ENTRYPOINT cd /build/app && \
  rm -rf node_modules && \
  mv /build/node_modules /build/app/ && \
  yarn build
      """, [:write])
      {stdOut, returnCode} = System.cmd "docker", ["build", "-t", imageNameTagged, "."]
    end
  end

  defp parse_args(args) do
    {options, _, _} = OptionParser.parse(args,
      switches: [name: :string]
    )
    options
  end

end
