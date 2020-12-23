# rubocop:disable all

module Day23
  def self.nice_display(numbers, round, current)
    round = round % numbers.size
    numbers = numbers.clone

    while numbers[round] != current
      numbers.push(numbers.shift)
    end

    numbers.each do |n|
      if n == current
        print "(#{n})"
      else
        print " #{n} "
      end
    end
    puts
  end

  def self.play(input, round)
    numbers = input.split("").map(&:to_i)

    current = numbers[0]
    to_move = numbers[1..3]

    des = current - 1
    loop do
      if des == 0
        des = 9
      end

      if to_move.include?(des)
        des = des - 1
        next
      end

      break
    end

    nice_display(numbers, round, current)

    res = numbers[4..-1] + [current]
    dest = res.find_index { |n| n == des }

    to_move.reverse.each do |n|
      res.insert(dest + 1, n)
    end

    res.join("")
  end

  def self.play_all(input, rounds)
    rounds.times.each do |i|
      input = play(input, i)
    end

    numbers = input.split("").map(&:to_i)

    while numbers[0] != 1
      numbers.push(numbers.shift)
    end

    numbers[1..-1].join("")
  end

  def self.solve1(input)
    puts play_all(input, 100)
  end
end

puts Day23.solve1("186524973")
