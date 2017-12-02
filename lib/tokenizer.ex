defmodule Tokenizer do

  def depToken() do

    fileName = getDepFileName()

    case File.read(fileName) do
      {:ok, data} -> hashData( data )
      {:error, reason} -> IO.puts("Unable to open '#{fileName}' because of following reason: #{reason}")
    end

  end

  def hashData( data ) do
    :crypto.hash(:sha, data) |> Base.encode16
  end

  def getDepFileName() do
    Npm.getDepFileName()
  end

end
