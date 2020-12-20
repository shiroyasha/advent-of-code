# rubocop:disable all

class Tile
  attr_reader :id
  attr_reader :map
  attr_reader :edges

  def initialize(id, map)
    @id = id
    @map = map
    @edges = find_edges
  end

  def find_edges
    edges = []

    edges << @map[0]
    edges << @map.map { |l| l[-1] }.join("")
    edges << @map[-1]
    edges << @map.map { |l| l[0] }.join("")

    edges
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
  end

  def put(tile)
  end

  def has?(tid)
    @used.include?(tid)
  end

  def set(x, y, tile)
    @used << tile.id
    @field[y+@size/2][x+@size/2] = tile
  end

  def as_map
    res = @field.reject { |l| l.all? { |e| e == nil } }
    res = res.transpose.reject { |l| l.all? { |e| e == nil } }.transpose

    res.map do |line|
      line.map { |t| t.id }
    end
  end
end

class Solver
  def initialize(tiles)
    @tiles = tiles

    @desk = Desk.new(6)
    @desk.set(0, 0, @tiles[0])
  end

  def solve
    tile = next_tile()

    @desk.put(tile)

    @desk.as_map
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

  def self.solve(tiles)
    Solver.new(tiles).solve
  end
end

if __FILE__ == $0
  Day20.solve(Day20.load(ARGV[0]))
end
