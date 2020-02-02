# frozen_string_literal: true

RSpec.describe '`spandx-index build nuget` command', type: :cli do
  it 'executes `spandx-index build help nuget` command successfully' do
    output = `spandx-index build help nuget`
    expected_output = <<~OUT
      Usage:
        spandx-index build nuget

      Options:
        -h, [--help], [--no-help]  # Display usage information\n\nCommand description...
    OUT

    expect(output).to eq(expected_output)
  end
end
