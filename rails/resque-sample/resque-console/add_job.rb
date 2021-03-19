require 'resque'
require_relative 'job'

Resque.redis="redis"

10.times {|n|
  Resque.enqueue( Job, n )
}
