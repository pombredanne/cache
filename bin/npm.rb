# frozen_string_literal: true
require 'bundler/inline'
require 'etc'

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

class NPM
  attr_reader :command, :queue

  def initialize
    @command = Command.new
    @queue = Queue.new
  end

  def run(cache)
    start = Time.now.to_i
    threads = Etc.nprocessors.times.map do |n|
      Thread.new do
        loop do
          item = queue.deq
          break if item == :stop

          cache.insert(item['name'], item['version'], [item['license']].compact)
        end
      end
    end

    command.run("https://replicate.npmjs.com/registry/_all_docs") do |response|
      json = Oj.load(response.body)
      json['rows'].each do |object|
        _id = object['id']
        key = object['key']
        command.run("https://replicate.npmjs.com/#{key}/") do |response|
          json = Oj.load(response.body)
          json['versions'].each do |version, data|
            queue.enq({
              'name' => data['name'],
              'version' => data['version'],
              'license' => data['license']
            })
          end
        end
      end
    end

    now = Time.now.to_i
    puts "Downloaded catalogue in #{now - start} seconds."

    previous = queue.size
    until queue.empty?
      tmp = queue.size
      puts "Drained: #{previous - tmp}/second. Remaining: #{tmp}"
      previous = tmp
      sleep 1
    end
    puts "Completed writes to disk in #{Time.now.to_i - now} seconds."
    threads.count.times { queue.enq(:stop) }
    threads.each(&:join)
    cache.rebuild_index
  end

  private

  def fetch_dependency(name, cache)
    command.run("https://replicate.npmjs.com/#{name}/") do |response|
      json = Oj.load(response.body)
      json['versions'].each do |version, data|
        puts [data['name'], data['version'], data['license']].inspect
      end
    end
  end
end

cache = Spandx::Core::Cache.new('npm', root: File.expand_path('.index'))
NPM.new.run(cache)
