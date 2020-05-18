# frozen_string_literal: true
require 'bundler/inline'

gemfile do
  source 'https://rubygems.org'

  gem 'oj'
  gem 'typhoeus'
  gem 'spandx'
end

trap("SIGINT") { exit(1) }

puts "Starting..."
Oj.default_options = { mode: :strict }

class Command
  attr_reader :hydra

  def initialize(hydra: Typhoeus::Hydra.hydra)
    @hydra = hydra
    @started = false
  end

  def run(url)
    request = Typhoeus::Request.new(url, followlocation: true, accept_encoding: 'gzip')
    request.on_complete do |response|
      yield response
    end
    hydra.queue(request)
    start!
  end

  private

  def start!
    return if @started

    @started = true
    hydra.run
  end
end

class Yarn
  attr_reader :command, :queue

  def initialize
    @command = Command.new
    @queue = Queue.new
  end

  def run(cache)
    ('a'..'z').each do |char|
      ('a'..'z').each do |other_char|
        command.run("https://registry.yarnpkg.com/-/v1/search?size=250&text=#{char}#{other_char}") do |response|
          json = Oj.load(response.body)
          json['objects'].each do |object|
            name = object['package']['name']
            version = object['package']['version']
            fetch_dependency(name, version, cache)
          end
        end
      end
    end
    cache.rebuild_index
  end

  private

  def fetch_dependency(name, version, cache)
    command.run("https://registry.yarnpkg.com/#{name}/#{version}") do |response|
      json = Oj.load(response.body)
      if json.nil?
        cache.insert(name, version, [])
      else
        x = json['versions'] ? json.fetch('versions').fetch(version) : json
        license = x['license']
        puts [name, version, license].inspect
        cache.insert(name, version, [license])
      end
    end
  end
end

cache = Spandx::Core::Cache.new('yarn', root: File.expand_path('.index'))
Yarn.new.run(cache)
