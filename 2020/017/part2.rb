# rubocop:disable all
require 'set'

class Grid
  def initialize(layer0)
    @map = Set.new()

    layer0.each.with_index do |line, y|
      line.each.with_index do |el, x|
        if el == '#'
          @map.add({x: x, y: y, z: 0, w: 0})
        end
      end
    end
  end

  def count
    @map.size
  end

  def show
    p @map.to_a
  end

  def nearby(el)
    d3 = [
      {x: el[:x]-1, y: el[:y]-1, z: el[:z]},
      {x: el[:x]+0, y: el[:y]-1, z: el[:z]},
      {x: el[:x]+1, y: el[:y]-1, z: el[:z]},
      {x: el[:x]-1, y: el[:y]+0, z: el[:z]},
      {x: el[:x]+0, y: el[:y]+0, z: el[:z]},
      {x: el[:x]+1, y: el[:y]+0, z: el[:z]},
      {x: el[:x]-1, y: el[:y]+1, z: el[:z]},
      {x: el[:x]+0, y: el[:y]+1, z: el[:z]},
      {x: el[:x]+1, y: el[:y]+1, z: el[:z]},

      {x: el[:x]-1, y: el[:y]-1, z: el[:z]+1},
      {x: el[:x]+0, y: el[:y]-1, z: el[:z]+1},
      {x: el[:x]+1, y: el[:y]-1, z: el[:z]+1},
      {x: el[:x]-1, y: el[:y]+0, z: el[:z]+1},
      {x: el[:x]+0, y: el[:y]+0, z: el[:z]+1},
      {x: el[:x]+1, y: el[:y]+0, z: el[:z]+1},
      {x: el[:x]-1, y: el[:y]+1, z: el[:z]+1},
      {x: el[:x]+0, y: el[:y]+1, z: el[:z]+1},
      {x: el[:x]+1, y: el[:y]+1, z: el[:z]+1},

      {x: el[:x]-1, y: el[:y]-1, z: el[:z]-1},
      {x: el[:x]+0, y: el[:y]-1, z: el[:z]-1},
      {x: el[:x]+1, y: el[:y]-1, z: el[:z]-1},
      {x: el[:x]-1, y: el[:y]+0, z: el[:z]-1},
      {x: el[:x]+0, y: el[:y]+0, z: el[:z]-1},
      {x: el[:x]+1, y: el[:y]+0, z: el[:z]-1},
      {x: el[:x]-1, y: el[:y]+1, z: el[:z]-1},
      {x: el[:x]+0, y: el[:y]+1, z: el[:z]-1},
      {x: el[:x]+1, y: el[:y]+1, z: el[:z]-1},
    ]

    w = el[:w]

    all = d3.map { |e| e.merge(w: w - 1) } + d3.map { |e| e.merge(w: w) } + d3.map { |e| e.merge(w: w + 1) }

    all - [el]
  end

  def transform
    counts = {}
    @map.to_a.each do |el|
      nearby(el).each do |n|
        counts[n] ||= 0
        counts[n] += 1
      end
    end

    new_map = Set.new()
    counts.each do |k, v|
      if @map.member?(k) && v == 2 || v == 3
        new_map.add(k)
      end

      if !@map.member?(k) && v == 3
        new_map.add(k)
      end
    end

    @map = new_map
  end
end

input = File.read(ARGV[0]).split("\n").map { |l| l.split("") }
g = Grid.new(input)

g.transform
p g.count

g.transform
p g.count

g.transform
p g.count

g.transform
p g.count

g.transform
p g.count

g.transform
p g.count

# g.show()
