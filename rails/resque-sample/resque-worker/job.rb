class Job
  @queue = :default

  def self.perform(time)
    puts "start #{time}"
    sleep time
    puts "end #{time}"
  end
end
