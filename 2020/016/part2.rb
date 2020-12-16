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

def discard_invalid_tickets(rules, tickets)
  tickets.select do |ticket|
    ticket.all? do |number|
      rules.any? do |rule|
        rule.range1.cover?(number) or rule.range2.cover?(number)
      end
    end
  end
end

def matches_all?(tickets, rule, position)
  tickets.all? do |ticket|
    rule.range1.cover?(ticket[position]) or rule.range2.cover?(ticket[position])
  end
end

def candidates(rules, tickets, except = [])
  rules.map do |rule|
    (0...tickets[0].size).select do |pos|
      !except.include?(pos) && matches_all?(tickets, rule, pos)
    end
  end
end

def order(rules, tickets)
  res = Array.new(rules.size)

  exceptions = []

  rules.size.times do
    candidate_list = candidates(rules, tickets, exceptions)

    # puts candidate_list.map(&:inspect)

    rule_index = candidate_list.find_index { |l| l.size == 1 }
    res_index = candidate_list[rule_index][0]

    res[res_index] = rules[rule_index]
    exceptions.push(res_index)
    # puts res_index

    # puts "----"
    # puts res
    # puts "----"
  end

  res
end

tickets = [input.mine] + input.others
tickets = discard_invalid_tickets(input.rules, tickets)

# puts input.rules
rules = order(input.rules, tickets)

sum = 1

rules.each.with_index do |rule, index|
  if rule.name.start_with?("departure")
    sum *= input.mine[index]
  end
end

puts rules
puts sum
