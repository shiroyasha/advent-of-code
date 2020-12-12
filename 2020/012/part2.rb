# rubocop:disable all

Instruction = Struct.new(:action, :amount)

class Vec
  attr_accessor :x, :y

  def initialize(x, y)
    @x = x
    @y = y
  end

  def manhattan
    @x.abs + @y.abs
  end

  def add(v)
    @x += v.x
    @y += v.y
  end

  def rotate_left
    @x, @y = -@y, @x
  end

  def rotate_right
    @x, @y = @y, -@x
  end

  def to_s
    "(#{@x}, #{@y})"
  end
end

class Ship
  attr_accessor :waypoint, :position

  def initialize(waypoint, position)
    @waypoint = waypoint
    @position = position
  end

  def run_all(instructions)
    instructions.each do |i|
      run(i)
    end
  end

  def run(instruction)
    case instruction.action
    when "F" then instruction.amount.times { @position.add(@waypoint) }
    when "L" then (instruction.amount / 90).times { @waypoint.rotate_left }
    when "R" then (instruction.amount / 90).times { @waypoint.rotate_right }
    when "N" then @waypoint.add(Vec.new(0, instruction.amount))
    when "S" then @waypoint.add(Vec.new(0, -instruction.amount))
    when "E" then @waypoint.add(Vec.new(instruction.amount, 0))
    when "W" then @waypoint.add(Vec.new(-instruction.amount, 0))
    end
  end
end

instructions = File.read(ARGV[0]).split("\n").map { |line| Instruction.new(line[0], line[1..-1].to_i) }

ship = Ship.new(Vec.new(10, 1), Vec.new(0, 0))
ship.run_all(instructions)

puts ship.position.manhattan
