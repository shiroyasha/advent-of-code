# rubocop:disable all

Rule = Struct.new(:name, :range1, :range2)
Input = Struct.new(:rules, :mine, :others)

def load_input(path)
  input = File.read(path)
  rules, mine, others = input.split("\n\n")

  rules = rules.split("\n").map do |rule|
    name, ranges = rule.split(": ")
    range1, range2 = ranges.split(" or ")

    range1 = range1.split("-").map(&:to_i)
    range1 = (range1[0]..range1[1])

    range2 = range2.split("-").map(&:to_i)
    range2 = (range2[0]..range2[1])

    Rule.new(name, range1, range2)
  end

  mine = mine.split("\n")[1].split(",").map(&:to_i)
  others = others.split("\n")[1..-1].map { |l| l.split(",").map(&:to_i) }

  Input.new(rules, mine, others)
end

input = load_input(ARGV[0])

puts input.rules
puts input.mine.inspect
puts input.others.inspect

def no_rule_matches?(rules, number)
  rules.none? do |rule|
    rule.range1.cover?(number) or rule.range2.cover?(number)
  end
end

invalid_sum = 0

input.others.each do |ticket|
  ticket.each do |number|
    if no_rule_matches?(input.rules, number)
      invalid_sum += number
    end
  end
end

puts invalid_sum
