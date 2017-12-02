defmodule Soos.Mixfile do
  use Mix.Project

  def project do
    [
      app: :soos,
      version: "0.1.0",
      elixir: "~> 1.5",
      escript: [main_module: Soos],
      start_permanent: Mix.env == :prod,
      deps: deps()
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      extra_applications: [:logger]
    ]
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      {:credo, "~> 0.3", only: [:dev, :test]},
      {:ini, git: "https://github.com/nathanjohnson320/ini.git"}
    ]
  end
end
