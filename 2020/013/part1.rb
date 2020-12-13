# rubocop:disable all

input = File.read(ARGV[0]).split("\n")

start = input[0].to_i
busses = input[1].split(",").reject { |entry| entry == "x" }.map(&:to_i)

el = busses.map { |bus| [bus - start % bus, bus] }.min_by { |el| el[0] }
puts el[0] * el[1]
