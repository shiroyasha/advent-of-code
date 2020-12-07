# rubocop:disable all

# input = File.read("input0.txt")
input = File.read("input1.txt")

bags = input.split("\n").map { |line| line.split(" ") }.map do |elements|
	name = elements[0] + " " + elements[1]

  elements.shift
  elements.shift
  elements.shift
  elements.shift

  rules = []

  while(elements.size > 0) do
    count = elements.shift
    n1 = elements.shift
    n2 = elements.shift
    elements.shift

    if count == "no"
      count = 0
    end

    rules.push([n1 + " " + n2, count.to_i])
  end

  [name, rules.to_h]
end.to_h

$bags = bags

#
# Part 1
#

$part1_memo = {}

def part1(name, search_name)
  return $part1_memo[name] if $part1_memo[name] != nil

  val = $bags[name].any? do |name, count|
    next if name == "other bags."

    direct = (name == search_name && count > 0)

    direct || part1(name, search_name)
  end

  $part1_memo[name] = val

  $part1_memo[name]
end

puts "Part 1: #{$bags.select { |name, rules| part1(name, "shiny gold") }.count}"

#
# Part 2
#

$part2_memo = {}

def part2(name)
  return $part2_memo[name] if $part2_memo[name] != nil

  vals = $bags[name].map do |name, count|
    next 0 if name == "other bags."

    count + count * part2(name)
  end

  $part2_memo[name] = vals.inject(:+)

  $part2_memo[name]
end

puts "Part 2: #{part2("shiny gold")}"
