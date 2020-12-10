# rubocop:disable all

input = File.read("input2.txt")

numbers = input.split("\n").map(&:to_i)

def find(numbers, adapter)
  return [] if numbers == []

  viable = numbers.select { |n| n - adapter <= 3 }
  p viable

  number = viable.min
  p number

  rest = numbers.reject { |n| n == number }
  p rest

  p "--------------"

  [ number ] + find(rest, number)
end

def find_diffs(numbers)
  res = numbers.each_cons(2).map { |a, b| b - a }

  [ res.select { |a| a == 1 }.count, res.select { |b| b == 3 }.count ]
end

def part1(numbers)
  adapters = find(numbers, 0)
  puts adapters.last + 3

  diffs = find_diffs([0] + adapters + [adapters.last+3])

  diffs[0] * diffs[1]
end

# puts "Result1: #{part1(numbers)}"


$memo = {}

def part2(numbers, index)
  if index == numbers.length - 1
    return 1
  end

  if $memo[index]
    return $memo[index]
  end

  numbers = numbers.sort

  current = numbers[index]

  res = 0

  for i in index+1...numbers.length do
    if numbers[i] - current <= 3
      # puts "Recurse: #{i}"
      res += part2(numbers, i)
    end
  end

  $memo[index] = res

  res
end

p part2([0] + numbers, 0)
