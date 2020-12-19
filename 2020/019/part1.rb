# rubocop:disable all

def assert_eq(a, b)
  unless a == b
    p "--------------------"
    p a
    p b
    raise "fail"
  end

  p "------------------------------------------------------------------------"
end

def parse_rules(rules)
  res = {}

  rules.split("\n").map do |line|
    number, definition = line.split(": ")
    number = number.to_i

    if definition[0] == '"'
      res[number] = definition[1]
    else
      res[number] = definition.split(" | ").map do |subrules|
        subrules.split(" ").map(&:to_i)
      end
    end
  end

  res
end

def match_str(rules, rule, message)
  if message[0] == rule
    [message[0], message[1..-1]]
  else
    ["", message]
  end
end

def match_arr(rules, rule, message)
  puts "#{rule}, #{message.inspect}"

  matched = ""
  rest = message

  rule.each do |ref|
    m, rest = match(rules, ref, rest)
    matched += m
  end

  [matched, rest]
end

def match(rules, rule_index, message)
  rule = rules[rule_index]

  puts "#{rule}, #{message.inspect}"

  if rule.class == String
    return match_str(rules, rule, message)
  end

  if rule.class == Array
    rule.each do |r|
      matched, rest = match_arr(rules, r, message)

      if rest == ""
        return matched, ""
      end
    end

    return "", message
  end
end

# ------------------------------

assert_eq(match_str([], "a", "a"), ["a", ""])
assert_eq(match_str([], "a", "ab"), ["a", "b"])
assert_eq(match_str([], "b", "ab"), ["", "ab"])

rules = [
  [[1, 1]],
  "a",
]

assert_eq(match(rules, 0, "bb"), ["", "bb"])
assert_eq(match(rules, 0, "aa"), ["aa", ""])

rules = [
  [[1, 2, 1]],
  "a",
  "b"
]

assert_eq(match(rules, 0, "aba"), ["aba", ""])
assert_eq(match(rules, 0, "abba"), ["", "abba"])

rules = [
  [[1, 2, 1], [1, 2, 2, 1]],
  "a",
  "b"
]

assert_eq(match(rules, 0, "aba"), ["aba", ""])
assert_eq(match(rules, 0, "abba"), ["abba", ""])

rules = [
  [[1, 2, 1], [3]],
  "a",
  "b",
  [[1, 2, 2, 1]]
]

assert_eq(match(rules, 0, "aba"), ["aba", ""])
assert_eq(match(rules, 0, "abba"), ["abba", ""])

rules = [
  [[3]],
  "a",
  "b",
  [[1, 1]],
  [[2, 2]],
]

assert_eq(match(rules, 0, "bb"), ["", "bb"])
assert_eq(match(rules, 0, "aa"), ["aa", ""])

rules = [
  [[1, 2], [2, 1]],
  "a",
  "b"
]

assert_eq(match(rules, 0, "ab"), ["ab", ""])
assert_eq(match(rules, 0, "ba"), ["ba", ""])

rules = [
  [[1]],
  [[2, 3], [3, 2]],
  "a",
  "b"
]

assert_eq(match(rules, 0, "ab"), ["ab", ""])
assert_eq(match(rules, 0, "ba"), ["ba", ""])

rules = [
  [[2, 1]],
  [[3, 3], [4, 4]],
  [[3, 4], [4, 3]],
  "a",
  "b"
]

assert_eq(match(rules, 0, "babb"), ["babb", ""])

# ------------------------------

input = File.read(ARGV[0])
rules, messages = input.split("\n\n")

rules = parse_rules(rules)
messages = messages.split("\n")

messages[0...1].each do |message|
  matched, rest = match(rules, 0, message)

  puts ""
  puts "#{message} => #{matched} #{rest == ""}"
  puts "----------------"
end
