# frozen_string_literal: true

require 'thor'

module Spandx
  module Index
    module Commands
      class Build < Thor
        namespace :build

        desc 'nuget', 'Command description...'
        method_option :help, aliases: '-h', type: :boolean,
                             desc: 'Display usage information'
        def nuget(*)
          if options[:help]
            invoke :help, ['nuget']
          else
            require_relative 'build/nuget'
            Spandx::Index::Commands::Build::Nuget.new(options).execute
          end
        end
      end
    end
  end
end
