# frozen_string_literal: true

=begin
Total pages: 10_094
Total items: 5_509_765
[5509765, "https://api.nuget.org/v3/catalog0/page10094.json"]

Template:
====================================
| count     | time (s) | guess (s) |
====================================
| 1_000     |          |
| 10_000    |          |
| 100_000   |          |
| 1_000_000 |          |
====================================
| 5_509_765 |          | N/A       |
====================================

Single Thread
====================================
| count     | time (s) | guess (s) |
====================================
| 1_000     | 1.527    | 8413
| 10_000    | 1.826    | 1006
| 100_000   | 2.952    | 162
| 1_000_000 | 19.314   | 106
====================================
| 5_509_765 | 91       | N/A       |
====================================
=end


puts "Installing gems..."
require 'bundler/inline'

gemfile do
  source 'https://rubygems.org'

  gem 'oj'
  gem 'typhoeus'
  gem 'spandx', git: 'https://github.com/mokhan/spandx.git', branch: 'concurrency'
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

class Nuget
  attr_reader :command, :queue

  def initialize
    @command = Command.new
    @queue = Queue.new
  end

  def run(cache)
    GC.start
    GC.disable
    threads = 8.times.map do |n|
      Thread.new do
        count = 0
        loop do
          item = queue.deq
          break if item == :stop

          count +=1
          puts [Thread.current.object_id, count].inspect
          cache.insert(item['id'] || item['nuget:id'], item['version'] || item['nuget:version'], [item['licenseExpression']])
        end
      end
    end
    command.run('https://api.nuget.org/v3/catalog0/index.json') do |response|
      json = Oj.load(response.body)
      0.upto(json['count']) do |n|
        fetch_page("https://api.nuget.org/v3/catalog0/page#{n}.json")
      end
    end
    GC.enable
    GC.start
    threads.count.times { queue.enq(:stop) }
    threads.each(&:join)
    cache.rebuild_index
  end

  def fetch_page(url)
    command.run(url) do |response|
      json = Oj.load(response.body)
      json['items'].each do |item|
        queue.enq(item)
      end
    rescue => error
      puts error.inspect
    end
  end
end

cache = Spandx::Core::Cache.new('nuget', root: File.expand_path('.index'))
Nuget.new.run(cache)
