# rubocop:disable all

class Memory
  def initialize
    @memory1 = {}
    @memory2 = {}
  end

  def remember(value, index)
    if @memory1[value] && @memory2[value]
      @memory1[value] = @memory2[value]
      @memory2[value] = index
    elsif @memory1[value]
      @memory2[value] = index
    else
      @memory1[value] = index
    end
  end

  def has_diff?(value)
    @memory1[value] && @memory2[value]
  end

  def get_diff(value)
    @memory2[value] - @memory1[value]
  end
end

def find_nth(input, nth)
  mem = Memory.new()

  numbers = input.split(",").map(&:to_i)
  last = numbers[numbers.length-1]

  numbers.each.with_index { |n, i| mem.remember(n, i) }

  i = numbers.length
  while i < nth
    puts i if i % 1000000 == 0

    if mem.has_diff?(last)
      current = mem.get_diff(last)

      mem.remember(current, i)
      last = current
    else
      last = 0

      mem.remember(0, i)
    end

    i += 1
  end

  last
end

# puts find_nth("0,3,6", 30000000)
# puts find_2020th("1,3,2")
puts find_nth("1,20,8,12,0,14", 30000000)
