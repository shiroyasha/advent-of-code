# rubocop:disable all

require 'set'

Vec = Struct.new(:x, :y) do
  def to_s
    "#{x}, #{y}"
  end
end

module Day24
  def self.find_coordinate(instructions)
    pos = Vec.new(0, 0)

    instructions.each do |i|
      case i
      when "e"
        pos.x += 1
      when "w"
        pos.x -= 1
      when "nw"
        pos.y += 1
      when "ne"
        pos.y += 1
        pos.x += 1
      when "sw"
        pos.y -= 1
        pos.x -= 1
      when "se"
        pos.y -= 1
      end
    end

    pos
  end

  def self.adjecent_pos(pos)
    [
      Vec.new(pos.x+1, pos.y),   # e
      Vec.new(pos.x-1, pos.y),   # w
      Vec.new(pos.x, pos.y+1),   # nw
      Vec.new(pos.x+1, pos.y+1), # ne
      Vec.new(pos.x-1, pos.y-1), # sw
      Vec.new(pos.x, pos.y-1),   # se
    ]
  end

  def self.create_map(path)
    lines = File.read(path).split("\n")
    input = lines.map { |l| parse_line(l) }

    map = Set.new()

    input.each do |ins|
      pos = find_coordinate(ins)

      if map.member?(pos)
        map.delete(pos)
      else
        map.add(pos)
      end
    end

    map
  end

  def self.calc_counts(map)
    counts = {}

    map.each do |pos1|
      counts[pos1] ||= 0

      adjecent_pos(pos1).each do |pos2|
        counts[pos2] ||= 0
        counts[pos2] += 1
      end
    end

    counts
  end

  def self.solve1(path)
    map = create_map(path)
    map.size
  end

  def self.solve2(path)
    map = create_map(path)

    puts ""

    100.times.each do |day|
      counts = calc_counts(map)

      counts.each do |pos, c|
        if map.member?(pos)
          if c == 0 || c > 2
            map.delete(pos)
          end
        else
          if c == 2
            map.add(pos)
          end
        end
      end

      puts "Day #{day + 1}: #{map.size}"
    end

    map.size
  end

  def self.parse_line(line)
    res = []

    while line.size > 0
      if line[0] == "e"
        res << "e"
        line = line[1..-1]
        next
      end

      if line[0] == "w"
        res << "w"
        line = line[1..-1]
        next
      end

      if line[0..1] == "nw"
        res << "nw"
        line = line[2..-1]
        next
      end

      if line[0..1] == "ne"
        res << "ne"
        line = line[2..-1]
        next
      end

      if line[0..1] == "sw"
        res << "sw"
        line = line[2..-1]
        next
      end

      if line[0..1] == "se"
        res << "se"
        line = line[2..-1]
        next
      end
    end

    res
  end
end

if __FILE__ == $0
  puts Day24.solve1(ARGV[0])
  puts Day24.solve2(ARGV[0])
end
