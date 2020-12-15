# rubocop:disable all

def all_indexes(arr, el)
  arr.each_index.select { |i| arr[i] == el }
end

def find_2020th(input)
  numbers = input.split(",").map(&:to_i)

  i = numbers.length

  while i < 2021
    last = numbers[i-1]
    indexes = all_indexes(numbers, last)

    if indexes.length == 1
      numbers.push(0)
    else
      numbers.push(indexes[-1] - indexes[-2])
    end

    i += 1
  end

  return numbers[2019]
end

puts find_2020th("0,3,6")
puts find_2020th("1,3,2")
puts find_2020th("1,20,8,12,0,14")
