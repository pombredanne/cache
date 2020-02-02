# frozen_string_literal: true

require_relative '../../command'

module Spandx
  module Index
    module Commands
      class Build
        class Nuget < Spandx::Index::Command
          def initialize(options)
            @options = options
          end

          def execute(output: $stdout)
            # Command logic goes here ...
            output.puts 'OK'
          end
        end
      end
    end
  end
end
