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
    [message[0]]
  else
    []
  end
end

def match_arr(rules, rule, message)
  return [""] if rule.size == 0

  res = []

  match(rules, rule[0], message).each do |match1|
    rest = message[match1.size..-1]

    match_arr(rules, rule[1..-1], rest).each do |match2|
      res << match1 + match2
    end
  end

  res
end

def match(rules, rule_index, message)
  rule = rules[rule_index]

  if rule.class == String
    return match_str(rules, rule, message)
  end

  if rule.class == Array
    return rule.map { |r| match_arr(rules, r, message) }.flatten
  end
end

# ------------------------------

assert_eq(match_str([], "a", "a"), ["a"])
assert_eq(match_str([], "a", "ab"), ["a"])
assert_eq(match_str([], "b", "ab"), [])

rules = [
  [[1, 1]],
  "a",
]

assert_eq(match(rules, 0, "bb"), [])
assert_eq(match(rules, 0, "aa"), ["aa"])

rules = [
  [[1, 2, 1]],
  "a",
  "b"
]

assert_eq(match(rules, 0, "aba"), ["aba"])
assert_eq(match(rules, 0, "abba"), [])

rules = [
  [[1, 2, 1], [1, 2, 2, 1]],
  "a",
  "b"
]

assert_eq(match(rules, 0, "aba"), ["aba"])
assert_eq(match(rules, 0, "abba"), ["abba"])

rules = [
  [[1, 2, 1], [3]],
  "a",
  "b",
  [[1, 2, 2, 1]]
]

assert_eq(match(rules, 0, "aba"), ["aba"])
assert_eq(match(rules, 0, "abba"), ["abba"])

rules = [
  [[3]],
  "a",
  "b",
  [[1, 1]],
  [[2, 2]],
]

assert_eq(match(rules, 0, "bb"), [])
assert_eq(match(rules, 0, "aa"), ["aa"])

rules = [
  [[1, 2], [2, 1]],
  "a",
  "b"
]

assert_eq(match(rules, 0, "ab"), ["ab"])
assert_eq(match(rules, 0, "ba"), ["ba"])

rules = [
  [[1]],
  [[2, 3], [3, 2]],
  "a",
  "b"
]

assert_eq(match(rules, 0, "ab"), ["ab"])
assert_eq(match(rules, 0, "ba"), ["ba"])

rules = [
  [[1, 2], [2, 1]],
  "a",
  "b"
]

assert_eq(match(rules, 0, "abbb"), ["ab"])
assert_eq(match(rules, 0, "babb"), ["ba"])

rules = [
  [[2, 1]],
  [[3, 3], [4, 4]],
  [[3, 4], [4, 3]],
  "a",
  "b"
]

assert_eq(match(rules, 0, "babb"), ["babb"])

# ------------------------------

input = File.read(ARGV[0])
rules, messages = input.split("\n\n")

rules = parse_rules(rules)
messages = messages.split("\n")

puts messages.count { |msg| match(rules, 0, msg).include?(msg) }
