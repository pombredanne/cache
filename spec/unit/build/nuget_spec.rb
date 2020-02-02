# frozen_string_literal: true

require 'spandx/index/commands/build/nuget'

RSpec.describe Spandx::Index::Commands::Build::Nuget do
  it 'executes `build nuget` command successfully' do
    output = StringIO.new
    options = {}
    command = described_class.new(options)

    command.execute(output: output)

    expect(output.string).to eq("OK\n")
  end
end
