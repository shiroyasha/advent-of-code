# rubocop:disable all

# input = File.read("input0.txt")
input = File.read("input1.txt")

Instruction = Struct.new(:name, :val)

instructions = input.split("\n").map do |ins|
  name, val = ins.split(" ")

  Instruction.new(name, val.to_i)
end

class Program
  def initialize(instructions)
    @instructions = instructions
    @current = 0
    @acc = 0

    @visited = {}
    @exit_code = 255
  end

  def exit_code
    @exit_code
  end

  def run
    # p @current
    # p @instructions[@current]

    if @visited[@current]
      @exit_code = 1
      return @acc
    end

    @visited[@current] = true

    ins = @instructions[@current]

    case ins.name
    when "acc"
      @acc += ins.val

      @current = @current + 1
    when "jmp"
      @current = @current + ins.val

    when "nop"
      @current = @current + 1
    end

    if @current >= @instructions.length
      @exit_code = 0
      @acc
    else
      run()
    end
  end
end

program = Program.new(instructions)
result = program.run()

puts "Result Part1: #{result}"

# ------------------------------------

def part2(instructions)
  instructions.each.with_index do |_, index|
    new_ins = instructions.map(&:clone)

    case new_ins[index].name
    when "jmp"
      new_ins[index].name = "nop"

      program = Program.new(new_ins)
      result = program.run()

      if program.exit_code() == 0
        return result
      end
    when "nop"
      new_ins[index].name = "jmp"

      program = Program.new(new_ins)
      result = program.run()

      if program.exit_code() == 0
        return result
      end
    when "acc"
      next

      puts program.exit_code()
    end
  end
end

puts "Result Part2: #{part2(instructions)}"
