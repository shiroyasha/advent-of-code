# rubocop:disable all

Instruction = Struct.new(:action, :amount)

instructions = File.read(ARGV[0]).split("\n").map { |line| Instruction.new(line[0], line[1..-1].to_i) }

DIRECTIONS = ["E", "N", "W", "S"]

class Pos
  attr_reader :direction
  attr_reader :x
  attr_reader :y

  def initialize(dir, x, y)
    @direction = dir
    @x = x
    @y = y
  end

  def move(direction, amount)
    case direction
    when "N" then @y += amount
    when "S" then @y -= amount
    when "E" then @x += amount
    when "W" then @x -= amount
    end
  end

  def left(angle)
    @direction = DIRECTIONS[(DIRECTIONS.index(@direction) + (angle / 90)) % 4]
  end

  def right(angle)
    @direction = DIRECTIONS[(DIRECTIONS.index(@direction) - (angle / 90)) % 4]
  end

  def manhattan
    @x.abs + @y.abs
  end
end

current = Pos.new("E", 0, 0)

instructions.each do |i|
  case i.action
  when "L"
    current.left(i.amount)
  when "R"
    current.right(i.amount)
  when "F"
    current.move(current.direction, i.amount)
  else
    current.move(i.action, i.amount)
  end

  p current
end

puts current.manhattan
