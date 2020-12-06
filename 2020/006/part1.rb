require 'set'

groups = File.read('input1.txt').split("\n\n")

# convert each person's answer list to a set
groups = groups.map { |g| g.split("\n").map { |p| Set.new(p.split('')) } }

# union on sets
puts "Part 1: #{groups.map { |g| g.inject(:+).size }.inject(:+)}"

# intersection on sets
puts "Part 2: #{groups.map { |g| g.inject(:&).size }.inject(:+)}"
