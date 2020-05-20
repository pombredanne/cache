# frozen_string_literal: true
require 'bundler/inline'
require 'etc'

gemfile do
  source 'https://rubygems.org'

  gem 'oj'
  gem 'spandx'
  gem 'typhoeus'
end

trap("SIGINT") { exit(1) }

puts "Starting..."
Oj.default_options = { mode: :strict }

class Command
  attr_reader :hydra

  def initialize(hydra: Typhoeus::Hydra.hydra)
    @hydra = hydra
  end

  def run(url)
    request = Typhoeus::Request.new(url, followlocation: true, accept_encoding: 'gzip')
    request.on_complete do |response|
      yield response
    end
    hydra.queue(request)
  end

  def start!
    hydra.run
  end
end

class NPM
  def self.run(cache)
    queue = Queue.new
    start = Time.now.to_i
    threads = Etc.nprocessors.times.map do |n|
      Thread.new do
        loop do
          item = queue.deq
          break if item == :stop

          cache.insert(item[:name], item[:version], [item[:license]].compact)
        end
      end
    end

    system("rm _all_docs") if File.exist?("_all_docs")
    exit 1 unless system('wget https://replicate.npmjs.com/registry/_all_docs')

    command = Command.new
    json = Oj.load(IO.read('_all_docs'))
    json.fetch('rows', []).each do |object|
      _id = object['id']
      key = object['key'].sub('/', '%2f')
      command.run("https://replicate.npmjs.com/#{key}") do |response|
        if response.success?
          json = Oj.load(response.body)
          json.fetch('versions', []).each do |version, data|
            queue.enq(name: data['name'], version: data['version'], license: data['license'])
          end
        elsif response.code == 0
          puts response.return_message
        else
          puts "ERROR: #{response.code} #{key}"
        end
      rescue
        puts "ERROR: https://replicate.npmjs.com/#{key}/"
      end
    end
    command.start!

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
end

NPM.run(Spandx::Core::Cache.new('npm', root: File.expand_path('.index')))
