# rubocop:disable all
require 'set'

class Tile
  attr_reader :id
  attr_reader :map
  attr_reader :edges

  def initialize(id, map)
    @id = id
    @map = map
  end

  def up
    @map[0]
  end

  def down
    @map[-1]
  end

  def left
    @map.map { |l| l[0] }.join("")
  end

  def right
    @map.map { |l| l[-1] }.join("")
  end

  def rotate
    @map = @map.map { |l| l.split("") }.transpose.map(&:reverse).map { |l| l.join("") }
    self
  end

  def flip
    @map = @map.map { |l| l.split("").reverse.join("") }
    self
  end
end

class Desk
  attr_reader :edges

  def initialize(size)
    @field = []
    @size = size

    size.times do
      line = []

      size.times do
        line << nil
      end

      @field << line
    end

    @used = []
    @edges = Set.new()
  end

  def can_fit?(x, y, tile)
    up = get(x, y+1)
    down = get(x, y-1)
    left = get(x-1, y)
    right = get(x+1, y)

    return false if up && tile.up != up.down
    return false if down && tile.down != down.up
    return false if left && tile.left != left.right
    return false if right && tile.right != right.left

    true
  end

  def has?(tid)
    @used.include?(tid)
  end

  def set(x, y, tile)
    @used << tile.id
    @field[y+@size/2][x+@size/2] = tile

    @edges.delete([x, y])

    @edges.add([x-1, y]) if empty?(x-1, y)
    @edges.add([x+1, y]) if empty?(x+1, y)
    @edges.add([x, y-1]) if empty?(x, y-1)
    @edges.add([x, y+1]) if empty?(x, y+1)
  end

  def empty?(x, y)
    @field[y+@size/2][x+@size/2] == nil
  end

  def get(x, y)
    @field[y+@size/2][x+@size/2]
  end

  def on_the_desk_count
    @used.size
  end

  def as_map
    res = @field.reject { |l| l.all? { |e| e == nil } }
    res = res.transpose.reject { |l| l.all? { |e| e == nil } }.transpose

    res.map do |line|
      line.map { |t| t.id }
    end
  end

  def image
    res = @field.reject { |l| l.all? { |e| e == nil } }
    res = res.transpose.reject { |l| l.all? { |e| e == nil } }.transpose

    res = res.map do |line|
      line.map do |t|
        t.map
      end
    end

    image = []

    res.each do |line|
      line[0].size.times.each do |index|
        image << line.map { |e| e[index] }.join("")
      end
    end

    image
  end
end

class Solver
  def initialize(tiles)
    @tiles = tiles

    @desk = Desk.new(@tiles.size * 3)
    @desk.set(0, 0, @tiles[0])
  end

  def solve
    @tiles.each do |tile|
      next if @desk.has?(tile.id)

      next if add(tile.rotate)
      next if add(tile.rotate)
      next if add(tile.rotate)

      next if add(tile.flip)

      next if add(tile.rotate)
      next if add(tile.rotate)
      next if add(tile.rotate)
    end

    if @desk.on_the_desk_count == @tiles.count
      @desk
    else
      solve
    end
  end

  def add(tile)
    edges = @desk.edges.to_a

    edges.each do |e|
      if @desk.can_fit?(e[0], e[1], tile)
        @desk.set(e[0], e[1], tile)
        return true
      end
    end

    false
  end

  def next_tile
    @tiles.find { |t| !@desk.has?(t.id) }
  end
end

module Day20
  def self.load(path)
    input = File.read(path)

    input.split("\n\n").map do |data|
      lines = data.split("\n")

      id = lines[0][5..-2].to_i
      map = lines[1..-1]

      Tile.new(id, map)
    end
  end

  def self.solve1(tiles)
    solution = Solver.new(tiles).solve.as_map

    puts solution

    p solution[0][0] * solution[0][-1] * solution[-1][0] * solution[-1][-1]
  end

  def self.solve2(tiles)
    solution = Solver.new(tiles).solve.image

    puts solution
  end
end

if __FILE__ == $0
  Day20.solve1(Day20.load(ARGV[0]))
  Day20.solve2(Day20.load(ARGV[0]))
end
