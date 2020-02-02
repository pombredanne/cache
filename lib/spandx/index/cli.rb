# frozen_string_literal: true

require 'thor'

module Spandx
  module Index
    # Handle the application command line parsing
    # and the dispatch to various command objects
    #
    # @api public
    class CLI < Thor
      # Error raised by this runner
      Error = Class.new(StandardError)

      desc 'version', 'spandx-index version'
      def version
        require_relative 'version'
        puts "v#{Spandx::Index::VERSION}"
      end
      map %w[--version -v] => :version

      require_relative 'commands/build'
      register Spandx::Index::Commands::Build, 'build', 'build [SUBCOMMAND]', 'Command description...'
    end
  end
end
