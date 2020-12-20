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

  def map_without_borders
    @map[1..-2].map do |line|
      line[1..-2]
    end
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
        t.map_without_borders
      end
    end

    image = []

    res.reverse.each do |line|
      line[0].size.times.each do |index|
        image << line.map { |e| e[index] }.join("")
      end
    end

    image.reverse
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

class SearchForMonstrer
  attr_reader :image

  def initialize(image)
    @image = image

    @monster = [
      "                  # ",
      "#    ##    ##    ###",
      " #  #  #  #  #  #   "
    ]
  end

  def find
    @image = Tile.new(0, @image).map
    find_on_image()

    @image = Tile.new(0, @image).rotate.map
    find_on_image()

    @image = Tile.new(0, @image).rotate.map
    find_on_image()

    @image = Tile.new(0, @image).rotate.map
    find_on_image()

    @image = Tile.new(0, @image).flip.map
    find_on_image()

    @image = Tile.new(0, @image).rotate.map
    find_on_image()

    @image = Tile.new(0, @image).rotate.map
    find_on_image()

    @image = Tile.new(0, @image).rotate.map
    find_on_image()
  end

  def find_on_image
    (0...@image.size).each do |y|
      (0...@image[0].size).each do |x|
        if is_monster(x, y)
          mark(x, y)
        end
      end
    end
  end

  def is_monster(x, y)
    return false if y > @image.size - @monster.size
    return false if x > @image[0].size - @monster[0].size

    l1 = @image[y][x...x + @monster[0].size].split("")
    l2 = @image[y+1][x...x + @monster[1].size].split("")
    l3 = @image[y+2][x...x + @monster[2].size].split("")

    @monster.size.times.each do |my|
      line = @image[y+my][x...x + @monster[my].size].split("")
      mline = @monster[my]

      line.each.with_index do |value, index|
        next if mline[index] != "#"

        return false if value != mline[index]
      end
    end

    true
  end

  def mark(x, y)
    @monster.each.with_index do |line, my|
      line.split("").each.with_index do |c, mx|
        if c == "#"
          @image[y+my][x+mx] = "O"
        end
      end
    end
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
    img = Solver.new(tiles).solve.image

    s = SearchForMonstrer.new(img)
    s.find

    puts s.image

    res = 0

    s.image.each do |line|
      line.split("").each do |c|
        if c == "#"
          res += 1
        end
      end
    end

    p res
  end
end

if __FILE__ == $0
  Day20.solve1(Day20.load(ARGV[0]))
  Day20.solve2(Day20.load(ARGV[0]))
end
