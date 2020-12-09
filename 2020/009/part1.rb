# rubocop:disable all

input = File.read("input0.txt")
preamble_size = 5

input = File.read("input1.txt")
preamble_size = 25

numbers = input.split("\n").map(&:to_i)

def sums(numbers)
  res = []

  numbers.each.with_index do |n1, i1|
    numbers.each.with_index do |n2, i2|
      next if i1 == i2

      res << n1 + n2
    end
  end

  res
end

def find1(numbers, preamble_size)
  numbers[preamble_size..-1].each.with_index do |num, index|

    all_pair_sums = sums(numbers[index...index+preamble_size])

    return num unless all_pair_sums.include?(num)
  end
end

def sum_of_min_and_max(numbers)
  numbers.min + numbers.max
end

def find2(numbers, target)
  (0...numbers.length).each do |from|
    res = 0

    (from...numbers.length).each do |till|
      res += numbers[till]

      next if res < target
      break if res > target

      return sum_of_min_and_max(numbers[from..till]) if res == target
    end
  end
end

res1 = find1(numbers, preamble_size)
puts "Result Part 1: #{res1}"

res2 = find2(numbers, res1)
puts "Result Part 2: #{res2}"
