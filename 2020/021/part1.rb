# rubocop:disable all

require 'set'

class Food
  attr_reader :ingreedients
  attr_reader :alergens

  def initialize(ingreedients, alergens)
    @ingreedients = ingreedients
    @alergens = alergens
  end
end

module Day21
  def self.parse(path)
    File.read(path).split("\n").map do |line|
      part1, part2 = line.split(" (contains ")

      part1 = part1.split(" ")
      alergens = part2[0..-2].split(", ")

      Food.new(part1, alergens)
    end
  end

  def self.candidates(food)
    map = {}

    alergens = food.map(&:alergens).flatten

    alergens.each do |i|
      map[i] = Set.new(food.map(&:ingreedients).flatten)
    end

    food.each do |f|
      map.each { |k, v| puts "#{k} #{v.inspect}" }
      puts "-----------"

      f.alergens.each do |a|
        map[a] = map[a] & Set.new(f.ingreedients)
      end
    end

    map
  end

  def self.solve1(path)
    food = parse(path)

    no_good_ingre = candidates(food)

    result = 0

    food.each do |f|
      f.ingreedients.each do |i|
        result += 1 unless no_good_ingre.member?(i)
      end
    end

    result
  end

  def self.solve2(path)
    food = parse(path)

    solution = {}
    m = candidates(food)

    loop do
      m.each { |k, v| puts "#{k} #{v.inspect}" }
      p solution
      puts "----"

      single = m.find { |k, v| v.size == 1 }
      break if single == nil

      ingreedient = single[1].to_a[0]
      solution[single[0]] = ingreedient

      m.each do |k, v|
        v.delete(ingreedient)
      end
    end

    solution.sort_by { |k, v| k }.map { |e| e[1] }.join(",")
  end
end

if __FILE__ == $0
  puts Day21.solve1(ARGV[0])
  puts Day21.solve2(ARGV[0])
end
