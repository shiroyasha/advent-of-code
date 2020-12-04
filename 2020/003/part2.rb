# rubocop:disable all

def countTrees(map, xSlope, ySlope)
  length = map[0].size
  height = map.size

  x = 0
  y = 0
  treeCount = 0

  loop do
    if map[y][x] == '#'
      treeCount += 1
    end

    x = (x+xSlope) % length
    y += ySlope

    if y >= height
      break
    end
  end

  treeCount
end

map = File.read("input1.txt").split("\n")

treeCounts =
  countTrees(map, 1, 1) *
  countTrees(map, 3, 1) *
  countTrees(map, 5, 1) *
  countTrees(map, 7, 1) *
  countTrees(map, 1, 2)

puts "Result: #{treeCounts}"
