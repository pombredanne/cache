# frozen_string_literal: true

RSpec.describe '`spandx-index build` command', type: :cli do
  it 'executes `spandx-index help build` command successfully' do
    output = `spandx-index help build`
    expected_output = <<~OUT
      Commands:
        spandx-index build help [COMMAND]  # Describe subcommands or one specific subcommand
        spandx-index build nuget           # Command description...

    OUT

    expect(output).to eq(expected_output)
  end
end
