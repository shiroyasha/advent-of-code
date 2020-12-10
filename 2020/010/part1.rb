# rubocop:disable all

input = File.read("input2.txt")

numbers = input.split("\n").map(&:to_i)

def find(numbers, adapter)
  return [] if numbers == []

  viable = numbers.select { |n| n - adapter <= 3 }
  number = viable.min
  rest = numbers.reject { |n| n == number }

  [ number ] + find(rest, number)
end

def find_diffs(numbers)
  res = numbers.each_cons(2).map { |a, b| b - a }

  [ res.select { |a| a == 1 }.count, res.select { |b| b == 3 }.count ]
end

def part1(numbers)
  adapters = find(numbers, 0)

  diffs = find_diffs([0] + adapters + [adapters.last+3])

  diffs[0] * diffs[1]
end


$memo = {}

def part2(numbers, index)
  return 1 if index == numbers.length - 1
  return $memo[index] if $memo[index]

  current = numbers[index]

  res = 0

  for i in index+1...numbers.length do
    if numbers[i] - current <= 3
      res += part2(numbers, i)
    end
  end

  $memo[index] = res

  res
end

puts "Result1: #{part1(numbers)}"
puts "Result2: #{part2([0] + numbers.sort, 0)}"
