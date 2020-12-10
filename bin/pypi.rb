# frozen_string_literal: true
require 'bundler/inline'
require 'etc'

gemfile do
  source 'https://rubygems.org'

  gem 'async'
  gem 'spandx', git: 'https://github.com/spandx/spandx.git'
  gem 'typhoeus'
end

trap("SIGINT") { exit(1) }

puts "Starting..."

class Http
  def get(url)
    Async do
      Typhoeus.get(url, followlocation: true, accept_encoding: 'gzip')
    end.wait
  end

  def ok?(response)
    response.success?
  end
end

#http = Spandx.http
http = Http.new
pypi = Spandx::Python::Pypi.new(http: http)
cache = Spandx::Core::Cache.new('pypi', root: File.expand_path('.index'))

start = Time.now.to_i
queue = Queue.new
threads = Etc.nprocessors.times.map do |n|
  Thread.new do
    loop do
      item = queue.deq
      break if item == :stop

      cache.insert(item[:name], item[:version], [item[:license]].compact)
    end
  end
end

pypi.each do |item|
  queue.enq(item)
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
