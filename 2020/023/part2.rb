# rubocop:disable all

require_relative "./circular_list.rb"

module Day23
  def self.play_all(list, current, rounds)
    rounds.times.each do |i|
      if i % 10000 == 0
        puts i
      end

      # list.nice_show(current.data)

      a = list.remove_next(current)
      b = list.remove_next(current)
      c = list.remove_next(current)

      des = current.data - 1

      loop do
        if des == 0
          des = list.length + 3
        end

        if [a, b, c].include?(des)
          des = des - 1
          next
        end

        break
      end

      node = list.find_fast(des)
      # if node == nil
      #   p list.length
      #   p des
      #   p a, b, c
      #   exit(0)
      # end

      list.insert_next(node, c)
      list.insert_next(node, b)
      list.insert_next(node, a)

      current = list.find_fast(current.data).next
    end

    # list.nice_show(1)

    node = list.find_first { |v| v.data == 1 }

    a = list.remove_next(node)
    b = list.remove_next(node)

    a * b
  end

  def self.solve1(input, len=100)
    list = CircularList.new()

    input.split("").map(&:to_i).each do |n|
      list.insert(n)
    end

    play_all(list, list.head, len)
  end

  def self.solve2(input)
    list = CircularList.new()

    last = nil
    input.split("").map(&:to_i).each do |n|
      last = list.insert(n)
    end

    (10..1000000).each do |n|
      last = list.insert_next(last, n)
    end

    p "------------"
    play_all(list, list.head, 10000000)
  end
end

if __FILE__ == $0
  # puts Day23.solve2("389125467")
  puts Day23.solve2("186524973")
end
