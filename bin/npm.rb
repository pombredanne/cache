# frozen_string_literal: true
require 'bundler/inline'
require 'etc'

gemfile do
  source 'https://rubygems.org'

  gem 'async-http'
  gem 'oj'
  gem 'spandx'
end

require 'async'
require 'async/barrier'
require 'async/http/internet'

trap("SIGINT") { exit(1) }

puts "Starting..."
Oj.default_options = { mode: :strict }

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

    if !File.exist?('_all_docs')
      exit 1 unless system('wget https://replicate.npmjs.com/registry/_all_docs')
    end

    Async do
      internet = Async::HTTP::Internet.new
      barrier = Async::Barrier.new
      headers = [['accept', 'application/json']]

      json = Oj.load(IO.read('_all_docs'))
      json.fetch('rows', []).each do |object|
        _id = object['id']
        key = object['key']

        barrier.async do
          begin
            response = internet.get("https://replicate.npmjs.com/#{key}/", headers)
            json = Oj.load(response.read)
            json.fetch('versions', []).each do |version, data|
              queue.enq(name: data['name'], version: data['version'], license: data['license'])
            end
          rescue
            puts "ERROR: https://replicate.npmjs.com/#{key}/"
          end
        end
      end
      barrier.wait
    ensure
      internet&.close
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
end

NPM.run(Spandx::Core::Cache.new('npm', root: File.expand_path('.index')))
