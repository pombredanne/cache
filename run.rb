# frozen_string_literal: true

puts "Installing gems..."
require 'bundler/inline'

gemfile do
  source 'https://rubygems.org'

  gem 'oj'
  gem 'typhoeus'
end

puts "Starting..."
Oj.default_options = { mode: :strict }

class Command
  attr_reader :hydra

  def initialize(hydra: Typhoeus::Hydra.hydra)
    @hydra = hydra
    @started = false
  end

  def run(url)
    request = Typhoeus::Request.new(url, followlocation: true)
    request.on_complete do |response|
      yield response, hydra
    end
    hydra.queue(request)
    start!
  end

  private

  def start!
    return if @started

    hydra.run
    @started = true
  end
end

class Nuget
  attr_reader :command

  def initialize
    @command = Command.new
  end

  def run
    command.run('https://api.nuget.org/v3/catalog0/index.json') do |response|
      json = Oj.load(response.body)
      json['items'].each do |item|
        fetch_page(item['@id'])
      end
    end
  end

  def fetch_page(url)
    command.run(url) do |response|
      json = Oj.load(response.body)
      json['items'].each do |item|
        puts item.inspect
      end
    end
  end
end

Nuget.new.run
