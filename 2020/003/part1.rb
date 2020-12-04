# rubocop:disable all

map = File.read("input1.txt").split("\n")

length = map[0].size
height = map.size

x = 0
y = 0
treeCount = 0

loop do
  p x, y

  if map[y][x] == '#'
    treeCount += 1
  end

  x = (x+3) % length
  y += 1

  if y == height
    break
  end
end

puts "Result: #{treeCount}"
