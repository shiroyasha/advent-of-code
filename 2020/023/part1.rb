# rubocop:disable all

module Day23
  def self.nice_display(numbers, round, current)
    round = round % numbers.size
    numbers = numbers.clone

    while numbers[round] != current
      numbers.push(numbers.shift)
    end

    numbers[0...round].join(" ") + " (" + numbers[round].to_s + ") " + numbers[round+1..-1].join(" ")
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

    puts "Move: #{round +1}"
    puts "Cups: #{nice_display(numbers, round, current)}"
    puts "Pick up: #{to_move}"
    puts "Destination: #{des}"

    res = numbers[4..-1] + [current]
    dest = res.find_index { |n| n == des }

    p res
    p dest

    to_move.reverse.each do |n|
      res.insert(dest + 1, n)
    end

    puts ""

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

  def self.solve2(input)

  end
end

if __FILE__ == $0
  puts Day23.solve1("186524973")
  puts Day23.solve2(ARGV[0])
end
